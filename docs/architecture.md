# Architecture

`craig-stars` is a monorepo containing code both for the frontend client and the backend server. The file structure is broken down as follows:

| path        | description                                                                                                                                                                                                             |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `/`         | The root folder contains the `main.go` entrypoint into the application.                                                                                                                                                 |
| `/cs`       | All core game logic and models are in the `cs` package. This is also what handles universe and turn generation logic. More details in the [cs](#cs) section                                                             |
| `/db`       | The `db` package handles serializing games to and from the database, as well as any UI database queries.                                                                                                                |
| `/server`   | The `server` package is where the webserver routes are configured. It also is the "glue" package that ties the game logic together with the database serialization. More details in the [server](#server) section       |
| `/cmd`      | The `cmd` package is where command line parsing is handled, as well as the entrypoint for serving the application.                                                                                                      |
| `/config`   | The `config` package is craig-stars configuration code lives. This config is loaded from `data/config/config.yaml` and is shared by the database and the server. The config defaults to settings for local development. |
| `/ai`       | The `ai` package is where the logic for ai players resides. The AI strives to be "just another player" with no special insight into the game world.                                                                     |
| `/test`     | The `test` package contains common testing utilities.                                                                                                                                                                   |
| `/frontend` | The `frontend` folder contains all the sveltekit front end code. More details in the [frontend](#frontend) section.                                                                                                     |

## cs

The [cs](/cs) package contains the models and game logic. It contains no serialization logic. There are two main interfaces into the game logic. The [Gamer](/cs/gamer.go) and the [Orderer](/cs/orderer.go). The `Gamer` is used to create games, generate universes and generate new turns. The `Orderer` is used to update player `Research`, `Planet`, `Fleet`, and `MineField` orders, as well as to handle in turn cargo transfers. The `server` makes uses of these interfaces to update game data and save changes back to the database.


## server

The [server](/server) package contains the backend [server](/server/server.go) and the [GameRunner](/server/gamerunner.go). The backend http server handles all `/api` requests from the front end. Routes for various resources are defined in their own files ([fleets.go](/server/fleets.go), [planets.go](/server/planets.go), etc).

The [GameRunner](/server/gamerunner.go) hooks up the [Gamer](/cs/gamer.go) with the database. The `GameRunner` hosts new games, updates game players, and generates new turns on when player's submit their turns. The `GameRunner` also loads a player's game information from the database.

```mermaid
graph TD;
    CS["CS\nGame Logic"]
    DB[(DB\nsqlite db)]
    Server["Server\ngo http server"]
    Cmd--serve-->Server
    Server--Host Game\nSubmit Turn\nGenerate Turn-->GameRunner;
    Server--reads data\nupdates player orders-->DB;
    Server--reads config-->Config;
    DB--reads config-->Config;
    GameRunner<--saves and loads games-->DB;
    GameRunner--generate universe/turn-->CS;
    GameRunner--runs ai player-->AI;
    Frontend--/api REST calls-->Server
    User["fa:fa-fa user"]-->Frontend
```

## frontend

The [frontend](/frontend) is a static SvelteKit site that loads game information from the backend `/api` endpoint and displays it to the user. The `frontend` also handles updating changes to player research, planet, fleet, and mine field orders by calling endpoints on the `/api`.

### Game View

The main [game view](</frontend/src/routes/(user)/games/(game)/[id]/(main)/Game.svelte>) is split into sections.

| section      | description                                                          |
| ------------ | -------------------------------------------------------------------- |
| Scanner      | The zoomable and pannable map display of the universe                |
| Command Pane | Tiles with information about the currently commanded planet or fleet |
| Summary Pane | A small summary about the currently selected map object              |

### frontend structure

| path     | description                                                                                                                                                                                                                                                                              |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `css`    | Contains tailwind `@apply` overrides and the various icon classes for techs and planets                                                                                                                                                                                                  |
| `lib`    | `lib` contains type definitions for server side models, Service classes for interacting with the backend `/api`, and any reusauble components (components that exist outside of a single page)                                                                                           |
| `routes` | `routes` contains all the user accessible routes for the application. Almost all routes require the user to be logged in and are in the `(user)` route group. The main game view is located under [src/routes/(user)/games/(game)/[id]](</frontend/src/routes/(user)/games/(game)/[id]>) |

### stores

Webapps are all about reactivity. When the user makes a change or some new data is loaded from the server, the front end needs to react to this new data and update a portion of the page. Svelte uses special [`$`](https://svelte.dev/docs/svelte-components#script-3-$-marks-a-statement-as-reactive) character to denote an inline store or the [store](https://svelte.dev/docs/svelte-store) types. `craig-stars` makes use of stores, along with [contexts](https://svelte.dev/docs/svelte#setcontext) to ensure the components update when their data change. For example, if cargo is transferred from a planet to a fleet, the planet minerals tile needs to update.

Every component under under the game route (`/games/[id]`) has access to a context:

```ts
type GameContext = {
  game: Readable<FullGame>;
  player: Readable<Player>;
  universe: Readable<Universe>;
  settings: Writable<PlayerSettings>;
  designs: Writable<ShipDesign[]>;
  messageNum: Writable<number>;
};

// accessed like this
const { game, player, universe, settings } = getGameContext();
```

This context contains reactive stores with the state of the game, for example a component like this would automatically update the game year when it was updated by the server.

```html
<script lang="ts">
  import { getGameContext } from '$lib/services/Contexts';

  const { game, player, universe } = getGameContext();
</script>

<h1>{$game.year}</h1>
```

`craig-stars` frontend code makes great use of contexts and stores to update the UI and keep components small.

The [Stores.ts](/frontend/src/lib/services/Stores.ts) defines some common stores used in the game (these should probably be included in the context...). These stores maintain the current state of selection.

Planets and Fleets owned by the player can be Commanded, at which point they show up in the Command Pane. All map objects can be Selected, at which point they show up in the Selection Summary.

A single click selects a map object. A second click commands the map object if it is owned by the player. Commanded map objects can have their orders changed by the player.

- CommandedPlanet - The currently commanded planet, or undefined if no planet is commanded.
- CommandedFleet - The currently commanded fleet, or undefined if no planet is commanded.
- SelectedMapObject - The currently selected map object, or undefined if no map object is selected.
