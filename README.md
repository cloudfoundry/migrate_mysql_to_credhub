# migrate_mysql_to_credhub

A Go tool to migrate broker data from MySQL to CredHub

This was used in the past for some kind of migration, and no longer maintained.

It is used in the [v5.0 branch of nfs-volume-release](https://github.com/cloudfoundry/nfs-volume-release/blob/6a89b454829a06479c63ab23e92108961f7b777d/.gitmodules#L9). Once this branch is deprecated, it would make sense to archive this repo.

- In 2019, some of the code in [service-broker-store](https://github.com/cloudfoundry/service-broker-store) that this repo depends on [was removed](https://github.com/cloudfoundry/service-broker-store/commit/8ce20271b626105189aaf2768e5c82fdff6807c4) on the basis that it was no longer needed
- In 2019, the errand that in [nfs-volume-release](https://github.com/cloudfoundry/nfs-volume-release) that ran this code [was removed](https://github.com/cloudfoundry/nfs-volume-release/commit/4e27c52f9f3413e51d2f4c972307468d0d566fcb)
