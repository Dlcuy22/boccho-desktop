# Boccho Desktop

A simple app to display animations and characters on your desktop.

## About This Project

This project is a simple desktop application for spawning animated characters on your desktop.

⚠ We did not include any model files in this repository.  
You can create your own models using the [boccho-toolkit](https://github.com/Dlcuy22/boccho-toolkit) or download/import character packs in `.bfk` format.

## Pack System (.bfk)

Boccho Desktop uses a custom pack format called **.bfk** (Boccho Frame Pack).
A `.bfk` file is essentially a ZIP-based archive that contains character folders with their respective animation frames (PNG/JPG).

### How to Import:

1. Open the application.
2. Click the **Import** button in the header.
3. Select **From .bfk** to browse for your pack file.
4. A preview modal will appear; click **Install** to extract the characters to your Frames directory.

## Requirements

This app requires SDL3 installed on your system.

- **Windows**: Download SDL3 and SDL3_image from the links below. You can place the DLLs in the same directory as the executable or add them to your `PATH`(e.g C:\Windows\System32).
  - [SDL3 Releases](https://github.com/libsdl-org/SDL/releases)
  - [SDL3_image Releases](https://github.com/libsdl-org/SDL_image/releases)

On linux you can install it using your package manager
Ubuntu

```bash
sudo apt install libsdl3-dev libsdl3-image-dev
```

Arch Linux

```bash
sudo pacman -S sdl3 sdl3_image
```

Fedora

```bash
sudo dnf install SDL3 SDL3_image
```

Later, I will add the DLLs to the Windows installer.

## Building

To build the application, run the following command:

```bash
wails build
```

To open dev mode (hot reload):

```bash
wails dev
```

## Simple - How does this app works?

This application uses the go-sdl3 library to create a window and render sprite frames.
Animations are implemented using a time-based update loop that cycles through frames at a fixed interval.

Each animation consists of a sequence of PNG images stored in a folder. These images are loaded into textures and rendered one by one to the SDL window, creating the illusion of motion.

The SDL window runs in a separate goroutine to ensure that the main Wails application remains responsive and does not freeze.
For implementation details, see AnimationEngine/Animation.go.

## Stack

- Golang
- React (TypeScript)
- Wails
