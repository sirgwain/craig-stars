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

## GCC
Additionally, the [go-sqlite3](https://github.com/mattn/go-sqlite3) package `craig-stars` relies on for databasing itself requires the [GNU compiler collection](https://gcc.gnu.org/) to run. You **will** need GCC installed and built to run `craig-stars` locally!

On Linux/Mac, you can simply download the latest version of the GCC compilers using whatever package software you installed previously and build it from there. However, on Windows, you'll need to install a Linux-like development interface (use [MinGW64](https://www.mingw-w64.org/) - Cygwin64 has been known to cause issues) to install/build the latest GCC version due to file type restrictions.

## Assets
You will need art assets for ships and planets - otherwise they'll just look like black boxes.Thankfully, you can now download the images with a single makefile command! (Since the image files are stored on the cloud, this requires an internet connection.) 
```bash
make images
```
The command will clear out the previous images folder before downloading and extracting `images.zip` (https://craig-stars.net/images/images.zip) into `frontend/static/images`

**You only ever need to download the image files once.** 
Of course, if any new items or icons come along that need to be downloaded, you'll need to update it with the new files if you want to see them in local host. 

## Installing Air
[Air](https://github.com/air-verse/air) is a Go utility that aids in automatic server restarting. While modifying the frontend code will trigger a hot reload of the local program (allowing for immediate confirmation of changes in real time), changes to the `golang` backend are only reflected the _next_ time the program is launched (requiring you to kill and restart the program each time). `Air` helps automate this "kill and restart" process by shutting down and reloading the server every time changes are detected. 

Installing air is **not required** to run craig-stars locally, but can be helpful if you plan on making frequent backup changes and want both real time confirmation and maximum laziness.

To install air, enter the following code into your terminal:
```bash
go install github.com/air-verse/air@latest
```

## Building and Running
After performing all that, go to your terminal and enter the following command:

```bash
make run
```

**Note** On first launch, this will create an empty database with a single `admin` user, password `admin`.

If done correctly, it should give a localhost link (http://localhost:5173/) representing the application being hosted locally on your machine. Go to that site to see a live reloading frontend proxied to the go server on port `:8080`. Updating Go code (backend) will re-launch the backend automatically (via air), while updating Svelte or Typescript code (frontend) will perform a hot reload with sveltekit/vite.

### Backend/Frontend only
To launch the backend separately from the frontend, you can call `air` directly (or equivalently, run the `dev_backend` makefile recipe which does just that). 

```bash
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

If one wants to exclusively launch the frontend, it can be run in development mode with npm:

```bash
make dev_frontend
```

# Visual Studio Code 
[VS Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened with VS Code in order to use frontend and backend plugins without issue in the same repo. It also comes with a built in terminal, debugging support, and an array of assorted bells and whistles useful for general software development.

## Testing
While manual local dev testing is good, software testing & debugging are also crucial to ensure things run (and continue to run) smoothly. After writing new or updating existing tests, there are several options as for how to run them:
* Run `make test` to run all the tests at once (great for overall checks to make sure everything works, bad for specific debugging)
* Run tests via command line (`go test` and `npm run test` for frontend/backend respectively, followed by the specific test file name(s) for specific coverage)
* Run/debug using the Testing panel in the activity bar - tests can be filtered by result, directory, etc. 
* Run/debug using the small buttons displayed in test files and next to test functions.

# Troubleshooting
See [Common Problems](faq#common-problems) in the FAQ for solutions to some common local dev issues.