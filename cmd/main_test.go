package main_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

const TestingKey = "t"

func TestFeaturesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:         "pretty",
			Paths:          []string{"features"},
			TestingT:       t, // Testing instance that will run subtests.
			DefaultContext: context.WithValue(context.Background(), TestingKey, t),
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func setEnvVar(ctx context.Context, key string, value string) error {
	v, ok := ctx.Value(TestingKey).(*testing.T)
	if !ok {
		return errors.New("failed to read testing.T")
	}
	v.Setenv(key, value)
	return nil
}

func checkFile(fileName string, volumeName string) error {
	volumePath, err := getVolumePath(volumeName)
	if err != nil {
		return err
	}
	_, err = os.Stat(volumePath + "/" + fileName)
	if err != nil {
		return err
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`a env var with key ([A-Za-z-\/0-9\_]+) and value (.*) is set$`, setEnvVar)
	sc.Step(`a file with name \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\' exists in volume \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\'`, checkFile)
	InitZfs(sc)
	InitStorage(sc)
	InitCmd(sc)
	InitConfig(sc)
}
