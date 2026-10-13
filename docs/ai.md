# AI players

AI players live in the [`ai`](../ai) package. This page describes how an AI decides its turn and what the difficulty levels change.

## How a turn works

The server runs each AI player once after the universe is generated and once after every turn is generated. The AI sees exactly what a human player in its seat would see: its own planets and fleets, and its intel on everything else. It gives orders through the same order functions the UI uses, so it can't do anything a human player can't.

The AI is rebuilt from the saved game every turn and keeps nothing in memory between turns. Anything it needs to remember is stored on the game objects themselves, as tags on fleets and production queue items. A fleet's `purpose` tag says what job it was built for, and its `idleSince` tag records the year it was first found with nothing to do.

Each turn is a fixed sequence of managers. Each one looks at the fleets that belong to its job, gives orders to the ones that are free, and asks for more ships if it needs them. Production runs near the end and turns those requests into queue items.

| Step | Manager           | What it does                                                                                                               |
| ---- | ----------------- | -------------------------------------------------------------------------------------------------------------------------- |
| 1    | Planning          | Sets the default production plan (auto factories, mines, terraforming, defenses) and research settings. Designs starbases. |
| 2    | War fleet sizing  | After the attack year, works out how many bombers, beam ships, torpedo ships and fuel ships make up a war fleet this year. |
| 3    | Scouting          | Sends each free scout to the nearest unexplored planet. Uses mineral packets as scouts for races that can.                 |
| 4    | Pickets           | Parks scouts with nothing left to explore over worlds outside planetary scanner range.                                     |
| 5    | Colonizing        | Loads colony ships and sends them to the best free world.                                                                  |
| 6    | Remote mining     | For AR, keeps a remote miner over each owned world.                                                                        |
| 7    | Mine laying       | Sends mine layers to owned worlds that have none. SD builds more as it grows.                                              |
| 8    | Invading, bombing | After the attack year, invades weak enemy worlds with colonists and sends war fleets to bomb.                              |
| 9    | Transport         | Moves colonists from full worlds to new colonies, then moves germanium with whatever transports are left.                  |
| 10   | Production        | Queues the ships the managers asked for, plus scanners and starbases.                                                      |
| 11   | Research          | Picks the next field from the research order.                                                                              |

Choices are deterministic. When two options are equally good, the lower planet number wins, so the same saved game always produces the same orders.

## Expansion

**Choosing a world.** A colony ship goes to the world with the best trade-off between travel time and habitability. Worlds with known hostile armed ships in orbit are skipped. A world that is only habitable after terraforming is used only when no habitable world is known.

**Loading colonists.** A source world keeps a reserve of its maximum population (25% on Normal) so it keeps growing quickly, and only the surplus is shipped. A ship takes enough for about 10% of the target's capacity, and at least 25 kT. An existing idle colony ship may take that minimum load early, before the source has a surplus, if the target is at least 70% habitable and the source has ten times the load.

