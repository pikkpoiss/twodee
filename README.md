#twodee


A library for doing 2d game stuff.  I'm not sure what format it will take,
except that it will use OpenGL and be basically sprite based.

My expectation is to use this for Ludum Dare competitions as I go.

## TODO

 - [x] Menus
 - [x] Sound
 - [x] Animations
 - [x] Flipping tiles vertically/horizontally
 - [x] Fullscreen toggle
 - [x] Keyboard / Mouse events
 - [x] Building on OSX / Linux / Windows
 - [x] Game grid + pathfinding
 - [x] Import from Tiled native file format (http://www.mapeditor.org/)
 - [x] Glow shaders
 - [ ] OSX / Windows packaging scripts
 - [ ] Very simple cinematics scripting
 - [ ] Collision detection

##Building (OSX)

Make sure Clang is your default compiler.

    export CC=clang
    export CXX=clang++

Get the environment set up:

    brew install go
    brew install hg
    brew install glew
    // Set up glfw3; see https://github.com/go-gl/glfw3.
    brew tap homebrew/versions
    brew install --build-bottle --static glfw3
    // Set up SDL.
    brew install sdl
    brew install libvorbis libogg sdl_mixer
    brew install sdl_image

Install go library deps:

    ./setup.sh

##Troubleshooting (OSX)


### Installing go-gl/gl.
Before Go 1.2, you may need to run with:

    CC=gcc CGO_CFLAGS=-ftrack-macro-expansion=0 go get github.com/go-gl/gl

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

### Issue with getting glfw3
If the go get github.com/go-gl/glfw3 command is throwing an error then you can go to your /usr/local/lib directory and create a symbolic link:

    ln -s libglfw.dylib libglfw3.dylib

### Installing kurrik/Go-SDL/sdl and kurrik/Go-SDL/mixer

    brew install gcc47
    CC=gcc-4.7 go get -u -v github.com/kurrik/Go-SDL/sdl
    CC=gcc-4.7 go get -u -v github.com/kurrik/Go-SDL/mixer

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

### Software Dependencies

#### Go 1.4.1
  * http://golang.org
  * Download Windows MSI
  * Install to `C:\Go\`
  * Add environment variable for user: `GOPATH = %USERPROFILE%\AppData\Local\Go`

#### MinGW-w64 4.8.2
  * http://mingw-w64.sourceforge.net/
  * MinGW-builds project
  * http://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/  
  * Downloaded `/4.8.2/threads-win32/seh/x86_64-4.8.2-release-win32-seh-rt_v3-rev3.7z`
  * Extract to `C:\mingw64`
  * Add `C:\mingw64\bin` to PATH

#### Git 1.9.0
  * Download http://git-scm.com/download/win
  * Install to `C:\git`
  * Git on PATH in installer

#### Mercurial 2.9.2
  * https://bitbucket.org/tortoisehg/files/downloads/
  * Downloaded `mercurial-2.9.2-x64.msi`

#### GTK+ 2.24
  * Contains pkg-config
  * http://ftp.acc.umu.se/pub/gnome/binaries/win32/gtk+/
  * Downloaded `/2.24/gtk+-bundle_2.24.10-20120208_win32.zip`
  * Extract to `C:\mingw64`
  * Add a new system variable `PKG_CONFIG_PATH, set to C:\mingw64\lib\pkgconfig`

#### Make 3.82.90
  * http://sourceforge.net/projects/mingw-w64/files/External%20binary%20packages%20%28Win64%20hosted%29/make/
  * Downloaded `make-3.82.90-20111115.zip`
  * Extract `bin_amd64/*` to `C:\mingw64\bin`

### Library Dependencies

#### GLEW 1.10.0
  * https://sourceforge.net/projects/glew
  * Downloaded `/files/glew/1.10.0/glew-1.10.0.zip`
  * Extract to `C:\src`
  * Git bash:

        cd /c/src/glew-1.10.0
        gcc -DGLEW_STATIC -DGLEW_NO_GLU -O2 -Wall -W -Iinclude -DGLEW_BUILD -o src/glew.o -c src/glew.c
        gcc -shared -Wl,-soname,libglew32.dll -Wl,--out-implib,lib/libglew32.dll.a -o lib/glew32.dll src/glew.o -LC:\mingw64\lib -lglu32 -lopengl32 -lgdi32 -luser32 -lkernel32
        # Create glew32.dll
        ar cr lib/libglew32.a src/glew.o
        cp lib/*.* ../../mingw64/lib
        cp -r include/GL ../../mingw64/include/

#### GLFW 3.0.4
  * http://www.glfw.org/download.html
  * 64-bit windows binaries
  * Downloaded `glfw-3.0.4.bin.WIN64.zip`
  * Extract to `C:\src`
  * Git bash:

        cd /c/src/glfw-3.0.4.bin.WIN64
        cp -r include/GLFW ../../mingw64/include/
        cp lib-mingw/*.* ../../mingw64/lib/
        cp lib-mingw/glfw3dll.a ../../mingw64/lib/libglfw3dll.a

#### SDL 1.2.15
  * http://www.libsdl.org/release/
  * Downloaded `SDL-1.2.15.zip`
  * Extract to `C:\src`
  * Git bash:

        cd /c/src/SDL-1.2.15
        ./configure
        make
        cp sdl.pc ../../mingw64/lib/pkgconfig
        mkdir ../../mingw64/include/SDL
        cp include/*.h ../../mingw64/include/SDL
        mkdir ../../mingw64/include/SDL/SDL
        cp include/SDL.h ../../mingw64/include/SDL/SDL
        cp build/.libs/* ../../mingw64/lib

#### SDL_image 1.2.12
  * http://www.libsdl.org/projects/SDL_image/release/
  * Downloaded `SDL_image-devel-1.2.12-VC.zip`
  * Extract to `C:\src`
  * Git bash:

        cd /c/src/SDL_image-1.2.12
        cp include/*.h ../../mingw64/include/
        cp lib/x64/*.dll ../../mingw64/lib/
        cp lib/x64/*.lib ../../mingw64/lib/

#### SDL_mixer 1.2.12
  * http://www.libsdl.org/projects/SDL_mixer/release/
  * Downloaded `SDL_mixer-devel-1.2.12-VC.zip`
  * Extract to `C:\src`
  * Git bash:

        cd /c/src/SDL_mixer-1.2.12
        cp include/*.h ../../mingw64/include/
        cp lib/x64/*.dll ../../mingw64/lib/
        cp lib/x64/*.lib ../../mingw64/lib/

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
