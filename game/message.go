package game

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PlayerMessage struct {
	ID                 uint                    `gorm:"primaryKey" json:"id"`
	CreatedAt          time.Time               `json:"createdAt"`
	UpdatedAt          time.Time               `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt          `gorm:"index" json:"deletedAt"`
	PlayerID           uint                    `json:"playerId"`
	Type               PlayerMessageType       `json:"type,omitempty"`
	Text               string                  `json:"text,omitempty"`
	TargetMapObjectNum uint                    `json:"targetMapObjectNum,omitempty"`
	TargetPlayerNum    uint                    `json:"targetPlayerNum,omitempty"`
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

func (m *messageClient) homePlanet(player *Player, planet *Planet) {
	text := fmt.Sprintf("Your home planet is %s. Your people are ready to leave the nest and explore the universe.  Good luck.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.ID})
}

func (m *messageClient) longMessage(player *Player) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text})
}

func (m *messageClient) minesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d mine(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltMine, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.ID})
}

func (m *messageClient) factoriesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d factory(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.ID})
}

func (m *messageClient) defensesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d defense(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetMapObjectNum: planet.ID})
}
