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
**Note**: The following instructions assume you have `go`, `make` and `npm` installled. `go` can be installed from [their website](https://go.dev/dl/]), while the easiest way to install the other 2 is to use a package manager: [chocolatey](https://chocolatey.org/install) for Windows, [Homebrew](https://brew.sh/) or apt-get/yum for Linux.

From there, you can install `make` and `node.js` (which includes npm) fairly easily. **Note**: make sure you install the **LTS** version of node. 

### GCC
Additionally, the [go-sqlite3](https://github.com/mattn/go-sqlite3) package `craig-stars` relies on for databasing itself requires the [GNU compiler collection](https://gcc.gnu.org/) to run. You **will** need GCC installed and built to run `craig-stars` locally!

On Linux/Mac, you can simply download the latest version of the GCC compilers using whatever package software you installed previously and build it from there. However, on Windows, you'll need to install a Linux-like development interface (use [MinGW64](https://www.mingw-w64.org/) - Cygwin64 has been known to cause issues) to install/build the latest GCC version due to file type restrictions.

## Assets
You will need art assets, and those can be downloaded from [https://craig-stars.net/images/images.zip](https://craig-stars.net/images/images.zip). Copy these images into the `frontend/static/images` folder (or make a new folder if it doesn't already exist).

## Installing Air
[Air](https://github.com/air-verse/air) is a Go utility that aids in automatic server restarting. While modifying the frontend code will trigger a hot reload of the local program (allowing for immediate confirmation of changes in real time), changes to the `golang` backend are only reflected the _next_ time the program is launched (requiring you to kill and restart the program each time). `Air` merely automates this "kill and restart" process by shutting down and reloading the server every time changes are detected. 

Installing air is **not required** to run `craig-stars` locally, but can be helpful if you plan on frequently making changes to the backend and want real time confirmation. 

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

If done correctly, it should give a localhost link (http://localhost:5173/) representing the application being hosted locally on your machine. Go to that site to see a live reloading frontend proxied to the go server on port `:8080`. Updating go code will relaunch the backend automatically (via air), while updating frontend code will do a hot reload with sveltekit/vite.

# Visual Studio Code 
[VS Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened with VS Code in order to use front end and backend plugins without issue in the same repo. It also comes with a built in terminal, debugging support, and about a thousand other bells and whistles useful for general software development. 


## backend
To launch the backend separately from the frontend, you can call `air` directly. 

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
"I try to click on the login button on localhost using the admin credentials and it does nothing! Worse, an error pops up in terminal!!!"

You probably aren't running the backend. Open a new terminal tab and type `air` to build the backend needed to handle all the nitty gritty logic stuff.

"When I run air, my computer complains about undefined Sqlite Drivers!"

See [Getting started](#Getting Started).

"Make is refusing to follow the instructions in the makefile!"
If worse comes to worst, try executing the instructions one by one (so do the instructions for make build, then make dev, etc etc). 
