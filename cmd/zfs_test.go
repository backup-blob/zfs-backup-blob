package main_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"os"
	"os/exec"
	"strings"
)

func createVolume(ctx context.Context, name string) (context.Context, error) {
	cmd := exec.Command("zfs", "create", "-u", name)
	cmd.Stderr = os.Stderr
	return ctx, cmd.Run()
}

func setField(ctx context.Context, name string, fieldName string, value string) (context.Context, error) {
	args := []string{"set", fieldName + "=" + value, name}
	cmd := exec.Command("zfs", args...)
	cmd.Stderr = os.Stderr
	return ctx, cmd.Run()
}

func mountVolume(ctx context.Context, name string) (context.Context, error) {
	cmd := exec.Command("zfs", "mount", name)
	cmd.Stderr = os.Stderr
	return ctx, cmd.Run()
}

func cleanVolumes(name string) error {
	cmd := exec.Command("zfs", "destroy", "-r", name)
	cmd.Stderr = os.Stderr
	cmd.Run() // ignore error
	return nil
}

func createSnapshot(name string) error {
	cmd := exec.Command("zfs", "snapshot", name)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func assertZfsField(zfsEntity string, fieldName string, value string) error {
	actualVal, err := getZfsField(fieldName, zfsEntity)
	if err != nil {
		return err
	}
	if actualVal != value {
		return fmt.Errorf("%s does not match %s", actualVal, value)
	}
	return nil
}

func datasetExists(name string) (bool, error) {
	cmd := exec.Command("zfs", "get", "-p", "all", name)
	buff := bytes.NewBuffer([]byte{})
	cmd.Stderr = buff
	cmd.Stdout = buff
	cmd.Run()

	return !strings.Contains(buff.String(), "dataset does not exist"), nil
}

func assertDatasetMissing(name string) error {
	exists, err := datasetExists(name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New(fmt.Sprintf("dataset %s exists, but should not", name))
	}

	return nil
}

func getZfsField(fieldName string, zfsEntity string) (string, error) {
	cmd := exec.Command("zfs", "get", "-o", "value", "-H", fieldName, zfsEntity)
	cmd.Stderr = os.Stderr
	buff := bytes.NewBuffer([]byte{})
	cmd.Stdout = buff
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	val := strings.ReplaceAll(buff.String(), "\n", "")
	return val, nil
}

func getVolumePath(volumeName string) (string, error) {
	return getZfsField("mountpoint", volumeName)
}

func createFileInVolume(volumeName string, fileName string, size int) error {
	volumePath, err := getVolumePath(volumeName)
	if err != nil {
		return err
	}
	f, err := os.Create(volumePath + "/" + fileName)
	if err != nil {
		return err
	}

	if err := f.Truncate(int64(size * 1024 * 1024)); err != nil {
		return err
	}
	return nil
}

func InitZfs(sc *godog.ScenarioContext) {
	sc.Step(`file in volume ([a-z-\/0-9]+) with name '([a-z-\/0-9\.]+)' and size (\d+)mb exists`, createFileInVolume)
	sc.Step(`the zfs entity \'(.*)\' has with field \'(.*)\' and value \'(.*)\'$`, assertZfsField)
	sc.Step(`the zfs entity \'(.*)\' does not exist`, assertDatasetMissing)
	sc.Step(`a snapshot \'([a-z-\/0-9@]+)\' exists$`, createSnapshot)
	sc.Step(`no child volume exists under parent ([\+\s<>_a-zA-Z-\/0-9@\.]+)$`, cleanVolumes)
	sc.Step(`they mount volume ([a-z-\/0-9]+)$`, mountVolume)
	sc.Step(`a volume with the name ([\+\s<>_a-zA-Z-\/0-9@\.]+) has field ([\+\s<>_a-zA-Z-\/0-9@\.\:]+) with value ([a-z-\/0-9]+)$`, setField)
	sc.Step(`a volume with the name ([\+\s<>_a-zA-Z-\/0-9@\.]+) exists$`, createVolume)
}
