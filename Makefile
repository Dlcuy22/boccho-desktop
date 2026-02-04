# Makefile
# Build and development automation for Boccho.
#
# Targets:
#   - run: runs the application locally
#   - build: builds the binary for windows amd64
#   - build-windows: builds the binary for windows amd64
#   - [LATER] build-all: builds binaries for Linux, Windows, and macOS (multiple architectures)

BINARY_NAME=Boccho
VERSION=1.0.0

run: 
	@echo Running Boccho...
	wails dev

build:
	@echo Building Boccho...
	wails build
	@echo Done, see ./build/bin

build-windows:
	@echo Building Boccho for Windows...
	wails build -nsis
	@echo Done, see ./build/windows

clean: 
	@echo Cleaning up...
	rm -rf build/bin/*
	@echo Done