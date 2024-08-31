package cs

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Target[T PlayerMessageTargetType | MapObjectType] struct {
	TargetType      T      `json:"targetType,omitempty"`
	TargetName      string `json:"targetName,omitempty"`
	TargetNum       int    `json:"targetNum,omitempty"`
	TargetPlayerNum int    `json:"targetPlayerNum,omitempty"`
}

// Throughout a turn various events will result in messages being sent to players.
// Messages have a type and a target (the target is focused in the UI when you click the Goto button)
// Messages also have a Spec that is used to store specific numbers for the UI to display on the message.
type PlayerMessage struct {
	Target[PlayerMessageTargetType]
	Type      PlayerMessageType `json:"type,omitempty"`
	Text      string            `json:"text,omitempty"`
	BattleNum int               `json:"battleNum,omitempty"`
	Spec      PlayerMessageSpec `json:"spec,omitempty"`
}

// The PlayerMessageSpec contains data specific to each message, like the amount of mines built
// of the field of research leveled up in.
type PlayerMessageSpec struct {
	// the thing being targeted by the message target, i.e. the planet for a fleet bombed a planet message
	Target[MapObjectType]
	Amount              int                             `json:"amount,omitempty"`
	Amount2             int                             `json:"amount2,omitempty"`
	PrevAmount          int                             `json:"prevAmount,omitempty"`
	SourcePlayerNum     int                             `json:"sourcePlayerNum,omitempty"`
	DestPlayerNum       int                             `json:"destPlayerNum,omitempty"`
	Name                string                          `json:"name,omitempty"`
	Cost                *Cost                           `json:"cost,omitempty"`
	Mineral             *Mineral                        `json:"mineral,omitempty"`
	Cargo               *Cargo                          `json:"cargo,omitempty"`
	QueueItemType       QueueItemType                   `json:"queueItemType,omitempty"`
	Field               TechField                       `json:"field,omitempty"`
	NextField           TechField                       `json:"nextField,omitempty"`
	TechGained          string                          `json:"techGained,omitempty"`
	LostTargetType      MapObjectType                   `json:"lostTargetType,omitempty"`
	Battle              BattleRecordStats               `json:"battle,omitempty"`
	Comet               *PlayerMessageSpecComet         `json:"comet,omitempty"`
	Bombing             *BombingResult                  `json:"bombing,omitempty"`
	MineralPacketDamage *MineralPacketDamage            `json:"mineralPacketDamage,omitempty"`
	MysteryTrader       *PlayerMessageSpecMysteryTrader `json:"mysteryTrader,omitempty"`
}

type PlayerMessageSpecComet struct {
	Size                          CometSize `json:"size,omitempty"`
	MineralsAdded                 Mineral   `json:"mineralsAdded,omitempty"`
	MineralConcentrationIncreased Mineral   `json:"mineralConcentrationIncreased,omitempty"`
	HabChanged                    Hab       `json:"habChanged,omitempty"`
	ColonistsKilled               int       `json:"colonistsKilled,omitempty"`
}

type PlayerMessageSpecMysteryTrader struct {
	MysteryTraderReward
	FleetNum int `json:"fleetNum" bson:"fleet_num"`
}

type PlayerMessageTargetType string

const (
	TargetNone          PlayerMessageTargetType = ""
	TargetPlanet        PlayerMessageTargetType = "Planet"
	TargetFleet         PlayerMessageTargetType = "Fleet"
	TargetWormhole      PlayerMessageTargetType = "Wormhole"
	TargetMineField     PlayerMessageTargetType = "MineField"
	TargetMysteryTrader PlayerMessageTargetType = "MysteryTrader"
	TargetMineralPacket PlayerMessageTargetType = "MineralPacket"
	TargetBattle        PlayerMessageTargetType = "Battle"
)

type PlayerMessageType int

