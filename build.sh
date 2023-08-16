#!/bin/bash

go build -o bin

mkdir DarkFlameMaster
cp -r bin/* DarkFlameMaster
cp -r data DarkFlameMaster
cp -r view DarkFlameMaster
cp -r conf DarkFlameMaster

tar czf DarkFlameMaster.tgz DarkFlameMaster

rm -rf DarkFlameMaster
