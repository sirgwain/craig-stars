# craig-stars

A web based Stars! clone.

## Tech Stack

- [gin-gonic](https://github.com/gin-gonic/gin)
- [cobra](https://github.com/spf13/cobra) (for cli)
- [viper](https://github.com/spf13/viper) (for config)
- [sqlx](https://github.com/jmoiron/sqlx) + sqlite (for session/user management)
- [sveltekit](https://kit.svelte.dev) (with static adaptor)
- [tailwindcss](https://tailwindcss.com)
- [daisyui](https://daisyui.com)

# dev

For development, install air `go install github.com/cosmtrek/air@latest`.

## backend

Start the go webserver:

```zsh
❯ air

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ , built with Go

watching .
watching cmd
watching config
watching db
!exclude dist
!exclude frontend
watching server
!exclude tmp
building...
running...
[GIN-debug] POST   /api/login                --> craig-stars/server.Login (5 handlers)
[GIN-debug] GET    /api/logout               --> craig-stars/server.Logout (5 handlers)
[GIN-debug] GET    /api/me                   --> craig-stars/server.Me (6 handlers)
[GIN-debug] GET    /api/ping                 --> craig-stars/server.Start.func2 (5 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

## frontend

Launch the frontend in development mode with

```zsh
yarn --cwd frontend run dev
```

Point your browser at [http://localhost:5173](http://localhost:5173) to see a live reloading frontend proxied to the go server on port `:8080`. Updating go code will relaunch the backend automatically. Updating frontend code will do a hot reload with sveltekit/vite.

# test

Run tests

```zsh
go test ./...
yarn --cwd frontend run test
```