const (
	PlayerMessageNone PlayerMessageType = iota
	PlayerMessageInfo
	PlayerMessageError
	PlayerMessagePlanetHomeworld
	PlayerMessagePlayerDiscovery
	PlayerMessagePlanetDiscovery
	PlayerMessagePlanetProductionQueueEmpty
	PlayerMessagePlanetProductionQueueComplete
	PlayerMessagePlanetBuiltMineralAlchemy
	PlayerMessagePlanetBuiltMine
	PlayerMessagePlanetBuiltFactory
	PlayerMessagePlanetBuiltDefense
	PlayerMessageFleetBuilt
	PlayerMessagePlanetBuiltStarbase
	PlayerMessagePlanetBuiltScanner
	PlayerMessagePlanetBuiltMineralPacket
	PlayerMessagePlanetBuiltTerraform
	PlayerMessageFleetOrdersComplete
	PlayerMessageFleetEngineFailure
	PlayerMessageFleetOutOfFuel
	PlayerMessageFleetGeneratedFuel
	PlayerMessageFleetScrapped
	PlayerMessageFleetMerged
	PlayerMessageFleetMergeInvalidNotFleet
	PlayerMessageFleetMergeInvalidUnowned
	PlayerMessageFleetPatrolTargeted
	PlayerMessageFleetRouteInvalidNotFriendlyPlanet
	PlayerMessageFleetRouteInvalidNotPlanet
	PlayerMessageFleetRouteInvalidNoRouteTarget
	PlayerMessageFleetTransportInvalid
	PlayerMessageFleetRoute
	PlayerMessageInvalid
	PlayerMessagePlanetColonized
	PlayerMessagePlayerGainTechLevel
	PlayerMessagePlanetBombed
	PlayerMessageUnused1
	PlayerMessageFleetBombedPlanet
	PlayerMessageUnused2
	PlayerMessagePlanetInvaded
	PlayerMessageFleetInvadedPlanet
	PlayerMessageBattle
	PlayerMessageFleetTransferredCargo
	PlayerMessageFleetSweptMines
	PlayerMessageFleetLaidMines
	PlayerMessageFleetMineFieldHit
	PlayerMessageFleetDumpedCargo
	PlayerMessageFleetStargateDamaged
	PlayerMessagePlanetPacketCaught
	PlayerMessagePlanetPacketDamage
	PlayerMessagePlanetPacketLanded
	PlayerMessageMineralPacketDiscovered
	PlayerMessageMineralPacketTargettingPlayerDiscovered
	PlayerMessagePlayerVictor
	PlayerMessageFleetReproduce
	PlayerMessagePlanetRandomMineralDeposit
	PlayerMessagePlanetPermaform
	PlayerMessagePlanetInstaform
	PlayerMessagePlanetPacketTerraform
	PlayerMessagePlanetPacketPermaform
	PlayerMessageFleetRemoteMined
	PlayerMessagePlayerTechGained
	PlayerMessageFleetTargetLost
	PlayerMessageFleetRadiatingEngineDieoff
	PlayerMessagePlanetDiedOff
	PlayerMessagePlanetEmptied
	PlayerMessagePlanetDiscoveryHabitable
	PlayerMessagePlanetDiscoveryTerraformable
	PlayerMessagePlanetDiscoveryUninhabitable
	PlayerMessagePlanetBuiltInvalidItem
	PlayerMessagePlanetBuiltInvalidMineralPacketNoMassDriver
	PlayerMessagePlanetBuiltInvalidMineralPacketNoTarget
	PlayerMessagePlanetPopulationDecreased
	PlayerMessagePlanetPopulationDecreasedOvercrowding
	PlayerMessagePlayerDead
	PlayerMessagePlayerNoPlanets
	PlayerMessagePlanetCometStrike
	PlayerMessagePlanetCometStrikeMyPlanet
	PlayerMessageFleetExceededSafeSpeed
	PlayerMessagePlanetBonusResearchArtifact
	PlayerMessageFleetTransferGiven
	PlayerMessageFleetTransferInvalidPlayer
	PlayerMessageFleetTransferInvalidColonists
	PlayerMessageFleetTransferInvalidGiveRefused
	PlayerMessageFleetTransferReceived
	PlayerMessageFleetTransferInvalidReceive
	PlayerMessageFleetTransferInvalidReceiveRefused
	PlayerMessagePlayerTechLevelGainedInvasion
	PlayerMessagePlayerTechLevelGainedScrapFleet
	PlayerMessagePlayerTechLevelGainedBattle
	PlayerMessageFleetDieoff
	PlayerMessageBattleAlly
	PlayerMessageBattleReports
	PlayerMessageMysteryTraderDiscovered
	PlayerMessageMysteryTraderChangedCourse
	PlayerMessageMysteryTraderAgain
	PlayerMessageMysteryTraderMetWithReward
	PlayerMessageMysteryTraderMetWithoutReward
	PlayerMessageMysteryTraderAlreadyRewarded
	PlayerMessagePlanetBuiltGenesisDevice
)

func newMessage(messageType PlayerMessageType) PlayerMessage {
	return PlayerMessage{Type: messageType}
}

// create a new message targeting a planet
func newPlanetMessage(messageType PlayerMessageType, target *Planet) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetName: target.Name, TargetNum: target.Num}}
}

// create a new message targeting a fleet
func newFleetMessage(messageType PlayerMessageType, target *Fleet) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetName: target.Name, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num}}
}

// create a new message targeting a minefield
func newMineFieldMessage(messageType PlayerMessageType, target *MineField) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: Target[PlayerMessageTargetType]{TargetType: TargetMineField, TargetName: target.Name, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num}}
}

// create a new message targeting a minefield
func newMineralPacketMessage(messageType PlayerMessageType, target *MineralPacket) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: Target[PlayerMessageTargetType]{TargetType: TargetMineralPacket, TargetName: target.Name, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num}}
}

// create a new message targeting a planet
func newMysteryTraderMessage(messageType PlayerMessageType, target *MysteryTrader) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: Target[PlayerMessageTargetType]{TargetType: TargetMysteryTrader, TargetNum: target.Num}}
}

// create a new message targeting a battle with the Name field as the location of the battle
func newBattleMessage(messageType PlayerMessageType, planet *Planet, battle *BattleRecord) PlayerMessage {
	planetNum := None
	targetType := TargetNone
	if planet != nil {
		planetNum = planet.Num
		targetType = TargetPlanet
	}

	return PlayerMessage{Type: messageType, Target: Target[PlayerMessageTargetType]{TargetType: targetType, TargetNum: planetNum}, BattleNum: battle.Num}
}

// use a spec in this message. spec.Name must be specified because the message details
// depend on it. Set it to the name of the PlayerMessage target
func (m PlayerMessage) withSpec(spec PlayerMessageSpec) PlayerMessage {
	m.Spec = spec
	return m
}

func (m PlayerMessage) withText(text string) PlayerMessage {
	m.Text = text
	return m
}

