package main_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/cmd/command"
	"github.com/cucumber/godog"
	"os"
	"strings"
)

func execCommand(ctx context.Context, args string) error {
	cmd := command.RootCmd
	buf := bytes.NewBuffer([]byte{})

	argsReplaced, err := replacePlaceholder(ctx, args)
	if err != nil {
		return err
	}

	cmd.SetOut(buf)
	cmd.SetErr(os.Stderr)
	cmd.SetArgs(strings.Split(argsReplaced, " "))

	errE := cmd.Execute()
	if errE != nil {
		return errE
	}

	setPlaceholder(ctx, "<exec_out>", buf.String())

	return nil
}

func expectStdout(ctx context.Context, expected string) error {
	actual, err := replacePlaceholder(ctx, "<exec_out>")
	if err != nil {
		return err
	}

	if expected != actual {
		return fmt.Errorf("%s != %s", expected, actual)
	}
	return nil
}

func InitCmd(sc *godog.ScenarioContext) {
	sc.Step(`expect stdout to equal`, expectStdout)
	sc.Step(`they execute the cli command \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\'`, execCommand)
}
