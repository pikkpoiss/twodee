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
include tests/CMakeFiles/empty.dir/depend.make

# Include the progress variables for this target.
include tests/CMakeFiles/empty.dir/progress.make

# Include the compile flags for this target's objects.
include tests/CMakeFiles/empty.dir/flags.make

tests/CMakeFiles/empty.dir/empty.c.o: tests/CMakeFiles/empty.dir/flags.make
tests/CMakeFiles/empty.dir/empty.c.o: tests/empty.c
	$(CMAKE_COMMAND) -E cmake_progress_report /Users/kurrik/workspace/twodee/build/glfw-3.1.1/CMakeFiles $(CMAKE_PROGRESS_1)
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Building C object tests/CMakeFiles/empty.dir/empty.c.o"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -o CMakeFiles/empty.dir/empty.c.o   -c /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests/empty.c

tests/CMakeFiles/empty.dir/empty.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/empty.dir/empty.c.i"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -E /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests/empty.c > CMakeFiles/empty.dir/empty.c.i

tests/CMakeFiles/empty.dir/empty.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/empty.dir/empty.c.s"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -S /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests/empty.c -o CMakeFiles/empty.dir/empty.c.s

tests/CMakeFiles/empty.dir/empty.c.o.requires:
.PHONY : tests/CMakeFiles/empty.dir/empty.c.o.requires

tests/CMakeFiles/empty.dir/empty.c.o.provides: tests/CMakeFiles/empty.dir/empty.c.o.requires
	$(MAKE) -f tests/CMakeFiles/empty.dir/build.make tests/CMakeFiles/empty.dir/empty.c.o.provides.build
.PHONY : tests/CMakeFiles/empty.dir/empty.c.o.provides

tests/CMakeFiles/empty.dir/empty.c.o.provides.build: tests/CMakeFiles/empty.dir/empty.c.o

tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o: tests/CMakeFiles/empty.dir/flags.make
tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o: deps/tinycthread.c
	$(CMAKE_COMMAND) -E cmake_progress_report /Users/kurrik/workspace/twodee/build/glfw-3.1.1/CMakeFiles $(CMAKE_PROGRESS_2)
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Building C object tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -o CMakeFiles/empty.dir/__/deps/tinycthread.c.o   -c /Users/kurrik/workspace/twodee/build/glfw-3.1.1/deps/tinycthread.c

tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/empty.dir/__/deps/tinycthread.c.i"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -E /Users/kurrik/workspace/twodee/build/glfw-3.1.1/deps/tinycthread.c > CMakeFiles/empty.dir/__/deps/tinycthread.c.i

tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/empty.dir/__/deps/tinycthread.c.s"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -S /Users/kurrik/workspace/twodee/build/glfw-3.1.1/deps/tinycthread.c -o CMakeFiles/empty.dir/__/deps/tinycthread.c.s

tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.requires:
.PHONY : tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.requires

tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.provides: tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.requires
	$(MAKE) -f tests/CMakeFiles/empty.dir/build.make tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.provides.build
.PHONY : tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.provides

tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.provides.build: tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o

# Object files for target empty
empty_OBJECTS = \
"CMakeFiles/empty.dir/empty.c.o" \
"CMakeFiles/empty.dir/__/deps/tinycthread.c.o"

# External object files for target empty
empty_EXTERNAL_OBJECTS =

tests/empty.app/Contents/MacOS/empty: tests/CMakeFiles/empty.dir/empty.c.o
tests/empty.app/Contents/MacOS/empty: tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o
tests/empty.app/Contents/MacOS/empty: tests/CMakeFiles/empty.dir/build.make
tests/empty.app/Contents/MacOS/empty: src/libglfw3.a
tests/empty.app/Contents/MacOS/empty: tests/CMakeFiles/empty.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --red --bold "Linking C executable empty.app/Contents/MacOS/empty"
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/empty.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
tests/CMakeFiles/empty.dir/build: tests/empty.app/Contents/MacOS/empty
.PHONY : tests/CMakeFiles/empty.dir/build

tests/CMakeFiles/empty.dir/requires: tests/CMakeFiles/empty.dir/empty.c.o.requires
tests/CMakeFiles/empty.dir/requires: tests/CMakeFiles/empty.dir/__/deps/tinycthread.c.o.requires
.PHONY : tests/CMakeFiles/empty.dir/requires

tests/CMakeFiles/empty.dir/clean:
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests && $(CMAKE_COMMAND) -P CMakeFiles/empty.dir/cmake_clean.cmake
.PHONY : tests/CMakeFiles/empty.dir/clean

tests/CMakeFiles/empty.dir/depend:
	cd /Users/kurrik/workspace/twodee/build/glfw-3.1.1 && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /Users/kurrik/workspace/twodee/build/glfw-3.1.1 /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests /Users/kurrik/workspace/twodee/build/glfw-3.1.1 /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests /Users/kurrik/workspace/twodee/build/glfw-3.1.1/tests/CMakeFiles/empty.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : tests/CMakeFiles/empty.dir/depend

