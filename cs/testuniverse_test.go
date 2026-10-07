//go:build !wasi && !wasm

package cs

import (
	"slices"
	"testing"
)

// testUniverse keeps fixture setup and phase execution consistent.
type testUniverse struct {
	t    *testing.T
	Game *FullGame
	turn *turnGenerator
}

func newTestUniverse(t *testing.T, s TestScenario) *testUniverse {
	t.Helper()
	game := BuildScenario(s)
	for _, player := range game.Players {
		player.Messages = nil
	}
	if err := game.computeSpecs(); err != nil {
		t.Fatal(err)
	}
	turn := newTurnGenerator(game)
	turn.log = testLogger
	return &testUniverse{t: t, Game: game, turn: &turn}
}

func (u *testUniverse) Player(num int) *Player {
	u.t.Helper()
	if num < 1 || num > len(u.Game.Players) {
		u.t.Fatalf("unknown player %d", num)
	}
	return u.Game.Players[num-1]
}

func (u *testUniverse) Planet(name string) *Planet {
	u.t.Helper()
	for _, p := range u.Game.Planets {
		if p.Name == name {
			return p
		}
	}
	u.t.Fatalf("unknown planet %q", name)
	return nil
}

func (u *testUniverse) Fleet(name string) *Fleet {
	u.t.Helper()
	var found *Fleet
	for _, f := range u.Game.Fleets {
		if f.Name == name {
			if found != nil {
				u.t.Fatalf("ambiguous fleet %q; use FleetFor", name)
			}
			found = f
		}
	}
	if found == nil {
		u.t.Fatalf("unknown fleet %q", name)
	}
	return found
}

func (u *testUniverse) FleetFor(playerNum int, name string) *Fleet {
	u.t.Helper()
	for _, f := range u.Game.Fleets {
		if f.PlayerNum == playerNum && f.Name == name {
			return f
		}
	}
	u.t.Fatalf("player %d: unknown fleet %q", playerNum, name)
	return nil
}

// TransferByHand uses the real order validation and the player's intel for
// foreign objects. An empty target means jettison/load from deep space.
func (u *testUniverse) TransferByHand(playerNum int, fleet, target string, cargo Cargo) {
	u.t.Helper()
	player := u.Player(playerNum)
	source := u.FleetFor(playerNum, fleet)
	var object CargoHolder
	if target != "" {
		for _, p := range u.Game.Planets {
			if p.Name == target {
				if p.PlayerNum == playerNum {
					object = p
				} else {
					object = player.GetPlanetIntel(p.Num)
				}
			}
		}
		for _, f := range u.Game.Fleets {
			if f.Name == target {
				if object != nil {
					u.t.Fatalf("ambiguous cargo target %q", target)
				}
				if f.PlayerNum == playerNum {
					object = f
				} else {
					intel := player.GetFleetIntel(f.PlayerNum, f.Num)
					if intel == nil {
						u.t.Fatalf("unknown fleet intel %q", target)
					}
					object = intel
				}
			}
		}
		for _, s := range u.Game.Salvages {
			if s.Name == target {
				intel := player.GetSalvageIntel(s.Num)
				if intel == nil {
					u.t.Fatalf("unknown salvage intel %q", target)
				}
				object = intel
			}
		}
		for _, p := range u.Game.MineralPackets {
			if p.Name == target {
				if p.PlayerNum == playerNum {
					object = p
				} else {
					intel := player.GetMineralPacketIntel(p.PlayerNum, p.Num)
					if intel == nil {
						u.t.Fatalf("unknown packet intel %q", target)
					}
					object = intel
				}
			}
		}
		if object == nil {
			u.t.Fatalf("unknown cargo target %q", target)
		}
	}
	if err := NewOrderer().TransferByHand(&u.Game.Rules, player, source, object, CargoTransferRequest{Cargo: cargo}); err != nil {
		u.t.Fatal(err)
	}
}

// GenerateTurn first drops objects deleted by the previous turn, like the
// server does when it saves, so tests can still check Delete flags after a turn.
func (u *testUniverse) GenerateTurn() {
	u.t.Helper()
	removeDeleted(u.Game.Universe)
	u.Game.buildMaps(u.Game.Players)
	if err := u.turn.generateTurn(); err != nil {
		u.t.Fatal(err)
	}
}

// GenerateTurns generates n turns in a row.
func (u *testUniverse) GenerateTurns(n int) {
	u.t.Helper()
	for range n {
		u.GenerateTurn()
	}
}

// removeDeleted drops objects deleted during a turn, like the server does when
// it saves a turn.
func removeDeleted(universe *Universe) {
	universe.Fleets = slices.DeleteFunc(universe.Fleets, func(o *Fleet) bool { return o.Delete })
	universe.Starbases = slices.DeleteFunc(universe.Starbases, func(o *Fleet) bool { return o.Delete })
	universe.Salvages = slices.DeleteFunc(universe.Salvages, func(o *Salvage) bool { return o.Delete })
	universe.Minefields = slices.DeleteFunc(universe.Minefields, func(o *Minefield) bool { return o.Delete })
	universe.MineralPackets = slices.DeleteFunc(universe.MineralPackets, func(o *MineralPacket) bool { return o.Delete })
	universe.MysteryTraders = slices.DeleteFunc(universe.MysteryTraders, func(o *MysteryTrader) bool { return o.Delete })
	universe.Wormholes = slices.DeleteFunc(universe.Wormholes, func(o *Wormhole) bool { return o.Delete })
}

func (u *testUniverse) Run(steps ...func(*turnGenerator)) {
	u.t.Helper()
	u.Game.buildMaps(u.Game.Players)
	for _, step := range steps {
		step(u.turn)
	}
}

func (u *testUniverse) RunE(steps ...func(*turnGenerator) error) {
	u.t.Helper()
	u.Game.buildMaps(u.Game.Players)
	for _, step := range steps {
		if err := step(u.turn); err != nil {
			u.t.Fatal(err)
		}
	}
}

func (u *testUniverse) Recompute() {
	u.t.Helper()
	if err := u.Game.computeSpecs(); err != nil {
		u.t.Fatal(err)
	}
}

func (u *testUniverse) Messages(playerNum int, typ PlayerMessageType) []PlayerMessage {
	u.t.Helper()
	var result []PlayerMessage
	for _, message := range u.Player(playerNum).Messages {
		if message.Type == typ {
			result = append(result, message)
		}
	}
	return result
}

// singleFleetScenario is a small declarative base for tests with one design.
// Callers customize scenario fields before the universe is built.
func singleFleetScenario(design ShipDesign) TestScenario {
	s := SingleUnitScenario()
	s.Players[0].Designs = Designs(design)
	s.Players[0].Fleets = []ScenarioFleet{{Design: design.Name, Waypoints: []ScenarioWaypoint{{Warp: 5}}}}
	return s
}

func twoPlayerFleetScenario(design ShipDesign) TestScenario {
	s := TwoPlayerScenario()
	s.Players[0].Designs = Designs(design)
	s.Players[0].Fleets = []ScenarioFleet{{Design: design.Name, Waypoints: []ScenarioWaypoint{{Warp: 5}}}}
	return s
}

// newSeededGame creates a game like Gamer.CreateGame, but with a fixed seed so
// generated universes and turns are the same on every run.
func newSeededGame(settings GameSettings) *Game {
	game := NewGamer().CreateGame(1, settings)
	game.Seed = 0
	game.Rules.ResetSeed(game.Seed)
	return game
}
