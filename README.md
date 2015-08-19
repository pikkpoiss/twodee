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

## Supported platforms

This library has been tested by developing games for:

 - OSX 10.10 Yosemite using x86_64 CPUs.
 - Windows 8.1 using x86_64 CPUs.

It should support:

 - Linux variants using x86_64 CPUs.

Other versions may work but have not been tested.  32-bit CPUs are not
supported.

## Requirements

 - [Windows](docs/requirements_win.md)
 - [OSX](docs/requirements_osx.md)

## Setup

This project is both a library and a set of support scripts which
will install dependencies needed to make the library work correctly.

To run the support scripts you must check out the library:

    git clone https://github.com/pikkpoiss/twodee.git
    cd twodee

This library depends on various C/C++ packages.  To try and get a stable
environment for building, sources have been included in the `lib` directory.
Running the following will build each package and attempt to install a
Golang wrapper:

    ./scripts/setup.sh

*On Windows* this will build shared libraries in `.dll` format and install
them to `build/usr/bin`.  These DLLs _must_ be included next to a packaged
executable in order for the built application to be portable.

*On OSX* this will build static libraries in `.a` format and install them
to `build/usr/lib`. Packaged executables linked against these libraries
should be portable without needing to package any shared libraries.

## Using

The twodee lib is not installed with `go get` like most Go libraries. Instead,
bootstrap a project using the provided scripts.

### Initializing a project

Once you have the libraries installed, run:

    ./scripts/setup_project.sh PATH_TO_PROJECT_DIRECTORY

This will set up a project structure, install support scripts, and copy
shared libraries if needed.

Once the project is set up, commit the files to git as needed and push to
origin.

If you rebuild the twodee libs, run the `setup_project.sh` script again.
Be careful though, as it will overwrite any modifications.

### Checking out the project

Other developers will need to follow the steps to set up the twodee library
before checking out the other project.  Once that is complete, they must
check out the project from git and run the following from the project root:

    git submodule init
    git submodule update

### Running the project

Run the project with the following command:

    make run

### Packaging the project

Build a distributable bundle with:

    make package

This will place a zip file in the `build/PLATFORM` directory with a bundled
executable.
