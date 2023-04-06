package cs

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PlayerMessage struct {
	Type            PlayerMessageType       `json:"type,omitempty"`
	Text            string                  `json:"text,omitempty"`
	TargetNum       int                     `json:"targetNum,omitempty"`
	TargetPlayerNum int                     `json:"targetPlayerNum,omitempty"`
	TargetType      PlayerMessageTargetType `json:"targetType,omitempty"`
}

type PlayerMessageTargetType string

const (
	TargetNone          PlayerMessageTargetType = ""
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
	PlayerMessageFleetInvalidTransport              PlayerMessageType = "FleetInvalidTransport"
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
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) longMessage(player *Player) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text})
}

func (m *messageClient) minesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d mine(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltMine, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) factoriesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d factory(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) defensesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d defense(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) fleetBuilt(player *Player, planet *Planet, fleet *Fleet, num int) {

	var text string

	if num == 1 {
		text = fmt.Sprintf("Your starbase at %s has built a new %s.", planet.Name, fleet.BaseName)
	} else {
		text = fmt.Sprintf("Your starbase at %s has built %d new %ss.", planet.Name, num, fleet.BaseName)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltShip, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetTransportedCargo(player *Player, fleet *Fleet, dest cargoHolder, cargoType CargoType, transferAmount int) {
	text := ""
	if cargoType == Colonists {
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has beamed %d %s from %s", fleet.Name, -transferAmount*100, cargoType, dest.getMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has beamed %d %s to %s", fleet.Name, transferAmount*100, cargoType, dest.getMapObject().Name)
		}
	} else {
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has loaded %d of %s from %s", fleet.Name, -transferAmount, cargoType, dest.getMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has unloaded %d of %s to %s", fleet.Name, transferAmount, cargoType, dest.getMapObject().Name)
		}
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageCargoTransferred, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetInvalidLoadCargo(player *Player, fleet *Fleet, dest cargoHolder, cargoType CargoType, transferAmount int) {
	text := fmt.Sprintf("%s attempted to load %dkT of %v from %s, but you do not own %s", fleet.Name, transferAmount, cargoType, dest.getMapObject().Name, dest.getMapObject().Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidTransport, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})

}

func (m *messageClient) fleetEngineFailure(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to engage it's engines due to balky equipment. Engineers think they have the problem fixed for the time being.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetEngineFailure, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetOutOfFuel(player *Player, fleet *Fleet, warpFactor int) {
	text := fmt.Sprintf("%s has run out of fuel. The fleet's speed has been decreased to Warp %d.", fleet.Name, warpFactor)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetOutOfFuel, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetGeneratedFuel(player *Player, fleet *Fleet, fuelGenerated int) {
	text := fmt.Sprintf("%s's ram scoops have produced %dmg of fuel from interstellar hydrogen.", fleet.Name, fuelGenerated)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetGeneratedFuel, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetScrapped(player *Player, fleet *Fleet, totalMinerals int, resources int, planet *Planet) {
	var text string
	if planet != nil {
		if planet.Spec.HasStarbase {
			text = fmt.Sprintf("%s has been dismantled for %dkT of minerals which have been deposited on %s.", fleet.Name, totalMinerals, planet.Name)
		} else {
			text = fmt.Sprintf("%s has been dismantled for %dkT of minerals at the starbase orbiting %s.", fleet.Name, totalMinerals, planet.Name)
		}
		if resources > 0 {
			text += fmt.Sprintf(" Ultimate recycling has also made %d resources available for immediate use (less if other ships were scrapped here this year).", resources)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		text = fmt.Sprintf("%s has been dismantled. The scrap was left in deep space.", fleet.Name)
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text})
	}
}

func (m *messageClient) fleetMerged(player *Player, fleet *Fleet, mergedInto *Fleet) {
	text := fmt.Sprintf("%s has been merged into %s.", fleet.Name, mergedInto.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetMerged, Text: text, TargetType: TargetFleet, TargetNum: mergedInto.Num, TargetPlayerNum: mergedInto.Num})
}

