# Local Development

craig-stars is a web based game. The backend logic and server is written in [Go](https://go.dev), while the frontend client is written in [TypeScript](https://www.typescriptlang.org) and powered by [`SvelteKit`](https://kit.svelte.dev).

## Prerequisites:

- Golang: 1.24 or higher, obtainable from [their website](https://go.dev/dl/)
- npm: [how to install](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- Respository forked and cloned on your device (instructions [here](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository))
- The [GNU compiler collection](https://gcc.gnu.org/) built locally on your device. Windows users can use [Mingw-w64](https://www.mingw-w64.org/), while linux/mac users can follow the [normal install instructions](https://gcc.gnu.org/install/index.html).

After all that, you'll also need to install 2 important Go dependencies:

- [Mage](https://github.com/magefile/mage), a make-like build tool/command executer helping to execute complex build commands.
- [Air](https://github.com/air-verse/air), a server utility aiding with automatic backend server restarting.
- [tygo](https://github.com/gzuidhof/tygo), a generator for creating typescript types for golang types.

Both can be installed using a single `go install` command:

```bash
go install github.com/magefile/mage@latest github.com/air-verse/air@latest github.com/gzuidhof/tygo@latest
```

(Mage should be included in go.mod regardless, but it never hurts to make sure it's there.)

## Assets

You will also need art assets for ships and planets - otherwise they'll just look like black boxes. Thankfully, you can now download the images with a single magefile command! (For obvious reasons, this requires an internet connection.)

```bash
mage images
```

This will clear out the previous images folder before downloading the zip file and extracting it into `frontend/static/images`. **You only ever need to download the image files once.**

## Building and Launching

After performing all that setup, you're good to go!
Head to the repo's root folder in your terminal and enter the following command to build and launch the server:

```bash
mage run
```

**Note** On first launch, this will create an empty database with a single `admin` user, password `admin`.

If setup correctly, you should get a localhost link from npm (http://localhost:5173/) representing the application being hosted locally on your machine. Go to that site to see a live-reloading frontend proxied to the go server on port `:8080`. Updating Go code (backend) will kill & restart the backend automatically via air, while updating Svelte or Typescript code (frontend) will perform a hot reload with sveltekit/vite.

### Launching Backend/Frontend only

If one wants to launch the backend or frontend separately (such as to have both processes in separate terminals), there are mage commands to launch them independently.

```bash
mage launch_frontend
```

```bash
mage launch_backend
```

# Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened with VS Code in order to use frontend and backend plugins without issue in the same repo, as well as [tasks.json](/..vscode/tasks.json) and [launch.json](/..vscode/tasks.json) files containing various prebuilt commands and debug configurations.
It also comes with a built in terminal, debugging support, and an array of assorted bells and whistles useful for general software development.

## Testing

While manual local dev testing is good, software testing & debugging are also crucial to ensure things run (and continue to run) smoothly. After writing new or updating existing tests, there are several options as for how to run them:

- Run `mage test` to run all the tests at once (great for overall checks to make sure everything works, bad for specific problem fixes)
- Run tests via command line (`go test` and `npm run test` for backend/frontend respectively, followed by the specific test file name(s) for specific coverage)
- Run/debug using the Testing panel in the activity bar - tests can be filtered by result, directory, etc.
- Run/debug using the small buttons displayed in test files and next to test functions.

NOTE: VSCode's Test Explorer has been known to adversely affect test performance. If your tests are failing due to timing out, try increasing the "Test timeout" variable in your settings.

# Troubleshooting

See the [Troubleshooting](troubleshooting) page for solutions to some common local dev issues.
