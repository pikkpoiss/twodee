#twodee


A library for doing 2d game stuff.  I'm not sure what format it will take,
except that it will use OpenGL and be basically sprite based.

My expectation is to use this for Ludum Dare competitions as I go.


##Building (OSX)

Make sure Clang is your default compiler.

    export CC=clang
    export CXX=clang++

Install deps:

    go get code.google.com/p/freetype-go/freetype
    go get github.com/Agon/googlmath
    go get github.com/go-gl/gl
    go get github.com/go-gl/glfw3

##Troubleshooting (OSX)


### Installing go-gl/gl.
Before Go 1.2, you may need to run with:

    CC=gcc CGO_CFLAGS=-ftrack-macro-expansion=0 \
    go get github.com/go-gl/gl

See http://stackoverflow.com/questions/16412644/using-opengl-from-go for background.

### Installing glfw
Maybe you need to install glfw's shared lib (TODO: see if there's a simple brew for this):

    git clone https://github.com/glfw/glfw
    mkdir glfw/build
    cd glfw/build
    cmake -DBUILD_SHARED_LIBS=1 ..
    make
    sudo make install

### Installing go-gl/glfw3
Might need to specify CFLAGS and LDFLAGS for deps:

    CGO_CFLAGS="-I/usr/include" \
    CGO_LDFLAGS="`pkg-config --libs glu x11 glfw3 xrandr xxf86vm xi xcursor` -lm" \
    go get github.com/go-gl/glfw3

### Running programs
Sometimes installed library paths are not in LD_LIBRARY_PATH. Try:

    LD_LIBRARY_PATH=/usr/local/lib go run *.go

## Building (Ubuntu Trusty)

Install deps:

    sudo apt-get install cmake libglu1-mesa-dev libxrandr-dev libxi-dev libxcursor-dev clang libglew-dev mercurial

Build glfw:

    git clone https://github.com/glfw/glfw
    mkdir glfw/build
    cd glfw/build
    cmake -DBUILD_SHARED_LIBS=1 ..
    make
    sudo make install
    cd /usr/local/lib
    sudo ln -s libglfw3.a libglfw.a

 Install go-gl/glfw3:

    CGO_CFLAGS="-I/usr/include" \
    CGO_LDFLAGS="`pkg-config --libs glu x11 glfw3 xrandr xxf86vm xi xcursor` -lm" \
    go get github.com/go-gl/glfw3

Install other deps:

    go get github.com/go-gl/gl
    go get code.google.com/p/freetype-go/freetype
    go get github.com/Agon/googlmath

-------------------------------------

##Old, old instructions

Ubuntu:

    sudo apt-get install freeglut3-dev
    sudo apt-get install libxmu-dev
    sudo apt-get install libxi-dev
    sudo apt-get install libxrandr-dev
    sudo apt-get install libglew-dev
    sudo apt-get install libglfw-dev

OSX:
  GLEW - http://glew.sourceforge.net/
  I used version 1.9.0

    tar xvzf glew-1.9.0.tgz
    cd glew-1.9.0
    make
    sudo make install

  I also needed
  libglfw - http://www.glfw.org/download.html
  I used version 2.7.6
  Unzip the source and cd to the base of the package.  Run (OSX):

    make cocoa-dist-install

  Need Mercurial

    brew install hg

  Make sure to use gcc for compiling go-gl/gl:

    CC=gcc go get -u github.com/go-gl/gl

Win:
  * Install Mercurial from http://mercurial.selenic.com/
  * Download binaries for GLFW http://sourceforge.net/projects/glfw/files/glfw/2.7.8/glfw-2.7.8.bin.WIN64.zip/download
  * Copy glfw-2.7.6.bin.WIN32/lib-mingw/x64 stuff to C:\MinGW64\lib
  * Copy dist include to gcc include
  * Copy glfw.dll to C:\Windows\System32

Then (all):

    go get github.com/go-gl/gl
    go get github.com/go-gl/glfw
    go get code.google.com/p/freetype-go/freetype
    go get code.google.com/p/freetype-go/freetype/truetype
    go get github.com/kurrik/gltext

I think cocoa-dist-install works now, but previously you had to explicitly
build the dylib using something like the following:

Unzip the source and cd to the base of the package.
You need to build the project as a shared lib, so I used the
following (OSX):

    cd lib/cocoa
    make -f Makefile.cocoa libglfw.dylib
    install -c -m 644 libglfw.dylib /usr/local/lib/

