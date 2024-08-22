---
title: zfs-backup-blob group
description: Actions related to volumes/fs in a group
---

Actions related to volumes/fs in a group

### Options

```
  -h, --help   help for group
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob](/cli/zfs-backup-blob/)	 - _
* [zfs-backup-blob group add-volume](/cli/zfs-backup-blob_group_add-volume/)	 - Add a volume to a group
* [zfs-backup-blob group list-volumes](/cli/zfs-backup-blob_group_list-volumes/)	 - List all volumes which belong to a group
* [zfs-backup-blob group snapshot](/cli/zfs-backup-blob_group_snapshot/)	 - Create snapshots for all zfs volume/fs belonging to a group
* [zfs-backup-blob group sync](/cli/zfs-backup-blob_group_sync/)	 - Sync snapshots of a group to remote
* [zfs-backup-blob group trim-local](/cli/zfs-backup-blob_group_trim-local/)	 - Trim local snapshots of a group
* [zfs-backup-blob group trim-remote](/cli/zfs-backup-blob_group_trim-remote/)	 - Trim remote backups of a group

