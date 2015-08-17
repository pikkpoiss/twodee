#!/usr/bin/env bash

. `git rev-parse --show-toplevel`/scripts/common.sh

if [[ "$PLATFORM" == "win" ]]; then
  export CC="/c/mingw/mingw32/bin/gcc.exe"
fi

BUILDROOT=$ROOT/build
PREFIX=$BUILDROOT/usr
INCDIR=$PREFIX/include
LIBDIR=$PREFIX/lib

export LDFLAGS="-L$LIBDIR"
export CPPFLAGS="-I$INCDIR $EXTRA_CPPFLAGS"
export CFLAGS="-I$INCDIR $EXTRA_CPPFLAGS"
export PKG_CONFIG_PATH="$PREFIX/lib/pkgconfig"

green "INIT" "Prefix is $PREFIX, Platform is $PLATFORM"

##### Folder setup #############################################################

if [ -n "$CLEAN" ]; then
  yellow "CLEAN" "Deleting the build path"
  rm -rf $BUILDROOT
fi

mkdir -p $BUILDROOT
cp lib/*.{zip,tar.gz} $BUILDROOT
cd $BUILDROOT

##### Helpers ##################################################################

function patch_makefile {
  if [[ "$PLATFORM" == "win" ]]; then
    yellow "BUILD" "patching Makefile"
    sed s/c:\\\//\\\/c\\\//g < Makefile > Makefile.new
    mv Makefile Makefile.old
    mv Makefile.new Makefile
  fi
}

##### Libraries ################################################################

if file_exists $PREFIX/lib/libglfw3.a; then
  green "EXISTS" "glfw"
else
  yellow "BUILD" "glfw"
  rm -rf glfw-3.1.1
  unzip -q glfw-3.1.1.zip
  cd glfw-3.1.1
  cmake \
    -G "Unix Makefiles" \
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
  unzip -q libogg-1.3.2.zip
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
  unzip -q libvorbis-1.3.5.zip
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
  yellow "BUILD" "unzip"
  unzip -q SDL2-2.0.3.zip
  cd SDL2-2.0.3
  if [[ "$PLATFORM" == "win" ]]; then
    yellow "BUILD" "patching"
    cd src
    git apply ../../../lib/SDL2-fix-gcc-compatibility.patch
    git apply ../../../lib/SDL2-prevent-duplicate-d3d11-declarations.patch
    cd ..
  fi
  yellow "BUILD" "configure"
  ./configure \
    --prefix=$PREFIX \
    --disable-shared
  yellow "BUILD" "make"
  make
  yellow "BUILD" "make install"
  make install
  cd ..
fi

if file_exists $PREFIX/lib/libSDL2_image.a; then
  green "EXISTS" "SDL2 image"
else
  yellow "BUILD" "SDL2 image"
  rm -rf SDL2_image-2.0.0
  if [[ "$PLATFORM" == "win" ]]; then
    tar -xf SDL2_image-devel-2.0.0-mingw.tar.gz
    cd SDL2_image-2.0.0/x86_64-w64-mingw32
    cp -r {include,lib} $PREFIX
    cd ../..
  else
    unzip -q SDL2_image-2.0.0.zip
    cd SDL2_image-2.0.0
    yellow "BUILD" "configure"
    ./configure \
      --disable-sdltest \
      --prefix=$PREFIX \
      --with-sdl-prefix=$PREFIX \
      --disable-png-shared \
      --disable-shared
    patch_makefile
    yellow "BUILD" "make"
    make
    yellow "BUILD" "make install"
    make install
    cd ..
  fi
fi

if file_exists $PREFIX/lib/libSDL2_mixer.a; then
  green "EXISTS" "SDL2 mixer"
else
  yellow "BUILD" "SDL2 mixer"
  rm -rf SDL2_mixer-2.0.0
  unzip -q SDL2_mixer-2.0.0.zip
  cd SDL2_mixer-2.0.0
  ./configure \
    --disable-sdltest \
    --prefix=$PREFIX \
    --with-sdl-prefix=$PREFIX \
    --disable-music-ogg-shared \
    --disable-shared
  make
  make install
  cd ..
fi

##### Go libraries #############################################################

export CGO_CFLAGS="-I$PREFIX/include"
export CGO_LDFLAGS="`$PREFIX/bin/sdl2-config --static-libs` -lvorbisfile -lvorbis -logg"

# Require libraries
go get -u -v github.com/scottferg/Go-SDL2/sdl
go get -u -v github.com/scottferg/Go-SDL2/mixer
go get -u -v github.com/go-gl/glfw/v3.1/glfw

# Do not require libraries
go get -u -v github.com/go-gl/gl/v3.3-core/gl
go get -u -v github.com/go-gl/mathgl/mgl32
go get -u -v code.google.com/p/freetype-go/freetype
go get -u -v code.google.com/p/freetype-go/freetype/raster
go get -u -v code.google.com/p/freetype-go/freetype/truetype
go get -u -v github.com/robertkrimen/otto

# Old
# go get -u github.com/kurrik/Go-SDL/mixer
# go get -u github.com/kurrik/Go-SDL/sdl
