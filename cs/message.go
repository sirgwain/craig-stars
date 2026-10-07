package cs

import (
	"fmt"
)

type Target[T PlayerMessageTargetType | MapObjectType] struct {
	TargetPosition  Vector `json:"targetPosition"`
	TargetType      T      `json:"targetType,omitempty"`
	TargetName      string `json:"targetName,omitempty"`
	TargetNum       int    `json:"targetNum,omitempty"`
	TargetPlayerNum int    `json:"targetPlayerNum,omitempty"`
}

type MapObjectTarget = Target[MapObjectType]
type PlayerMessageTarget = Target[PlayerMessageTargetType]

func (t Target[T]) PrettyString() string {
	return fmt.Sprintf("Target: %s Type: %s Player: %d Num: %d", t.TargetName, t.TargetType, t.TargetPlayerNum, t.TargetNum)
}

func (t Target[T]) Targeting(mo MapObject) bool {
	return string(mo.Type) == string(t.TargetType) && mo.PlayerNum == t.TargetPlayerNum && mo.Num == t.TargetNum
}

// Throughout a turn various events will result in messages being sent to players.
// Messages have a type and a target (the target is focused in the UI when you click the Goto button)
// Messages also have a Spec that is used to store specific numbers for the UI to display on the message.
type PlayerMessage struct {
	Target[PlayerMessageTargetType]
	Type      PlayerMessageType `json:"type"`
	Text      string            `json:"text,omitempty"` // Legacy rendered text from saved games.
	BattleNum int               `json:"battleNum,omitempty"`
	Spec      PlayerMessageSpec `json:"spec"`
}

