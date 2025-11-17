#!/bin/bash
set -e

VERSION="1.0.0"
NAME="terraform-provider-envsend"

PLATFORMS=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
  OS=$(echo $PLATFORM | cut -d " " -f1)
  ARCH=$(echo $PLATFORM | cut -d " " -f2)
  OUTPUT="${NAME}_${VERSION}_${OS}_${ARCH}"

  if [ "$OS" == "windows" ]; then
    OUTPUT="${OUTPUT}.exe"
  fi

  echo "Building $OUTPUT ..."
  GOOS=$OS GOARCH=$ARCH go build -o "dist/${OUTPUT}"
done
