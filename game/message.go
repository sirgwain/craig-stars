package game

import (
	"fmt"
	"time"
)

type PlayerMessage struct {
	ID                 int64                   `json:"id"`
	CreatedAt          time.Time               `json:"createdAt"`
	UpdatedAt          time.Time               `json:"updatedAt"`
	PlayerID           int64                   `json:"playerId"`
	Type               PlayerMessageType       `json:"type,omitempty"`
	Text               string                  `json:"text,omitempty"`
	TargetMapObjectNum int                     `json:"targetMapObjectNum,omitempty"`
	TargetPlayerNum    int                     `json:"targetPlayerNum,omitempty"`
	TargetType         PlayerMessageTargetType `json:"targetType,omitempty"`
}

type PlayerMessageTargetType string

const (
	TargetNone          PlayerMessageTargetType = "None"
	TargetPlanet        PlayerMessageTargetType = "Planet"
	TargetFleet         PlayerMessageTargetType = "Fleet"
	TargetWormhole      PlayerMessageTargetType = "Wormhole"
	TargetMineField     PlayerMessageTargetType = "MineField"
	TargetMysteryTrader PlayerMessageTargetType = "MysteryTrader"
	TargetBattle        PlayerMessageTargetType = "Battle"
)

type PlayerMessageType string

const (
	PlayerMessageInfo                               PlayerMessageType = "Info"
	PlayerMessageError                              PlayerMessageType = "Error"
	PlayerMessageHomePlanet                         PlayerMessageType = "HomePlanet"
	PlayerMessagePlayerDiscovery                    PlayerMessageType = "PlayerDiscovery"
	PlayerMessagePlanetDiscovery                    PlayerMessageType = "PlanetDiscovery"
	PlayerMessagePlanetProductionQueueEmpty         PlayerMessageType = "PlanetProductionQueueEmpty"
	PlayerMessagePlanetProductionQueueComplete      PlayerMessageType = "PlanetProductionQueueComplete"
	PlayerMessageBuiltMine                          PlayerMessageType = "BuiltMine"
	PlayerMessageBuiltFactory                       PlayerMessageType = "BuiltFactory"
	PlayerMessageBuiltDefense                       PlayerMessageType = "BuiltDefense"
	PlayerMessageBuiltShip                          PlayerMessageType = "BuiltShip"
	PlayerMessageBuiltStarbase                      PlayerMessageType = "BuiltStarbase"
	PlayerMessageBuiltMineralPacket                 PlayerMessageType = "BuiltMineralPacket"
	PlayerMessageBuiltTerraform                     PlayerMessageType = "BuiltTerraform"
	PlayerMessageFleetOrdersComplete                PlayerMessageType = "FleetOrdersComplete"
	PlayerMessageFleetEngineFailure                 PlayerMessageType = "FleetEngineFailure"
	PlayerMessageFleetOutOfFuel                     PlayerMessageType = "FleetOutOfFuel"
	PlayerMessageFleetGeneratedFuel                 PlayerMessageType = "FleetGeneratedFuel"
	PlayerMessageFleetScrapped                      PlayerMessageType = "FleetScrapped"
	PlayerMessageFleetMerged                        PlayerMessageType = "FleetMerged"
	PlayerMessageFleetInvalidMergeNotFleet          PlayerMessageType = "FleetInvalidMergeNotFleet"
	PlayerMessageFleetInvalidMergeUnowned           PlayerMessageType = "FleetInvalidMergeUnowned"
	PlayerMessageFleetPatrolTargeted                PlayerMessageType = "FleetPatrolTargeted"
	PlayerMessageFleetInvalidRouteNotFriendlyPlanet PlayerMessageType = "FleetInvalidRouteNotFriendlyPlanet"
	PlayerMessageFleetInvalidRouteNotPlanet         PlayerMessageType = "FleetInvalidRouteNotPlanet"
	PlayerMessageFleetInvalidRouteNoRouteTarget     PlayerMessageType = "FleetInvalidRouteNoRouteTarget"
	PlayerMessageFleetRoute                         PlayerMessageType = "FleetRoute"
	PlayerMessageInvalid                            PlayerMessageType = "Invalid"
	PlayerMessagePlanetColonized                    PlayerMessageType = "PlanetColonized"
	PlayerMessageGainTechLevel                      PlayerMessageType = "GainTechLevel"
	PlayerMessageMyPlanetBombed                     PlayerMessageType = "MyPlanetBombed"
	PlayerMessageMyPlanetRetroBombed                PlayerMessageType = "MyPlanetRetroBombed"
	PlayerMessageEnemyPlanetBombed                  PlayerMessageType = "EnemyPlanetBombed"
	PlayerMessageEnemyPlanetRetroBombed             PlayerMessageType = "EnemyPlanetRetroBombed"
	PlayerMessageMyPlanetInvaded                    PlayerMessageType = "MyPlanetInvaded"
	PlayerMessageEnemyPlanetInvaded                 PlayerMessageType = "EnemyPlanetInvaded"
	PlayerMessageBattle                             PlayerMessageType = "Battle"
	PlayerMessageCargoTransferred                   PlayerMessageType = "CargoTransferred"
	PlayerMessageMinesSwept                         PlayerMessageType = "MinesSwept"
	PlayerMessageMinesLaid                          PlayerMessageType = "MinesLaid"
	PlayerMessageMineFieldHit                       PlayerMessageType = "MineFieldHit"
	PlayerMessageFleetDumpedCargo                   PlayerMessageType = "FleetDumpedCargo"
	PlayerMessageFleetStargateDamaged               PlayerMessageType = "FleetStargateDamaged"
	PlayerMessageMineralPacketCaught                PlayerMessageType = "MineralPacketCaught"
	PlayerMessageMineralPacketDamage                PlayerMessageType = "MineralPacketDamage"
	PlayerMessageMineralPacketLanded                PlayerMessageType = "MineralPacketLanded"
	PlayerMessageVictor                             PlayerMessageType = "Victor"
	PlayerMessageFleetReproduce                     PlayerMessageType = "FleetReproduce"
	PlayerMessageRandomMineralDeposit               PlayerMessageType = "RandomMineralDeposit"
	PlayerMessagePermaform                          PlayerMessageType = "Permaform"
	PlayerMessageInstaform                          PlayerMessageType = "Instaform"
	PlayerMessagePacketTerraform                    PlayerMessageType = "PacketTerraform"
	PlayerMessagePacketPermaform                    PlayerMessageType = "PacketPermaform"
	PlayerMessageRemoteMined                        PlayerMessageType = "RemoteMined"
)

