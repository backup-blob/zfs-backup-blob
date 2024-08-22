package driver

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"os/exec"
	"strings"
)

type zfs struct {
	binaryPath string
	cmd        domain.Commander
	logger     domain.LogDriver
}

func NewZfs(binaryPath string, cmd domain.Commander, logger domain.LogDriver) domain.ZfsDriver {
	return &zfs{
		binaryPath: binaryPath,
		cmd:        cmd,
		logger:     logger,
	}
}

func NewZfsFromConfig(conf *config.ZfsConfig, logger domain.LogDriver) domain.ZfsDriver {
	binaryPath := "zfs"
	defaultCmd := exec.Command

	if conf.ZfsPath != "" {
		binaryPath = conf.ZfsPath
	}

	return NewZfs(binaryPath, defaultCmd, logger)
}

func (z *zfs) Send(p *domain.SendParameters) *exec.Cmd {
	var args []string
	args = append(args, "send", "--raw")

	if p.WithParameters {
		args = append(args, "-p")
	}

	if p.PreviousSnapshotName != "" {
		args = append(args, "-I", p.PreviousSnapshotName)
	}

	args = append(args, p.SnapshotName)

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	return z.cmd(z.binaryPath, args...)
}

func (z *zfs) Receive(p *domain.ReceiveParameters) *exec.Cmd {
	var args []string

	doNotMount := "-u"
	args = append(args, "receive", doNotMount, p.TargetName)

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	return z.cmd(z.binaryPath, args...)
}

func (z *zfs) List(p *domain.ListParameters) *exec.Cmd {
	// TODO: allow filter by volume
	var args []string
	args = append(args, "list") //nolint:gocritic // more readable with append
	args = append(args, "-t", strings.Join(p.Type, ","))
	args = append(args, "-o", strings.Join(p.Fields, ","))

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	return z.cmd(z.binaryPath, args...)
}

func (z *zfs) Snapshot(name string) *exec.Cmd {
	var args []string
	args = append(args, "snapshot", name)

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	return z.cmd(z.binaryPath, args...)
}

func (z *zfs) SetField(fieldName, value, zfsEntity string) error {
	args := []string{"set", fieldName + "=" + value, zfsEntity}

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	cmd := z.cmd(z.binaryPath, args...)

	return cmd.Run()
}

func (z *zfs) GetField(fieldName, zfsEntity string) *exec.Cmd {
	args := []string{"get", "-o", "value", "-H", fieldName, zfsEntity}

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	return z.cmd(z.binaryPath, args...)
}

func (z *zfs) Destroy(zfsEntity string) error {
	args := []string{"destroy", zfsEntity}

	z.logger.Debugf("invoke %s %v", z.binaryPath, args)

	cmd := z.cmd(z.binaryPath, args...)

	return cmd.Run()
}
