#!/bin/bash
set -euo pipefail

APPNAME=equeselfgo

[ -f ./go.mod ] || exit 0

# install fyne-cross
go install github.com/fyne-io/fyne-cross@v1.6.1 # or develop or master
~/go/bin/fyne-cross version

# clean up
rm -rf ./bin ./fyne-cross ./tmp-pkg ./fyne_metadata_init.go ./fyne.syso ./${APPNAME} ./${APPNAME}.app ./${APPNAME}.exe ./*.tar.xz ./*.apk

# linux and freebsd
for os in linux freebsd; do
	echo "Building $os..."
	~/go/bin/fyne-cross $os -pull -arch=amd64,arm64 -release
	echo "--------------------------------------------------------------------------------"
done

# darwin
# https://developer.apple.com/download/all/?q=Xcode%2012.5.1
echo "Building darwin..."
XCODE_CMD_TOOLS="Command_Line_Tools_for_Xcode_12.5.1.dmg"
if [ ! -d ~/Downloads/SDKs/MacOSX.sdk ] && [ -f ~/Downloads/${XCODE_CMD_TOOLS} ]; then
	~/go/bin/fyne-cross darwin-sdk-extract -engine docker -pull -xcode-path ~/Downloads/${XCODE_CMD_TOOLS}
fi
if [ -d ~/Downloads/SDKs/MacOSX.sdk ]; then
	~/go/bin/fyne-cross darwin -pull -arch=arm64 -macosx-sdk-path ~/Downloads/SDKs/MacOSX.sdk
fi
echo "--------------------------------------------------------------------------------"

# windows
echo "Building windows..."
~/go/bin/fyne-cross windows -pull -arch=amd64,arm64
echo "--------------------------------------------------------------------------------"

# android
echo "Building android..."
~/go/bin/fyne-cross android -pull -release
echo "--------------------------------------------------------------------------------"

# export
echo "Copying files..."
mkdir -p ./bin
[ -f ./fyne-cross/bin/darwin-arm64/${APPNAME} ] && cp ./fyne-cross/bin/darwin-arm64/${APPNAME} ./bin/${APPNAME}-darwin-arm64.app
[ -f ./fyne-cross/bin/linux-amd64/${APPNAME} ] && cp ./fyne-cross/bin/linux-amd64/${APPNAME} ./bin/${APPNAME}-linux-amd64.bin
[ -f ./fyne-cross/bin/linux-arm64/${APPNAME} ] && cp ./fyne-cross/bin/linux-arm64/${APPNAME} ./bin/${APPNAME}-linux-arm64.bin
[ -f ./fyne-cross/bin/freebsd-amd64/${APPNAME} ] && cp ./fyne-cross/bin/freebsd-amd64/${APPNAME} ./bin/${APPNAME}-freebsd-amd64.bin
[ -f ./fyne-cross/bin/freebsd-arm64/${APPNAME} ] && cp ./fyne-cross/bin/freebsd-arm64/${APPNAME} ./bin/${APPNAME}-freebsd-arm64.bin
[ -f ./fyne-cross/bin/windows-amd64/${APPNAME}.exe ] && cp ./fyne-cross/bin/windows-amd64/${APPNAME}.exe ./bin/${APPNAME}-windows-amd64.exe
[ -f ./fyne-cross/bin/windows-arm64/${APPNAME}.exe ] && cp ./fyne-cross/bin/windows-arm64/${APPNAME}.exe ./bin/${APPNAME}-windows-arm64.exe
[ -f ./fyne-cross/dist/android/${APPNAME}.apk ] && cp ./fyne-cross/dist/android/${APPNAME}.apk ./bin/${APPNAME}-android.apk

cd ./bin
zip -r ./${APPNAME}.zip ./
cd ..
echo "--------------------------------------------------------------------------------"

# Build AppImage
./build-appimage.sh
echo "--------------------------------------------------------------------------------"

# Build deb package
# ./build-deb.sh
echo "--------------------------------------------------------------------------------"

# Build rpm package
# ./build-rpm.sh
echo "--------------------------------------------------------------------------------"

ls -alh ./bin/*
echo "Build completed successfully!"
