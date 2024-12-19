#!/bin/bash
set -e

supported_platforms=(
    darwin-amd64
    darwin-arm64
    linux-386
    linux-amd64
    linux-arm
    linux-arm64
    windows-386
    windows-amd64
    windows-arm64
)

VERSION="$1"
BUILD_DATE="$(date -u "+%Y-%m-%d %H:%M:%S UTC")"

echo "VERSION=${VERSION} BUILD_DATE=${BUILD_DATE}"

for ((j = 0; j < ${#supported_platforms[@]}; j++)); do
    p="${supported_platforms[$j]}"
    goos="${p%-*}"
    goarch="${p#*-}"

    ext=""
    if [ "$goos" = "windows" ]; then
        ext=".exe"
    fi

    cgo_enabled=0
    export CC=""
    export CXX=""

    # Enable CGO and set cross-compilation toolchain for macOS
    if [ "$goos" = "darwin" ]; then
        cgo_enabled=1
        export CC="/home/runner/work/osxcross/target/bin/o64-clang"
        export CXX="/home/runner/work/osxcross/target/bin/o64-clang++"
    fi

    echo "go build ( $((j + 1)) / ${#supported_platforms[@]} ): GOOS=${goos} GOARCH=${goarch} CGO_ENABLED=${cgo_enabled} CC=${CC} CXX=${CXX} -o dist/bing-wallpaper-changer_${VERSION}_${p}${ext}"

    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="$cgo_enabled" go build \
        -trimpath \
        -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -extldflags=-static" \
        -tags="osusergo netgo static_build" \
        -o "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed" \
        "cmd/bing-wallpaper-changer/main.go"

    (
        upx --best -q -q -v "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed" -o "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}" &&
            rm "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed"
    ) || mv "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed" "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}"

    sha256sum "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}" | awk '{print $1}' >"dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.sha256"
done
