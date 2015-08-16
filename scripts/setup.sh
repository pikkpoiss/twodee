#!/usr/bin/env bash

. `git rev-parse --show-toplevel`/scripts/common.sh

BUILDROOT=$ROOT/build
PREFIX=$BUILDROOT/usr

green "INIT" "Prefix is $PREFIX"

##### Folder setup #############################################################

if [ -n "$CLEAN" ]; then
  yellow "CLEAN" "Deleting the build path"
  rm -rf $BUILDROOT
fi

mkdir -p $BUILDROOT
cp lib/*.zip $BUILDROOT
cd $BUILDROOT

##### Libraries ################################################################

if file_exists $BUILDROOT/foo; then
  green "EXISTS" "glfw"
else
  yellow "BUILD" "glfw"
  unzip glfw-3.1.1.zip
  cd glfw-3.1.1
  cmake \
    -DBUILD_SHARED_LIBS=OFF \
    -DCMAKE_INSTALL_PREFIX:PATH=$PREFIX \
    .
  make
  make install
fi
