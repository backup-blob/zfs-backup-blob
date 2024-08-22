---
title: zfs-backup-blob group add-volume
description: Add a volume to a group
---

Add a volume to a group

```
zfs-backup-blob group add-volume [flags]
```

### Options

```
  -g, --group string    Volume group name (default "default")
  -h, --help            help for add-volume
      --volume string   Volume name (Example: pool/vol1)
```

### Options inherited from parent commands

```
  -c, --configPath string   Path to the config file (default "~/.bbackup.yaml")
  -l, --logLevel string     LogLevel = debug|disabled (default "disabled")
```

### SEE ALSO

* [zfs-backup-blob group](/cli/zfs-backup-blob_group/)	 - Actions related to volumes/fs in a group

