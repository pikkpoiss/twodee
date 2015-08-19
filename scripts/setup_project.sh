#!/usr/bin/env bash

. `git rev-parse --show-toplevel`/scripts/common.sh

if [ $# -ne 1 ]; then
  echo "Usage: $0 PATH"
  exit 1
fi

DESTPATH=$1
CWD=`pwd`

if [ ! -d "$DESTPATH" ]; then
  echo "$DESTPATH must be a directory"
  exit 1
fi

echo -n "Path: "
green $DESTPATH
echo -n "Does this look OK? [Press 'y'] "
read -n 1 CONFIRM
echo

if [ "$CONFIRM" != "y" ]; then
  red "Stopping"
  exit 1
fi
echo

##### Setup ####################################################################

if [ ! -d "$DESTPATH/.git" ]; then
  echo -n "Dest is not a git repo, running git init... "
  cd $DESTPATH
  git init --quiet
  cd $CWD
else
  echo -n "Dest is a git repo, doing nothing... "
fi
green "OK"

echo -n "Copying files... "
cp -r $ROOT/scripts/project/* $DESTPATH
green "OK"

WIN_DLLS=$ROOT/build/usr/bin/*.dll
if [ -n "$WIN_DLLS" ]; then
  echo -n "Copying windows dlls... "
  cp $WIN_DLLS $DESTPATH/lib/win32-x64/
else
  echo -n "Could not find windows dlls, doing nothing... "
fi
green "OK"

GIT_REPO=`git config --get remote.origin.url`
echo -n "Adding subproject '$GIT_REPO'... "
cd $DESTPATH/lib
git submodule add --quiet $GIT_REPO
cd $CWD
green "OK"
