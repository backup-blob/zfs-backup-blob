---
title: zfs-backup-blob backup restore
description: Restore a backup
---

Restore a backup

```
zfs-backup-blob backup restore [flags]
```

### Options

```
  -b, --blob-key string   S3 Key to the backup to restore (excluding prefix)
  -h, --help              help for restore
  -r, --restoreAll        Restores all incremental snapshots including the full backup
  -t, --target string     Path to zfs pool/<volume/filesystem>
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob backup](/cli/zfs-backup-blob_backup/)	 - Actions related to backups

