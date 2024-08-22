---
title: FAQ
description: Frequently asked questions
---

### How to execute the cli at certain times

It's recommended to run the cli with the
[linux crond daemon](https://www.redhat.com/sysadmin/automate-linux-tasks-cron).

### How to make sure the state file for a volume is not corrupted.

Make sure zfs-backup-blob is not executed concurrently (for example 2 cronjobs operate on the same group in parallel).

You can build a DIY locking into your cronjob if needed.

### An error occurred while uploading to the blob storage. What to do?

Retry the upload. Retries are safe operations, they continue from the point where the error occurred.

It will restart the upload from zero, which is wasteful but since an upload error should be rare this should be
acceptable.

### How to have different configurations for different groups.

Assuming you have two groups of volumes which require a different frequency and a different blob storage as target.

With the `--configPath` on every cli command you are specify the config you want to use.