**How many at once.** The number of colony ships in flight or waiting is capped by difficulty, and never exceeds the number of known free worlds. A developed shipyard whose orders would take less than a year of its resources has capacity going unused, so each one raises the cap by eight (scaled by the difficulty's pace). Orders for colony ships that no longer have a destination are removed from the queue. A loaded colony ship with nowhere left to go returns its colonists to the nearest owned world.

**Filling colonies.** Transport fleets carry surplus colonists from worlds above the reserve to colonies below it. Loads under 25 kT aren't sent.

**Minerals.** Factories need germanium, and young colonies rarely mine enough. A colony's first starbase or ships can wait years for ironium or boranium too. Transports with no colonists to move carry all three from worlds with a surplus to the worlds waiting for them. A source keeps what it needs itself plus 500 kT of each.

**AR.** AR races keep less in reserve, because their income grows with the square root of population and spreading out pays sooner. They build remote miners for their own worlds.

## Production

Each manager's request becomes ship orders at a suitable shipyard. A planet must have a starbase with a dock, and for most ship types its factories and mines must be at least half built. Colony ships and transports are built only where there are colonists to spare.

Scouts and colony ships are queued ahead of freighters, and freighters ahead of warships. New orders go to the top of the queue. Mineral packets are moved ahead of ships so a cheap packet isn't left waiting behind them.

A ship or starbase that is short of minerals stops everything behind it in the queue, including the mines that would supply those minerals. When the next order is stuck like that, the AI puts a quarter of the planet's resources into mines ahead of it.

Scanners and starbases are added when they can be finished quickly, or sooner if the planet is threatened.

A world's standard plan builds all its factories and mines before it terraforms. Once a world is half developed, the AI puts a quarter of its resources into terraforming ahead of the rest, because a more habitable world grows faster and holds more.

## War

Nothing offensive happens before the attack year, which depends on difficulty and on the game's start mode.

**War fleets.** A war fleet is bombers escorted by beam and torpedo ships, with one fuel tanker for every six of them. Bombers are heavy, and a fleet without tankers only has the fuel to crawl. The mix of beam and torpedo ships follows which design is the better value, and the fleet size grows as the game goes on.

**Mustering.** The ships of a fleet are often finished at different shipyards, years apart. After a year waiting, loose warships travel to whichever shipyard already holds the most, and form a fleet there once 60% of a full fleet has gathered or nothing more is on its way.

**How many at once.** The AI builds one war fleet per shipyard that is developed enough, and pauses when the warships already on order would cost more than three years of 30% of its income.

**Targets.** A fleet picks the enemy world that is quickest to reach on the fuel it has and smallest, avoiding starbases where it can, and ignores anything more than ten years away. On Normal it only attacks defenders it at least matches in rated strength, and on Hard it wants a quarter more. Fleets with nothing they can take on merge with others waiting at the same world. A fleet whose target stops being an enemy world picks another.

**Refuelling.** A fleet below a quarter of its fuel goes to the nearest of its starbases before taking another target.

**Invasion.** An idle transport fleet invades an explored enemy world that has no starbase when it can land at least twice the defenders, after allowing for the invaders its planetary defenses would kill. It prefers a world the AI's bombers are already over, then the nearest. It is recalled if the world gains a starbase, changes hands or grows past that ratio.

## Idle and surplus ships

- **Scouts** are built in proportion to what is left to explore: one per 10 unexplored planets. A scout with less than a quarter of its fuel is treated as spent: it no longer counts against the limit and its target is left for another scout, so a fresh one takes over. Once there is nothing to explore, they watch worlds outside planetary scanner range, enemy worlds first. A scout sitting inside the range of a better penetrating scanner is scrapped.
- **Transports** idle for five years are scrapped at the nearest starbase, keeping one fleet per world with colonists to spare.
- **Warships** that have waited five years for a fleet that never filled are sent as a partial fleet if they make up 60% of one and include bombers. Ships whose design has a replacement at least 1.5 times as strong are scrapped.

## Race doctrines

Each primary racial trait can change the standard play through a doctrine, defined in [`ai/doctrine.go`](../ai/doctrine.go). A doctrine is a small set of numbers and switches that the managers read. It applies at every difficulty.

| Race | What its doctrine changes                                                                                                                                                                              |
| ---- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| HE   | Twice as many colony ships at once, and idle colony ships seed worlds from 40% habitability instead of 70%.                                                                                            |
| SS   | Scouts are designed with a cloak as well as a scanner.                                                                                                                                                 |
| WM   | Attacks five years earlier and builds half again as many war fleets at once.                                                                                                                           |
| CA   | Judges a free world by its habitability after terraforming, since it terraforms on arrival, and sends half again as many colony ships.                                                                 |
| IS   | Builds five defenses a year on every world, ahead of terraforming.                                                                                                                                     |
| SD   | Keeps a mine layer over every owned world that has a starbase. Detonates a standard minefield while at least five armed enemy ships are inside it and they outnumber its own ships there three to one. |
| PP   | Gives a colony a cheap orbital fort with a mass driver once it is a fifth developed and can afford one within five years.                                                                              |
| IT   | The same early orbital fort, carrying a stargate. Fleets use gates automatically when both ends have one.                                                                                              |
| AR   | Keeps less population in reserve, builds remote miners for its own worlds, and researches energy first.                                                                                                |
| JoaT | Plays the standard game.                                                                                                                                                                               |

Most doctrines don't raise the score against another AI, and a few cost a little: defenses and minefields aren't scored and an AI opponent rarely tests them. Measured over eight maps, IT's is worth about +2%, HE's and CA's are level, WM's, IS's and PP's cost about 2%, and SD's costs about 5%.

## AI races

The standard AI races are in [`ai/races.go`](../ai/races.go), one for each primary racial trait. Each is a legal race: it costs no more race points than a player gets, and a test checks that. HE, SS and PP use the original game's standard-level computer races. The others are this project's own designs.

The original game made its computer players harder by giving them better races at higher levels, well outside the legal point budget at its two top levels. Here the difficulty levels change how the AI plays instead, and only Cheater uses illegal races.

## Difficulty levels

Every difficulty follows the same rules and gives the same kinds of orders. Difficulty changes how quickly the AI commits to expansion and war. Easy is slower and more cautious; it doesn't make deliberate mistakes.

| Setting                                                                | Easy                    | Normal                  | Hard                                               | Cheater        |
| ---------------------------------------------------------------------- | ----------------------- | ----------------------- | -------------------------------------------------- | -------------- |
| Pace (scales scouts, war fleets, miners and mine layers built at once) | 0.5×                    | 1×                      | 1.5×                                               | 1.5×           |
| Colony ships at once                                                   | 2                       | 4                       | 8                                                  | 8              |
| Population a source world keeps                                        | 33%                     | 25%                     | 25%                                                | 25%            |
| Population an AR source world keeps                                    | 25%                     | 15%                     | 10%                                                | 10%            |
| First attack, normal start                                             | Year 35                 | Year 25                 | Year 20                                            | Year 20        |
| First attack, accelerated start                                        | Year 30                 | Year 20                 | Year 15                                            | Year 15        |
| First attack, maximum start                                            | Year 10                 | Year 0                  | Year 0                                             | Year 0         |
| Shipyard development needed to build bombers                           | 90%                     | 90%                     | 65%                                                | 65%            |
| Research                                                               | Leftover resources only | Leftover resources only | 15% reserved on established worlds, plus leftovers | Same as Hard   |
| Research order                                                         | Standard                | Standard                | Race-specific goals first                          | Same as Hard   |
| Settles worlds that need terraforming, as a last resort                | No                      | Yes                     | Yes                                                | Yes            |
| Prefers mineral-rich worlds when colonizing                            | No                      | No                      | Yes                                                | Yes            |
| Strength it needs over known defenders before attacking                | Doesn't check           | Equal                   | A quarter more                                     | A quarter more |
| IS builds mine layers (SD always does)                                 | No                      | No                      | Yes                                                | Yes            |
| Races                                                                  | Standard pool           | Standard pool           | Standard pool                                      | Cheater pool   |

**Easy** expands with half as many colony ships, keeps more population at home, and doesn't attack until ten years after Normal would.

**Normal** is the baseline the other levels are defined against.

**Hard** commits to more at once, attacks five years earlier, and reserves part of its income for research. It researches the technology its race depends on before the standard order: for example weapons for WM, the terraforming fields for CA, and stargate technology for IT. It only attacks a world when its fleet is rated at least 1.25 times the defenses it knows about.

**Cheater** plays exactly like Hard. The only difference is the race: instead of one of the ten standard AI races, it gets one of two races (a JoaT and an HE) that cost about 3,000 race points more than a player is allowed.

## Limits

Things the AI doesn't do yet:

- Defend its worlds with fleets, or respond to incoming mineral packets.
- Predict battles. Strength comparisons use design power ratings.
- Ship minerals other than germanium between its worlds.
- Use the deeper race tricks: SS mineral theft, PP packet attacks and mineral flinging, CA remote terraforming, or IS growing colonists aboard freighters.
- Use diplomacy. Other players are either enemies or not.

## Testing

`TestAIRegression` in [`server/ai_regression_test.go`](../server/ai_regression_test.go) plays each standard AI race on Normal against a JoaT for 80 turns, and compares what it measures with the recording in [`server/ai-regression`](../server/ai-regression). A difference means AI behavior changed. When the change is intended, record it and review the diff:

```sh
AI_REGRESSION_UPDATE=1 go test ./server -run '^TestAIRegression$' -count=1
```

Each recording has one row per metric and one column per snapshot year (every ten years to 2480), covering the economy, the homeworld, research, exploration, ships owned, idle and queued by type, race traits such as minefields and packets in flight, and the war: fleets attacking, battles, and worlds taken, emptied and lost.

`AI_REGRESSION_WITHOUT=<behavior>` records the games with a behavior switched off, to see what it changes over time.

`TestAIGamesAreReproducible` plays one game twice, a year at a time, and stops at the first year the two differ, listing what differs. The recordings depend on a game always playing out the same way, so run this first if a recording changes for no reason.

### Comparing a behavior

`TestAICompare` measures whether a behavior helps. Behaviors that can be compared are wrapped in `ai.enabled("name")` in the code. For every map and AI race it plays two games against the JoaT opponent, one where only the race has the behavior and one where only the opponent does, and reports the score of the side with it against the side without:

```sh
AI_COMPARE=spare-colony go test ./server -run '^TestAICompare$' -count=1 -timeout 30m -v
```

`AI_COMPARE_SEEDS`, `AI_COMPARE_TURNS` (80 by default) and `AI_COMPARE_SIZE` change the maps, length and universe size. Use at least 80 turns: economic changes often show nothing by turn 50. Score isn't the whole story, either. A change that explores less can score higher only because it finds fewer enemy worlds to build war fleets against.

To look at these games in the app, start the server in AI mode:

```sh
go run main.go serve --ai-mode
```

Like `--test-mode`, this uses a temporary database that is deleted when the server stops. Before serving, it plays the same ten games the regression test plays, which takes about fifteen seconds. Start the frontend as usual, log in as `admin` (password `admin`) and open All Games from the menu. Each game is named for its races, such as "AI SS vs JoaT", and stands at its final year. The game menu's View As switches between the two players.

`TestAIStrategyThirtyTurns` plays the ten races together. Set `AI_STRATEGY_ALL_MODES=1` to run every difficulty for the full thirty turns.
