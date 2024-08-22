package usecase_test

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
)

type TestCase struct {
	name             string
	policy           config.RemoteTrimPolicy
	backupMap        map[string]domain.BackupRecord
	backupsRemaining []string
	backupsDeleted   []string
}

var testCases = []TestCase{
	{
		name:   "delete one full backup due to policy",
		policy: "IIIFF",
		backupMap: map[string]domain.BackupRecord{
			"1": {Type: domain.Full},
			"2": {Type: domain.Full},
			"3": {Type: domain.Full},
			"4": {Type: domain.Incremental, ParentBackupKey: "3"},
			"5": {Type: domain.Incremental, ParentBackupKey: "4"},
			"6": {Type: domain.Incremental, ParentBackupKey: "5"},
			"7": {Type: domain.Incremental, ParentBackupKey: "6"},
		},
		backupsRemaining: []string{"2", "3", "4", "5", "6", "7"},
		backupsDeleted:   []string{"1"},
	},
	{
		name:   "full backup remains due to dependency",
		policy: "III",
		backupMap: map[string]domain.BackupRecord{
			"1": {Type: domain.Full},
			"2": {Type: domain.Full},
			"3": {Type: domain.Full},
			"4": {Type: domain.Incremental, ParentBackupKey: "3"},
			"5": {Type: domain.Incremental, ParentBackupKey: "4"},
			"6": {Type: domain.Incremental, ParentBackupKey: "5"},
			"7": {Type: domain.Incremental, ParentBackupKey: "6"},
		},
		backupsRemaining: []string{"3", "4", "5", "6", "7"},
		backupsDeleted:   []string{"1", "2"},
	},
	{
		name:   "delete incremental and full backup since due to policy",
		policy: "IIIF",
		backupMap: map[string]domain.BackupRecord{
			"1": {Type: domain.Full},
			"2": {Type: domain.Incremental, ParentBackupKey: "1"},
			"3": {Type: domain.Incremental, ParentBackupKey: "2"},
			"4": {Type: domain.Full},
			"5": {Type: domain.Incremental, ParentBackupKey: "4"},
			"6": {Type: domain.Incremental, ParentBackupKey: "5"},
			"7": {Type: domain.Incremental, ParentBackupKey: "6"},
		},
		backupsRemaining: []string{"4", "5", "6", "7"},
		backupsDeleted:   []string{"1", "2", "3"},
	},
}
