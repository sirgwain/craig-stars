# craig-stars

A web based Stars! clone.

## Tech Stack

- [go-chi](https://github.com/chi/go-chi)
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
watching ai
watching ansible
watching ansible/group_vars
watching ansible/inventory
watching ansible/roles
watching ansible/roles/craig-stars
watching ansible/roles/craig-stars/defaults
watching ansible/roles/craig-stars/handlers
watching ansible/roles/craig-stars/tasks
watching ansible/roles/craig-stars/templates
watching ansible/roles/root
watching ansible/roles/root/defaults
watching ansible/roles/root/handlers
watching ansible/roles/root/tasks
watching ansible/roles/root/templates
watching ansible/roles/service
watching ansible/roles/service/handlers
watching ansible/roles/service/tasks
watching ansible/roles/service/templates
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
11:06AM DBG Debug logging enabled
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
