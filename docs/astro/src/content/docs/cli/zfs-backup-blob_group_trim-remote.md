---
title: zfs-backup-blob group trim-remote
description: Trim remote backups of a group
---

Trim remote backups of a group

```
zfs-backup-blob group trim-remote [flags]
```

### Options

```
  -d, --dry-run        Don't actually remove backup blob
  -g, --group string   Volume group name (default "default")
  -h, --help           help for trim-remote
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob group](/cli/zfs-backup-blob_group/)	 - Actions related to volumes/fs in a group

