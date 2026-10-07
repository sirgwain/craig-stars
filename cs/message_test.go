package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_message_eventSnapshots(t *testing.T) {
	player := &Player{Num: 1}
	fleet := &Fleet{MapObject: MapObject{Type: MapObjectTypeFleet, PlayerNum: 1, Num: 2, Name: "Original"}}
	planet := &Planet{MapObject: MapObject{Type: MapObjectTypePlanet, Num: 3, Name: "Source"}, Hab: Hab{49, 60, 70}}
	dest := &Fleet{MapObject: MapObject{Type: MapObjectTypeFleet, PlayerNum: 1, Num: 4, Name: "Survivor"}}

	messager.fleetMerged(player, fleet, dest)
	fleet.Name = "Changed"
	assert.Equal(t, "Original", player.Messages[0].Spec.Name)
	assert.Equal(t, "Survivor", player.Messages[0].TargetName)
	assert.Equal(t, dest.Num, player.Messages[0].TargetNum)

	for _, habType := range HabTypes {
		messager.planetTerraform(player, planet, habType, -1)
		message := player.Messages[len(player.Messages)-1]
		assert.Equal(t, FromHabType(habType), message.Spec.HabType)
		assert.Equal(t, -1, message.Spec.Amount)
		assert.Equal(t, planet.Hab.Get(habType), message.Spec.Amount2)
		planet.Hab.Set(habType, 99)
		assert.NotEqual(t, planet.Hab.Get(habType), message.Spec.Amount2)
	}

	messager.fleetTransportedCargo(player, fleet, planet, Colonists, -12)
	transfer := player.Messages[len(player.Messages)-1]
	require.NotNil(t, transfer.Spec.CargoTransfer)
	assert.Equal(t, Colonists, transfer.Spec.CargoTransfer.CargoType)
	assert.Equal(t, -12, transfer.Spec.CargoTransfer.Transfered)
	assert.Equal(t, "Source", transfer.Spec.TargetName)
	planet.Name = "Changed"
	assert.Equal(t, "Source", transfer.Spec.TargetName)

	for _, message := range player.Messages {
		assert.Empty(t, message.Text)
	}
}

func Test_message_stargateSnapshots(t *testing.T) {
	player := &Player{Num: 1}
	fleet := &Fleet{MapObject: MapObject{PlayerNum: 1, Num: 2, Name: "Jumper"}}
	source := Waypoint{MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1, TargetName: "Earth"}, Position: Vector{10, 20}}
	dest := Waypoint{MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 2, TargetName: "Mars"}, Position: Vector{30, 40}}
	cargo := Cargo{Ironium: 2, Boranium: 3, Germanium: 4, Colonists: 5}
	messager.fleetStargateDumpedCargo(player, fleet, source, dest, cargo)
	messager.fleetStargateDamaged(player, fleet, source, dest, 25, 3, 4)
	messager.fleetStargateDestroyed(player, fleet, source, dest)
	messager.fleetStargateInvalidRange(player, fleet, source, dest, 123.45)
	cargo.Colonists = 0
	source.TargetName = "Changed"
	dest.TargetName = "Changed"

	for _, message := range player.Messages {
		assert.Empty(t, message.Text)
		assert.Equal(t, "Earth", message.Spec.TargetName)
		assert.Equal(t, source.Position, message.Spec.TargetPosition)
		require.NotNil(t, message.Spec.RouteTarget)
		assert.Equal(t, "Mars", message.Spec.RouteTarget.TargetName)
		assert.Equal(t, dest.Position, message.Spec.RouteTarget.TargetPosition)
	}
	assert.Equal(t, PlayerMessageFleetDumpedCargo, player.Messages[0].Type)
	assert.Equal(t, 5, player.Messages[0].Spec.Cargo.Colonists)
	assert.Equal(t, 25, player.Messages[1].Spec.Amount)
	assert.Equal(t, 7, player.Messages[1].Spec.Amount2)
	assert.Equal(t, PlayerMessageFleetStargateDestroyed, player.Messages[2].Type)
	assert.Equal(t, 123.45, player.Messages[3].Spec.Distance)
}

func Test_message_packetImpactSnapshot(t *testing.T) {
	player := &Player{Num: 1}
	planet := &Planet{MapObject: MapObject{Num: 1, Name: "Earth"}}
	planet.Spec.HasStarbase = true
	planet.Starbase = &Fleet{}
	planet.Starbase.Spec.HasMassDriver = true
	packet := &MineralPacket{Cargo: Cargo{Ironium: 100}}
	damage := MineralPacketDamage{Killed: 1200, DefensesDestroyed: 3, Uncaught: 50}
	messager.planetPacketDamage(player, planet, packet, damage)
	packet.Cargo = Cargo{}
	planet.Spec.HasStarbase = false
	planet.Cargo.Colonists = 10
	damage.Killed = 0

	message := player.Messages[0]
	assert.Empty(t, message.Text)
	assert.Equal(t, 100, message.Spec.Amount)
	assert.True(t, message.Spec.HasMassDriver)
	assert.True(t, message.Spec.PlanetEmptied)
	require.NotNil(t, message.Spec.MineralPacketDamage)
	assert.Equal(t, 1200, message.Spec.MineralPacketDamage.Killed)
	assert.Equal(t, 3, message.Spec.MineralPacketDamage.DefensesDestroyed)
	assert.Equal(t, 50, message.Spec.MineralPacketDamage.Uncaught)
}

func Test_message_invasionStartingPopulations(t *testing.T) {
	attacker, defender := &Player{Num: 1}, &Player{Num: 2}
	planet := &Planet{MapObject: MapObject{Num: 3, Name: "Earth", PlayerNum: 2}}
	for _, player := range []*Player{attacker, defender} {
		messager.planetInvaded(player, planet, "Invaders", attacker, defender, 12000, 10000, 8700, 10000, true)
		message := player.Messages[0]
		assert.Equal(t, 12000, message.Spec.Amount)
		assert.Equal(t, 10000, message.Spec.Amount2)
		assert.Equal(t, 8700, message.Spec.Invasion.AttackersKilled)
		assert.Equal(t, 10000, message.Spec.Invasion.DefendersKilled)
		assert.Equal(t, "Earth", message.TargetName)
	}
	assert.Equal(t, PlayerMessageFleetInvadedPlanet, attacker.Messages[0].Type)
	assert.Equal(t, PlayerMessagePlanetInvaded, defender.Messages[0].Type)
}
