@echo off
echo Creating exact directory structure...

mkdir client
mkdir ui
mkdir db
mkdir config
mkdir models
mkdir errors
mkdir logger
mkdir internal
mkdir examples

echo.
echo Structure created successfully!
echo.
echo Created directories:
dir /b /ad
echo.
pause