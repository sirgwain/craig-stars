//go:build !wasi && !wasm

package cs

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// byHandScenario has player 1's freighters, player 2's freighter and a salvage, all at player 1's planet.
// Everything starts at the same location:
//
//	Planet 1 (player 1)      500kT ironium, boranium and germanium
//	Teamster #1 (player 1)   50kT ironium, 20kT boranium   (420kT hold, 900mg fuel, full)
//	Teamster #2 (player 1)   empty                         (420kT hold, 900mg fuel, full)
//	Teamster #3 (player 1)   70kT germanium                (420kT hold, 900mg fuel, full)
//	Bystander #1 (player 1)  100kT ironium, never transfers anything
//	Hauler #1 (player 2)     empty                         (210kT hold, 450mg fuel, full)
//	Salvage #1               100kT ironium, boranium and germanium
//
// In u.TransferByHand, positive cargo is loaded into the fleet and negative cargo is unloaded from it.
// Transfers with player 1's own planet and fleets happen right away. Transfers with the salvage and
// player 2's fleet only change player 1's intel until fleetByHandTransfers settles them.
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
					{Name: "Bystander #1", Design: "Teamster", At: "Planet 1", Cargo: Cargo{Ironium: 100}},
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
		// Teamster #2 loads 100kT ironium from the planet (planet 500 -> 400), takes Teamster #1's 50kT
		// (Teamster #1 50 -> 0), then unloads all 150kT back to the planet (planet 400 -> 550)
		u.TransferByHand(1, "Teamster #2", "Planet 1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #1", Cargo{Ironium: 50})
		u.TransferByHand(1, "Teamster #2", "Planet 1", Cargo{Ironium: -150})

		u.turn.fleetByHandTransfers()

		// settling leaves everything as the player left it
		assert.Equal(t, Cargo{Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, 550, u.Planet("Planet 1").Cargo.Ironium)
	})

	t.Run("another player loaded the salvage first", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// both players see 100kT ironium in the salvage. Player 2 loads 60kT. Player 1 loads all 100kT into
		// Teamster #2, then moves it to Teamster #3 (70kT germanium -> 70kT germanium, 100kT ironium).
		// Together they took 160kT from a salvage that has 100kT
		u.TransferByHand(2, "Hauler #1", "Salvage #1", Cargo{Ironium: 60})
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #3", Cargo{Ironium: -100})

		u.turn.fleetByHandTransfers()

		// player 1 settles first and gets all 100kT, so nothing changes for them
		assert.Equal(t, Cargo{Germanium: 70, Ironium: 100}, u.FleetFor(1, "Teamster #3").Cargo)
		assert.Empty(t, byHandIncompleteMessages(u.Player(1)))

		// the salvage is empty when player 2 settles. The Hauler gets nothing and player 2 gets a message
		// saying 0 of the 60kT loaded
		assert.Equal(t, Cargo{}, u.FleetFor(2, "Hauler #1").Cargo)
		assert.Equal(t, 0, u.Game.Salvages[0].Cargo.Ironium)
		messages := byHandIncompleteMessages(u.Player(2))
		if assert.Len(t, messages, 1) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 0, Wanted: 60}, *messages[0].Spec.CargoTransfer)
		}
	})

	t.Run("shortfall comes out of the fleet we moved the cargo to", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// player 1 sees 100kT ironium in the salvage and loads it all into Teamster #2, then moves it to
		// Teamster #3 (70kT germanium -> 70kT germanium, 100kT ironium). Teamster #2 is empty again
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #3", Cargo{Ironium: -100})
		// someone else took 40kT before player 1's load is settled. The salvage has 60kT
		u.Game.Salvages[0].Cargo.Ironium = 60

		u.turn.fleetByHandTransfers()

		// like the original game, the transfers are replayed in order: Teamster #2 only loads 60kT, so it
		// can only move 60kT to Teamster #3 (100kT -> 60kT). Teamster #1's 50kT is left alone. The old code
		// created the missing 40kT out of nothing
		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, Cargo{Germanium: 70, Ironium: 60}, u.FleetFor(1, "Teamster #3").Cargo)
		assert.Equal(t, Cargo{Ironium: 50, Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Equal(t, 0, u.Game.Salvages[0].Cargo.Ironium)
		// one message for the short load, and one for the short move to Teamster #3
		messages := byHandIncompleteMessages(u.Player(1))
		if assert.Len(t, messages, 2) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 60, Wanted: 100}, *messages[0].Spec.CargoTransfer)
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 60, Wanted: 100, Status: CargoTransferStatusCargo}, *messages[1].Spec.CargoTransfer)
		}
	})

	t.Run("fleets that didn't transfer don't lose cargo", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// Teamster #1 holds the most ironium (200kT) but never transfers anything
		u.FleetFor(1, "Teamster #1").Cargo = Cargo{Ironium: 200}
		// Teamster #2 loads 100kT ironium from the salvage and unloads it all on the planet (500kT -> 600kT)
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Planet 1", Cargo{Ironium: -100})
		// someone else took 40kT before player 1's load is settled. The salvage has 60kT
		u.Game.Salvages[0].Cargo.Ironium = 60

		u.turn.fleetByHandTransfers()

		// only 60kT reached the planet (600kT -> 560kT). Teamster #1 keeps its 200kT
		assert.Equal(t, Cargo{Ironium: 200}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Equal(t, 560, u.Planet("Planet 1").Cargo.Ironium)
	})

	t.Run("shortfalls follow cargo along a chain of fleets", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// Teamster #2 loads 100kT ironium from the salvage and passes it to Teamster #3, which passes it to
		// Teamster #1 (50kT -> 150kT ironium). Then Teamster #1 unloads 120kT on the planet (500kT -> 620kT),
		// keeping 30kT
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Teamster #3", Cargo{Ironium: -100})
		u.TransferByHand(1, "Teamster #3", "Teamster #1", Cargo{Ironium: -100})
		u.TransferByHand(1, "Teamster #1", "Planet 1", Cargo{Ironium: -120})
		// someone else took 40kT before player 1's load is settled. The salvage has 60kT
		u.Game.Salvages[0].Cargo.Ironium = 60

		u.turn.fleetByHandTransfers()

		// replayed in order: Teamster #2 only loads 60kT and passes 60kT along, so Teamster #1 has 110kT and
		// its 120kT unload only moves 110kT. Teamster #1 ends with no ironium instead of 30kT, and the planet
		// has 610kT instead of 620kT
		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, Cargo{Germanium: 70}, u.FleetFor(1, "Teamster #3").Cargo)
		assert.Equal(t, Cargo{Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Equal(t, 610, u.Planet("Planet 1").Cargo.Ironium)
	})

	t.Run("shortfall reduces cargo we gave away", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// player 1 loads 100kT ironium from the salvage into Teamster #2 and gives all 100kT to player 2's
		// Hauler. Teamster #2 is empty again
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		u.TransferByHand(1, "Teamster #2", "Hauler #1", Cargo{Ironium: -100})
		// someone else took 40kT before player 1's load is settled. The salvage has 60kT
		u.Game.Salvages[0].Cargo.Ironium = 60

		u.turn.fleetByHandTransfers()

		// player 1 only got 60kT, so they can only give 60kT. The 40kT shortfall comes out of the gift to the
		// Hauler instead of Teamster #1's 50kT ironium
		assert.Equal(t, Cargo{}, u.FleetFor(1, "Teamster #2").Cargo)
		assert.Equal(t, Cargo{Ironium: 60}, u.FleetFor(2, "Hauler #1").Cargo)
		// one message for the short load, one for the short gift
		assert.Len(t, byHandIncompleteMessages(u.Player(1)), 2)
	})

	t.Run("cargo a target can't hold comes back", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// Teamster #1 gives all 50kT of its ironium to player 2's Hauler (Teamster #1 -> 20kT boranium)
		u.TransferByHand(1, "Teamster #1", "Hauler #1", Cargo{Ironium: -50})
		// player 2 loaded 190kT germanium into the Hauler before player 1's gift is settled. The Hauler holds
		// 210kT, so it only has room for 20kT
		hauler := u.FleetFor(2, "Hauler #1")
		hauler.Cargo = Cargo{Germanium: hauler.Spec.CargoCapacity - 20}

		u.turn.fleetByHandTransfers()

		// the Hauler takes 20kT and the other 30kT goes back to Teamster #1 (30kT ironium, 20kT boranium)
		assert.Equal(t, Cargo{Germanium: hauler.Spec.CargoCapacity - 20, Ironium: 20}, hauler.Cargo)
		assert.Equal(t, Cargo{Ironium: 30, Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		messages := byHandIncompleteMessages(u.Player(1))
		if assert.Len(t, messages, 1) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Ironium, Transfered: 20, Wanted: 50, Status: CargoTransferStatusDestCargoCapacity}, *messages[0].Spec.CargoTransfer)
		}
	})

	t.Run("cargo for a fleet that's gone comes back", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		// Teamster #1 gives all 50kT of its ironium to player 2's Hauler (Teamster #1 -> 20kT boranium)
		u.TransferByHand(1, "Teamster #1", "Hauler #1", Cargo{Ironium: -50})
		// the Hauler is gone (merged or scrapped) before the gift is settled
		u.FleetFor(2, "Hauler #1").Delete = true

		u.turn.fleetByHandTransfers()

		// all 50kT goes back to Teamster #1
		assert.Equal(t, Cargo{Ironium: 50, Boranium: 20}, u.FleetFor(1, "Teamster #1").Cargo)
		assert.Len(t, byHandIncompleteMessages(u.Player(1)), 1)
	})

	t.Run("fuel given to another player's fleet is delivered", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		fleet, hauler := u.FleetFor(1, "Teamster #1"), u.FleetFor(2, "Hauler #1")
		hauler.Fuel = 0
		player.GetFleetIntel(hauler.PlayerNum, hauler.Num).Fuel = 0
		startFuel := fleet.Fuel
		// Teamster #1 (900mg) gives 50mg to the Hauler, which starts with an empty 450mg tank.
		// Teamster #1 drops to 850mg right away. The Hauler gets the fuel when it's settled
		if err := NewOrderer().TransferByHand(&u.Game.Rules, player, fleet, player.GetFleetIntel(hauler.PlayerNum, hauler.Num), CargoTransferRequest{Fuel: -50}); err != nil {
			t.Fatal(err)
		}

		u.turn.fleetByHandTransfers()

		assert.Equal(t, startFuel-50, fleet.Fuel)
		assert.Equal(t, 50, hauler.Fuel)
		assert.Empty(t, byHandIncompleteMessages(player))
	})

	t.Run("fuel a full tank can't hold comes back", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		fleet, hauler := u.FleetFor(1, "Teamster #1"), u.FleetFor(2, "Hauler #1")
		hauler.Fuel = 0
		player.GetFleetIntel(hauler.PlayerNum, hauler.Num).Fuel = 0
		startFuel := fleet.Fuel
		// Teamster #1 (900mg) gives 50mg to the Hauler (Teamster #1 -> 850mg)
		if err := NewOrderer().TransferByHand(&u.Game.Rules, player, fleet, player.GetFleetIntel(hauler.PlayerNum, hauler.Num), CargoTransferRequest{Fuel: -50}); err != nil {
			t.Fatal(err)
		}
		// player 2 refueled the Hauler to 430mg of 450mg before the gift is settled, so it only has room for 20mg
		hauler.Fuel = hauler.Spec.FuelCapacity - 20

		u.turn.fleetByHandTransfers()

		// the Hauler fills up with 20mg and the other 30mg goes back to Teamster #1 (850mg -> 880mg)
		assert.Equal(t, startFuel-20, fleet.Fuel)
		assert.Equal(t, hauler.Spec.FuelCapacity, hauler.Fuel)
		messages := byHandIncompleteMessages(player)
		if assert.Len(t, messages, 1) {
			assert.Equal(t, PlayerMessageSpecCargoTransfer{CargoType: Fuel, Transfered: 20, Wanted: 50, Status: CargoTransferStatusDestCargoCapacity}, *messages[0].Spec.CargoTransfer)
		}
	})

	t.Run("fuel only moves between our fleets and other fleets", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		fleet, hauler := u.FleetFor(1, "Teamster #1"), u.FleetFor(2, "Hauler #1")
		orderer := NewOrderer()
		// these are all rejected when the player makes them, so nothing is recorded to settle.
		// can't take fuel from another player's fleet, even with an empty tank
		assert.Error(t, orderer.TransferByHand(&u.Game.Rules, player, u.FleetFor(1, "Teamster #2").withFuel(0), player.GetFleetIntel(hauler.PlayerNum, hauler.Num), CargoTransferRequest{Fuel: 10}))
		// can't jettison fuel or give it to a planet
		assert.Error(t, orderer.TransferByHand(&u.Game.Rules, player, fleet, nil, CargoTransferRequest{Fuel: -10}))
		assert.Error(t, orderer.TransferByHand(&u.Game.Rules, player, fleet, u.Planet("Planet 1"), CargoTransferRequest{Fuel: -10}))
		assert.Empty(t, player.CargoTransfers)
	})

	t.Run("colonists can't be put in space", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		fleet := u.FleetFor(1, "Teamster #1").withCargo(Cargo{Colonists: 10})
		orderer := NewOrderer()
		// Teamster #1 has 10kT of colonists. Jettisoning them or putting them in the salvage is rejected, and
		// they stay aboard
		assert.Error(t, orderer.TransferByHand(&u.Game.Rules, player, fleet, nil, CargoTransferRequest{Cargo: Cargo{Colonists: -10}}))
		assert.Error(t, orderer.TransferByHand(&u.Game.Rules, player, fleet, player.GetSalvageIntel(u.Game.Salvages[0].Num), CargoTransferRequest{Cargo: Cargo{Colonists: -10}}))
		assert.Equal(t, Cargo{Colonists: 10}, fleet.Cargo)
		assert.Empty(t, player.CargoTransfers)
	})

	t.Run("colonists recorded in space come back", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		fleet := u.FleetFor(1, "Teamster #1")
		// orders reject this now, but a game in progress could have a jettison of 10kT of colonists from
		// Teamster #1 recorded before that. Teamster #1 already lost the colonists when it was made
		player.CargoTransfers = CargoTransfers{fleet.Position.String(): []ByHandCargoTransfer{
			{SourceFleetNum: fleet.Num, MapObjectTarget: MapObjectTarget{TargetPosition: fleet.Position}, Cargo: Cargo{Colonists: 10}},
		}}

		u.turn.fleetByHandTransfers()

		// the colonists go back to Teamster #1 instead of dying in a salvage
		assert.Equal(t, 10, fleet.Cargo.Colonists)
		messages := byHandIncompleteMessages(player)
		if assert.Len(t, messages, 1) {
			assert.Equal(t, CargoTransferStatusDeepSpace, messages[0].Spec.CargoTransfer.Status)
		}
	})

	t.Run("transfers follow merged fleets", func(t *testing.T) {
		u := newByHandTestUniverse(t)
		player := u.Player(1)
		// Teamster #2 loads 100kT ironium from the salvage, then merges into Teamster #3
		// (Teamster #3 -> 70kT germanium, 100kT ironium). The salvage load now belongs to Teamster #3
		u.TransferByHand(1, "Teamster #2", "Salvage #1", Cargo{Ironium: 100})
		fleet2, fleet3 := u.FleetFor(1, "Teamster #2"), u.FleetFor(1, "Teamster #3")
		if _, err := NewOrderer().Merge(&u.Game.Rules, player, []*Fleet{fleet3, fleet2}); err != nil {
			t.Fatal(err)
		}
		// someone else took 30kT before the load is settled. The salvage has 70kT
		u.Game.Salvages[0].Cargo.Ironium = 70

		u.turn.fleetByHandTransfers()

		// the 30kT shortfall comes out of Teamster #3 (100kT -> 70kT), and the message is about Teamster #3
		assert.Equal(t, Cargo{Germanium: 70, Ironium: 70}, fleet3.Cargo)
		messages := byHandIncompleteMessages(player)
		if assert.Len(t, messages, 1) {
			assert.Equal(t, fleet3.Num, messages[0].Target.TargetNum)
		}
	})
}