func (spec PlayerMessageSpec) withTargetFleet(fleet *Fleet) PlayerMessageSpec {
	spec.Target = Target[MapObjectType]{
		TargetType:      MapObjectTypeFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
		TargetName:      fleet.Name,
	}
	return spec
}

func (spec PlayerMessageSpec) withTargetPlanet(planet *Planet) PlayerMessageSpec {
	if planet == nil {
		return spec
	}
	spec.Target = Target[MapObjectType]{
		TargetType:      MapObjectTypePlanet,
		TargetPlayerNum: planet.PlayerNum,
		TargetNum:       planet.Num,
		TargetName:      planet.Name,
	}
	return spec
}

func (spec PlayerMessageSpec) withTargetMinefield(mineField *MineField) PlayerMessageSpec {
	spec.Target = Target[MapObjectType]{
		TargetType:      MapObjectTypeMineField,
		TargetPlayerNum: mineField.PlayerNum,
		TargetNum:       mineField.Num,
		TargetName:      mineField.Name,
	}
	return spec
}

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

func (m *messageClient) battle(player *Player, planet *Planet, battle *BattleRecord) {
	location := fmt.Sprintf("Space (%0f, %0f)", battle.Position.X, battle.Position.Y)
	if planet != nil {
		location = planet.Name
	}

	// create a new message targeting a battle
	player.Messages = append(player.Messages, newBattleMessage(PlayerMessageBattle, planet, battle).
		withSpec(PlayerMessageSpec{Name: location, Battle: battle.Stats}))
}

func (m *messageClient) battleAlly(player *Player, planet *Planet, battle *BattleRecord) {
	location := fmt.Sprintf("Space (%0f, %0f)", battle.Position.X, battle.Position.Y)
	if planet != nil {
		location = planet.Name
	}

	// create a new message targeting a battle
	player.Messages = append(player.Messages, newBattleMessage(PlayerMessageBattleAlly, planet, battle).
		withSpec(PlayerMessageSpec{Name: location, Battle: battle.Stats}))
}

func (mc *messageClient) battleReports(player *Player) {
	// Battle report messages are always the first message of the year (other than victory)
	player.Messages = append([]PlayerMessage{{Type: PlayerMessageBattleReports}}, player.Messages...)
}

/*
 * Fleet Messages
 */

func (m *messageClient) fleetBombedPlanet(player *Player, fleet *Fleet, planet *Planet, bombing BombingResult) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetBombedPlanet, fleet).
		withSpec(PlayerMessageSpec{Bombing: &bombing}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetBuilt(player *Player, planet *Planet, fleet *Fleet, numBuilt int) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetBuilt, fleet).
		withSpec(PlayerMessageSpec{Name: fleet.BaseName, Amount: numBuilt}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetColonizeNonPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to colonize, but is not currently orbiting a planet. The order has been canceled.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})
}

func (m *messageClient) fleetColonizeOwnedPlanet(player *Player, planet *Planet, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to colonize %s, but %s is already populated. The order has been canceled.", fleet.Name, planet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})

}

func (m *messageClient) fleetColonizeWithNoModule(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to colonize a planet without a colonization module. The order has been canceled.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})

}

func (m *messageClient) fleetColonizeWithNoColonists(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to colonize a planet, but has failed to bring along any colonists. The order has been canceled.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})
}

func (m *messageClient) fleetCompletedAssignedOrders(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has completed its assigned orders.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetOrdersComplete, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetDieOff(player *Player, fleet *Fleet, death int) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetDieoff, fleet).withSpec(
		PlayerMessageSpec{Amount: death},
	))
}

func (m *messageClient) fleetEngineFailure(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to engage its engines due to balky equipment. Engineers think they have the problem fixed for the time being.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetEngineFailure, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})
}

func (m *messageClient) fleetExceededSafeSpeed(player *Player, fleet *Fleet, explodedShips int) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetExceededSafeSpeed, fleet).withSpec(
		PlayerMessageSpec{Amount: explodedShips},
	))
}

func (m *messageClient) fleetGeneratedFuel(player *Player, fleet *Fleet, fuelGenerated int) {
	hasRamScoop := false
	for _, token := range fleet.Tokens {
		if token.design.Spec.Engine.FreeSpeed > 1 {
			hasRamScoop = true
			break
		}
	}
	if hasRamScoop {
		text := fmt.Sprintf("%s's ramscoops have produced %dmg of fuel from interstellar hydrogen.", fleet.Name, fuelGenerated)
	} else {
		text := fmt.Sprintf("%s's engines have produced %dmg of fuel from interstellar hydrogen.", fleet.Name, fuelGenerated)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetGeneratedFuel, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})
}

func (m *messageClient) fleetMerged(player *Player, fleet *Fleet, mergedInto *Fleet) {
	text := fmt.Sprintf("%s has been merged into %s.", fleet.Name, mergedInto.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetMerged, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: mergedInto.Num, TargetPlayerNum: mergedInto.PlayerNum}})
}

