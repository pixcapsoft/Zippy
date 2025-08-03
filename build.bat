@echo off
setlocal

REM ====== CONFIGURE METADATA HERE ======
set VERSION=1.0.0
set AUTHOR=PixCap Soft

REM ====== CLEAN OLD BUILDS ======
rmdir /s /q build 2>nul

REM ====== WINDOWS BUILDS ======
set GOOS=windows
for %%A in (amd64 386) do (
    set GOARCH=%%A
    call :winbuild %%A
)
goto afterwin

:winbuild
set OUTDIR=build\windows_%1
mkdir %OUTDIR% 2>nul
go build -o %OUTDIR%\zippy.exe -ldflags "-X \"main.ZippyVersion=%VERSION%\" -X \"main.ZippyAuthor=%AUTHOR%\""
goto :eof

:afterwin

REM ====== LINUX BUILDS ======
set GOOS=linux
for %%A in (amd64 386) do (
    set GOARCH=%%A
    call :linbuild %%A
)
goto afterlin

:linbuild
set OUTDIR=build\linux_%1
mkdir %OUTDIR% 2>nul
go build -o %OUTDIR%\zippy
goto :eof

:afterlin

REM ====== MACOS BUILDS ======
set GOOS=darwin
for %%A in (amd64) do (
    set GOARCH=%%A
    call :macbuild %%A
)
goto aftermac

:macbuild
set OUTDIR=build\macos_%1
mkdir %OUTDIR% 2>nul
go build -o %OUTDIR%\zippy
goto :eof

:aftermac

echo.
echo All builds complete! Check the build\ folders.
endlocal
pause