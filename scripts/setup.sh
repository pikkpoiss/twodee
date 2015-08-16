#!/usr/bin/env bash

. `git rev-parse --show-toplevel`/scripts/common.sh

BUILDROOT=$ROOT/build
PREFIX=$BUILDROOT/usr

green "INIT" "Prefix is $PREFIX"

##### Folder setup #############################################################

if [ -n "$CLEAN" ]; then
  yellow "CLEAN" "Deleting the build path"
  rm -rf $BUILDROOT
  yellow "CLEAN" "Deleting $GOPATH/src and $GOPATH/pkg"
  rm -rf $GOPATH/src
  rm -rf $GOPATH/pkg
fi

mkdir -p $BUILDROOT
cp lib/*.zip $BUILDROOT
cd $BUILDROOT

##### Libraries ################################################################

if file_exists $PREFIX/lib/libglfw3.a; then
  green "EXISTS" "glfw"
else
  yellow "BUILD" "glfw"
  rm -rf glfw-3.1.1
  unzip glfw-3.1.1.zip
  cd glfw-3.1.1
  cmake \
    -DBUILD_SHARED_LIBS=OFF \
    -DCMAKE_INSTALL_PREFIX:PATH=$PREFIX \
    .
  make
  make install
  cd ..
fi

if file_exists $PREFIX/lib/libogg.a; then
  green "EXISTS" "ogg"
else
  yellow "BUILD" "ogg"
  rm -rf libogg-1.3.2
  unzip libogg-1.3.2.zip
  cd libogg-1.3.2
  ./configure \
    --prefix=$PREFIX \
    --disable-shared
  make
  make install
  cd ..
fi

if file_exists $PREFIX/lib/libvorbis.a; then
  green "EXISTS" "vorbis"
else
  yellow "BUILD" "vorbis"
  rm -rf libvorbis-1.3.5
  unzip libvorbis-1.3.5.zip
  cd libvorbis-1.3.5
  ./configure \
    --prefix=$PREFIX \
    --disable-shared
  make
  make install
  cd ..
fi

if file_exists $PREFIX/lib/libSDL2.a; then
  green "EXISTS" "SDL2"
else
  yellow "BUILD" "SDL2"
  rm -rf SDL2-2.0.3
  unzip SDL2-2.0.3.zip
  cd SDL2-2.0.3
  ./configure \
    --prefix=$PREFIX \
    --disable-shared
  make
  make install
  cd ..
fi

if file_exists $PREFIX/lib/libSDL2_mixer.a; then
  green "EXISTS" "SDL2 mixer"
else
  yellow "BUILD" "SDL2 mixer"
  rm -rf SDL2_mixer-2.0.0
  unzip SDL2_mixer-2.0.0.zip
  cd SDL2_mixer-2.0.0
  ./configure \
    --prefix=$PREFIX \
    --with-sdl-prefix=$PREFIX \
    --disable-music-ogg-shared \
    --disable-shared
  make
  make install
  cd ..
fi

##### Go libraries #############################################################

CGO_CFLAGS="-I$PREFIX/include"

go get -u code.google.com/p/freetype-go/freetype
go get -u code.google.com/p/freetype-go/freetype/raster
go get -u code.google.com/p/freetype-go/freetype/truetype
go get -u github.com/go-gl/gl/v3.3-core/gl
go get -u github.com/go-gl/glfw/v3.1/glfw
go get -u github.com/go-gl/mathgl/mgl32
go get -u github.com/kurrik/Go-SDL/mixer
go get -u github.com/kurrik/Go-SDL/sdl
go get -u github.com/robertkrimen/otto
