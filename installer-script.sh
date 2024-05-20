#!/bin/bash

REPO="szeretlekblackfire/animedrive-go"
VERSION="v1.0.0"

if [ "$(uname)" == "Darwin" ]; then
    OS="darwin"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    OS="linux"
elif [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ]; then
    OS="windows"
else
    echo "Unsupported OS"
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$VERSION/animedrive-dl-$OS-amd64"

if [ "$OS" == "windows" ]; then
    curl -L $URL -o animedrive-dl.exe
    mv animedrive-dl.exe /usr/local/bin/animedrive-dl.exe
else
    curl -L $URL -o animedrive-dl
    chmod +x animedrive-dl
    mv animedrive-dl /usr/local/bin/animedrive-dl
fi

echo "animedrive-dl installed successfully."
