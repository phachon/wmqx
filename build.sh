#!/bin/sh

# start go mod
export GO111MODULE=on
# set goproxy
export GOPROXY=https://goproxy.cn

PROJECT_NAME="wmqx"
BUILD_DIR="release"
ROOT_DIR=`pwd`

# windows .exe
if [ "${GOOS}" = "windows" ]; then
    PROJECT_NAME=${PROJECT_NAME}".exe"
fi

rm -rf ${BUILD_DIR}

build_app() {
    mkdir -p ${ROOT_DIR}/${BUILD_DIR}/conf
    mkdir -p ${ROOT_DIR}/${BUILD_DIR}/logs
    mkdir -p ${ROOT_DIR}/${BUILD_DIR}/bin

    /bin/cp -r ${ROOT_DIR}/conf/default.toml ${ROOT_DIR}/${BUILD_DIR}/conf/
    /bin/cp -r ${ROOT_DIR}/README.md ${ROOT_DIR}/${BUILD_DIR}
    /bin/cp -r ${ROOT_DIR}/README_CN.md ${ROOT_DIR}/${BUILD_DIR}
    /bin/cp -r ${ROOT_DIR}/LICENSE ${ROOT_DIR}/${BUILD_DIR}
    /bin/cp -r ${ROOT_DIR}/scripts/*.sh ${ROOT_DIR}/${BUILD_DIR}

    chmod -R 755 ${ROOT_DIR}/${BUILD_DIR}/conf/
    chmod -R 755 ${ROOT_DIR}/${BUILD_DIR}/logs/
    chmod -R 755 ${ROOT_DIR}/${BUILD_DIR}/bin/
    chmod -R 755 ${ROOT_DIR}/${BUILD_DIR}/*.sh

    go build -o ${PROJECT_NAME} -ldflags "-s -w" ./

    if [ -f "${ROOT_DIR}/${PROJECT_NAME}"  ]; then
        mv ${ROOT_DIR}/${PROJECT_NAME} ${ROOT_DIR}/${BUILD_DIR}/bin/
    fi
    return
}

echo ">> WMQX start build!"
build_app
