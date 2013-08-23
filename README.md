twodee
======

A library for doing 2d game stuff.  I'm not sure what format it will take, 
except that it will use OpenGL and be basically sprite based.

I'm planning on including a scene graph implementation and high level
functions for handling events and input, so this should be suitable for
making quick games in Go.

My expectation is to use this for Ludum Dare competitions as I go.


Dependencies
------------
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

Old instructions
----------------
I think cocoa-dist-install works now, but previously you had to explicitly
build the dylib using something like the following:

Unzip the source and cd to the base of the package.
You need to build the project as a shared lib, so I used the
following (OSX):

    cd lib/cocoa
    make -f Makefile.cocoa libglfw.dylib
    install -c -m 644 libglfw.dylib /usr/local/lib/

Thanks
------
Silkscreen font by Jason Kottke
http://kottke.org/plus/type/silkscreen/
