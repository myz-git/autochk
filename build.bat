@echo off
rem -----------------------------------------------------------------
rem build.bat : build autochk.exe (Windows) and autochk (linux/amd64)
rem -----------------------------------------------------------------
chcp 65001 >nul
setlocal EnableDelayedExpansion

rem -------- obtain git information ---------------------------------
where git >nul 2>nul
if errorlevel 1 (
  set "VERSION=0.0.0-dev"
  set "COMMIT=unknown"
) else (
  for /f "usebackq delims=" %%i in (`git describe --tags --always 2^>nul`) do set "VERSION=%%i"
  if "!VERSION!"=="" set "VERSION=0.0.0-dev"
  for /f "usebackq delims=" %%i in (`git rev-parse --short HEAD 2^>nul`) do set "COMMIT=%%i"
  if "!COMMIT!"=="" set "COMMIT=unknown"
)

rem -------- build date ---------------------------------------------
for /f "usebackq delims=" %%i in (`powershell -NoProfile -Command "$([DateTimeOffset]::Now.ToString('yyyy-MM-ddTHH:mm:ssK'))"`) do set "BUILD_DATE=%%i"

set "LDFLAGS=-s -w -X main.Version=!VERSION! -X main.Commit=!COMMIT! -X main.BuildDate=!BUILD_DATE!"

rem -------- Windows binary -----------------------------------------
echo === Building Windows binary ===
go build -ldflags="%LDFLAGS%" -o autochk.exe main.go
if errorlevel 1 goto :fail

rem -------- Linux binary -------------------------------------------
echo === Building Linux binary ===
set "GOOS=linux"
set "GOARCH=amd64"
set "CGO_ENABLED=0"
go build -ldflags="%LDFLAGS%" -o autochk main.go
if errorlevel 1 goto :fail

rem restore env
set "GOOS="
set "GOARCH="
set "CGO_ENABLED="

rem copy with version tag
copy /y autochk.exe "autochk-!VERSION!.exe" >nul
copy /y autochk     "autochk-!VERSION!"     >nul

echo Build succeeded. version=!VERSION! commit=!COMMIT! date=!BUILD_DATE!
exit /b 0

:fail
echo Build FAILED.    version=!VERSION! commit=!COMMIT!
exit /b 1
