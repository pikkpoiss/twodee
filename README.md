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

Then (both):

    go get github.com/banthar/gl
    go get github.com/jteeuwen/glfw

Setup
-----
To run the examples from the local source, run:

    ./setup_devel.sh

From the project root.  This will symlink the twodee source folder into your
$GOPATH, so the examples can be built from a local checkout.


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

