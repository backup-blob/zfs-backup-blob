---
title: zfs-backup-blob backup incremental
description: Create a incremental backup of a zfs snapshot
---

Create a incremental backup of a zfs snapshot

```
zfs-backup-blob backup incremental [flags]
```

### Options

```
  -b, --base string   The base snapshot name to base this increment on (Example: pool/vol@snapshot0)
  -h, --help          help for incremental
  -s, --snap string   Snapshot name (Example: pool/vol@snapshot1)
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob backup](/cli/zfs-backup-blob_backup/)	 - Actions related to backups

