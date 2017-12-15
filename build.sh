#!/usr/bin/env bash

workDir=$(pwd)

outputAppPrefix=${workDir}/bin/gvm

srcFile=${workDir}/gvm.go

platform=$(uname -s)

echo "Current operating systems:    ${platform}"

rm -rf ${workDir}/bin/*

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${outputAppPrefix} -ldflags "-s -w" ${srcFile}
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${outputAppPrefix}-darwin -ldflags "-s -w" ${srcFile}
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${outputAppPrefix}-amd64 -ldflags "-s -w" ${srcFile}
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${outputAppPrefix}-windows.exe -ldflags "-s -w" ${srcFile}

if [ $platform = "Darwin" ];then
    if [ -n `which upx` ];then
        upx ${outputAppPrefix}
    fi
fi