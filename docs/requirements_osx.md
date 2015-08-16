# Requirements (OSX)

 - Go 1.4.2 or later (From [golang.org](http://golang.org))
 - Mercurial (`brew install hg` or equivalent)
 - Clang (Available with Xcode on OSX)
 - Xcode tools `xcode-select --install`

## Troubleshooting

### Clang

Compiles should use Clang.  You may need to use the following if you
have a different compiler as default:

    export CC=clang
    export CXX=clang++

### Libraries not found

This shouldn't happen if you ran `./scripts/setup.sh` and built static
libraries.  If you use shared libs instead (Brew) then sometimes
installed library paths are not in LD_LIBRARY_PATH. Try:

    LD_LIBRARY_PATH=/usr/local/lib go run *.go

### glfw3

Shouldn't happen if you use `./scripts/setup.sh`.  If you don't, you
might need to specify CFLAGS and LDFLAGS for deps:

    CGO_CFLAGS="-I/usr/include" \
    CGO_LDFLAGS="`pkg-config --libs glu x11 glfw3 xrandr xxf86vm xi xcursor` -lm" \
    go get github.com/go-gl/glfw3

## Brew
Note: this is *not* the suggested way of installing dependencies.  You should
install the dependencies listed above.  However if you absolutely need to use
brew, then the following may work:

    brew install go
    brew install hg
    brew install glew
    // Set up glfw3; see https://github.com/go-gl/glfw3.
    brew tap homebrew/versions
    brew install --build-bottle --static glfw3
    // Set up SDL.
    brew install sdl2
    brew install libvorbis libogg sdl2_mixer
    brew install sdl2_image

However, the risk with this is that you will not build static libraries,
meaning distributing the application will require managing shared libraries
for each platform you support.
