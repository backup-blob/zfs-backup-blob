Feature: Group Commands

  Background:
    Given no child volume exists under parent zfs-pool2/unencrypted/tmp
    And no child volume exists under parent zfs-pool2/encrypted/tmp
    And a volume with the name zfs-pool2/unencrypted/tmp exists
    And a volume with the name zfs-pool2/encrypted/tmp exists
    And create storage
    And the s3 key id is set at 'AWS_ACCESS_KEY_ID'
    And the 'AWS_DEFAULT_REGION' env var is set to 'us-east-1'
    And the s3 key secret is set at 'AWS_SECRET_ACCESS_KEY'
    And the s3 base url is set as placeholder '<path>'
    And a bucket with the name 'test-bucket' exists
    And the config is loaded from 'example_config_trim.yaml'
    And the placeholders are replaced in config
    And the config is persisted at '<cfgpath>'

  Scenario: The one where a full backup is deleted
    * a volume with the name zfs-pool2/unencrypted/tmp/50 exists
    * a volume with the name zfs-pool2/unencrypted/tmp/50 has field backup_blob::group with value group2
    * they mount volume zfs-pool2/unencrypted/tmp/50
    * a env var with key BB_FLAG_TIME and value 1300000000 is set
    * they execute the cli command 'group snapshot -g group2 -t full -c <cfgpath>'
    * a env var with key BB_FLAG_TIME and value 1400000000 is set
    * they execute the cli command 'group snapshot -g group2 -t full -c <cfgpath>'
    * a env var with key BB_FLAG_TIME and value 1500000000 is set
    * they execute the cli command 'group snapshot -g group2 -t incremental -c <cfgpath>'
    * they execute the cli command 'group sync -g group2 -c <cfgpath> -l debug'
    * they execute the cli command 'group trim-remote -g group2 -c <cfgpath> -l debug'
    * the key 'zfs-pool2/unencrypted/tmp/50/backup_blob_2011-03-13T07-06-40Z' does not exists in bucket 'test-bucket'
    * the key 'zfs-pool2/unencrypted/tmp/50/backup_blob_2014-05-13T16-53-20Z' exists in bucket 'test-bucket'
    * the key 'zfs-pool2/unencrypted/tmp/50/backup_blob_2017-07-14T02-40-00Z' exists in bucket 'test-bucket'

  Scenario: The one where no backup is deleted due to dry-run flag
    * a volume with the name zfs-pool2/unencrypted/tmp/50 exists
    * a volume with the name zfs-pool2/unencrypted/tmp/50 has field backup_blob::group with value group2
    * they mount volume zfs-pool2/unencrypted/tmp/50
    * a env var with key BB_FLAG_TIME and value 1300000000 is set
    * they execute the cli command 'group snapshot -g group2 -t full -c <cfgpath>'
    * a env var with key BB_FLAG_TIME and value 1400000000 is set
    * they execute the cli command 'group snapshot -g group2 -t full -c <cfgpath>'
    * a env var with key BB_FLAG_TIME and value 1500000000 is set
    * they execute the cli command 'group snapshot -g group2 -t incremental -c <cfgpath>'
    * they execute the cli command 'group sync -g group2 -c <cfgpath> -l debug'
    * they execute the cli command 'group trim-remote -d -g group2 -c <cfgpath> -l debug'
    * the key 'zfs-pool2/unencrypted/tmp/50/backup_blob_2011-03-13T07-06-40Z' exists in bucket 'test-bucket'
    * the key 'zfs-pool2/unencrypted/tmp/50/backup_blob_2014-05-13T16-53-20Z' exists in bucket 'test-bucket'
    * the key 'zfs-pool2/unencrypted/tmp/50/backup_blob_2017-07-14T02-40-00Z' exists in bucket 'test-bucket'

  Scenario: Can create snapshots of volumes in a group
    * a volume with the name zfs-pool2/unencrypted/tmp/49 exists
    * a volume with the name zfs-pool2/unencrypted/tmp/49 has field backup_blob::group with value group1
    * a env var with key BB_FLAG_TIME and value 1405544146 is set
    * they execute the cli command 'group snapshot -g group1 -t full -c <cfgpath> -l debug'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/49@backup_blob_2014-07-16T20-55-46Z' has with field 'backup_blob::group' and value 'group1'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/49@backup_blob_2014-07-16T20-55-46Z' has with field 'backup_blob::type' and value 'full'

  Scenario: Can sync snapshots of a group to remote
    * a volume with the name zfs-pool2/unencrypted/tmp/54 exists
    * a volume with the name zfs-pool2/unencrypted/tmp/54 has field backup_blob::group with value group3
    * they mount volume zfs-pool2/unencrypted/tmp/54
    * a env var with key BB_FLAG_TIME and value 1400000000 is set
    * they execute the cli command 'group snapshot -g group3 -t full -c <cfgpath> -l debug'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/54@backup_blob_2014-05-13T16-53-20Z' has with field 'backup_blob::group' and value 'group3'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/54@backup_blob_2014-05-13T16-53-20Z' has with field 'backup_blob::type' and value 'full'
    * a env var with key BB_FLAG_TIME and value 1500000000 is set
    * file in volume zfs-pool2/unencrypted/tmp/54 with name 'foo-file.txt' and size 26mb exists
    * they execute the cli command 'group snapshot -g group3 -t incremental -c <cfgpath>'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/54@backup_blob_2017-07-14T02-40-00Z' has with field 'backup_blob::group' and value 'group3'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/54@backup_blob_2017-07-14T02-40-00Z' has with field 'backup_blob::type' and value 'incremental'
    * they execute the cli command 'group sync -g group3 -c <cfgpath> -l debug'
    * the key 'zfs-pool2/unencrypted/tmp/54/backup_blob_2014-05-13T16-53-20Z' exists in bucket 'test-bucket'
    * the key 'zfs-pool2/unencrypted/tmp/54/backup_blob_2017-07-14T02-40-00Z' exists in bucket 'test-bucket'

  Scenario: The one where volumes within a group are listed
    * a volume with the name zfs-pool2/unencrypted/tmp/51 exists
    * a volume with the name zfs-pool2/unencrypted/tmp/51 has field backup_blob::group with value group1
    * a volume with the name zfs-pool2/unencrypted/tmp/52 exists
    * they execute the cli command 'group list-volumes -c <cfgpath>'
    * expect stdout to equal
      """
+------------------------------+--------+
| VOLUME                       | GROUP  |
+------------------------------+--------+
| zfs-pool2/unencrypted/tmp/51 | group1 |
+------------------------------+--------+

      """

  Scenario: The one where a volume is added to a group
    * a volume with the name zfs-pool2/unencrypted/tmp/52 exists
    * a env var with key BB_FLAG_TIME and value 1400000000 is set
    * they execute the cli command 'group add-volume --volume zfs-pool2/unencrypted/tmp/52 -g group1 -c <cfgpath>'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/52@backup_blob_2014-05-13T16-53-20Z' has with field 'backup_blob::group' and value 'group1'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/52@backup_blob_2014-05-13T16-53-20Z' has with field 'backup_blob::type' and value 'full'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/52' has with field 'backup_blob::group' and value 'group1'

  Scenario: The one where local snapshots are trimmed
    * a volume with the name zfs-pool2/unencrypted/tmp/53 exists
    * a volume with the name zfs-pool2/unencrypted/tmp/53 has field backup_blob::group with value group1
    * a env var with key BB_FLAG_TIME and value 1405544146 is set
    * they execute the cli command 'group snapshot -g group1 -t full -c <cfgpath> -l debug'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/53@backup_blob_2014-07-16T20-55-46Z' has with field 'backup_blob::group' and value 'group1'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/53@backup_blob_2014-07-16T20-55-46Z' has with field 'backup_blob::type' and value 'full'
    * a env var with key BB_FLAG_TIME and value 1406544146 is set
    * they execute the cli command 'group snapshot -g group1 -t full -c <cfgpath> -l debug'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/53@backup_blob_2014-07-28T10-42-26Z' has with field 'backup_blob::group' and value 'group1'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/53@backup_blob_2014-07-28T10-42-26Z' has with field 'backup_blob::type' and value 'full'
    * they execute the cli command 'group trim-local -g group1 -c <cfgpath> -l debug'
    * the zfs entity 'zfs-pool2/unencrypted/tmp/53@backup_blob_2014-07-16T20-55-46Z' does not exist