func (m *messageClient) fleetMineFieldHit(player *Player, fleet *Fleet, fleetPlayer *Player, mineField *MineField, damage int, shipsDestroyed int) {
	var text string
	if fleet.PlayerNum == player.Num {
		// it's our fleet, it must be someone else's minefield
		if fleet.Spec.TotalShips <= shipsDestroyed {
			text = fmt.Sprintf("%s has been annihilated in a %s mine field at %v.",
				fleet.Name, mineField.Type, mineField.Position)
		} else {
			text = fmt.Sprintf("%s has been stopped in a %s mine field at %v.",
				fleet.Name, mineField.Type, mineField.Position)
			if damage > 0 {
				if shipsDestroyed > 0 {
					text += fmt.Sprintf(" Your fleet has taken %d damage points and %d ships were destroyed.",
						damage, shipsDestroyed)
				} else {
					text += fmt.Sprintf(" Your fleet has taken %d damage points, but none of your ships were destroyed.",
						damage)
				}
			} else {
				text = fmt.Sprintf("%s has been stopped in a %s mine field at %v.",
					fleet.Name, mineField.Type, mineField.Position)
			}
		}
	} else {
		// it's not our fleet, it must be our minefield
		if fleet.Spec.TotalShips <= shipsDestroyed {
			text = fmt.Sprintf("%s %s has been annihilated in your %s mine field at %v.",
				fleetPlayer.Race.PluralName, fleet.Name, mineField.Type, mineField.Position)
		} else {
			text = fmt.Sprintf("%s %s has been stopped in your %s mine field at %v.",
				fleetPlayer.Race.PluralName, fleet.Name, mineField.Type, mineField.Position)
			if damage > 0 {
				if shipsDestroyed > 0 {
					text += fmt.Sprintf(" Your mines have inflicted %d damage points and destroyed %d ships.",
						damage, shipsDestroyed)
				} else {
					text += fmt.Sprintf(" Your mines have inflicted %d damage points, but you didn't manage to destroy any ships.",
						damage)
				}
			} else {
				text = fmt.Sprintf("%s has been stopped in your %s mine field at %v.",
					fleet.Name, mineField.Type, mineField.Position)
			}
		}
	}

	player.Messages = append(player.Messages, PlayerMessage{
		Type:   PlayerMessageFleetMineFieldHit,
		Text:   text,
		Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num},
	})

}

func (m *messageClient) fleetMineFieldSwept(player *Player, fleet *Fleet, mineField *MineField, numMinesSwept int) {
	var text string

	if fleet.PlayerNum == player.Num {
		text = fmt.Sprintf("%s has swept %d mines from a mine field at %v.", fleet.Name, numMinesSwept, mineField.Position)
	} else {
		text = fmt.Sprintf("Someone has swept %d mines from your mine field at %v.", numMinesSwept, mineField.Position)
	}

	// this will be removed if the mines are gone, so target the fleet
	if mineField.NumMines <= 10 {
		if fleet.PlayerNum == player.Num {
			player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetSweptMines, fleet).
				withSpec(PlayerMessageSpec{Amount: numMinesSwept}.withTargetMinefield(mineField)).
				withText(text))
		}
	} else {
		player.Messages = append(player.Messages, newMineFieldMessage(PlayerMessageFleetSweptMines, mineField).
			withSpec(PlayerMessageSpec{Amount: numMinesSwept}).
			withText(text))
	}

}

func (m *messageClient) fleetMinesLaidFailed(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to lay mines, but has no mine layers. The order has been canceled.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetMinesLaid(player *Player, fleet *Fleet, mineField *MineField, numMinesLaid int) {
	var text string
	if mineField.NumMines == numMinesLaid {
		text = fmt.Sprintf("%s has dispensed %d mines.", fleet.Name, numMinesLaid)
	} else {
		text = fmt.Sprintf("%s has increased a minefield by %d mines.", fleet.Name, numMinesLaid)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetLaidMines, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetOutOfFuel(player *Player, fleet *Fleet, warpSpeed int) {
	text := fmt.Sprintf("%s has run out of fuel. The fleet's speed has been decreased to Warp %d.", fleet.Name, warpSpeed)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetOutOfFuel, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})
}

func (m *messageClient) fleetPatrolTargeted(player *Player, fleet *Fleet, target *FleetIntel) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetPatrolTargeted, fleet).withSpec(
		PlayerMessageSpec{
			Name:   fleet.Name,
			Target: Target[MapObjectType]{TargetType: MapObjectTypeFleet, TargetName: target.Name, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num},
		},
	))
}

func (m *messageClient) fleetInvalidMergeNotFleet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to complete its merge orders as the waypoint destination wasn't a fleet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetMergeInvalidNotFleet, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetInvalidMergeNotOwned(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to complete its merge orders as the destination fleet wasn't one of yours.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetMergeInvalidUnowned, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetInvalidRouteNotPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s could not be routed because it is not orbiting a planet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetRouteInvalidNotPlanet, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetInvalidRouteNotFriendlyPlanet(player *Player, fleet *Fleet, planet *Planet) {
	text := fmt.Sprintf("%s could not be routed because you are not friends with the inhabitants of %s.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetRouteInvalidNotFriendlyPlanet, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetInvalidRouteNoRouteTarget(player *Player, fleet *Fleet, planet *Planet) {
	text := fmt.Sprintf("%s could not be routed at %s as the planet has no route set.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetRouteInvalidNoRouteTarget, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetRadiatingEngineDieoff(player *Player, fleet *Fleet, colonistsKilled int) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetRadiatingEngineDieoff, fleet).
		withSpec(PlayerMessageSpec{Amount: colonistsKilled}))
}

func (m *messageClient) fleetReproduce(player *Player, fleet *Fleet, colonistsGrown int, planet *Planet, over int) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetReproduce, fleet).
		withSpec(PlayerMessageSpec{Amount: colonistsGrown, Amount2: over}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRemoteMineNoMiners(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s has orders to remote mine %s, but the fleet doesn't have any remote mining modules. The order has been canceled.", fleet.Name, planet.Name)))
}

func (m *messageClient) fleetRemoteMineInhabited(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s has orders to remote mine %s, but the planet is already inhabited. The order has been canceled.", fleet.Name, planet.Name)))
}

func (m *messageClient) fleetRemoteMineDeepSpace(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s has orders to remote mine in deep space. The order has been canceled.", fleet.Name)))
}

func (m *messageClient) fleetRemoteMined(player *Player, fleet *Fleet, planet *Planet, mineral Mineral) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetRemoteMined, fleet).
		withSpec(PlayerMessageSpec{Mineral: &mineral}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRouted(player *Player, fleet *Fleet, planet *Planet, target string) {
	text := fmt.Sprintf("%s has been routed by the citizens of %s to %s.", fleet.Name, planet.Name, target)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetRoute, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num}})
}

