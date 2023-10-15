#!/bin/bash

mkdir DarkFlameMaster

COMMIT_ID=$(git log | head -n 1 | awk '{print $2}')
AUTHOR=$(git log | head -n 2 | tail -n 1 | awk '{print $2}')
BRANCH=$(git branch | grep \* | awk '{print $2}')
BUILD_TIME=$(date "+%Y-%m-%d|%H:%M:%S")
GO_VERSION=$(go version | awk '{print $3}')
BUILD_MODE="release"
SERVER_INFO="$COMMIT_ID;$AUTHOR;$BRANCH;$BUILD_TIME;$GO_VERSION;$BUILD_MODE"

go build -ldflags "-X DarkFlameMaster/serverinfo.ServerInfo=$SERVER_INFO" -o DarkFlameMaster

cp -r bin/* DarkFlameMaster
cp -r data DarkFlameMaster
cp -r view DarkFlameMaster
cp -r conf DarkFlameMaster

mkdir DarkFlameMaster/tools
mkdir DarkFlameMaster/tools/seatmaker

cd ./tools/dumper || exit
go build -o dumper
mv dumper ../../DarkFlameMaster/tools

cd ../manager || exit
go build -o manager
mv manager ../../DarkFlameMaster/tools

cd ../seatmaker || exit
go build -o seatmaker
mv seatmaker ../../DarkFlameMaster/tools/seatmaker
cp -r view ../../DarkFlameMaster/tools/seatmaker

cd ../..

tar czf DarkFlameMaster.tgz DarkFlameMaster

rm -rf DarkFlameMaster