type Messager interface {
	homePlanet(player *Player, planet *Planet)
}

type messageClient struct {
}

var messager = messageClient{}

func (m *messageClient) error(player *Player, err error) {
	text := fmt.Sprintf("Something went wrong on the server. Please contact the administrator, %v", err)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageError, Text: text})
}

func (m *messageClient) homePlanet(player *Player, planet *Planet) {
	text := fmt.Sprintf("Your home planet is %s. Your people are ready to leave the nest and explore the universe.  Good luck.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.Num})
}

func (m *messageClient) longMessage(player *Player) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text})
}

func (m *messageClient) minesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d mine(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltMine, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.Num})
}

func (m *messageClient) factoriesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d factory(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.Num})
}

func (m *messageClient) defensesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d defense(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.Num})
}

func (m *messageClient) fleetBuilt(player *Player, planet *Planet, fleet *Fleet, num int) {

	var text string

	if num == 1 {
		text = fmt.Sprintf("Your starbase at %s has built a new %s.", planet.Name, fleet.BaseName)
	} else {
		text = fmt.Sprintf("Your starbase at %s has built %d new %ss.", planet.Name, num, fleet.BaseName)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltShip, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) fleetTransportedCargo(player *Player, fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) {
	text := ""
	if cargoType == Colonists {
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has beamed %d %s from %s", fleet.Name, -transferAmount*100, cargoType, dest.GetMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has beamed %d %s to %s", fleet.Name, transferAmount*100, cargoType, dest.GetMapObject().Name)
		}
	} else {
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has loaded %d %s from %s", fleet.Name, -transferAmount, cargoType, dest.GetMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has unloaded %d %s to %s", fleet.Name, transferAmount, cargoType, dest.GetMapObject().Name)
		}
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageCargoTransferred, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) fleetEngineFailure(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to engage it's engines due to balky equipment. Engineers think they have the problem fixed for the time being.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetEngineFailure, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) fleetOutOfFuel(player *Player, fleet *Fleet, warpFactor int) {
	text := fmt.Sprintf("%s has run out of fuel. The fleet's speed has been decreased to Warp %d.", fleet.Name, warpFactor)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetOutOfFuel, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) fleetGeneratedFuel(player *Player, fleet *Fleet, fuelGenerated int) {
	text := fmt.Sprintf("%s's ram scoops have produced %dmg of fuel from interstellar hydrogen.", fleet.Name, fuelGenerated)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetGeneratedFuel, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) colonizeNonPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a waypoint with no Planet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) colonizeOwnedPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet that is already inhabited.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})

}

func (m *messageClient) colonizeWithNoModule(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet without a colonization module.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})

}

func (m *messageClient) colonizeWithNoColonists(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet without bringing any colonists.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetMapObjectNum: fleet.Num})
}

func (m *messageClient) planetColonized(player *Player, planet *Planet) {
	text := "Your colonists are now in control of {planet.Name}"
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetColonized, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.Num})
}
