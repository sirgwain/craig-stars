//go:build !wasi && !wasm

package cs

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// byHandScenario has player 1's freighters, player 2's freighter and a salvage, all at player 1's planet
func byHandScenario() TestScenario {
	second := AIPlayer("Player 2")
	second.Designs = Designs(DesignTeamster)
	second.Fleets = []ScenarioFleet{{Name: "Hauler #1", Design: "Teamster", At: "Planet 1"}}
	return TestScenario{
		Name: "By Hand Transfers",
		Players: []ScenarioPlayer{
			{
				Designs: Designs(DesignTeamster),
				Fleets: []ScenarioFleet{
					{Design: "Teamster", At: "Planet 1", Quantity: 2, Cargo: Cargo{Ironium: 50, Boranium: 20}},
					{Design: "Teamster", At: "Planet 1", Quantity: 2},
					{Design: "Teamster", At: "Planet 1", Quantity: 2, Cargo: Cargo{Germanium: 70}},
				},
				Salvages: []Salvage{{MapObject: MapObject{Position: Vector{}}, Cargo: Cargo{Ironium: 100, Boranium: 100, Germanium: 100}}},
			},
			second,
		},
		Planets: []ScenarioPlanet{
			{Name: "Planet 1", Owner: 1, Cargo: Cargo{Ironium: 500, Boranium: 500, Germanium: 500, Colonists: 2500}},
			{Name: "Planet 2", Owner: 2, Position: Vector{X: 100}, Cargo: Cargo{Colonists: 2500}},
		},
	}
}

func newByHandTestUniverse(t *testing.T) *testUniverse {
	u := newTestUniverse(t, byHandScenario())
	for _, player := range u.Game.Players {
		discoverer := newDiscoverer(testLogger, player)
		discoverer.discoverSalvage(u.Game.Salvages[0])
		for _, fleet := range u.Game.Fleets {
			if fleet.PlayerNum != player.Num {
				discoverer.discoverFleet(fleet, true)
				discoverer.discoverFleetCargo(fleet)
			}
		}
	}
	return u
}

func byHandIncompleteMessages(player *Player) []PlayerMessage {
	messages := []PlayerMessage{}
	for _, message := range player.Messages {
		if message.Type == PlayerMessageFleetByHandTransferIncomplete {
			messages = append(messages, message)
		}
	}
	return messages
}