// The PlayerMessageSpec contains data specific to each message, like the amount of mines built
// or the field of research leveled up in.
type PlayerMessageSpec struct {
	// the thing being targeted by the message target, i.e. the planet for a fleet bombed a planet message
	MapObjectTarget
	Amount              int                             `json:"amount,omitempty"`
	Amount2             int                             `json:"amount2,omitempty"`
	Battle              *BattleRecordStats              `json:"battle,omitempty"`
	Bombing             *BombingResult                  `json:"bombing,omitempty"`
	Cargo               *Cargo                          `json:"cargo,omitempty"`
	CargoTransfer       *PlayerMessageSpecCargoTransfer `json:"cargoTransfer,omitempty"`
	Comet               *PlayerMessageSpecComet         `json:"comet,omitempty"`
	Cost                *Cost                           `json:"cost,omitempty"`
	DestPlayerNum       int                             `json:"destPlayerNum,omitempty"`
	Distance            float64                         `json:"distance,omitempty"`
	Error               string                          `json:"error,omitempty"`
	Field               TechField                       `json:"field,omitempty"`
	HabType             TerraformHabType                `json:"habType,omitempty"`
	HasMassDriver       bool                            `json:"hasMassDriver,omitempty"`
	Invasion            *PlayerMessageSpecInvasion      `json:"invasion,omitempty"`
	LostTargetType      MapObjectType                   `json:"lostTargetType,omitempty"`
	MinefieldDamage     *MinefieldDamage                `json:"minefieldDamage,omitempty"`
	Mineral             *Mineral                        `json:"mineral,omitempty"`
	MineralPacketDamage *MineralPacketDamage            `json:"mineralPacketDamage,omitempty"`
	MysteryTrader       *PlayerMessageSpecMysteryTrader `json:"mysteryTrader,omitempty"`
	Name                string                          `json:"name,omitempty"`
	NextField           TechField                       `json:"nextField,omitempty"`
	PlanetEmptied       bool                            `json:"planetEmptied,omitempty"`
	PrevAmount          int                             `json:"prevAmount,omitempty"`
	QueueItemType       QueueItemType                   `json:"queueItemType,omitempty"`
	RouteTarget         *MapObjectTarget                `json:"routeTarget,omitempty"`
	SourcePlayerNum     int                             `json:"sourcePlayerNum,omitempty"`
	TechGained          string                          `json:"techGained,omitempty"`
	TerraformAmount     Hab                             `json:"terraformAmount,omitempty"`
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

type PlayerMessageSpecInvasion struct {
	FleetName         string `json:"fleetName,omitempty"`
	AttackerPlayerNum int    `json:"attackerPlayerNum"`
	DefenderPlayerNum int    `json:"defenderPlayerNum"`
	AttackersKilled   int    `json:"attackersKilled"`
	DefendersKilled   int    `json:"defendersKilled"`
	Successful        bool   `json:"successful"`
}

type PlayerMessageSpecCargoTransfer struct {
	CargoType  CargoType           `json:"cargoType"`
	Transfered int                 `json:"transfered"`
	Wanted     int                 `json:"wanted"`
	Status     CargoTransferStatus `json:"status"`
}

type PlayerMessageTargetType string

const (
	TargetNone          PlayerMessageTargetType = ""
	TargetPlanet        PlayerMessageTargetType = "Planet"
	TargetFleet         PlayerMessageTargetType = "Fleet"
	TargetWormhole      PlayerMessageTargetType = "Wormhole"
	TargetMinefield     PlayerMessageTargetType = "Minefield"
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
	PlayerMessageFleetMinefieldSweptMines
	PlayerMessageFleetLaidMines
	PlayerMessageFleetMinefieldHit
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
	PlayerMessagePlayerAcquirablePartGainedScrapFleet
	PlayerMessagePlayerAcquirablePartGainedBattle
	PlayerMessageFleetByHandTransferIncomplete
	PlayerMessagePlanetBuiltBeyondMaximum
	PlayerMessagePlanetBuiltInvalidShip
	PlayerMessageFleetColonizeInvalidNotPlanet
	PlayerMessageFleetColonizeInvalidOwnedPlanet
	PlayerMessageFleetColonizeInvalidNoModule
	PlayerMessageFleetColonizeInvalidNoColonists
	PlayerMessageFleetLayMinesInvalidNoMineLayers
	PlayerMessageFleetRemoteMineInvalidNoMiners
	PlayerMessageFleetRemoteMineInvalidInhabited
	PlayerMessageFleetRemoteMineInvalidDeepSpace
	PlayerMessageFleetStargateInvalidSource
	PlayerMessageFleetStargateInvalidSourceOwner
	PlayerMessageFleetStargateInvalidDest
	PlayerMessageFleetStargateInvalidDestOwner
	PlayerMessageFleetStargateInvalidRange
	PlayerMessageFleetStargateInvalidMass
	PlayerMessageFleetStargateInvalidColonists
	PlayerMessagePlanetInvadeInvalidEmpty
	PlayerMessagePlanetInvadeInvalidStarbase
	PlayerMessageFleetStargateDestroyed
	PlayerMessageFleetEngineStrainDestroyed
	PlayerMessagePlanetRemoteTerraform
)

func newMessage(messageType PlayerMessageType) PlayerMessage {
	return PlayerMessage{Type: messageType}
}

func newPlayerMessage(messageType PlayerMessageType, targetPlayerNum int) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: PlayerMessageTarget{TargetPlayerNum: targetPlayerNum}}
}

// create a new message targeting a planet
func newPlanetMessage(messageType PlayerMessageType, target *Planet) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: PlayerMessageTarget{TargetType: TargetPlanet, TargetName: target.Name, TargetNum: target.Num}}
}

// create a new message targeting a fleet
func newFleetMessage(player *Player, messageType PlayerMessageType, target *Fleet) PlayerMessage {
	// don't expose the actual fleet name to other players
	targetName := target.Name
	if player.Num != target.PlayerNum {
		targetName = fmt.Sprintf("%s #%d", target.Tokens[0].design.Hull, target.Num)
	}
	return PlayerMessage{Type: messageType, Target: PlayerMessageTarget{TargetType: TargetFleet, TargetName: targetName, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num}}
}

// create a new message targeting a minefield
func newMineralPacketMessage(messageType PlayerMessageType, target *MineralPacket) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: PlayerMessageTarget{TargetType: TargetMineralPacket, TargetName: target.Name, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num}}
}

// create a new message targeting a planet
func newMysteryTraderMessage(messageType PlayerMessageType, target *MysteryTrader) PlayerMessage {
	return PlayerMessage{Type: messageType, Target: PlayerMessageTarget{TargetType: TargetMysteryTrader, TargetNum: target.Num}}
}

