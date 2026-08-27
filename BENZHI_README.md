# lZ989_8642

基于 Go 实现的Web项目，一款业务流程工具，完成接收、校验、保存、展示。

## Standard commands

```bash
go build ./...
go test -count=1 ./...
```

## Run

```bash
go run ./cmd/server
```

## Frontend

```bash
(cd web && npm run build)
```

## Docker validation

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh my-go-task linux/arm64
./build_benzhi_docker.sh my-go-task linux/amd64
docker run -it my-go-task:latest
```

## Known initial failures

Initial validation failures are retained in the package command output and run logs; they are not copied into the project repository.