func Test_turn_fleetByHandTransfers(t *testing.T) {
	t.Run("transfers between our own objects aren't replayed", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		u.TransferByHand(1, "Teamster #2", "Planet 1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #1", Cargo{Ironium: 50})
		u.TransferByHand(1, "Teamster #2", "Planet 1", Cargo{Ironium: -150})

		assert.Empty(t, u.Player(1).CargoTransfers)

		u.turn.fleetByHandTransfers()

		assert.Equal(t, Cargo{Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, 550, u.Planet("Planet 1").Cargo.Ironium)
	})

	t.Run("another player loaded the salvage first", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// player 2 loads 60kT, player 1 loads all 100kT, then moves it to another fleet
		u.TransferByHand(2, "Hauler #1", "Salvage #1", Cargo{Ironium: 60})
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #3", Cargo{Ironium: -100})

		u.turn.fleetByHandTransfers()

		// player 1 settles first and gets everything it took
		assert.Equal(t, Cargo{Germanium: 70, Ironium: 100}, u.FleetFor(1, "Teamster #3").Cargo)
		assert.Empty(t, byHandIncompleteMessages(u.Player(1)))

		// player 2 only gets what's left. It doesn't make any
		assert.Equal(t, Cargo{}, u.FleetFor(2, "Hauler #1").Cargo)
		assert.Equal(t, 0, u.Game.Salvages[0].Cargo.Ironium)
		messages := byHandIncompleteMessages(u.Player(2))
		if assert.Len(t, messages, 1) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 0, Wanted: 60}, *messages[0].Spec.CargoTransfer)
		}
	})

	t.Run("shortfall comes out of the fleet we moved the cargo to", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// player 1 sees 100kT and loads it all, then moves it to another fleet
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #3", Cargo{Ironium: -100})
		// someone else took 40kT before player 1's load is settled
		u.Game.Salvages[0].Cargo.Ironium = 60

		u.turn.fleetByHandTransfers()

		// no cargo is created
		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, Cargo{Germanium: 70, Ironium: 60}, u.FleetFor(1, "Teamster #3").Cargo)
		assert.Equal(t, 0, u.Game.Salvages[0].Cargo.Ironium)
		messages := byHandIncompleteMessages(u.Player(1))
		if assert.Len(t, messages, 1) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 60, Wanted: 100}, *messages[0].Spec.CargoTransfer)
		}
	})

	t.Run("shortfall reduces cargo we gave away", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// player 1 loads 100kT from the salvage and gives it all to player 2
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Hauler #1", Cargo{Ironium: -100})
		u.Game.Salvages[0].Cargo.Ironium = 60

		u.turn.fleetByHandTransfers()

		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, Cargo{Ironium: 60}, u.FleetFor(2, "Hauler #1").Cargo)
		assert.Len(t, byHandIncompleteMessages(u.Player(1)), 2)
	})

	t.Run("cargo a target can't hold comes back", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		u.TransferByHand(1, "Teamster #1", "Hauler #1", Cargo{Ironium: -50})
		// player 2 filled their freighter before player 1's unload is settled
		hauler := u.FleetFor(2, "Hauler #1")
		hauler.Cargo = Cargo{Germanium: hauler.Spec.CargoCapacity - 20}

		u.turn.fleetByHandTransfers()

		assert.Equal(t, Cargo{Germanium: hauler.Spec.CargoCapacity - 20, Ironium: 20}, hauler.Cargo)
		assert.Equal(t, Cargo{Ironium: 30, Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		messages := byHandIncompleteMessages(u.Player(1))
		if assert.Len(t, messages, 1) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 20, Wanted: 50, Status: CargoTransferStatusDestCargoCapacity}, *messages[0].Spec.CargoTransfer)
		}
	})

	t.Run("cargo for a fleet that's gone comes back", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		u.TransferByHand(1, "Teamster #1", "Hauler #1", Cargo{Ironium: -50})
		u.FleetFor(2, "Hauler #1").Delete = true

		u.turn.fleetByHandTransfers()

		assert.Equal(t, Cargo{Ironium: 50, Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Len(t, byHandIncompleteMessages(u.Player(1)), 1)
	})

	t.Run("transfers follow merged fleets", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		fleet2, fleet3 := u.FleetFor(1, "Teamster #2"), u.FleetFor(1, "Teamster #3")
		if _, err := NewOrderer().Merge(&u.Game.Rules, player, []*Fleet{fleet3, fleet2}); err != nil {
			t.Fatal(err)
		}
		u.Game.Salvages[0].Cargo.Ironium = 70

		u.turn.fleetByHandTransfers()

		assert.Equal(t, Cargo{Germanium: 70, Ironium: 70}, fleet3.Cargo)
		messages := byHandIncompleteMessages(player)
		if assert.Len(t, messages, 1) {
			assert.Equal(t, fleet3.Num, messages[0].Target.TargetNum)
		}
	})
}

// cargo totals everything that holds minerals at the by hand scenario's location
func byHandMinerals(u *testUniverse) Cargo {
	total := Cargo{}
	for _, fleet := range u.Game.Fleets {
		if !fleet.Delete {
			total = total.Add(fleet.Cargo)
		}
	}
	for _, salvage := range u.Game.Salvages {
		if !salvage.Delete {
			total = total.Add(salvage.Cargo)
		}
	}
	total = total.Add(u.Planet("Planet 1").Cargo)
	total.Colonists = 0
	return total
}