// create a new message targeting a battle with the Name field as the location of the battle
func newBattleMessage(messageType PlayerMessageType, planet *Planet, battle *BattleRecord) PlayerMessage {
	planetNum := None
	targetType := TargetNone
	if planet != nil {
		planetNum = planet.Num
		targetType = TargetPlanet
	}

	return PlayerMessage{Type: messageType, Target: PlayerMessageTarget{TargetType: targetType, TargetNum: planetNum}, BattleNum: battle.Num}
}

// Use a spec for event data. Amount and Amount2 have message-specific meanings;
// the primary object's name is stored in the message target.
func (m PlayerMessage) withSpec(spec PlayerMessageSpec) PlayerMessage {
	m.Spec = spec
	return m
}

func (spec PlayerMessageSpec) withTargetFleet(fleet *Fleet) PlayerMessageSpec {
	spec.MapObjectTarget = MapObjectTarget{
		TargetType:      MapObjectTypeFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
		TargetName:      fleet.Name,
		TargetPosition:  fleet.Position,
	}
	return spec
}

func (spec PlayerMessageSpec) withTargetPlanet(planet *Planet) PlayerMessageSpec {
	if planet == nil {
		return spec
	}
	spec.MapObjectTarget = MapObjectTarget{
		TargetType:      MapObjectTypePlanet,
		TargetPlayerNum: planet.PlayerNum,
		TargetNum:       planet.Num,
		TargetName:      planet.Name,
		TargetPosition:  planet.Position,
	}
	return spec
}

func (spec PlayerMessageSpec) withTargetMinefield(minefield *Minefield) PlayerMessageSpec {
	spec.MapObjectTarget = MapObjectTarget{
		TargetType:      MapObjectTypeMinefield,
		TargetPlayerNum: minefield.PlayerNum,
		TargetNum:       minefield.Num,
		TargetName:      minefield.Name,
		TargetPosition:  minefield.Position,
	}
	return spec
}

// Snapshot waypoint names and positions, including when the destination no longer exists.
func (spec PlayerMessageSpec) withStargateSource(wp Waypoint) PlayerMessageSpec {
	spec.MapObjectTarget = wp.MapObjectTarget
	spec.MapObjectTarget.TargetPosition = wp.Position
	return spec
}

func (spec PlayerMessageSpec) withStargateTargets(source, dest Waypoint) PlayerMessageSpec {
	spec = spec.withStargateSource(source)
	target := dest.MapObjectTarget
	target.TargetPosition = dest.Position
	spec.RouteTarget = &target
	return spec
}

type messageClient struct {
}

var messager = messageClient{}

func (m *messageClient) error(player *Player, err error) {
	player.Messages = append(player.Messages, newMessage(PlayerMessageError).
		withSpec(PlayerMessageSpec{Error: err.Error()}))
}

func (m *messageClient) battle(player *Player, planet *Planet, battle *BattleRecord) {
	location := fmt.Sprintf("Space (%d, %d)", battle.Position.X, battle.Position.Y)
	if planet != nil {
		location = planet.Name
	}

	// create a new message targeting a battle
	player.Messages = append(player.Messages, newBattleMessage(PlayerMessageBattle, planet, battle).
		withSpec(PlayerMessageSpec{Name: location, Battle: &battle.Stats}))
}

func (m *messageClient) battleAlly(player *Player, planet *Planet, battle *BattleRecord) {
	location := fmt.Sprintf("Space (%d, %d)", battle.Position.X, battle.Position.Y)
	if planet != nil {
		location = planet.Name
	}

	// create a new message targeting a battle
	player.Messages = append(player.Messages, newBattleMessage(PlayerMessageBattleAlly, planet, battle).
		withSpec(PlayerMessageSpec{Name: location, Battle: &battle.Stats}))
}

func (mc *messageClient) battleReports(player *Player) {
	// Battle report messages are always the first message of the year (other than victory)
	player.Messages = append([]PlayerMessage{{Type: PlayerMessageBattleReports}}, player.Messages...)
}

/*
 * Fleet Messages
 */

