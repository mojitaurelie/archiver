#!/bin/bash

platforms=("darwin/amd64" "darwin/arm64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mipsle" "linux/mips64" "linux/mips64le" "linux/ppc64le" "linux/riscv64")

if [[ -d "./build" ]]
then
    rm -r ./build
fi

mkdir build
cd build

for platform in "${platforms[@]}"
do
    echo "* Compiling for $platform..."
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='archiver_'$GOOS'_'$GOARCH
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name -a ..
done