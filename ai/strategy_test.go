package ai

import (
	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDifficultyUsesDifferentStrategies(t *testing.T) {
	game := cs.NewGame()
	configs := map[cs.AIDifficulty]playerConfig{}
	for _, difficulty := range []cs.AIDifficulty{cs.AIDifficultyEasy, cs.AIDifficultyNormal, cs.AIDifficultyHard, cs.AIDifficultyCheater} {
		p := cs.NewPlayer(1, cs.NewRace().WithSpec(&game.Rules))
		p.AIDifficulty = difficulty
		configs[difficulty] = NewAIPlayer(game, &cs.StaticTechStore, p, cs.PlayerMapObjects{}).config
	}
	assert.Less(t, configs[cs.AIDifficultyEasy].colonyPipeline, configs[cs.AIDifficultyNormal].colonyPipeline)
	assert.Less(t, configs[cs.AIDifficultyNormal].colonyPipeline, configs[cs.AIDifficultyHard].colonyPipeline)
	assert.Greater(t, configs[cs.AIDifficultyEasy].startAttackingYear, configs[cs.AIDifficultyNormal].startAttackingYear)
	assert.Greater(t, configs[cs.AIDifficultyNormal].startAttackingYear, configs[cs.AIDifficultyHard].startAttackingYear)
	assert.Equal(t, configs[cs.AIDifficultyHard], configs[cs.AIDifficultyCheater])
}

func TestUnavailableDesignCanBeReconsidered(t *testing.T) {
	game := cs.NewGame()
	p := cs.NewPlayer(1, cs.NewRace().WithSpec(&game.Rules))
	ai := NewAIPlayer(game, &cs.StaticTechStore, p, cs.PlayerMapObjects{})
	ai.designsByPurpose[cs.ShipDesignPurposeTorpedoFighter] = nil
	_, err := ai.designShip("Torpedo", cs.ShipDesignPurposeTorpedoFighter, cs.FleetPurposeFighter)
	require.NoError(t, err)
}

func TestInvaderTargetsWeakEnemyWorld(t *testing.T) {
	ai, source, target, f := colonizationFixture(t)
	design, err := cs.DesignShip(&ai.game.Rules, ai.Player, &cs.Privateer, "Freighter", 200, 0, cs.ShipDesignPurposeColonistFreighter, cs.FleetPurposeColonistFreighter)
	require.NoError(t, err)
	ai.Designs = append(ai.Designs, design)
	f.Tokens = []cs.ShipToken{{DesignNum: design.Num, Quantity: 1}}
	f.InjectDesigns(ai.Designs)
	f.Spec = cs.ComputeFleetSpec(&ai.game.Rules, ai.Player, f)
	require.False(t, f.Spec.Colonizer)
	f.Purpose = cs.FleetPurposeColonistFreighter
	f.SetTag(cs.TagPurpose, string(cs.FleetPurposeColonistFreighter))
	ai.Relations = []cs.PlayerRelationship{{}, {Relation: cs.PlayerRelationEnemy}}
	target.PlayerNum = 2
	target.Cargo.Colonists = 50
	source.Cargo.Colonists = source.Spec.MaxPopulation / 100

	require.NoError(t, ai.invade())
	require.Len(t, f.Waypoints, 2)
	assert.Equal(t, cs.FleetPurposeInvader, f.Purpose)
	assert.Equal(t, target.Num, f.Waypoints[1].TargetNum)
	assert.Equal(t, cs.TransportActionUnloadAll, f.Waypoints[1].TransportTasks.Colonists.Action)
	assert.Equal(t, f.Spec.CargoCapacity, f.Cargo.Colonists)

	// a starbase makes the world too dangerous, come home
	target.Spec.HasStarbase = true
	require.NoError(t, ai.invade())
	assert.Equal(t, source.Num, f.Waypoints[1].TargetNum)
	assert.NotEqual(t, cs.FleetPurposeInvader, f.Purpose)
}

func TestARColonyShipKeepsItsColonyHull(t *testing.T) {
	game := cs.NewGame()
	race := cs.NewRace()
	race.PRT = cs.AR
	p := cs.NewPlayer(1, race.WithSpec(&game.Rules))
	p.TechLevels = cs.TechLevel{Energy: 5, Weapons: 5, Propulsion: 5, Construction: 5, Electronics: 5, Biotechnology: 5}
	ai := NewAIPlayer(game, &cs.StaticTechStore, p, cs.PlayerMapObjects{})
	design, err := ai.designShip("Colony", cs.ShipDesignPurposeColonizer, cs.FleetPurposeColonizer)
	require.NoError(t, err)
	require.NotNil(t, design)
	assert.True(t, design.Spec.Colonizer)
	assert.Equal(t, cs.ColonyShip.Name, design.Hull)
}
