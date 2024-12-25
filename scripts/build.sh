#!/bin/bash
set -eux

##################################
# クロスコンパイルするスクリプト
##################################

# 定数ライク
APP_NAME=na2me

# 変数ライク
APP_VERSION=v0.0.0  # タグ
APP_COMMIT=00000000 # 短縮ハッシュ

cd `dirname $0`
cd ..

# ================

function is_git_repo {
    echo `git rev-parse --is-inside-work-tree`
}

if [ $(is_git_repo) = "true" ]; then
    APP_VERSION=`git describe --tag --abbrev=0`
    APP_COMMIT=`git rev-parse --short HEAD`
fi

# ================

BUILD_INFO=$(cat <<-EOF
built at `date +%Y-%m-%d`
commit $APP_COMMIT
EOF
          )

# cmd <output> <GOOS> <GOARCH> <CGO>
cmd() {
    output=$1
    goos=$2
    goarch=$3
    cgo=$4

    docker run \
           --rm \
           -v $PWD:/work \
           -w /work \
           --env CGO_ENABLED=$cgo \
           --env GOCACHE=/work/.cache \
           --env GOPATH=/work/.cache \
           --env GOOS=$goos \
           --env GOARCH=$goarch \
           na2mebuilder \
           go build -o build/$output -buildvcs=false -ldflags "-X github.com/kijimaD/na2me/lib/consts.AppVersion=$APP_VERSION -X 'github.com/kijimaD/na2me/lib/consts.BuildInfo=$BUILD_INFO'" ./main.go
}

start() {
    docker build --target builder -t na2mebuilder .

    cmd "${APP_NAME}_linux_amd64" linux amd64 1
    # cmd "${APP_NAME}_linux_arm64" linux arm64 1
    cmd "${APP_NAME}_windows_amd64.exe" windows amd64 0
    cmd "${APP_NAME}_windows_arm64.exe" windows arm64 0
    cmd "${APP_NAME}.wasm" js wasm 0
}

start