func (m *messageClient) fleetBombedPlanet(player *Player, fleet *Fleet, planet *Planet, bombing BombingResult) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetBombedPlanet, fleet).
		withSpec(PlayerMessageSpec{Bombing: &bombing}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetBuilt(player *Player, planet *Planet, fleet *Fleet, numBuilt int, routeTarget *MapObject) {
	var target MapObjectTarget
	if routeTarget != nil {
		target = MapObjectTarget{
			TargetPosition:  routeTarget.Position,
			TargetType:      routeTarget.Type,
			TargetNum:       routeTarget.Num,
			TargetPlayerNum: routeTarget.PlayerNum,
			TargetName:      routeTarget.Name,
		}
	}
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetBuilt, fleet).
		withSpec(PlayerMessageSpec{Name: fleet.BaseName, Amount: numBuilt, RouteTarget: &target}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetColonizeNonPlanet(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetColonizeInvalidNotPlanet, fleet))
}

func (m *messageClient) fleetColonizeOwnedPlanet(player *Player, planet *Planet, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetColonizeInvalidOwnedPlanet, fleet).
		withSpec(PlayerMessageSpec{}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetColonizeWithNoModule(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetColonizeInvalidNoModule, fleet))
}

func (m *messageClient) fleetColonizeWithNoColonists(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetColonizeInvalidNoColonists, fleet))
}

func (m *messageClient) fleetCompletedAssignedOrders(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetOrdersComplete, fleet))
}

func (m *messageClient) fleetDieOff(player *Player, fleet *Fleet, death int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetDieoff, fleet).withSpec(
		PlayerMessageSpec{Amount: death},
	))
}

func (m *messageClient) fleetEngineFailure(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetEngineFailure, fleet))
}

func (m *messageClient) fleetExceededSafeSpeed(player *Player, fleet *Fleet, explodedShips int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetExceededSafeSpeed, fleet).withSpec(
		PlayerMessageSpec{Amount: explodedShips},
	))
}

func (m *messageClient) fleetEngineStrainDestroyed(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetEngineStrainDestroyed, fleet))
}

func (m *messageClient) fleetGeneratedFuel(player *Player, fleet *Fleet, fuelGenerated int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetGeneratedFuel, fleet).withSpec(
		PlayerMessageSpec{Amount: fuelGenerated},
	))
}

func (m *messageClient) fleetMerged(player *Player, fleet *Fleet, mergedInto *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetMerged, mergedInto).
		withSpec(PlayerMessageSpec{Name: fleet.Name}))
}

func (m *messageClient) fleetMinefieldHit(player *Player, fleet *Fleet, minefield *Minefield, minefieldDamage MinefieldDamage) {
	player.Messages = append(player.Messages, newFleetMessage(player,
		PlayerMessageFleetMinefieldHit, fleet).withSpec(PlayerMessageSpec{MinefieldDamage: &minefieldDamage}.withTargetMinefield(minefield)))
}

func (m *messageClient) fleetMinefieldSwept(player *Player, fleet *Fleet, minefield *Minefield, numMinesSwept int) {
	player.Messages = append(player.Messages, newFleetMessage(player,
		PlayerMessageFleetMinefieldSweptMines, fleet).withSpec(PlayerMessageSpec{Amount: numMinesSwept}.withTargetMinefield(minefield)))
}

func (m *messageClient) fleetMinesLaidFailed(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetLayMinesInvalidNoMineLayers, fleet))
}

func (m *messageClient) fleetMinesLaid(player *Player, fleet *Fleet, minefield *Minefield, numMinesLaid int) {
	player.Messages = append(player.Messages, newFleetMessage(player,
		PlayerMessageFleetLaidMines, fleet).withSpec(PlayerMessageSpec{Amount: numMinesLaid}.withTargetMinefield(minefield)))
}

// Amount is the resulting warp speed.
func (m *messageClient) fleetOutOfFuel(player *Player, fleet *Fleet, warpSpeed int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetOutOfFuel, fleet).
		withSpec(PlayerMessageSpec{Amount: warpSpeed}))
}

func (m *messageClient) fleetPatrolTargeted(player *Player, fleet *Fleet, target *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetPatrolTargeted, fleet).withSpec(
		PlayerMessageSpec{
			Name:            fleet.Name,
			MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypeFleet, TargetName: target.Name, TargetPlayerNum: target.PlayerNum, TargetNum: target.Num},
		},
	))
}

