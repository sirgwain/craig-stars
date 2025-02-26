# Local Development

craig-stars is a web based game. The backend logic and server is written in [Go](https://go.dev), while the frontend client is written in [TypeScript](https://www.typescriptlang.org) and powered by [`SvelteKit`](https://kit.svelte.dev).

## Prerequisites:

- Golang: 1.24 or higher, obtainable from [their website](https://go.dev/dl/)
- npm: [how to install](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- Respository forked and cloned on your device (instructions [here](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository))
- The [GNU compiler collection](https://gcc.gnu.org/) built locally on your device. Windows users can use [Mingw-w64](https://www.mingw-w64.org/), while linux/mac users can follow the [normal install instructions](https://gcc.gnu.org/install/index.html).

### Go Deps

After all that, you'll also need to install [Mage](https://github.com/magefile/mage), a make-like build tool/command executer helping to execute complex build commands.
It's included in the project's `go.mod` dependency tracker anyways, but using `go install` allows us to run it from the command line directly.

```bash
go install github.com/magefile/mage@latest
```

## Assets

You will also need art assets for ships and planets - otherwise they'll just look like black boxes. Thankfully, you can now download the images with a single magefile command! (For obvious reasons, this requires an internet connection.)

```bash
mage images
```

This will clear out the previous images folder before downloading the zip file and extracting it into `frontend/static/images`. **You only ever need to download the image files once.**

## Building and Launching

After performing all that setup, you should be good to go!
You have 2 methods to launch the server:

1. (Recommended) In VS Code, run the "Run Build Task" command (default keybinding `Ctrl+Shift+B`). This builds the server before launching the frontend and backend in separate terminals.
2. Run `mage run` from your terminal. This does essentially the same thing, but launches them inside the same terminal within separate goroutines. (_Note_: Don't worry if Mage complains about exceeding cleanup deadlines.)

On first launch, this will create an empty database in `./data` with a single `admin` user (password `admin`). (If it fails, try clearing the data folder and trying again.)

With some luck, you should get a localhost link from npm (http://localhost:5173/) representing the application being hosted locally on your machine. Go to that site to see a live-reloading frontend proxied to the go server on port `:8080`. Updating Go code (backend) will kill & restart the backend automatically via air, while updating Svelte or Typescript code (frontend) will perform a hot reload with sveltekit/vite.

### Launching Backend/Frontend only

If one wants to launch the backend or frontend separately (such as to have both processes in separate terminals), there are mage commands to launch them separately.

```bash
mage launch_frontend
```

```bash
mage launch_backend
```

(For those curious, this is how VS Code launches the server.)

# Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened with VS Code in order to use frontend and backend plugins without issue in the same repo.
It also contains [tasks.json](/.vscode/tasks.json) and [launch.json](/.vscode/tasks.json) files containing various prebuilt commands and debug configurations. (There's a build task to build & launch the entire server in 1 button press.)
It also comes with a built in terminal, debugging support, and an array of assorted bells and whistles useful for general software development.

## Running Tests

While manual local dev testing is good, software testing & debugging are also crucial to ensure things run (and continue to run) smoothly.\
`craig-stars` makes use of 3 different software testing providers:

- [gotestsum](https://github.com/gotestyourself/gotestsum) for backend Golang unit tests. This runs `go test` under the hood
  - _Note_: `VSCode-Go` doesn't currently support running alternate test tools, so running tests from within VS Code's UI will just use regular old `go test`.
- [Vitest](https://vitest.dev/guide/cli.html) for frontend unit tests.
- [Playwright](https://playwright.dev/docs/running-tests) for end-to-end integration tests.

After writing new or updating existing tests, there are several options as for how to run them:

- Run tests from the command line:
  - `mage test` to run everything at once. Great for overall checks to make sure everything works, bad for specific problem fixes.
  - `mage test_golang`, `mage test_vitest` and `mage test_playwright` to run tests for a given test provider. Each passes their arguments directly to the test provider.
  - Protip: to run test functions matching a regex, run `mage test_backend --run="XXX"`, `mage test-frontend -- XXX`.
- Run tests using th
- Run/debug using the Testing panel in the activity bar - tests can be filtered by result, directory, etc.
- Run/debug using the small buttons displayed in test files and next to test functions.

NOTE: VSCode's Test Explorer has been known to adversely affect test performance. If your tests are failing due to timing out, try increasing the "Test timeout" variable in your settings.

# Troubleshooting

See the [Troubleshooting](troubleshooting) page for solutions to some common local dev issues.
