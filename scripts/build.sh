#!/usr/bin/env bash

set -eu

package=pr
platforms=(
	"windows/amd64"
	"darwin/amd64"
	"darwin/arm64"
	"linux/amd64"
	"linux/arm64"
)

for platform in "${platforms[@]}"; do
	echo "Building for $platform"

	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}

	# Build the binary
	env GOOS=$GOOS GOARCH=$GOARCH go build -o $package

	# Zip the binary
	mkdir -p dist
	output_name="$package-$GOOS-$GOARCH.tar.gz"
	tar czf dist/$output_name $package

	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
