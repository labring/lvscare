# build.sh v3.0.2
set -x
COMMIT_SHA1=$(git rev-parse --short HEAD || echo "0.0.0")
go build -o lvscare -mod vendor -ldflags "-X github.com/fanux/lvscare/cmd.Version=$1 -X github.com/fanux/lvscare/cmd.Githash=$COMMIT_SHA1 -X github.com/fanux/lvscare/cmd.Author=goreleaser" main.go
