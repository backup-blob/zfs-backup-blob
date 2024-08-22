---
title: zfs-backup-blob group trim-local
description: Trim local snapshots of a group
---

Trim local snapshots of a group

```
zfs-backup-blob group trim-local [flags]
```

### Options

```
  -d, --dry-run        Don't actually remove backup blob
  -g, --group string   Volume group name (default "default")
  -h, --help           help for trim-local
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob group](/cli/zfs-backup-blob_group/)	 - Actions related to volumes/fs in a group

