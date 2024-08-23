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
**Note**: The following instructions assume you have `go`, `make` and `npm` installled. `go` can be installed quite easily from [their website](https://go.dev/dl/]), while the easiest way to install the other 2 is to use a package manager: [chocolatey](https://chocolatey.org/install) for Windows, [Homebrew](https://brew.sh/) or apt-get/yum for Linux (& Chromebook by extension).

From there, you can install `make` and `node.js` (which includes npm) fairly easily. NOTE: If `node.js` doesn't co-operate after being installed, try uninstalling it and re-installing the **LTS** version from their website instead.

After cloning the repo (easily done with Github Desktop), build `craig-stars` locally once. This will build the frontend and then the backend. The backend embeds the frontend resources into its binary to support a single binary deployment. 

### Installing go-sqlite3
If you don't already have it, you'll also need to install [go-sqlite3](https://github.com/mattn/go-sqlite3) for things to work. However, this itself requires the [GCC compiler collection](https://gcc.gnu.org/) to function. 

For Linux/Mac, you can simply install the latest version of the GCC compilers using whatever package software you installed previously.

However, on Windows, you'll need to install a Linux-like development interface (use [MinGW](https://www.mingw-w64.org/) - Cygwin64 WILL NOT WORK) to install & build the latest GCC version.

Once that's done, enter `go install github.com/mattn/gosqlite3` into your computer's native terminal (or the terminal in VS Code) to install the driver.

## Assets
You will need art assets, and those can be downloaded from [https://craig-stars.net/images/images.zip](https://craig-stars.net/images/images.zip). Copy these images into the `frontend/static/images` folder (or make a new folder if it doesn't already exist).

## Building and Running
Finally, go to your terminal and enter the following commands:

```bash
make build
```

Install air for automatic server restarts while developing: 

```bash
go install github.com/air-verse/air@latest
```

Launch the frontend and backend at the same time with `make`
```bash
make dev
```

**Note** On first launch, this will create an empty database with a single `admin` user, password `admin`.

Enter `npm run preview` to view the localhost link (localhost:XXXX) representing the application being hosted locally on your machine. Go to that site to see a live reloading frontend proxied to the go server on port `:8080`. Updating go code will relaunch the backend automatically, while updating frontend code will do a hot reload with sveltekit/vite. (For obvious reasons, this only lasts while the terminal is actually _open_; closing it will close the program.)

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
