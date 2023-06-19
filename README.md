# craig-stars

A web based Stars! clone.

## Tech

- [gin-gonic](https://github.com/gin-gonic/gin)
- [cobra](https://github.com/spf13/cobra) (for cli)
- [viper](https://github.com/spf13/viper) (for config)
- [gorm](https://gorm.io) + sqlite (for session/user management)
- [sveltekit](https://kit.svelte.dev) (with static adaptor)
- [tailwindcss](https://tailwindcss.com)
- [daisyui](https://daisyui.com)
- [svelteui](https://www.svelteui.org)

# build

To build both the front end and the go binary, call:

```zsh
make
```

# run

To run the app

```zsh
make run
```

Point your browser at [http://localhost:8080](http://localhost:8080)

# user management

To create a new user, call the cli

```zsh
> ./dist/craig-stars-darwin-universal create user -u admin -p admin
> ./dist/craig-stars-darwin-universal list users

Users

│ USERNAME │
│──────────│
│ admin    │
```

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
yarn --cwd run dev
```

Point your browser at [http://localhost:3000](http://localhost:3000) to see a live reloading frontend proxied to the go server on port `:8080`. Updating go code will relaunch the backend automatically. Updating frontend code will do a hot reload with sveltekit/vite.

# test

## frontend

Install playwright `npx playwright install`

Run tests

```zsh
yarn --cwd frontend run test
```
