#!/usr/bin/env bash

set -euo pipefail

VERSION="${VERSION:-1.7.5}"
VERSION="${VERSION#v}"
BUILD_DIR="${BUILD_DIR:-build}"
PACKAGE="git.taurusxin.com/taurusxin/ncmdump-go"
LDFLAGS="-w -s -X main.version=${VERSION}"

mkdir -p "${BUILD_DIR}"
rm -f \
    "${BUILD_DIR}/ncmdump-go" \
    "${BUILD_DIR}/ncmdump-go.exe" \
    "${BUILD_DIR}/ncmdump-go_linux_amd64_${VERSION}.tar.gz" \
    "${BUILD_DIR}/ncmdump-go_linux_arm64_${VERSION}.tar.gz" \
    "${BUILD_DIR}/ncmdump-go_darwin_amd64_${VERSION}.tar.gz" \
    "${BUILD_DIR}/ncmdump-go_darwin_arm64_${VERSION}.tar.gz" \
    "${BUILD_DIR}/ncmdump-go_windows_amd64_${VERSION}.zip"

create_tar_archive() {
    local os="$1"
    local arch="$2"
    local binary="${BUILD_DIR}/ncmdump-go"
    local archive="${BUILD_DIR}/ncmdump-go_${os}_${arch}_${VERSION}.tar.gz"
    local tar_version

    echo "Building for ${os} ${arch}..."
    CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" \
        go build -trimpath -ldflags="${LDFLAGS}" -o "${binary}" "${PACKAGE}"

    # Builds are sometimes produced from Git Bash on Windows, where chmod does
    # not change NTFS permissions. Force the executable bit into GNU tar's
    # header so macOS/Linux users can run the extracted binary immediately.
    chmod 0755 "${binary}"
    tar_version="$(tar --version 2>/dev/null || true)"
    if [[ "${tar_version}" == *"GNU tar"* ]]; then
        tar --mode=0755 -zcf "${archive}" -C "${BUILD_DIR}" ncmdump-go
    else
        tar -zcf "${archive}" -C "${BUILD_DIR}" ncmdump-go
    fi

    if ! tar -tvzf "${archive}" | awk 'NR == 1 { exit(substr($1, 1, 10) == "-rwxr-xr-x" ? 0 : 1) }'; then
        echo "Archive validation failed: ncmdump-go is not executable in ${archive}" >&2
        exit 1
    fi

    rm "${binary}"
}

create_tar_archive linux amd64
create_tar_archive linux arm64
create_tar_archive darwin amd64
create_tar_archive darwin arm64

# Windows amd64
echo "Building for Windows amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
    go build -trimpath -ldflags="${LDFLAGS}" -o "${BUILD_DIR}/ncmdump-go.exe" "${PACKAGE}"
zip -q -j "${BUILD_DIR}/ncmdump-go_windows_amd64_${VERSION}.zip" "${BUILD_DIR}/ncmdump-go.exe"
rm "${BUILD_DIR}/ncmdump-go.exe"

echo "Release archives written to ${BUILD_DIR}/"
