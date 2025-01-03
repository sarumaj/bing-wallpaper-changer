#!/bin/bash
set -e

supported_platforms=(
    darwin-amd64
    darwin-arm64
    linux-amd64
    linux-arm64
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
    extldflags="-static"
    build_tags="osusergo netgo static_build"

    if [ "${CGO_ENABLED:-0}" = "1" ]; then
        build_tags="osusergo netgo"
        ext="-cgo"
        extldflags=""
    fi

    ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -extldflags '${extldflags}'"
    if [ "$goos" = "windows" ]; then
        ext+=".exe"
        ldflags+=" -H windowsgui"
    fi

    echo "go build ( $((j + 1)) / ${#supported_platforms[@]} ): GOOS=${goos} GOARCH=${goarch} CGO_ENABLED=${CGO_ENABLED:-0} -o dist/bing-wallpaper-changer_${VERSION}_${p}${ext}"

    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="${CGO_ENABLED:-0}" go build \
        -trimpath \
        -ldflags="${ldflags}" \
        -tags="${build_tags}" \
        -o "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed" \
        "cmd/bing-wallpaper-changer/main.go"

    (
        upx --best -q -q -v "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed" -o "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}" &&
            rm "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed"
    ) || mv "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.uncompressed" "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}"

    sha256sum "dist/bing-wallpaper-changer_${VERSION}_${p}${ext}" | awk '{print $1}' >"dist/bing-wallpaper-changer_${VERSION}_${p}${ext}.sha256"
done
