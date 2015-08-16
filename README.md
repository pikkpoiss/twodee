#twodee

A library for 2d games using OpenGL and Go.

Under heavy development, we have been using this for Ludum Dare competitions
so it changes from time to time.

## Features

 - Menus
 - Sound
 - Animations
 - Fullscreen mode
 - Keyboard / Mouse / Gamepad events
 - Building on OSX / Linux / Windows
 - Game grid + pathfinding
 - Import from Tiled native file format (http://www.mapeditor.org/)
 - Some effects shaders (like Glow)

## Requirements

 - [Windows](docs/requirements_win.md)
 - [OSX](docs/requirements_osx.md)

## Setup

This library depends on various C/C++ packages.  To try and get a stable
environment for building, sources have been included in the `lib` directory.
Running the following will build each package and attempt to install a
Golang wrapper:

    ./scripts/setup.sh


# The following is out of date and will be removed shortly.

## Building (Ubuntu Trusty)

Install deps:

    sudo apt-get install cmake libglu1-mesa-dev libxrandr-dev libxi-dev libxcursor-dev clang libglew-dev mercurial
    sudo apt-get install libsdl1.2-dev libsdl-mixer1.2-dev libsdl-image1.2-dev

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

Install SDL stuff:

    CGO_CFLAGS="-I/usr/include/SDL" go get github.com/kurrik/Go-SDL/sdl
    CGO_CFLAGS="-I/usr/include/SDL" go get github.com/kurrik/Go-SDL/mixer

## Building (Windows 8.1)

### Go Library Dependencies
(may need to put `CGO_CFLAGS="-I C:\mingw64\include" CGO_LDFLAGS="-L C:\mingw64\lib"` in front of some/all of these but it doesn't seem to need it any more for all of them.

    go get -u github.com/go-gl/gl/v3.3-core/gl
    go get -u github.com/go-gl/glfw/v3.1/glfw
    go get -u github.com/go-gl/mathgl/mgl32
    go get -u github.com/robertkrimen/otto
    go get -u github.com/go-gl/glfw/v3.1/glfw
    go get code.google.com/p/freetype-go/freetype
    go get -u github.com/kurrik/Go-SDL/sdl
    CGO_CFLAGS="-I C:\mingw64\include" CGO_LDFLAGS="-L C:\mingw64\lib" go get -u github.com/kurrik/Go-SDL/mixer

### Twodee

    git clone git@github.com:kurrik/twodee-examples.git
    cd twodee-examples
    git submodule init
    git submodule update
    cd examples/basic
    PATH="$PATH;C:\mingw64\lib" go run *.go
    PATH="/c/mingw64/lib:$PATH" go run *.go (cygwin)
