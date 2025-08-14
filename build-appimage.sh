#!/bin/bash
set -euo pipefail

APPNAME=equeselfgo

[ -f ./go.mod ] || exit 0
[ -f ./bin/${APPNAME}-linux-amd64.bin ] || exit 0

echo -n "Building appimage..."
if [ ! -f ~/Downloads/appimagetool-x86_64.AppImage ]; then
	wget -q -O ~/Downloads/appimagetool-x86_64.AppImage https://github.com/AppImage/appimagetool/releases/download/continuous/appimagetool-x86_64.AppImage
	chmod +x ~/Downloads/appimagetool-x86_64.AppImage
fi

rm -rf ./build-appimage
mkdir -p ./build-appimage/opt/${APPNAME}

cp ./bin/${APPNAME}-linux-amd64.bin ./build-appimage/opt/${APPNAME}/${APPNAME}
cp ./Icon.png ./build-appimage/opt/${APPNAME}/Icon.png
cp ./Icon.png ./build-appimage/ico.png
cp ./${APPNAME}.desktop ./build-appimage/${APPNAME}.desktop
cp ./appimage-apprun.sh ./build-appimage/AppRun

~/Downloads/appimagetool-x86_64.AppImage ./build-appimage ./bin/${APPNAME}-x86_64.AppImage
