package domain

import (
	"os/exec"
	"strings"
)

type Commander = func(name string, arg ...string) *exec.Cmd

type SendParameters struct {
	WithParameters       bool
	PreviousSnapshotName string
	SnapshotName         string
}

type ReceiveParameters struct {
	TargetName string
}

type ListParameters struct {
	Type   []string
	Fields []string
}

type ZfsDriver interface {
	Send(p *SendParameters) *exec.Cmd
	Snapshot(name string) *exec.Cmd
	List(p *ListParameters) *exec.Cmd
	Receive(p *ReceiveParameters) *exec.Cmd
	GetField(fieldName string, zfsEntity string) *exec.Cmd
	SetField(fieldName, value, zfsEntity string) error
	Destroy(zfsEntity string) error
}

type ZfsVolume struct {
	Name      string
	GroupName string
}

type ZfsSnapshot struct {
	Name       string
	VolumeName string
}

type ZfsSnapshotWithType struct {
	ZfsSnapshot
	Type BackupType
}

func NewZfsSnapshot(fullName string) *ZfsSnapshot {
	splitted := strings.Split(fullName, "@")
	if len(splitted) != 2 {
		return nil
	}

	return &ZfsSnapshot{
		Name:       splitted[1],
		VolumeName: splitted[0],
	}
}

func (z ZfsSnapshot) FullName() string {
	return z.VolumeName + "@" + z.Name
}

func (z ZfsSnapshot) NormalizedFullPath() string {
	return z.VolumeName + "/" + z.Name
}
