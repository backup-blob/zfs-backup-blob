package driver_test

import (
	"errors"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestSpecConfig(t *testing.T) {

	Convey("Given the Load function is called", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should load config", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureConfig))

				So(err, ShouldBeNil)
				So(confLoaded.GetMiddlewares(), ShouldNotBeEmpty)
				So(confLoaded.GetStorageDriver(), ShouldNotBeNil)
				So(confLoaded.GetZfsDriver(), ShouldNotBeNil)
				So(confLoaded.GetConfig(), ShouldNotBeNil)
			})
		})
		Convey("When Unmarshal fails", func() {
			Convey("It should error", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureConfigErrUnmarshal))

				So(err.Error(), ShouldContainSubstring, "configRepo.loadFileConfig failed unmarshalFile failed yaml")
				So(confLoaded, ShouldBeNil)
			})
		})
		Convey("When the stage type cannot be extracted", func() {
			Convey("It should error", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureFaultyStageType))

				So(err.Error(), ShouldContainSubstring, "config type for stage: middleware not found")
				So(confLoaded, ShouldBeNil)
			})
		})
		Convey("When the remote is not found", func() {
			Convey("It should error", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureConfigRemoteNotFound))

				So(err.Error(), ShouldContainSubstring, "remote not found: middleware2")
				So(confLoaded, ShouldBeNil)
			})
		})
		Convey("When middleware is invalid", func() {
			Convey("It should error", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureConfigInvalidMiddleware))

				So(err.Error(), ShouldContainSubstring, "invalid middleware on stage index 1")
				So(confLoaded, ShouldBeNil)
			})
		})
		Convey("When the stage cannot be decoded", func() {
			Convey("It should error", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureFaultyStage))

				So(err.Error(), ShouldContainSubstring, "failed to read stage: middleware")
				So(confLoaded, ShouldBeNil)
			})
		})
		Convey("When loading of the reader fails", func() {
			Convey("It should error", func() {
				params := getLoadParams(fixtureConfig)
				params.ConfigReader = &errorReader{err: errors.New("error")}
				confLoaded, err := driver.NewConfigDriver(params)

				So(err, ShouldBeError, "configRepo.loadFileConfig failed error")
				So(confLoaded, ShouldBeNil)
			})
		})
		Convey("When required fields are missing", func() {
			Convey("It should error", func() {
				confLoaded, err := driver.NewConfigDriver(getLoadParams(fixtureMissingRequiredField))
				So(err, ShouldBeError, "validation failed: Key: 'Config.Stages[middleware].Password' Error:Field validation for 'Password' failed on the 'required' tag")
				So(confLoaded, ShouldBeNil)
			})
		})
	})
}

func getLoadParams(configFixture string) *config.LoadParams {
	return &config.LoadParams{
		ConfigReader: strings.NewReader(configFixture),
		StageMapping: map[string]func() config.ConfigStage{
			"source":     config.NewZfsConfig,
			"sink":       config.NewS3Config,
			"middleware": config.NewCryptConfig,
			"throttle":   config.NewThrottleConfig,
		},
		StorageDriverFunc: func(s *config.S3Config) (domain.StorageDriver, error) {
			return mockStorageDriver{}, nil
		},
		ZfsDriverFunc: func(c *config.ZfsConfig) domain.ZfsDriver {
			return mockZfsDriver{}
		},
		ToMiddleware: func(c config.ConfigStage) domain.Middleware {
			switch v := c.(type) {
			case *config.CryptConfig:
				return driver.NewCrypt(v)
			default:
				return nil
			}
		},
	}
}

type mockStorageDriver struct {
	domain.StorageDriver
}

type mockZfsDriver struct {
	domain.ZfsDriver
}

var fixtureConfig = `
stages:
  sink:
    type: sink
    bucket: test-bucket
    baseEndpoint: <path>
    usePathStyle: true
  middleware:
    type: middleware
    password: "hello"
    remote: sink
  source:
    type: source
    remote: middleware
`

var fixtureConfigInvalidMiddleware = `
stages:
  sink:
    type: sink
    bucket: test-bucket
    baseEndpoint: <path>
    usePathStyle: true
  middleware:
    type: throttle
    writeSpeed: 10
    readSpeed: 10
    remote: sink
  source:
    type: source
    remote: middleware
`

var fixtureConfigRemoteNotFound = `
stages:
  sink:
    type: sink
    bucket: test-bucket
    baseEndpoint: <path>
    usePathStyle: true
  middleware:
    type: middleware
    password: "hello"
    remote: sink
  source:
    type: source
    remote: middleware2
`

var fixtureFaultyStageType = `
stages:
  sink:
    type: sink
    bucket: test-bucket
    baseEndpoint: <path>
    usePathStyle: true
  middleware:
    typed: middleware
    password: "hello"
    remote: sink
  source:
    type: source
    remote: middleware
`

var fixtureFaultyStage = `
stages:
  sink:
    type: sink
    bucket: test-bucket
    baseEndpoint: <path>
    usePathStyle: true
  middleware: []
  source:
    type: source
    remote: middleware
`

var fixtureMissingRequiredField = `
stages:
  sink:
    type: sink
    bucket: test-bucket
    baseEndpoint: <path>
    usePathStyle: true
  middleware:
    type: middleware
    remote: sink
  source:
    type: source
    remote: middleware
`

var fixtureConfigErrUnmarshal = `
$$$$
`

type errorReader struct {
	err error
}

// Read implements the io.Reader interface.
func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
