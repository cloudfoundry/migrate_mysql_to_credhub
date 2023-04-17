# migrate_mysql_to_credhub

An (arvhived) Go tool to migrate broker data from MySQL to CredHub

This was used in the past for some kind of migration, but it's no longer in use, and no longer maintained.

- In 2019, some of the code in [service-broker-store](https://github.com/cloudfoundry/service-broker-store) that this repo depends on [was removed](https://github.com/cloudfoundry/service-broker-store/commit/8ce20271b626105189aaf2768e5c82fdff6807c4) on the basis that it was no longer needed
- In 2019, the errand that in [nfs-volume-release](https://github.com/cloudfoundry/nfs-volume-release) that ran this code [was removed](https://github.com/cloudfoundry/nfs-volume-release/commit/4e27c52f9f3413e51d2f4c972307468d0d566fcb)