func (m *messageClient) fleetInvalidMergeNotFleet(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetMergeInvalidNotFleet, fleet))
}

func (m *messageClient) fleetInvalidMergeNotOwned(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetMergeInvalidUnowned, fleet))
}

func (m *messageClient) fleetInvalidRouteNotPlanet(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRouteInvalidNotPlanet, fleet))
}

func (m *messageClient) fleetInvalidRouteNotFriendlyPlanet(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRouteInvalidNotFriendlyPlanet, fleet).
		withSpec(PlayerMessageSpec{}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetInvalidRouteNoRouteTarget(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRouteInvalidNoRouteTarget, fleet).
		withSpec(PlayerMessageSpec{}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRadiatingEngineDieoff(player *Player, fleet *Fleet, colonistsKilled int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRadiatingEngineDieoff, fleet).
		withSpec(PlayerMessageSpec{Amount: colonistsKilled}))
}

func (m *messageClient) fleetReproduce(player *Player, fleet *Fleet, colonistsGrown int, planet *Planet, over int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetReproduce, fleet).
		withSpec(PlayerMessageSpec{Amount: colonistsGrown, Amount2: over}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRemoteMineNoMiners(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRemoteMineInvalidNoMiners, fleet).
		withSpec(PlayerMessageSpec{}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRemoteMineInhabited(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRemoteMineInvalidInhabited, fleet).
		withSpec(PlayerMessageSpec{}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRemoteMineDeepSpace(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRemoteMineInvalidDeepSpace, fleet))
}

func (m *messageClient) fleetRemoteMined(player *Player, fleet *Fleet, planet *Planet, mineral Mineral) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRemoteMined, fleet).
		withSpec(PlayerMessageSpec{Mineral: &mineral}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetRouted(player *Player, fleet *Fleet, planet *Planet, target *MapObject) {
	routeTarget := target.ToTarget()
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetRoute, fleet).
		withSpec(PlayerMessageSpec{RouteTarget: &routeTarget}.withTargetPlanet(planet)))
}

func (m *messageClient) fleetScrapped(player *Player, fleet *Fleet, cost Cost, planet *Planet) {
	if planet != nil {
		player.Messages = append(player.Messages, newPlanetMessage(PlayerMessageFleetScrapped, planet).
			withSpec(PlayerMessageSpec{Cost: &cost, Cargo: &fleet.Cargo}.withTargetFleet(fleet)))
	} else {
		player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetScrapped, fleet))
	}
}
func (m *messageClient) fleetStargateInvalidSource(player *Player, fleet *Fleet, wp0 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidSource, fleet).
		withSpec(PlayerMessageSpec{}.withStargateSource(wp0)))
}

func (m *messageClient) fleetStargateInvalidSourceOwner(player *Player, fleet *Fleet, wp0 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidSourceOwner, fleet).
		withSpec(PlayerMessageSpec{}.withStargateSource(wp0)))
}

func (m *messageClient) fleetStargateInvalidDest(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidDest, fleet).
		withSpec(PlayerMessageSpec{}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetStargateInvalidDestOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidDestOwner, fleet).
		withSpec(PlayerMessageSpec{}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetStargateInvalidRange(player *Player, fleet *Fleet, wp0, wp1 Waypoint, totalDist float64) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidRange, fleet).
		withSpec(PlayerMessageSpec{Distance: totalDist}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetStargateInvalidMass(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidMass, fleet).
		withSpec(PlayerMessageSpec{}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetStargateInvalidColonists(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateInvalidColonists, fleet).
		withSpec(PlayerMessageSpec{}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetStargateDumpedCargo(player *Player, fleet *Fleet, wp0, wp1 Waypoint, cargo Cargo) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetDumpedCargo, fleet).
		withSpec(PlayerMessageSpec{Cargo: &cargo}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetStargateDestroyed(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateDestroyed, fleet).
		withSpec(PlayerMessageSpec{}.withStargateTargets(wp0, wp1)))
}

// Amount is damage points; Amount2 is the total number of ships lost.
func (m *messageClient) fleetStargateDamaged(player *Player, fleet *Fleet, wp0, wp1 Waypoint, damage int, shipsLostToDamage int, shipsLostToTheVoid int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetStargateDamaged, fleet).
		withSpec(PlayerMessageSpec{Amount: damage, Amount2: shipsLostToDamage + shipsLostToTheVoid}.withStargateTargets(wp0, wp1)))
}