func (m *messageClient) fleetScrapped(player *Player, fleet *Fleet, cost Cost, planet *Planet) {
	if planet != nil {
		player.Messages = append(player.Messages, newPlanetMessage(PlayerMessageFleetScrapped, planet).
			withSpec(PlayerMessageSpec{Cost: &cost}.withTargetFleet(fleet)))
	} else {
		player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetScrapped, fleet))
	}
}
func (m *messageClient) fleetStargateInvalidSource(player *Player, fleet *Fleet, wp0 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s, but no stargate exists there.", fleet.Name, wp0.TargetName)))
}

func (m *messageClient) fleetStargateInvalidSourceOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s, but could not because the starbase is not owned by you or your allies.", fleet.Name, wp0.TargetName)))
}

func (m *messageClient) fleetStargateInvalidDest(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but no stargate could be detected at the destination.", fleet.Name, wp0.TargetName, wp1.TargetName)))
}

func (m *messageClient) fleetStargateInvalidDestOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but could not because the destination starbase is not owned by you or your allies.", fleet.Name, wp0.TargetName, wp1.TargetName)))
}

func (m *messageClient) fleetStargateInvalidRange(player *Player, fleet *Fleet, wp0, wp1 Waypoint, totalDist float64) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but the distance of %.1f Ly. was far beyond the max range of the stargates.", fleet.Name, wp0.TargetName, wp1.TargetName, totalDist)))
}

func (m *messageClient) fleetStargateInvalidMass(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but your ships are far too massive for the gate's limits.", fleet.Name, wp0.TargetName, wp1.TargetName)))
}

func (m *messageClient) fleetStargateInvalidColonists(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but you are carrying colonists and can't drop them off as you don't own the planet.", fleet.Name, wp0.TargetName, wp1.TargetName)))
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

	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageInvalid, fleet).
		withText(text))
}

func (m *messageClient) fleetStargateDestroyed(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetStargateDamaged, fleet).
		withText(fmt.Sprintf("Heedless to the danger, %s attempted to use the stargate at %s to reach %s. The fleet never arrived. The distance or mass must have been too great.", fleet.Name, wp0.TargetName, wp1.TargetName)))
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
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing %d ships to the unforgiving void. Exceeding the capability of your stargates can be dangerous...", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	} else if totalShipsLost >= 10 && totalShipsLost <= 50 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing %d ships to the great unknown. Such disregard for stargates' capabilities is not recommended...", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	} else if totalShipsLost >= 50 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing an unbelievable %d ships to the cosmic ocean. The jump was far in excess of the capabilities of stargates involved...", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	}

	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetStargateDamaged, fleet).
		withText(text))
}

func (m *messageClient) fleetTransferGiven(player *Player, fleet *Fleet, targetPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetTransferGiven, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num, DestPlayerNum: targetPlayer.Num, Name: fleet.BaseName}))
}

func (m *messageClient) fleetTransferInvalidColonists(player *Player, fleet *Fleet, targetPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetTransferInvalidColonists, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num, DestPlayerNum: targetPlayer.Num}))
}

func (m *messageClient) fleetTransferInvalidGiveRefused(player *Player, fleet *Fleet, targetPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetTransferInvalidGiveRefused, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num, DestPlayerNum: targetPlayer.Num, Name: fleet.Name}))
}

func (m *messageClient) fleetTransferInvalidPlayer(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetTransferInvalidPlayer, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num}))
}

func (m *messageClient) fleetTransferInvalidReceiveRefused(player *Player, fleet *Fleet, givingPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetTransferInvalidReceiveRefused, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: givingPlayer.Num, DestPlayerNum: player.Num, Name: fleet.Name}))
}

func (m *messageClient) fleetTransferReceived(player *Player, fleet *Fleet, givingPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(PlayerMessageFleetTransferReceived, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: givingPlayer.Num, DestPlayerNum: player.Num, Name: fleet.BaseName}))
}

func (m *messageClient) fleetTransportedCargo(player *Player, fleet *Fleet, dest cargoHolder, cargoType CargoType, transferAmount int) {
	text := ""
	if cargoType == Colonists {
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has beamed %d colonists from %s.", fleet.Name, -transferAmount*100, dest.getMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has beamed %d colonists to %s.", fleet.Name, transferAmount*100, dest.getMapObject().Name)
		}
	} else {
		units := "kT"
		if cargoType == Fuel {
			units = "mg"
		}
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has loaded %d%s of %v from %s.", fleet.Name, -transferAmount, units, cargoType, dest.getMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has unloaded %d%s of %v to %s.", fleet.Name, transferAmount, units, cargoType, dest.getMapObject().Name)
		}
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetTransferredCargo, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})
}

