package main

import (
	"os"

	"github.com/jessevdk/go-flags"

	"code.cloudfoundry.org/lager/lagerflags"
)

var opts struct {
	DBDriver string `long:"dbDriver" description:"Database driver name when using SQL to store broker state" required:"true"`

	DBHostname string `long:"dbHostname" description:"Database hostname when using SQL to store broker state" required:"true"`

	DBPort string `long:"dbPort" description:"Database port when using SQL to store broker state" required:"true"`

	DBName string `long:"dbName" description:"Database name when using SQL to store broker state" required:"true"`

	DBUsername string `long:"dbUsername" description:"Database username when using SQL to store broker state" required:"true"`

	DBPassword string `long:"dbPassword" description:"Database password when using SQL to store broker state" required:"true"`

	DBCACertPath string `long:"dbCACertPath" description:"Path to CA Cert for database SSL connection"`

	CredhubURL string `long:"credhubURL" description:"CredHub server URL when using CredHub to store broker state" required:"true"`

	CredhubCACertPath string `long:"credhubCACertPath" description:"Path to CA Cert for CredHub"`

	UAAClientID string `long:"uaaClientID" description:"UAA client ID when using CredHub to store broker state" required:"true"`

	UAAClientSecret string `long:"uaaClientSecret" description:"UAA client secret when using CredHub to store broker state" required:"true"`

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

}