func (m *messageClient) fleetTransferGiven(player *Player, fleet *Fleet, targetPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferGiven, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num, DestPlayerNum: targetPlayer.Num, Name: fleet.BaseName}))
}

func (m *messageClient) fleetTransferInvalidColonists(player *Player, fleet *Fleet, targetPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferInvalidColonists, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num, DestPlayerNum: targetPlayer.Num}))
}

func (m *messageClient) fleetTransferInvalidGiveRefused(player *Player, fleet *Fleet, targetPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferInvalidGiveRefused, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num, DestPlayerNum: targetPlayer.Num, Name: fleet.Name}))
}

func (m *messageClient) fleetTransferInvalidPlayer(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferInvalidPlayer, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: player.Num}))
}

func (m *messageClient) fleetTransferInvalidReceiveRefused(player *Player, fleet *Fleet, givingPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferInvalidReceiveRefused, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: givingPlayer.Num, DestPlayerNum: player.Num, Name: fleet.Name}))
}

func (m *messageClient) fleetTransferReceived(player *Player, fleet *Fleet, givingPlayer *Player) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferReceived, fleet).
		withSpec(PlayerMessageSpec{SourcePlayerNum: givingPlayer.Num, DestPlayerNum: player.Num, Name: fleet.BaseName}))
}

func (m *messageClient) fleetTransportedCargo(player *Player, fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransferredCargo, fleet).
		withSpec(PlayerMessageSpec{
			MapObjectTarget: dest.GetMapObject().ToTarget(),
			CargoTransfer:   &PlayerMessageSpecCargoTransfer{CargoType: cargoType, Transfered: transferAmount},
		}))
}

func (m *messageClient) fleetByHandTransferIncomplete(player *Player, fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int, wanted int, status CargoTransferStatus) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetByHandTransferIncomplete, fleet).
		withSpec(PlayerMessageSpec{
			MapObjectTarget: dest.GetMapObject().ToTarget(),
			CargoTransfer:   &PlayerMessageSpecCargoTransfer{CargoType: cargoType, Transfered: transferAmount, Wanted: wanted, Status: status}}))
}

func (m *messageClient) fleetTransportInvalid(player *Player, fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int, wanted int, status CargoTransferStatus) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTransportInvalid, fleet).
		withSpec(PlayerMessageSpec{
			MapObjectTarget: dest.GetMapObject().ToTarget(),
			CargoTransfer:   &PlayerMessageSpecCargoTransfer{CargoType: cargoType, Transfered: transferAmount, Wanted: wanted, Status: status},
		}),
	)

}

func (m *messageClient) fleetTargetLost(player *Player, fleet *Fleet, targetName string, targetType MapObjectType) {
	player.Messages = append(player.Messages, newFleetMessage(player, PlayerMessageFleetTargetLost, fleet).
		withSpec(PlayerMessageSpec{LostTargetType: targetType, Name: targetName}))
}

/*
 * MineralPacket Messages
 */