func (m *messageClient) fleetTransportInvalid(player *Player, fleet *Fleet, dest cargoHolder, cargoType CargoType, transferAmount int) {
	text := fmt.Sprintf("%s attempted to load %dkT of %v from %s, but you do not own %s. The order has been canceled.", fleet.Name, -transferAmount, cargoType, dest.getMapObject().Name, dest.getMapObject().Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetTransportInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}})

}

func (m *messageClient) fleetTargetLost(player *Player, fleet *Fleet, targetName string, targetType MapObjectType) {
	text := ""
	if targetType == MapObjectTypeFleet {
		text = fmt.Sprintf("The fleet you were tracking with %s, %s, appears to have outrun the range of your scanners. Orders for your fleet have been changed to go to the last known location of that fleet.", fleet.Name, targetName)
	} else {
		text = fmt.Sprintf("The %s that you were tracking with %s, appears to have disappeared. Orders for your fleet have been changed to go to the last known location of the target.", targetName, fleet.Name)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetTargetLost, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum}, Spec: PlayerMessageSpec{LostTargetType: targetType}})
}

/*
 * MineralPacket Messages
 */

func (m *messageClient) planetBuiltMineralPacket(player *Player, planet *Planet, packet *MineralPacket, target string) {
	player.Messages = append(player.Messages, newMineralPacketMessage(PlayerMessagePlanetBuiltMineralPacket, packet).
		withSpec(PlayerMessageSpec{Amount: packet.Cargo.Total()}.withTargetPlanet(planet)))
}

func (m *messageClient) mineralPacketDiscovered(player *Player, packet *MineralPacket, target *Planet) {
	player.Messages = append(player.Messages, newMineralPacketMessage(PlayerMessageMineralPacketDiscovered, packet).withSpec(PlayerMessageSpec{}.withTargetPlanet(target)))
}

func (m *messageClient) mineralPacketDiscoveredTargettingPlayer(player *Player, packet *MineralPacket, target *Planet, damage MineralPacketDamage) {
	player.Messages = append(player.Messages, newMineralPacketMessage(PlayerMessageMineralPacketTargettingPlayerDiscovered, packet).withSpec(PlayerMessageSpec{MineralPacketDamage: &damage}.withTargetPlanet(target)))
}

/*
 * Planet Messages
 */

func (m *messageClient) planetHomeworld(player *Player, planet *Planet) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetHomeworld, planet))
}

func (m *messageClient) planetBombed(player *Player, planet *Planet, fleet *Fleet, bombing BombingResult) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBombed, planet).
		withSpec(PlayerMessageSpec{Bombing: &bombing}.withTargetFleet(fleet)))
}

func (m *messageClient) planetBonusResearchArtifact(player *Player, planet *Planet, amount int, field TechField) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBonusResearchArtifact, planet).withSpec(
		PlayerMessageSpec{Amount: amount, Field: field},
	))
}

func (m *messageClient) planetBuiltDefenses(player *Player, planet *Planet, numBuilt int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltDefense, planet).
		withSpec(PlayerMessageSpec{Amount: numBuilt}))
}

func (m *messageClient) planetBuiltFactories(player *Player, planet *Planet, numBuilt int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltFactory, planet).
		withSpec(PlayerMessageSpec{Amount: numBuilt}))
}

func (m *messageClient) planetBuiltMineralAlchemy(player *Player, planet *Planet, numBuilt int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltMineralAlchemy, planet).withSpec(PlayerMessageSpec{Amount: numBuilt}))
}

func (m *messageClient) planetBuiltMines(player *Player, planet *Planet, numBuilt int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltMine, planet).
		withSpec(PlayerMessageSpec{Amount: numBuilt}))
}

func (m *messageClient) planetBuiltScanner(player *Player, planet *Planet, scanner string) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltScanner, planet).
		withSpec(PlayerMessageSpec{Name: scanner}))
}

func (m *messageClient) planetBuiltGenesisDevice(player *Player, planet *Planet) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltGenesisDevice, planet))
}

func (m *messageClient) planetBuiltStarbase(player *Player, planet *Planet, fleet *Fleet) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltStarbase, planet).
		withSpec(PlayerMessageSpec{Name: fleet.BaseName}))
}

func (m *messageClient) planetColonized(player *Player, planet *Planet) {
	text := fmt.Sprintf("Your colonists are now in control of %s.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetColonized, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetComet(player *Player, planet *Planet, size CometSize, mineralsAdded Mineral, mineralConcentrationIncreased Mineral, habChanged Hab, colonistsKilled int) {
	if planet.PlayerNum == player.Num {
		player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetCometStrikeMyPlanet, planet).withSpec(
			PlayerMessageSpec{
				Comet: &PlayerMessageSpecComet{
					Size:                          size,
					MineralsAdded:                 mineralsAdded,
					MineralConcentrationIncreased: mineralConcentrationIncreased,
					HabChanged:                    habChanged,
					ColonistsKilled:               colonistsKilled,
				},
			},
		))
	} else {
		player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetCometStrike, planet).withSpec(
			PlayerMessageSpec{
				Comet: &PlayerMessageSpecComet{
					Size: size,
				},
			},
		))
	}
}

func (m *messageClient) planetDiedOff(player *Player, planet *Planet) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetDiedOff, planet))
}

