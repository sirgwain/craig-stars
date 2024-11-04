# Development
craig-stars is a web based game. The backend is written in [golang](https://go.dev). The frontend is written in [typescript](https://www.typescriptlang.org) and powered by [`SvelteKit`](https://kit.svelte.dev).

## Architecture
For detailed information about `craig-stars` architecture, check out the [architecture](architecture.md) page.

## Tech Stack
`craig-stars` is built on top of the following excellent technologies:

- [golang](https://go.dev)
- [sveltekit](https://kit.svelte.dev) (with static adaptor)
- [sqlx](https://github.com/jmoiron/sqlx) + [sqlite](https://www.sqlite.org)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goverter](https://github.com/jmattheis/goverter)
- [go-chi](https://github.com/chi/go-chi)
- [zerolog](https://github.com/rs/zerolog)
- [disgo](https://github.com/disgoorg/disgo)
- [go-pkgz/auth](https://github.com/go-pkgz/auth)
- [cobra](https://github.com/spf13/cobra) (for cli)
- [viper](https://github.com/spf13/viper) (for config)
- [tailwindcss](https://tailwindcss.com)
- [daisyui](https://daisyui.com)

Icons are either hand crafted, taken from the original Stars! files or from the wonderful [game-icons.net](https://game-icons.net) and [heroicons.com](https://heroicons.com).

# Getting Started
**Note**: The following instructions assume you have `go`, `make` and `npm` installled. `go` can be installed from [their website](https://go.dev/dl/]), while the easiest way to install the other 2 is to use a package manager: [chocolatey](https://chocolatey.org/install) for Windows, [Homebrew](https://brew.sh/) for Mac or apt-get/yum for Linux.

From there, you can install `make` and `node.js` (which includes npm) fairly easily. **Note**: make sure you install the **LTS** version of node. 

If you're on Windows, *make sure to have Powershell installed and in your PATH* for the makefile to function correctly. In 99% of cases, this shouldn't be an issue (powershell is bundled with all currently supported Windows versions to date and has been the default shell since Windows 10), but if your machine somehow doesn't have it, use winget in command prompt (`winget install --id Microsoft.PowerShell --source winget`) to install the latest version. (On Mac/Linux, any terminal should be fine as long as it supports the `mkdir`, `rm`, `cp` and `date` instructions.)

### GCC
Additionally, the [go-sqlite3](https://github.com/mattn/go-sqlite3) package `craig-stars` relies on for databasing itself requires the [GNU compiler collection](https://gcc.gnu.org/) to run. You **will** need GCC installed and built to run `craig-stars` locally!

On Linux/Mac, you can simply download the latest version of the GCC compilers using whatever package software you installed previously and build it from there. However, on Windows, you'll need to install a Linux-like development interface (use [MinGW64](https://www.mingw-w64.org/) - Cygwin64 has been known to cause issues) to install/build the latest GCC version due to file type restrictions.

## Assets
You will need art assets for ships and such, and those can be downloaded from [https://craig-stars.net/images/images.zip](https://craig-stars.net/images/images.zip). Copy these images into the `frontend/static/images` folder (or make a new folder if it doesn't already exist).

## Installing Air
[Air](https://github.com/air-verse/air) is a Go utility that aids in automatic server restarting. While modifying the frontend code will trigger a hot reload of the local program (allowing for immediate confirmation of changes in real time), changes to the `golang` backend are only reflected the _next_ time the program is launched (requiring you to kill and restart the program each time). `Air` helps automate this "kill and restart" process by shutting down and reloading the server every time changes are detected. 

Installing air is **not required** to run craig-stars locally, but can be helpful if you plan on frequently making changes to the backend and want real time confirmation.

To install air, enter the following code into your terminal:
```bash
go install github.com/air-verse/air@latest
```

## Building and Running
After performing all that, go to your terminal and enter the following commands:

```bash
make build
make dev
```

**Note** On first launch, this will create an empty database with a single `admin` user, password `admin`.

If done correctly, it should give a localhost link (http://localhost:5173/) representing the application being hosted locally on your machine. Go to that site to see a live reloading frontend proxied to the go server on port `:8080`. Updating go code will re-launch the backend automatically (via air), while updating frontend code will do a hot reload with sveltekit/vite.

# Visual Studio Code 
[VS Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened with VS Code in order to use frontend and backend plugins without issue in the same repo. It also comes with a built in terminal, debugging support, and an array of assorted bells and whistles useful for general software development.

## backend
To launch the backend separately from the frontend, you can call `air` directly (or equivalently, run the `dev_backend` makefile recipe). 

```zsh
❯ air

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ , built with Go

watching .
watching ai
watching cmd
watching config
watching cs
!exclude data
watching db
!exclude dist
!exclude frontend
watching server
watching test
!exclude tmp
!exclude vendor
building...
running...
7:47AM DBG Debug logging enabled
```

## frontend
Launch the frontend in development mode with npm:

```zsh
cd frontend
npm run dev
```

## test
Run tests:

```zsh
make test
```

# Troubleshooting
"I try to click on the login button on localhost using the admin credentials and it does nothing! Also, an error pops up in my terminal!"

You probably aren't running the backend. Open a new terminal tab and type `air` to build the backend needed to handle all the nitty gritty logic stuff.

"When I run air, my computer complains about undefined Sqlite Drivers!"

See [Getting started](#Getting Started).

"When I run make build, I get an obscure error about 'executable not found in %PATH%' or 'build target excluding all files in XXX'!"
What's probably happening is you're trying to generate go files or build the server with the incorrect GOARCH and GOOS settings. Try using `go env -u GOOS GOARCH` to reset them to their defaults and see if the problems persist.

"I tried to do all of the above, but I'm still getting errors in the command line!"
Consult this ordered checklist of vague general suggestions:
1. Read the error message to try and figure out why it's failing. Make (as it's being used here) effectively just copy-pastes its commands into the terminal one by one (expanding variables here and there), so syntax errors in the terminal can oftentimes reflect issues with the makefile itself rather than any device-related issues. 
2. Try and search online for the error message or similar problems to see if others may have found solutions to the problem for you. 
3. If all else fails, reach out in the #stars-clones or #craig-stars channels in the discord (ideally with images/text of the commands used and/or resulting error messages).