// byHandFuel totals the fuel in every fleet
func byHandFuel(u *testUniverse) int {
	total := 0
	for _, fleet := range u.Game.Fleets {
		if !fleet.Delete {
			total += fleet.Fuel
		}
	}
	return total
}

// byHandMinerals totals the minerals in everything that holds them at the by hand scenario's location
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
			// a fleet that never transfers anything never loses or gains anything
			bystander := u.FleetFor(1, "Bystander #1")
			orderer := NewOrderer()
			// player 2's freighter has room for fuel
			hauler.Fuel = 0
			haulerIntel := player.GetFleetIntel(hauler.PlayerNum, hauler.Num)
			fuelGiven := 0
			start := byHandMinerals(u)
			startFuel := byHandFuel(u)

			playerFleets := func() []*Fleet {
				fleets := []*Fleet{}
				for _, fleet := range u.Game.Fleets {
					if fleet.PlayerNum == player.Num && !fleet.Delete && fleet != bystander {
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
					request := CargoTransferRequest{Cargo: Cargo{}.WithCargo(CargoTypes[rng.Intn(3)], rng.Intn(161)-80)}
					if _, ok := dest.(*Fleet); ok && rng.Intn(3) == 0 {
						request = CargoTransferRequest{Fuel: rng.Intn(161) - 80}
						if dest == haulerIntel && fuelGiven-request.Fuel > hauler.Spec.FuelCapacity {
							// we don't know their tank size, but keep it from overflowing so nothing comes back
							continue
						}
					}
					// invalid transfers are rejected, just like in the UI
					if err := orderer.TransferByHand(rules, player, fleet, dest, request); err == nil && dest == haulerIntel {
						fuelGiven -= request.Fuel
					}
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
			wantHaulerFuel := fuelGiven
			wantFleetFuel := make([]int, len(fleets))
			for i, fleet := range fleets {
				wantFleetFuel[i] = fleet.Fuel
			}
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
				hauler.Fuel += rng.Intn(hauler.availableFuelSpace() + 1)
			}

			u.turn.fleetByHandTransfers()

			got := byHandMinerals(u)
			assert.Equal(t, start.Add(outside), got, "no cargo is created or destroyed")
			for _, fleet := range append(playerFleets(), hauler) {
				assert.False(t, fleet.Cargo.HasNegative())
				assert.LessOrEqual(t, fleet.Cargo.Total(), fleet.Spec.CargoCapacity)
				assert.GreaterOrEqual(t, fleet.Fuel, 0)
				assert.LessOrEqual(t, fleet.Fuel, fleet.Spec.FuelCapacity)
			}
			assert.Empty(t, player.CargoTransfers)
			assert.Equal(t, Cargo{Ironium: 100}, bystander.Cargo)
			assert.Equal(t, bystander.Spec.FuelCapacity, bystander.Fuel)

			if interfere {
				return
			}

			assert.Empty(t, byHandIncompleteMessages(player))
			assert.Equal(t, startFuel, byHandFuel(u), "no fuel is created or destroyed")
			assert.Equal(t, wantHaulerFuel, hauler.Fuel)
			for i, fleet := range fleets {
				assert.Equal(t, wantFleetCargo[i], fleet.Cargo, fleet.Name)
				assert.Equal(t, wantFleetFuel[i], fleet.Fuel, fleet.Name)
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
