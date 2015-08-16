# Requirements (Windows)

## ConEmu

This is optional, but having a good console with copy &amp; paste support,
etc will really help speed things up.  ConEmu has good integration with Git
Bash (you can open a window using Git Bash and keep all the nice terminal
features).

 - http://conemu.github.io/

## 7z

You need this for unpacking some projects and bundling.

 - Project http://www.7-zip.org/7z.html

## Git 1.9.5

In general, you just need git.  This also installs Git Bash which makes
running the support scripts much easier on Windows.

 - Download http://git-scm.com/download/win
 - Install to `C:\git`
 - Choose an option which puts git on `PATH` during installation
 - Follow https://help.github.com/articles/generating-ssh-keys/ for keys.
 - Follow https://help.github.com/articles/working-with-ssh-key-passphrases/
   so you don't have to keep entering passwords all the time.

## Mercurial 3.5.0

You need this to go-get some dependencies.

 - Download https://bitbucket.org/tortoisehg/files/downloads/
 - `mercurial-3.5.0-x64.msi`

## MinGW-w64 5.1.0

Contains a compiler and various tools.

 - Project http://mingw-w64.org/
 - http://mingw-w64.org/doku.php/download/mingw-builds
 - Links to SourceForge (ugh) http://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/
 - Get installer `mingw-w64-install.exe`
 - Extract to `C:\mingw`
 - Add `C:\mingw\mingw32\bin` to `PATH`

## GTK+ 2.24

Contains pkg-config.

 - Download http://ftp.acc.umu.se/pub/gnome/binaries/win32/gtk+/
 - Downloaded `/2.24/gtk+-bundle_2.24.10-20120208_win32.zip`
 - Extract to `C:\mingw\mingw32`
 - Don't overwrite files
 - Add a new system variable `PKG_CONFIG_PATH`, set to `C:\mingw\mingw32\lib\pkgconfig`

## Make 3.82.90

Allows the bundling scripts to work.

 - Download http://sourceforge.net/projects/mingw-w64/files/External%20binary%20packages%20%28Win64%20hosted%29/make/
  * Downloaded `make-3.82.90-20111115.zip`
  * Extract `bin_amd64/*` to `C:\mingw\mingw32\bin`