func (m *messageClient) planetBuiltMineralPacket(player *Player, planet *Planet, packet *MineralPacket) {
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

func (m *messageClient) planetProductionQueueStatus(player *Player, planet *Planet, completed bool) {
	messageType := PlayerMessagePlanetProductionQueueEmpty
	if completed {
		messageType = PlayerMessagePlanetProductionQueueComplete
	}
	player.Messages = append(player.Messages, newPlanetMessage(messageType, planet))
}

// PrevAmount and Amount snapshot the planet owner's habitability before and after
// the fleet's work. Amount2 is 1 for improving and -1 for degrading the planet;
// TerraformAmount records the actual changes to each habitat axis.
func (m *messageClient) planetRemoteTerraform(player, planetPlayer *Player, planet *Planet, fleet *Fleet, initialHab Hab, deterraform bool) {
	direction := 1
	if deterraform {
		direction = -1
	}
	spec := PlayerMessageSpec{
		PrevAmount:      planetPlayer.Race.GetPlanetHabitability(initialHab),
		Amount:          planetPlayer.Race.GetPlanetHabitability(planet.Hab),
		Amount2:         direction,
		TerraformAmount: planet.Hab.Subtract(initialHab),
		SourcePlayerNum: fleet.PlayerNum,
	}.withTargetFleet(fleet)
	// Use the same public fleet name as other fleet reports for the planet owner.
	spec.TargetName = newFleetMessage(player, PlayerMessagePlanetRemoteTerraform, fleet).TargetName
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetRemoteTerraform, planet).withSpec(spec))
}

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
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetColonized, planet))
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
	terraformAmount := terraformer.GetTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	habTerraformed := player.Race.GetPlanetHabitability(planet.Hab.Add(terraformAmount))

	if player.Race.Spec.Instaforming {
		hab = habTerraformed // CAs instantly terraform any planet they inhabit, so referring to "hab after terraforming" is a bit disingenuous
	}

	if hab >= 0 {
		messageType = PlayerMessagePlanetDiscoveryHabitable
	} else if habTerraformed > 0 {
		messageType = PlayerMessagePlanetDiscoveryTerraformable
	} else {
		messageType = PlayerMessagePlanetDiscoveryUninhabitable
	}

	player.Messages = append(player.Messages, newPlanetMessage(messageType, planet))
}

func (m *messageClient) planetInstaform(player *Player, planet *Planet, terraformAmount Hab) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetInstaform, planet).
		withSpec(PlayerMessageSpec{TerraformAmount: terraformAmount}))
}

// Amount and Amount2 snapshot the initial attacking and defending populations in colonists.
func (m *messageClient) planetInvaded(player *Player, planet *Planet, fleetName string, attacker, defender *Player, attackers int, defenders int, attackersKilled int, defendersKilled int, successful bool) {
	invasion := PlayerMessageSpecInvasion{
		FleetName:         fleetName,
		AttackerPlayerNum: attacker.Num,
		DefenderPlayerNum: defender.Num,
		AttackersKilled:   attackersKilled,
		DefendersKilled:   defendersKilled,
		Successful:        successful,
	}
	if player.Num == attacker.Num {
		player.Messages = append(player.Messages, newPlanetMessage(PlayerMessageFleetInvadedPlanet, planet).
			withSpec(PlayerMessageSpec{Amount: attackers, Amount2: defenders, Invasion: &invasion}))
	} else {
		player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetInvaded, planet).
			withSpec(PlayerMessageSpec{Amount: attackers, Amount2: defenders, Invasion: &invasion}))
	}
}

func (m *messageClient) planetInvadeEmpty(player *Player, planet *Planet, fleet *Fleet) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetInvadeInvalidEmpty, planet).
		withSpec(PlayerMessageSpec{}.withTargetFleet(fleet)))
}
func (m *messageClient) planetInvadeStarbase(player *Player, planet *Planet, fleet *Fleet) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetInvadeInvalidStarbase, planet).
		withSpec(PlayerMessageSpec{}.withTargetFleet(fleet)))
}

// Amount is the packet mineral mass in kT.
func (m *messageClient) planetPacketArrived(player *Player, planet *Planet, packet *MineralPacket) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPacketLanded, planet).
		withSpec(PlayerMessageSpec{Amount: packet.Cargo.Total()}))
}

// Amount is the packet mineral mass in kT.
func (m *messageClient) planetPacketCaught(player *Player, planet *Planet, packet *MineralPacket) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPacketCaught, planet).
		withSpec(PlayerMessageSpec{Amount: packet.Cargo.Total()}))
}

// Amount is the packet mineral mass in kT.
func (m *messageClient) planetPacketDamage(player *Player, planet *Planet, packet *MineralPacket, damage MineralPacketDamage) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPacketDamage, planet).
		withSpec(PlayerMessageSpec{
			Amount: packet.Cargo.Total(), MineralPacketDamage: &damage,
			HasMassDriver: planet.Spec.HasStarbase && planet.Starbase.Spec.HasMassDriver,
			PlanetEmptied: planet.GetPopulation() == 0,
		}))
}

// Amount is the signed change; Amount2 is the resulting raw habitat value.
func (m *messageClient) planetPacketPermaform(player *Player, planet *Planet, habType HabType, change int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPacketPermaform, planet).
		withSpec(PlayerMessageSpec{HabType: FromHabType(habType), Amount: change, Amount2: planet.Hab.Get(habType)}))
}

