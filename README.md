# Boccho Desktop

A simple app to display animations and characters on your desktop.

## About This Project

This project is a simple desktop application for spawning animated characters on your desktop.

⚠ We did not include any model files in this repository.  
You need to create your own model using the [boccho-toolkit](https://github.com/Dlcuy22/boccho-toolkit) repository.  
The model consists of frames with transparent backgrounds placed in the `./Frames` folder.

This app require SDL3 Installed in your system, if you are in windows go ahead and download SDL3 and SDL3_image here 
[SDL](https://github.com/libsdl-org/SDL_image/releases)
[SDL_image](https://github.com/libsdl-org/SDL/releases)
You can either install them to your PATH (for example C:\Windows\System32) or place the DLL files in the same directory as the executable.


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
````
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

* Golang
* React (TypeScript)
* Wails

