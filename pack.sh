#!/bin/bash
VER=$1
if [ "$VER" = "" ]; then
    echo 'please input pack version!'
    exit 1
fi
RELEASE="release-${VER}"
rm -rf release-*
mkdir ${RELEASE}

# windows amd64
echo 'Start pack windows amd64...'
GOOS=windows GOARCH=amd64 go build  
tar -czvf "${RELEASE}/wmqx-windows-amd64.tar.gz" wmqx.exe config.toml message.json log/.gitignore License.md README.md README_CN.md
rm -rf wmqx.exe

echo 'Start pack windows X386...'
GOOS=windows GOARCH=386 go build 
tar -czvf "${RELEASE}/wmqx-windows-386.tar.gz" wmqx.exe config.toml message.json log/.gitignore License.md README.md README_CN.md
rm -rf wmqx.exe

echo 'Start pack linux amd64'
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
tar -czvf "${RELEASE}/wmqx-linux-amd64.tar.gz" wmqx.exe config.toml message.json log/.gitignore License.md README.md README_CN.md
rm -rf wmqx

echo 'Start pack linux 386'
GOOS=linux GOARCH=386 go build -ldflags "-s -w"
tar -czvf "${RELEASE}/wmqx-linux-386.tar.gz" wmqx.exe config.toml message.json log/.gitignore License.md README.md README_CN.md
rm -rf wmqx

echo 'Start pack mac amd64'
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w"
tar -czvf "${RELEASE}/wmqx-mac-amd64.tar.gz" wmqx.exe config.toml message.json log/.gitignore License.md README.md README_CN.md
rm -rf wmqx

echo 'END'