// Amount is the signed change; Amount2 is the resulting raw habitat value.
func (m *messageClient) planetPacketTerraform(player *Player, planet *Planet, habType HabType, change int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPacketTerraform, planet).
		withSpec(PlayerMessageSpec{HabType: FromHabType(habType), Amount: change, Amount2: planet.Hab.Get(habType)}))
}

// Amount is the signed change; Amount2 is the resulting raw habitat value.
func (m *messageClient) planetPermaform(player *Player, planet *Planet, habType HabType, change int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPermaform, planet).
		withSpec(PlayerMessageSpec{HabType: FromHabType(habType), Amount: change, Amount2: planet.Hab.Get(habType)}))
}

func (m *messageClient) planetPopulationDecreased(player *Player, planet *Planet, prevAmount int, amount int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPopulationDecreased, planet).
		withSpec(PlayerMessageSpec{PrevAmount: prevAmount, Amount: amount}))
}

func (m *messageClient) planetPopulationDecreasedOvercrowding(player *Player, planet *Planet, amount int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetPopulationDecreasedOvercrowding, planet).
		withSpec(PlayerMessageSpec{Amount: amount}))
}

// Amount is the signed change; Amount2 is the resulting raw habitat value.
func (m *messageClient) planetTerraform(player *Player, planet *Planet, habType HabType, change int) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlanetBuiltTerraform, planet).
		withSpec(PlayerMessageSpec{HabType: FromHabType(habType), Amount: change, Amount2: planet.Hab.Get(habType)}))
}

/*
 * Player Messages
 */

func (m *messageClient) playerDiscovered(player *Player, otherPlayer *Player) {
	player.Messages = append(player.Messages, newPlayerMessage(PlayerMessagePlayerDiscovery, otherPlayer.Num).
		withSpec(PlayerMessageSpec{Name: otherPlayer.Race.PluralName}))
}

// Amount is the research level reached.
func (m *messageClient) playerGainTechLevel(player *Player, field TechField, level int, nextField TechField) {
	player.Messages = append(player.Messages, newMessage(PlayerMessagePlayerGainTechLevel).
		withSpec(PlayerMessageSpec{Field: field, NextField: nextField, Amount: level}))
}

func (m *messageClient) playerTechGained(player *Player, field TechField, tech *Tech) {
	player.Messages = append(player.Messages, newMessage(PlayerMessagePlayerTechGained).
		withSpec(PlayerMessageSpec{Field: field, TechGained: tech.Name}))
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

func (m *messageClient) playerAcquirablePartGainedBattle(player *Player, planet *Planet, record *BattleRecord, tech string) {
	player.Messages = append(player.Messages, newBattleMessage(PlayerMessagePlayerAcquirablePartGainedBattle, planet, record).
		withSpec(PlayerMessageSpec{TechGained: tech}))
}

func (m *messageClient) playerAcquirablePartGainedScrappedFleet(player *Player, planet *Planet, fleetName string, tech string) {
	player.Messages = append(player.Messages, newPlanetMessage(PlayerMessagePlayerAcquirablePartGainedScrapFleet, planet).
		withSpec(PlayerMessageSpec{TechGained: tech, Name: fleetName}))
}

// tell a player they are dead. This always appears as the first message
func (mc *messageClient) playerDead(player, deadPlayer *Player) {
	player.Messages = append([]PlayerMessage{newPlayerMessage(PlayerMessagePlayerDead, deadPlayer.Num)}, player.Messages...)
}

// tell a player they have no planets but still have colonists. This always appears as the first message
func (mc *messageClient) playerNoPlanets(player *Player, numColonists int) {
	player.Messages = append([]PlayerMessage{newMessage(PlayerMessagePlayerNoPlanets).withSpec(PlayerMessageSpec{Amount: numColonists})}, player.Messages...)
}

func (m *messageClient) playerVictory(player *Player, victor *Player) {
	// Victory messages are always the first message of the year.
	message := newPlayerMessage(PlayerMessagePlayerVictor, victor.Num).
		withSpec(PlayerMessageSpec{Name: victor.Race.Name})
	player.Messages = append([]PlayerMessage{message}, player.Messages...)
}
