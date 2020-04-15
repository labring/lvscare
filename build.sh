# build.sh v3.0.2
set -x
COMMIT_SHA1=$(git rev-parse --short HEAD || echo "0.0.0")
go build -o lvscare -mod vendor -ldflags "-X github.com/fanux/lvscare/version.Version=$1 -X github.com/fanux/lvscare/version.Build=${COMMIT_SHA1}" main.go
