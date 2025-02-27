# Local Development

craig-stars is a web based game. The backend logic and server is written in [Go](https://go.dev), while the frontend client is written in [TypeScript](https://www.typescriptlang.org) and powered by [`SvelteKit`](https://kit.svelte.dev).

## Prerequisites:

- Golang: 1.24 or higher, obtainable from [their website](https://go.dev/dl/)
- npm: [how to install](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- Respository forked and cloned on your device (instructions [here](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository))
- The [GNU compiler collection](https://gcc.gnu.org/) built locally on your device. Windows users can use [Mingw-w64](https://www.mingw-w64.org/), while linux/mac users can follow the [normal install instructions](https://gcc.gnu.org/install/index.html).

### Go Deps

After all that, you'll also need to install [Mage](https://github.com/magefile/mage), a make-like build tool/command executer written in Go for execution of complex build commands.
It's included in the project's `go.mod` dependency tracker anyways, but using `go install` allows us to run it from the command line directly.

```bash
go install github.com/magefile/mage@latest
```

**Disclaimer**: Magefile commands must be run from the repository root (otherwise mage won't find the files).

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
2. Run `mage run` from your terminal inside the root folder. This does essentially the same thing, but launches them inside the same terminal within separate goroutines. (_Note_: Don't worry if Mage complains about exceeding cleanup deadlines.)

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

(For those curious, this is how the aforementioned build task launches the server.)

# Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened inside VS Code in order to use frontend and backend plugins without issue in the same repo. The repository also contains [tasks.json](/.vscode/tasks.json) and [launch.json](/.vscode/tasks.json) files containing various prebuilt commands and debug configurations.

# Testing

<!--? Do we need to move this to its own section? -->

While manual local dev testing is certainly valuable, software testing & debugging are also crucial to ensure things run (and continue to run) smoothly.\
`craig-stars` makes use of 3 different automated software testing providers:

- [gotestsum](https://github.com/gotestyourself/gotestsum) for backend Golang unit tests. This runs `go test` under the hood and does fancy formatting on the output.
- [Vitest](https://vitest.dev/) for frontend unit tests.
- [Playwright](https://playwright.dev/) for end-to-end integration/UI tests.

Each provider comes with its [own](../gotestsum) [config](../frontend/vite.config.ts) [files](../frontend/playwright.config.ts), with varying settings for CI and non-CI runs.

## Running & Debugging tests

After writing new or updating existing tests, there are several options as for how to run them.

- Run `mage test` to run everything at once. Great for overall checks to make sure everything works, bad for specific problem fixes.
- Run `mage test_golang`, `mage test_vitest` and `mage test_playwright` to run tests for a given test provider at a time. Each passes their arguments directly to the test provider, so you can pass all the same arguments to them as you would to `go test` or `vitest`. (Test reports are saved to `tmp/test-results` as JUnit XML files.)
  - Protip: To test only files matching a specific file name or regex, you can use the `--run=` flag for `go test` or simply enter the test file name for vitest & playwright.
- Run the various test tasks inside `tasks.json` (the green ones with icons). There's 4 in total, one for each of the above mage commands.
- Run tests from VS Code's UI, via either the Test Explorer panel or the small buttons displayed within test files.
  - Unfortunately, `vscode-go` doesn't currently support running alternate test tools for UI commands, so running backend tests this way will just use plain old `go test`.

_NOTE_: Slower devices may have trouble running backend tests within the default timeout of 30s, especially ones inside `./server` involving repeated serialization to & from the database. If your tests are routinely timing out while succeeding on CI, consider increasing the "Go: Test Timeout" variable in your local settings.

# Troubleshooting

See the [Troubleshooting](troubleshooting) page for solutions to some common local dev issues.