// Test_turn_fleetByHandTransfersRandom makes random by hand transfers, splits and merges, then settles them.
// When nobody else touches the cargo, settling doesn't change anything the player saw. When someone does,
// no cargo is created or destroyed
func Test_turn_fleetByHandTransfersRandom(t *testing.T) {
	for seed := int64(0); seed < 300; seed++ {
		interfere := seed%2 == 1
		t.Run(fmt.Sprintf("seed %d interfere %t", seed, interfere), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			u := newByHandTestUniverse(t)
			rules := &u.Game.Rules
			player := u.Player(1)
			planet := u.Planet("Planet 1")
			salvage := u.Game.Salvages[0]
			hauler := u.FleetFor(2, "Hauler #1")
			orderer := NewOrderer()
			start := byHandMinerals(u)

			playerFleets := func() []*Fleet {
				fleets := []*Fleet{}
				for _, fleet := range u.Game.Fleets {
					if fleet.PlayerNum == player.Num && !fleet.Delete {
						fleets = append(fleets, fleet)
					}
				}
				return fleets
			}

			for i := 0; i < 30; i++ {
				fleets := playerFleets()
				fleet := fleets[rng.Intn(len(fleets))]
				switch op := rng.Intn(10); {
				case op < 7:
					// transfer some of one mineral with a random target
					var dest CargoHolder
					switch rng.Intn(5) {
					case 0:
						dest = planet
					case 1:
						dest = fleets[rng.Intn(len(fleets))]
					case 2:
						dest = player.GetSalvageIntel(salvage.Num)
					case 3:
						dest = player.GetFleetIntel(hauler.PlayerNum, hauler.Num)
					case 4:
						// jettison
					}
					if dest == fleet {
						continue
					}
					cargo := Cargo{}.WithCargo(CargoTypes[rng.Intn(3)], rng.Intn(161)-80)
					// invalid transfers are rejected, just like in the UI
					_ = orderer.TransferByHand(rules, player, fleet, dest, CargoTransferRequest{Cargo: cargo})
				case op < 8:
					if len(fleets) < 2 {
						continue
					}
					other := fleets[rng.Intn(len(fleets))]
					if other == fleet {
						continue
					}
					if _, err := orderer.Merge(rules, player, []*Fleet{fleet, other}); err != nil {
						t.Fatal(err)
					}
				default:
					newFleets, err := orderer.SplitAll(rules, player, fleets, fleet)
					if err != nil {
						t.Fatal(err)
					}
					for _, newFleet := range newFleets {
						u.Game.Fleets = append(u.Game.Fleets, newFleet)
						if err := u.Game.addFleet(newFleet); err != nil {
							t.Fatal(err)
						}
					}
				}
			}

			// what the player sees before the turn is generated
			fleets := playerFleets()
			wantFleetCargo := make([]Cargo, len(fleets))
			for i, fleet := range fleets {
				wantFleetCargo[i] = fleet.Cargo
			}
			wantPlanetCargo := planet.Cargo
			wantHaulerCargo := player.GetFleetIntel(hauler.PlayerNum, hauler.Num).Cargo
			wantSalvageCargo := player.GetSalvageIntel(salvage.Num).Cargo.Add(player.getByHandTransfer(MapObjectTarget{TargetPosition: salvage.Position}))

			// another player takes from the salvage and loads their freighter
			outside := Cargo{}
			if interfere {
				taken := Cargo{}.WithCargo(CargoTypes[rng.Intn(3)], rng.Intn(salvage.Cargo.Ironium+1))
				taken = taken.Add(Cargo{}.WithCargo(CargoTypes[rng.Intn(3)], rng.Intn(salvage.Cargo.Germanium+1)))
				if !salvage.Cargo.CanTransfer(taken) {
					taken = Cargo{}
				}
				salvage.Cargo = salvage.Cargo.Subtract(taken)
				loaded := Cargo{Boranium: rng.Intn(hauler.availableCargoSpace() + 1)}
				hauler.Cargo = hauler.Cargo.Add(loaded)
				outside = loaded.Subtract(taken)
			}

			u.turn.fleetByHandTransfers()

			got := byHandMinerals(u)
			assert.Equal(t, start.Add(outside), got, "no cargo is created or destroyed")
			for _, fleet := range playerFleets() {
				assert.False(t, fleet.Cargo.HasNegative())
				assert.LessOrEqual(t, fleet.Cargo.Total(), fleet.Spec.CargoCapacity)
			}
			assert.Empty(t, player.CargoTransfers)

			if interfere {
				return
			}

			assert.Empty(t, byHandIncompleteMessages(player))
			for i, fleet := range fleets {
				assert.Equal(t, wantFleetCargo[i], fleet.Cargo, fleet.Name)
			}
			assert.Equal(t, wantPlanetCargo, planet.Cargo)
			assert.Equal(t, wantHaulerCargo, hauler.Cargo)
			salvageCargo := Cargo{}
			for _, s := range u.Game.Salvages {
				if !s.Delete && s.Position == salvage.Position {
					salvageCargo = salvageCargo.Add(s.Cargo)
				}
			}
			assert.Equal(t, wantSalvageCargo, salvageCargo)
			assert.False(t, slices.ContainsFunc(u.Game.Salvages, func(s *Salvage) bool { return s.Cargo.HasNegative() }))
		})
	}
}
