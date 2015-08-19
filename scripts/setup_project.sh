#!/usr/bin/env bash

. `git rev-parse --show-toplevel`/scripts/common.sh

echo -n "Enter the path to the project dir: "
read PATH

echo -n "Path: "
green $PATH
echo -n "Does this look OK? [Press 'y'] "
read -n 1 CONFIRM
echo

if [ "$CONFIRM" != "y" ]; then
  red "Stopping"
  exit 1
fi
echo

echo -n "Copying files... "
cp -r $ROOT/scripts/project/* $PATH
green "OK"

echo -n "Copying dlls... "
cp $ROOT/build/usr/bin/*.dll $PATH/lib/win32-x64/
green "OK"
