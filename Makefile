PHONY: install lint format

install:
	@asdf install
	@go install github.com/google/go-licenses@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/smartystreets/goconvey@latest
	@go install go.uber.org/mock/mockgen@latest
	@make -C docs install

lint:
	@test -z $(gofmt -l .)
	@golangci-lint run ./...
	@make -C docs lint

license-check:
	@go-licenses report ./... --ignore github.com/backup-blob/zfs-backup-blob,golang.org/x/sys/unix

format:
	@go fmt ./...

test:
	@go test ./... -cover -short -count=1

coverage:
	@go test ./... -cover -short -coverprofile=coverage.out
	@go tool cover -html=coverage.out

test-acceptance:
	@go test ./... -cover -run Integration -count=1

gen-mocks:
	@mockgen -source=internal/domain/zfs.go -destination=internal/domain/mocks/zfs.go -package mocks
	@mockgen -source=internal/domain/snapshot.go -destination=internal/domain/mocks/snapshot.go -package mocks
	@mockgen -source=internal/domain/storage.go -destination=internal/domain/mocks/storage.go -package mocks
	@mockgen -source=internal/domain/backup.go -destination=internal/domain/mocks/backup.go -package mocks
	@mockgen -source=internal/domain/volume.go -destination=internal/domain/mocks/volume.go -package mocks
	@mockgen -source=internal/domain/render.go -destination=internal/domain/mocks/render.go -package mocks
	@mockgen -source=internal/domain/backup_state.go -destination=internal/domain/mocks/backup_state.go -package mocks
	@mockgen -source=internal/domain/config/config.go -destination=internal/domain/config/mocks/config.go -package mocks_config