func (m *messageClient) fleetInvalidMergeNotFleet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to complete it's merge orders as the waypoint destination wasn't a fleet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidMergeNotFleet, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidMergeNotOwned(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to complete it's merge orders as the destination fleet wasn't one of yours.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidMergeUnowned, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetPatrolTargeted(player *Player, fleet *Fleet, target *Fleet) {
	text := fmt.Sprintf("Your patrolling %s has targeted %s for intercept.", fleet.Name, target.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetPatrolTargeted, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidRouteNotPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s could not be routed because it is not at a planet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidRouteNotPlanet, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidRouteNotFriendlyPlanet(player *Player, fleet *Fleet, planet *Planet) {
	text := fmt.Sprintf("%s could not be routed because you are not friends with the owners of %s", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidRouteNotFriendlyPlanet, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidRouteNoRouteTarget(player *Player, fleet *Fleet, planet *Planet) {
	text := fmt.Sprintf("%s could not be routed because %s has no route set.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidRouteNoRouteTarget, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetRouted(player *Player, fleet *Fleet, planet *Planet, target string) {
	text := fmt.Sprintf("%s has been routed by the citizens of %s to %s", fleet.Name, planet.Name, target)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetRoute, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetBuiltForComposition(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetStargateInvalidSource(player *Player, fleet *Fleet, wp0 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s, but no stargate exists there.", fleet.Name, wp0.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidSourceOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s, but could not because the starbase is not owned by you or a friend of yours.", fleet.Name, wp0.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidDest(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but no stargate could be detected at the destination.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidDestOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but could not because the destination starbase is not owned by you or a friend of yours.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidRange(player *Player, fleet *Fleet, wp0, wp1 Waypoint, totalDist float64) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but the distance of %.1f l.y. was outside the max range of the stargates.", fleet.Name, wp0.TargetName, wp1.TargetName, totalDist),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidMass(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but your ships are too massive.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidColonists(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but you are carrying colonists and can't drop them off as you don't own the planet.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateDumpedCargo(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint, cargo Cargo) {
	var text string
	if cargo.HasColonists() && cargo.HasMinerals() {
		text = fmt.Sprintf("%s has unloaded %d colonists and %dkt of minerals in preparation for jumping through the stargate at %s to reach %s.", fleet.Name, cargo.Colonists*100, cargo.Total()-cargo.Colonists, wp0.TargetName, wp1.TargetName)
	} else if cargo.HasColonists() {
		text = fmt.Sprintf("%s has unloaded %d colonists in preparation for jumping through the stargate at %s to reach %s.", fleet.Name, cargo.Colonists*100, wp0.TargetName, wp1.TargetName)
	} else {
		text = fmt.Sprintf("%s has unloaded %dkt of minerals in preparation for jumping through the stargate at %s to reach %s.", fleet.Name, cargo.Total(), wp0.TargetName, wp1.TargetName)
	}
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            text,
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateDestroyed(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageFleetStargateDamaged,
		Text:            fmt.Sprintf("Heedless to the danger, %s attempted to use the stargate at %s to reach %s. The fleet never arrived. The distance or mass must have been too great.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateDamaged(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint, damage int, startingShips int, shipsLostToDamage int, shipsLostToTheVoid int) {
	totalShipsLost := shipsLostToDamage + shipsLostToTheVoid
	var text string
	if totalShipsLost == 0 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing no ships but suffering %d dp of damage. They exceeded the capability of the gates.", fleet.Name, wp0.TargetName, wp1.TargetName, damage)
	} else if totalShipsLost < 5 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing only %d ship%s to the treacherous void. They were fortunate. They exceeded the capability of the gates.", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost, func() string {
			if totalShipsLost == 1 {
				return ""
			} else {
				return "s"
			}
		}())
	} else if totalShipsLost >= 5 && totalShipsLost <= 10 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing %d ships to the unforgiving void. Exceeding the capability of your stargates is not recommended.", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	} else if totalShipsLost >= 10 && totalShipsLost <= 50 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s unfortunately losing %d ships to the great unknown. Exceeding the capability of your stargates is dangerous.", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	} else if totalShipsLost >= 50 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing an unbelievable %d ships. The jump was far in excess of the capabilities of starbases involved..", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	}
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageFleetStargateDamaged,
		Text:            text,
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetReproduce(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetCompletedAssignedOrders(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetHitMineField(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) colonizeNonPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a waypoint with no Planet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) colonizeOwnedPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet that is already inhabited.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})

}

func (m *messageClient) colonizeWithNoModule(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet without a colonization module.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})

}

func (m *messageClient) colonizeWithNoColonists(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet without bringing any colonists.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) planetColonized(player *Player, planet *Planet) {
	text := fmt.Sprintf("Your colonists are now in control of %s", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetColonized, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) planetInvadeEmpty(player *Player, planet *Planet, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to invade %s, but the planet is uninhabited.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) planetInvaded(player *Player, planet *Planet, fleet *Fleet, planetOwner string, fleetOwner string, attackersKilled int, defendersKilled int) {
	var text string

	// use this formatter to get commas on the text
	p := message.NewPrinter(language.English)
	if player.Num == fleet.PlayerNum {
		if planet.PlayerNum == fleet.PlayerNum {
			// we invaded and won
			text = p.Sprintf("Your %s has successfully invaded %s planet %s killing off all colonists", fleet.Name, planetOwner, planet.Name)
		} else {
			// we invaded and lost
			text = p.Sprintf("Your %s tried to invade %s, but all of your colonists were killed by %s. You valiant fighters managed to kill %d of their colonists.", fleet.Name, planet.Name, planetOwner, defendersKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageEnemyPlanetInvaded, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		if planet.PlayerNum == fleet.PlayerNum {
			// we were invaded, and lost
			text = p.Sprintf("%s %s has successfully invaded your planet %s, killing off all of your colonists", fleetOwner, fleet.Name, planet.Name)
		} else {
			// we were invaded, and lost
			text = p.Sprintf("%s %s tried to invade %s, but you were able to fend them off. You lost %d colonists in the invasion.", fleetOwner, fleet.Name, planet.Name, defendersKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMyPlanetInvaded, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	}
}

func (m *messageClient) techLevel(player *Player, field TechField, level int, nextField TechField) {
	text := fmt.Sprintf("Your scientists have completed research into Tech Level %d for %v.  They will continue their efforts in the %v field.", level, field, nextField)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text})
}

func (m *messageClient) playerDiscovered(player *Player, otherPlayer *Player) {
	text := fmt.Sprintf("You have discovered a new species, the %s. You are not alone in the universe!", otherPlayer.Race.PluralName)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlayerDiscovery, Text: text})
}

func (m *messageClient) permaform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := "decreased"
	if change > 0 {
		changeText = "increased"
	}
	text := fmt.Sprintf("Your race has permanently %s the %s on %s.", changeText, habType, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePermaform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) terraform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := "decreased"
	if change > 0 {
		changeText = "increased"
	}

	var newValueText string
	newValue := planet.Hab.Get(habType)
	switch habType {
	case Grav:
		newValueText = gravString(newValue)
	case Temp:
		newValueText = tempString(newValue)
	case Rad:
		newValueText = radString(newValue)
	}

	text := fmt.Sprintf("Your terraforming efforts on %s have %s the %s to %s", planet.Name, changeText, habType, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltTerraform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}
