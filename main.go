package main

import (
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerflags"
	"code.cloudfoundry.org/migrate_mysql_to_credhub/migrator"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims"
	"github.com/go-sql-driver/mysql"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	DBDriver string `long:"dbDriver" description:"Database driver name when using SQL to store broker state" required:"true"`

	DBHostname string `long:"dbHostname" description:"Database hostname when using SQL to store broker state" required:"true"`

	DBPort string `long:"dbPort" description:"Database port when using SQL to store broker state" required:"true"`

	DBName string `long:"dbName" description:"Database name when using SQL to store broker state" required:"true"`

	DBUsername string `long:"dbUsername" description:"Database username when using SQL to store broker state" required:"true"`

	DBPassword string `long:"dbPassword" description:"Database password when using SQL to store broker state" required:"true"`

	DBCACertPath string `long:"dbCACertPath" description:"Path to CA Cert for database SSL connection"`

	DBSkipHostnameValidation bool `long:"dbSkipHostnameValidation" description:"Skip DB server hostname validation when connecting over TLS"`

	CredhubURL string `long:"credhubURL" description:"CredHub server URL when using CredHub to store broker state" required:"true"`

	CredhubCACertPath string `long:"credhubCACertPath" description:"Path to CA Cert for CredHub"`

	UAAClientID string `long:"uaaClientID" description:"UAA client ID when using CredHub to store broker state" required:"true"`

	UAAClientSecret string `long:"uaaClientSecret" description:"UAA client secret when using CredHub to store broker state" required:"true"`

	UAACACertPath string `long:"uaaCACertPath" description:"Path to CA Cert for UAA used for CredHub authorization"`

	StoreID string `long:"storeID" description:"Store ID used to namespace instance details and bindings (credhub only)" required:"true"`

	MinLogLevel string `long:"logLevel" default:"info" description:"Log level: debug, info, error or fatal"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		panic(err)
	}

	logger, _ := lagerflags.NewFromConfig("migrate_mysql_to_credhub", lagerflags.LagerConfig{LogLevel: opts.MinLogLevel})
	logger.Info("migrating")
	defer logger.Info("ends")

	var dbCACert string
	if opts.DBCACertPath != "" {
		b, err := ioutil.ReadFile(opts.DBCACertPath)
		if err != nil {
			logger.Fatal("cannot-read-db-ca-cert", err, lager.Data{"path": opts.DBCACertPath})
		}
		dbCACert = string(b)
	}

	var credhubCACert string
	if opts.CredhubCACertPath != "" {
		b, err := ioutil.ReadFile(opts.CredhubCACertPath)
		if err != nil {
			logger.Fatal("cannot-read-credhub-ca-cert", err, lager.Data{"path": opts.CredhubCACertPath})
		}
		credhubCACert = string(b)
	}

	var uaaCACert string
	if opts.UAACACertPath != "" {
		b, err := ioutil.ReadFile(opts.UAACACertPath)
		if err != nil {
			logger.Fatal("cannot-read-credhub-ca-cert", err, lager.Data{"path": opts.UAACACertPath})
		}
		uaaCACert = string(b)
	}

	dbStore, err := brokerstore.NewSqlStore(
		logger,
		opts.DBDriver,
		opts.DBUsername,
		opts.DBPassword,
		opts.DBHostname,
		opts.DBPort,
		opts.DBName,
		dbCACert,
		opts.DBSkipHostnameValidation,
	)
	if err != nil {
		if HandleSQLStoreError(err) != nil {
			logger.Fatal("failed-to-initialize-sql-store", err)
		}

		logger.Info("missing-sql-database")
		return
	}

	credhubShim, err := credhub_shims.NewCredhubShim(
		opts.CredhubURL,
		credhubCACert,
		opts.UAAClientID,
		opts.UAAClientSecret,
		uaaCACert,
		&credhub_shims.CredhubAuthShim{},
	)
	if err != nil {
		logger.Fatal("failed-to-create-credhub-shim", err)
	}
	credhubStore := brokerstore.NewCredhubStore(
		logger,
		credhubShim,
		opts.StoreID,
	)
	if err != nil {
		logger.Fatal("failed-to-initialize-credhub-store", err)
	}

	migrator := migrator.NewMigrator(logger)
	err = migrator.Migrate(dbStore, credhubStore)
	if err != nil {
		logger.Fatal("failed-to-migrate", err)
	}
}

func HandleSQLStoreError(err error) error {
	if err == nil {
		return nil
	}

	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1049 {
			return nil
		}
	}

	return err
}