func (m *messageClient) planetDiscovered(player *Player, planet *Planet) {
	messageType := PlayerMessagePlanetDiscovery
	hab := player.Race.GetPlanetHabitability(planet.Hab)

	terraformer := NewTerraformer()
	terraformAmount := terraformer.getTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	habTerraformed := player.Race.GetPlanetHabitability(planet.Hab.Add(terraformAmount))

	if hab >= 0 {
		messageType = PlayerMessagePlanetDiscoveryHabitable
	} else if habTerraformed > 0 {
		messageType = PlayerMessagePlanetDiscoveryTerraformable
	} else {
		messageType = PlayerMessagePlanetDiscoveryUninhabitable
	}

	player.Messages = append(player.Messages, newPlanetMessage(messageType, planet))
}

func (m *messageClient) planetEmptied(player *Player, planet *Planet) {
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetEmptied, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetInstaform(player *Player, planet *Planet, terraformAmount Hab) {
	text := fmt.Sprintf("Your race has instantly terraformed %s up to optimal conditions. Its value is now %d%.", planet.Name, planet.value)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetInstaform, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetInvaded(player *Player, planet *Planet, fleet *Fleet, planetOwner string, fleetOwner string, attackersKilled int, defendersKilled int, successful bool) {
	var text string

	// use this formatter to get commas on the text
	p := message.NewPrinter(language.English)
	if player.Num == fleet.PlayerNum {
		if successful {
			// we invaded and won
			text = p.Sprintf("Your troops beaming down from %s have successfully wrested %s from %s control, killing off all their colonists with only %d causalties.", fleet.Name, planet.Name, planetOwner, attackersKilled)
		} else {
			// we invaded and lost
			text = p.Sprintf("Your troops beaming down from %s tried to invade %s, but all of them were massacred by the %s. Your valiant fighters managed to kill %d of their colonists in return.", fleet.Name, planet.Name, planetOwner, defendersKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvadedPlanet, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
	} else {
		if successful {
			// we were invaded, and lost
			text = p.Sprintf("%s %s has successfully invaded %s and wrested it from your control. Your colonists managed to defeat %d of their invaders before being overrun.", fleetOwner, fleet.Name, planet.Name, attackersKilled)
		} else {
			// we were invaded, and lost
			text = p.Sprintf("%s %s tried to invade %s, but your troops were able to fend them off. You lost %d colonists in the process.", fleetOwner, fleet.Name, planet.Name, defendersKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetInvaded, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
	}
}

func (m *messageClient) planetInvadeEmpty(player *Player, planet *Planet, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to invade %s, but the planet is uninhabited. The order has been canceled.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}
func (m *messageClient) planetInvadeStarbase(player *Player, planet *Planet, fleet *Fleet) {
	text := fmt.Sprintf("%s has orders to invade %s, but the planet is protected by a starbase. The order has been canceled.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPacketArrived(player *Player, planet *Planet, packet *MineralPacket) {
	text := fmt.Sprintf("Your mineral packet containing %dkT of minerals has arrived at %s.", packet.Cargo.Total(), planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetPacketLanded, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPacketCaught(player *Player, planet *Planet, packet *MineralPacket) {
	text := fmt.Sprintf("Your mass accelerator at %s has successfully captured a packet containing %dkT of minerals.", planet.Name, packet.Cargo.Total())
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetPacketCaught, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPacketDamage(player *Player, planet *Planet, packet *MineralPacket, colonistsKilled, defensesDestroyed int) {
	var text string
	if planet.Spec.HasStarbase && planet.Starbase.Spec.HasMassDriver {
		if defensesDestroyed == 0 {
			text = fmt.Sprintf("Your mass accelerator at %s was partially successful at capturing a %dkT mineral packet. Unable to completely slow the packet, %d of your colonists were killed in the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled)
		} else {
			text = fmt.Sprintf("Your mass accelerator at %s was partially successful at capturing a %dkT mineral packet. Unfortunately, %d of your colonists and %d of your defenses were destroyed in the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled, defensesDestroyed)
		}
	} else {
		if planet.population() == 0 {
			text = fmt.Sprintf("%s was annihilated by a mineral packet. All of your colonists were killed.", planet.Name)
		} else if defensesDestroyed == 0 {
			text = fmt.Sprintf("%s was bombarded with a %dkT mineral packet. %d of your colonists were killed in the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled)
		} else {
			text = fmt.Sprintf("%s was bombarded with a %dkT mineral packet. %d of your colonists and %d of your defenses were destroyed in the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled, defensesDestroyed)
		}
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetPacketDamage, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPacketPermaform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := ""
	if change > 0 {
		changeText = "increased"
	} else {
		changeText = "decreased"
	}
	newValueText := ""
	newValue := planet.Hab.Get(habType)
	switch habType {
	case Grav:
		newValueText = gravString(newValue)
	case Temp:
		newValueText = tempString(newValue)
	case Rad:
		newValueText = radString(newValue)
	}
	text := fmt.Sprintf("Your mineral packet has permanently %s the %s on %s to %s.", changeText, habType, planet.Name, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetPacketPermaform, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPacketTerraform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := ""
	if change > 0 {
		changeText = "increased"
	} else {
		changeText = "decreased"
	}

	newValueText := ""
	newValue := planet.Hab.Get(habType)
	switch habType {
	case Grav:
		newValueText = gravString(newValue)
	case Temp:
		newValueText = tempString(newValue)
	case Rad:
		newValueText = radString(newValue)
	}

	text := fmt.Sprintf("Your mineral packet hitting %s has %s its %s to %s.", planet.Name, changeText, habType, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetPacketTerraform, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPermaform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := "decreased"
	if change > 0 {
		changeText = "increased"
	}
	newValueText := ""
	newValue := planet.Hab.Get(habType)
	switch habType {
	case Grav:
		newValueText = gravString(newValue)
	case Temp:
		newValueText = tempString(newValue)
	case Rad:
		newValueText = radString(newValue)
	}
	text := fmt.Sprintf("Your colonists have permanently %s the %s on %s to %s.", changeText, habType, planet.Name, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetPermaform, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

func (m *messageClient) planetPopulationDecreased(player *Player, planet *Planet, prevAmount int, amount int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPopulationDecreased, planet).
		withSpec(PlayerMessageSpec{PrevAmount: prevAmount, Amount: amount}))
}

func (m *messageClient) planetPopulationDecreasedOvercrowding(player *Player, planet *Planet, amount int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPopulationDecreasedOvercrowding, planet).
		withSpec(PlayerMessageSpec{Amount: amount}))
}

func (m *messageClient) planetTerraform(player *Player, planet *Planet, habType HabType, change int) {
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

	text := fmt.Sprintf("Your terraforming efforts on %s have %s the %s to %s.", planet.Name, changeText, habType, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetBuiltTerraform, Text: text, Target: Target[PlayerMessageTargetType]{TargetType: TargetPlanet, TargetNum: planet.Num}})
}

/*
 * Player Messages
 */

func (m *messageClient) playerDiscovered(player *Player, otherPlayer *Player) {
	text := fmt.Sprintf("You have discovered a new species, the %s. You are not alone in this universe!", otherPlayer.Race.PluralName)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlayerDiscovery, Text: text})
}

func (m *messageClient) playerGainTechLevel(player *Player, field TechField, level int, nextField TechField) {
	text := fmt.Sprintf("Your scientists have completed research into Tech Level %d for %v. They will continue their efforts in the %v field.", level, field, nextField)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlayerGainTechLevel, Text: text, Spec: PlayerMessageSpec{Field: field, NextField: nextField, Amount: level}})
}

func (m *messageClient) playerTechGained(player *Player, field TechField, tech *Tech) {
	var text string
	switch tech.Category {
	case TechCategoryShipHull:
		fallthrough
	case TechCategoryStarbaseHull:
		text = fmt.Sprintf("Your recent breakthrough in %v has also given you the %s hull type. To build ships with this design, go to Commands -> Ship Designer and select Create New Design.", field, tech.Name)
	case TechCategoryPlanetaryDefense:
		text = fmt.Sprintf("Your recent breakthrough in %v has also taught you how to build %s defenses. All existing planetary defenses have been upgraded to the new technology.", field, tech.Name)
	case TechCategoryPlanetaryScanner:
		text = fmt.Sprintf("Your recent breakthrough in %v has also taught you how to build the %s scanner. All existing planetary scanners have been upgraded to the new technology.", field, tech.Name)
	default:
		text = fmt.Sprintf("Your recent breakthrough in %v has also given you the %s benefit.", field, tech.Name)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlayerTechGained, Text: text, Spec: PlayerMessageSpec{Field: field, TechGained: tech.Name}})
}

func (m *messageClient) playerTechGainedBattle(player *Player, planet *Planet, record *BattleRecord, field TechField) {
	player.Messages = append(player.Messages, newBattleMessage(PlayerMessagePlayerTechLevelGainedBattle, planet, record).
		withSpec(PlayerMessageSpec{Field: field}))
}

func (m *messageClient) playerTechGainedInvasion(player *Player, planet *Planet, field TechField) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlayerTechLevelGainedInvasion, planet).
		withSpec(PlayerMessageSpec{Field: field}))
}

func (m *messageClient) playerTechGainedScrappedFleet(player *Player, planet *Planet, fleetName string, field TechField) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlayerTechLevelGainedScrapFleet, planet).
		withSpec(PlayerMessageSpec{Field: field, Name: fleetName}))
}

// tell a player they are dead. This always appears as the first message
func (mc *messageClient) playerDead(player, deadPlayer *Player) {
	player.Messages = append([]PlayerMessage{newMessage(PlayerMessagePlayerDead).withSpec(PlayerMessageSpec{Target: Target[MapObjectType]{TargetPlayerNum: deadPlayer.Num}})}, player.Messages...)
}

// tell a player they have no planets but still have colonists. This always appears as the first message
func (mc *messageClient) playerNoPlanets(player *Player, numColonists int) {
	player.Messages = append([]PlayerMessage{newMessage(PlayerMessagePlayerNoPlanets).withSpec(PlayerMessageSpec{Amount: numColonists})}, player.Messages...)
}

func (mc *messageClient) playerVictory(player *Player, victor *Player) {
	var text string
	if player.Num == victor.Num {
		text = "You have been declared the winner of this grand game. You may continue to play though, if you wish to really rub your nose in everyone else's face."
	} else {
		text = fmt.Sprintf("The %s have been declared the winner of this game. You are advised to accept their supremacy, though you may continue the fight regardless.", victor.Race.Name)
	}
	// Victory messages are always the first message of the year
	player.Messages = append([]PlayerMessage{{Type: PlayerMessagePlayerVictor, Text: text}}, player.Messages...)
}
