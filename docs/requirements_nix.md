# Requirements (Ubuntu Trusty)

The following is from an older version of the library but may provide some hints
until we can test with a linux distro again.

## (OLD) Building (Ubuntu Trusty)

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
