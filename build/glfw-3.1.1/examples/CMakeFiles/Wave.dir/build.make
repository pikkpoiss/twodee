# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.0

#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:

# Remove some rules from gmake that .SUFFIXES does not remove.
SUFFIXES =

.SUFFIXES: .hpux_make_needs_suffix_list

# Suppress display of executed commands.
$(VERBOSE).SILENT:

# A target that is always out of date.
cmake_force:
.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/local/Cellar/cmake/3.0.2/bin/cmake

# The command to remove a file.
RM = /usr/local/Cellar/cmake/3.0.2/bin/cmake -E remove -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /Users/kurrik/workspace/twodee/build/glfw-3.1.1

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /Users/kurrik/workspace/twodee/build/glfw-3.1.1

# Include any dependencies generated for this target.
include examples/CMakeFiles/Wave.dir/depend.make

# Include the progress variables for this target.
include examples/CMakeFiles/Wave.dir/progress.make

# Include the compile flags for this target's objects.
include examples/CMakeFiles/Wave.dir/flags.make

examples/CMakeFiles/Wave.dir/wave.c.o: examples/CMakeFiles/Wave.dir/flags.make
examples/CMakeFiles/Wave.dir/wave.c.o: examples/wave.c
	$(CMAKE_COMMAND) -E cmake_progress_report /Users/kurrik/workspace/twodee/build/glfw-3.1.1/CMakeFiles $(CMAKE_PROGRESS_1)
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Building C object examples/CMakeFiles/Wave.dir/wave.c.o"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -o CMakeFiles/Wave.dir/wave.c.o   -c /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples/wave.c

examples/CMakeFiles/Wave.dir/wave.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/Wave.dir/wave.c.i"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -E /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples/wave.c > CMakeFiles/Wave.dir/wave.c.i

examples/CMakeFiles/Wave.dir/wave.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/Wave.dir/wave.c.s"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -S /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples/wave.c -o CMakeFiles/Wave.dir/wave.c.s

examples/CMakeFiles/Wave.dir/wave.c.o.requires:
.PHONY : examples/CMakeFiles/Wave.dir/wave.c.o.requires

examples/CMakeFiles/Wave.dir/wave.c.o.provides: examples/CMakeFiles/Wave.dir/wave.c.o.requires
	$(MAKE) -f examples/CMakeFiles/Wave.dir/build.make examples/CMakeFiles/Wave.dir/wave.c.o.provides.build
.PHONY : examples/CMakeFiles/Wave.dir/wave.c.o.provides

examples/CMakeFiles/Wave.dir/wave.c.o.provides.build: examples/CMakeFiles/Wave.dir/wave.c.o

# Object files for target Wave
Wave_OBJECTS = \
"CMakeFiles/Wave.dir/wave.c.o"

# External object files for target Wave
Wave_EXTERNAL_OBJECTS =

examples/Wave.app/Contents/MacOS/Wave: examples/CMakeFiles/Wave.dir/wave.c.o
examples/Wave.app/Contents/MacOS/Wave: examples/CMakeFiles/Wave.dir/build.make
examples/Wave.app/Contents/MacOS/Wave: src/libglfw3.a
examples/Wave.app/Contents/MacOS/Wave: examples/CMakeFiles/Wave.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --red --bold "Linking C executable Wave.app/Contents/MacOS/Wave"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/Wave.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
examples/CMakeFiles/Wave.dir/build: examples/Wave.app/Contents/MacOS/Wave
.PHONY : examples/CMakeFiles/Wave.dir/build

examples/CMakeFiles/Wave.dir/requires: examples/CMakeFiles/Wave.dir/wave.c.o.requires
.PHONY : examples/CMakeFiles/Wave.dir/requires

examples/CMakeFiles/Wave.dir/clean:
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples && $(CMAKE_COMMAND) -P CMakeFiles/Wave.dir/cmake_clean.cmake
.PHONY : examples/CMakeFiles/Wave.dir/clean

examples/CMakeFiles/Wave.dir/depend:
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1 && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /Users/kurrik/workspace/twodee/build/glfw-3.1.1 /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples /Users/kurrik/workspace/twodee/build/glfw-3.1.1 /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples /Users/kurrik/workspace/twodee/build/glfw-3.1.1/examples/CMakeFiles/Wave.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : examples/CMakeFiles/Wave.dir/depend

