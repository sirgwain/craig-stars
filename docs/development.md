# Development
craig-stars is a web based game. The backend is written in [golang](https://go.dev). The frontend is written in [typescript](https://www.typescriptlang.org) and powered by [`SvelteKit`](https://kit.svelte.dev).

## Tech Stack
`craig-stars` is built on top of the following excellent technologies:

- [golang](https://go.dev)
- [sveltekit](https://kit.svelte.dev) (with static adaptor)
- [sqlx](https://github.com/jmoiron/sqlx) + [sqlite](https://www.sqlite.org)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [go-chi](https://github.com/chi/go-chi)
- [zerolog](https://github.com/rs/zerolog)
- [disgo](https://github.com/disgoorg/disgo)
- [go-pkgz/auth](https://github.com/go-pkgz/auth)
- [cobra](https://github.com/spf13/cobra) (for cli)
- [viper](https://github.com/spf13/viper) (for config)
- [tailwindcss](https://tailwindcss.com)
- [daisyui](https://daisyui.com)

Icons are either hand crafted or from the wonderful [game-icons.net](https://game-icons.net) and [heroicons.com](https://heroicons.com)

# Getting Started
**Note**: The following instructions assume you have `make`, `go`, and `npm` installled.

After cloning the repo, build `craig-stars` locally once. This will build the frontend and then the backend. The backend embeds the frontend resources into its binary to support a single binary deployment.

```bash
make build
```

Install air for automatic server restarts while developing: 

```bash
go install github.com/cosmtrek/air@latest`
```

Launch the frontend and backend at the same time with `make`
```bash
make dev
```

Point your browser at [http://localhost:5173](http://localhost:5173) to see a live reloading frontend proxied to the go server on port `:8080`. Updating go code will relaunch the backend automatically. Updating frontend code will do a hot reload with sveltekit/vite.

# Visual Studio Code 
[VS Code](https://code.visualstudio.com) is highly recommended for development. `craig-stars` comes with a [cs.code-workspace](/cs.code-workspace) file that can be opened with VS Code in order to use front end and backend plugins without issue in the same repo.


## backend
To launch the front end and backend separately, you can call `air` directly.

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

# test
Run tests

```zsh
make test
```
