---
title: zfs-backup-blob group snapshot
description: Create snapshots for all zfs volume/fs belonging to a group
---

Create snapshots for all zfs volume/fs belonging to a group

```
zfs-backup-blob group snapshot [flags]
```

### Options

```
  -g, --group string   Volume group name (default "default")
  -h, --help           help for snapshot
  -t, --type string    Backup type (full|incremental)
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob group](/cli/zfs-backup-blob_group/)	 - Actions related to volumes/fs in a group

