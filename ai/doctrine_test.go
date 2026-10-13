package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDoctrineTestAI(prt cs.PRT) *aiPlayer {
	game := cs.NewGame()
	race := cs.NewRace()
	race.PRT = prt
	player := cs.NewPlayer(1, race.WithSpec(&game.Rules))
	player.Num = 1
	return NewAIPlayer(game, &cs.StaticTechStore, player, cs.PlayerMapObjects{})
}

func TestDoctrineFollowsTheRace(t *testing.T) {
	joat, he, wm := newDoctrineTestAI(cs.JoaT), newDoctrineTestAI(cs.HE), newDoctrineTestAI(cs.WM)

	assert.Equal(t, doctrine{}, joat.doctrine, "JoaT plays the standard game")
	assert.Greater(t, he.colonyLimit(), joat.colonyLimit())
	assert.Less(t, he.seedHabitability(), joat.seedHabitability())
	assert.Equal(t, joat.config.startAttackingYear-5, wm.config.startAttackingYear)
}

func TestDoctrineCanBeSwitchedOffForATest(t *testing.T) {
	Experiment = func(playerNum int, behavior string) bool { return behavior == "doctrine-HE" }
	t.Cleanup(func() { Experiment = nil })
	assert.Equal(t, doctrine{}, newDoctrineTestAI(cs.HE).doctrine)
	assert.NotEqual(t, doctrine{}, newDoctrineTestAI(cs.WM).doctrine)
}

func TestCAJudgesAWorldByWhatItWillBecome(t *testing.T) {
	world := cs.NewPlanet()
	world.Spec.Habitability = -10
	world.Spec.TerraformedHabitability = 60

	ca := newDoctrineTestAI(cs.CA)
	require.True(t, ca.Race.Spec.Instaforming)
	assert.Equal(t, 60, ca.colonyHabitability(world))
	assert.Equal(t, -10, newDoctrineTestAI(cs.JoaT).colonyHabitability(world))
}

func TestISPlansDefensesAheadOfTerraforming(t *testing.T) {
	is := newDoctrineTestAI(cs.IS)
	require.NoError(t, is.plan())
	items := is.PlayerPlans.ProductionPlans[0].Items
	defenses, terraform := -1, -1
	for i, item := range items {
		if item.Type == cs.QueueItemTypeAutoDefenses && defenses == -1 {
			defenses = i
		}
		if item.Type == cs.QueueItemTypeAutoMaxTerraform {
			terraform = i
		}
	}
	require.NotEqual(t, -1, defenses)
	require.NotEqual(t, -1, terraform)
	assert.Less(t, defenses, terraform)
	assert.Equal(t, doctrines[cs.IS].defensesPerYear, items[defenses].Quantity)
}

func TestSSScoutsCarryACloak(t *testing.T) {
	ss := newDoctrineTestAI(cs.SS)
	ss.TechLevels = cs.TechLevel{Energy: 5, Weapons: 5, Propulsion: 5, Construction: 5, Electronics: 5, Biotechnology: 5}
	design, err := ss.designShip("Scout", cs.ShipDesignPurposeScout, cs.FleetPurposeScout)
	require.NoError(t, err)
	require.NotNil(t, design)
	assert.Equal(t, cs.ShipDesignPurposeScout, design.Purpose, "still a scout to the rest of the AI")

	joat := newDoctrineTestAI(cs.JoaT)
	joat.TechLevels = ss.TechLevels
	plain, err := joat.designShip("Scout", cs.ShipDesignPurposeScout, cs.FleetPurposeScout)
	require.NoError(t, err)
	assert.Greater(t, design.Spec.CloakPercent, plain.Spec.CloakPercent)
}

func TestSDDetonatesOnlyWhenTheEnemyIsInTheField(t *testing.T) {
	sd := newDoctrineTestAI(cs.SD)
	require.True(t, sd.Race.Spec.CanDetonateMinefields)
	sd.Relations = []cs.PlayerRelationship{{}, {Relation: cs.PlayerRelationEnemy}}

	field := &cs.Minefield{MinefieldType: cs.MinefieldTypeStandard, NumMines: 10_000}
	field.Position = cs.Vector{X: 100, Y: 100}
	sd.Minefields = []*cs.Minefield{field}

	// an armed enemy design we have seen
	warship := cs.NewShipDesign(2, 1)
	warship.Spec.PowerRating = 100
	sd.ShipDesignIntels = append(sd.ShipDesignIntels, warship)
	raiders := &cs.Fleet{Tokens: []cs.ShipToken{{DesignNum: 1, Quantity: detonateEnemyShips}}}
	raiders.PlayerNum = 2
	raiders.Position = field.Position
	sd.FleetIntels = []*cs.Fleet{raiders}

	require.NoError(t, sd.detonateMinefields())
	assert.True(t, field.Detonate, "a raiding force is inside and none of our ships are")

	// our own fleet in the field would be hit too
	ours := &cs.Fleet{}
	ours.Position = field.Position
	ours.Spec.TotalShips = detonateEnemyShips
	sd.Fleets = []*cs.Fleet{ours}
	require.NoError(t, sd.detonateMinefields())
	assert.False(t, field.Detonate)

	// and it stops once the enemy has gone
	sd.Fleets, sd.FleetIntels = nil, nil
	field.Detonate = true
	require.NoError(t, sd.detonateMinefields())
	assert.False(t, field.Detonate)
}

func TestITOutpostIsAnOrbitalFortWithAGate(t *testing.T) {
	it := newDoctrineTestAI(cs.IT)
	it.TechLevels = cs.TechLevel{Energy: 5, Weapons: 5, Propulsion: 5, Construction: 5, Electronics: 5, Biotechnology: 5}
	require.NoError(t, it.designStarbases())
	require.NotNil(t, it.outpostDesign)
	assert.NotEmpty(t, it.outpostDesign.Spec.Stargate)
	assert.Equal(t, cs.OrbitalFort.Name, it.outpostDesign.Hull)

	joat := newDoctrineTestAI(cs.JoaT)
	joat.TechLevels = it.TechLevels
	require.NoError(t, joat.designStarbases())
	assert.Nil(t, joat.outpostDesign)
}
