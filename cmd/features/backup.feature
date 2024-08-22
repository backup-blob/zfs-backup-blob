Feature: Backup Commands

  Background:
    Given no child volume exists under parent zfs-pool2/unencrypted/tmp
    And no child volume exists under parent zfs-pool2/encrypted/tmp
    And a volume with the name zfs-pool2/unencrypted/tmp exists
    And a volume with the name zfs-pool2/encrypted/tmp exists
    And create storage
    And the s3 key id is set at 'AWS_ACCESS_KEY_ID'
    And the s3 key secret is set at 'AWS_SECRET_ACCESS_KEY'
    And the s3 base url is set as placeholder '<path>'
    And the 'AWS_DEFAULT_REGION' env var is set to 'us-east-1'
    And a bucket with the name 'test-bucket' exists
    And the config is loaded from 'example_config.yaml'
    And the placeholders are replaced in config
    And the config is persisted at '<cfgpath>'

  Scenario: Can list backups on remote for volume/fs
    * a volume with the name zfs-pool2/unencrypted/tmp/45 exists
    * they mount volume zfs-pool2/unencrypted/tmp/45
    * a snapshot 'zfs-pool2/unencrypted/tmp/45@test1' exists
    * a env var with key BB_FLAG_FIX_SIZE and value YES is set
    * they execute the cli command 'backup full -s zfs-pool2/unencrypted/tmp/45@test1 -c <cfgpath>'
    * they execute the cli command 'backup list --volume zfs-pool2/unencrypted/tmp/45 -c <cfgpath>'
    * expect stdout to equal
      """
+------------------------------------+------+------+
| KEY                                | TYPE | SIZE |
+------------------------------------+------+------+
| zfs-pool2/unencrypted/tmp/45/test1 | full | 1MB  |
+------------------------------------+------+------+

      """

  Scenario: Can create a full backup and restore it
    * a volume with the name zfs-pool2/unencrypted/tmp/46 exists
    * they mount volume zfs-pool2/unencrypted/tmp/46
    * file in volume zfs-pool2/unencrypted/tmp/46 with name 'foo-file.txt' and size 26mb exists
    * a snapshot 'zfs-pool2/unencrypted/tmp/46@test1' exists
    * they execute the cli command 'backup full -s zfs-pool2/unencrypted/tmp/46@test1 -c <cfgpath>'
    * they execute the cli command 'backup restore -b zfs-pool2/unencrypted/tmp/46/test1 -t zfs-pool2/unencrypted/tmp/hello3 -c <cfgpath> -restoreAll=false'
    * they mount volume zfs-pool2/unencrypted/tmp/hello3
    * a file with name 'foo-file.txt' exists in volume 'zfs-pool2/unencrypted/tmp/hello3'

  Scenario: Can create incremental backup and restore it
    * a volume with the name zfs-pool2/unencrypted/tmp/47 exists
    * they mount volume zfs-pool2/unencrypted/tmp/47
    * file in volume zfs-pool2/unencrypted/tmp/47 with name 'one.txt' and size 26mb exists
    * a snapshot 'zfs-pool2/unencrypted/tmp/47@one' exists
    * file in volume zfs-pool2/unencrypted/tmp/47 with name 'two.txt' and size 14mb exists
    * a snapshot 'zfs-pool2/unencrypted/tmp/47@two' exists
    * they execute the cli command 'backup full -s zfs-pool2/unencrypted/tmp/47@one -c <cfgpath>'
    * they execute the cli command 'backup incremental -b zfs-pool2/unencrypted/tmp/47@one -s zfs-pool2/unencrypted/tmp/47@two -c <cfgpath>'
    * they execute the cli command 'backup restore -b zfs-pool2/unencrypted/tmp/47/one -t zfs-pool2/unencrypted/tmp/47-restored -c <cfgpath> -restoreAll=false'
    * they execute the cli command 'backup restore -b zfs-pool2/unencrypted/tmp/47/two -t zfs-pool2/unencrypted/tmp/47-restored -c <cfgpath> -restoreAll=false'
    * they mount volume zfs-pool2/unencrypted/tmp/47-restored
    * a file with name 'one.txt' exists in volume 'zfs-pool2/unencrypted/tmp/47-restored'
    * a file with name 'two.txt' exists in volume 'zfs-pool2/unencrypted/tmp/47-restored'

  # volume is not mounted due to encryption key needs to be passed manually
  @encrypted
  Scenario: Can create incremental backup and restore on a encrypted volume
    * a volume with the name zfs-pool2/encrypted/tmp/vol1 exists
    * they mount volume zfs-pool2/encrypted/tmp/vol1
    * file in volume zfs-pool2/encrypted/tmp/vol1 with name 'one.txt' and size 26mb exists
    * a snapshot 'zfs-pool2/encrypted/tmp/vol1@one' exists
    * file in volume zfs-pool2/encrypted/tmp/vol1 with name 'two.txt' and size 14mb exists
    * a snapshot 'zfs-pool2/encrypted/tmp/vol1@two' exists
    * they execute the cli command 'backup full -s zfs-pool2/encrypted/tmp/vol1@one -c <cfgpath>'
    * they execute the cli command 'backup incremental -b zfs-pool2/encrypted/tmp/vol1@one -s zfs-pool2/encrypted/tmp/vol1@two -c <cfgpath>'
    * they execute the cli command 'backup restore -b zfs-pool2/encrypted/tmp/vol1/one -t zfs-pool2/unencrypted/tmp/vol1-restored -c <cfgpath> -restoreAll=false'
    * they execute the cli command 'backup restore -b zfs-pool2/encrypted/tmp/vol1/two -t zfs-pool2/unencrypted/tmp/vol1-restored -c <cfgpath> -restoreAll=false'

  Scenario: Can restore all incremental backups
    * a volume with the name zfs-pool2/unencrypted/tmp/48 exists
    * they mount volume zfs-pool2/unencrypted/tmp/48
    * file in volume zfs-pool2/unencrypted/tmp/48 with name 'one.txt' and size 26mb exists
    * a snapshot 'zfs-pool2/unencrypted/tmp/48@one' exists
    * file in volume zfs-pool2/unencrypted/tmp/48 with name 'two.txt' and size 14mb exists
    * a snapshot 'zfs-pool2/unencrypted/tmp/48@two' exists
    * they execute the cli command 'backup full -s zfs-pool2/unencrypted/tmp/48@one -c <cfgpath>'
    * they execute the cli command 'backup incremental -b zfs-pool2/unencrypted/tmp/48@one -s zfs-pool2/unencrypted/tmp/48@two -c <cfgpath>'
    * they execute the cli command 'backup restore -b zfs-pool2/unencrypted/tmp/48/two -t zfs-pool2/unencrypted/tmp/48-restore -c <cfgpath>'
    * they mount volume zfs-pool2/unencrypted/tmp/48-restore
    * a file with name 'one.txt' exists in volume 'zfs-pool2/unencrypted/tmp/48-restore'
    * a file with name 'two.txt' exists in volume 'zfs-pool2/unencrypted/tmp/48-restore'