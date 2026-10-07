package cs

import (
	"fmt"
	"log/slog"
	"slices"
)

// ByHandCargoTransfers are cargo transfers the player made in the UI with something they don't own. They
// are settled when a turn is generated. When a fleet is deleted or merged its by hand transfers go to the
// fleet that took its place.
type ByHandCargoTransfer struct {
	MapObjectTarget
	SourceFleetNum int   `json:"sourceFleetNum,omitempty"`
	Cargo          Cargo `json:"cargo"`
	Fuel           int   `json:"fuel,omitempty"`
}

// CargoTransfers are a player's ByHandCargoTransfers, keyed by the location they were made. Each location
// is a shared cargo bucket for the player, see settleByHandLoads and settleByHandUnloads.
type CargoTransfers map[string][]ByHandCargoTransfer

type cargoTransferer struct {
	log     *slog.Logger
	invader invader
	game    *FullGame
}

// CargoTransferStatus will alert the user if a CargoTransfer didn't go through due to insufficient capacity or available cargo
type CargoTransferStatus int

const (
	CargoTransferStatusNone CargoTransferStatus = iota
	CargoTransferStatusOwned
	CargoTransferStatusCargo
	CargoTransferStatusCargoCapacity
	CargoTransferStatusDestCargo
	CargoTransferStatusDestCargoCapacity
	// if a starbase is present, you cannot drop invaders
	CargoTransferStatusDestStarbase
	CargoTransferStatusDestUnowned
	// colonists can't survive in deep space or salvage
	CargoTransferStatusDeepSpace
)

func (r CargoTransferStatus) String() string {
	switch r {
	case CargoTransferStatusNone:
		return "None"
	case CargoTransferStatusOwned:
		return "Owned"
	case CargoTransferStatusCargo:
		return "Insufficient Cargo"
	case CargoTransferStatusCargoCapacity:
		return "Insufficient Cargo Capacity"
	case CargoTransferStatusDestCargo:
		return "Insufficient Destination Cargo"
	case CargoTransferStatusDestCargoCapacity:
		return "Insufficient Destination Cargo Capacity"
	case CargoTransferStatusDestStarbase:
		return "Destination Has Starbase"
	case CargoTransferStatusDestUnowned:
		return "Destination Unowned"
	case CargoTransferStatusDeepSpace:
		return "Deep Space"
	default:
		return fmt.Sprintf("Unknown %d", r)
	}
}

// cargoTransferResult is the result of a single CargoType cargo transfer to a dest
type cargoTransferResult struct {
	status      CargoTransferStatus // if transfer fails, this is the reason
	fleet       *Fleet
	dest        CargoHolder
	cargoType   CargoType
	transferred int
	wanted      int
}

func newCargoTransferer(log *slog.Logger, game *FullGame) cargoTransferer {
	return cargoTransferer{log: log, game: game, invader: newInvader()}
}

func (cargoTransfers CargoTransfers) getTransfers(position Vector) []ByHandCargoTransfer {
	return cargoTransfers[position.String()]
}

// getByHandTransfer sums all cargo load/unloads for this position to determine the total amount of by hand cargo here
func (cargoTransfers CargoTransfers) getByHandTransfer(target MapObjectTarget) Cargo {
	cargo := Cargo{}
	if cargoTransfers == nil {
		return cargo
	}
	key := target.TargetPosition.String()

	for _, transfer := range cargoTransfers[key] {
		if transfer.MapObjectTarget != target {
			continue
		}
		cargo = cargo.Add(transfer.Cargo)
	}
	return cargo
}

// transferByHand adds a byHand transfer to a target. Cargo and fuel given to the target are positive
func (cargoTransfers CargoTransfers) transferByHand(fleet *Fleet, target MapObjectTarget, cargo Cargo, fuel int) {
	// add the new cargo transfer to the player
	key := fleet.Position.String()
	transfers := cargoTransfers[key]

	// if the last transfer is the same fleet/target, just update it instead of creating another one
	var lastTransfer *ByHandCargoTransfer
	if len(transfers) > 0 {
		lastTransfer = &transfers[len(transfers)-1]
	}
	if lastTransfer != nil && lastTransfer.SourceFleetNum == fleet.Num && lastTransfer.MapObjectTarget == target {
		lastTransfer.Cargo = lastTransfer.Cargo.Add(cargo)
		lastTransfer.Fuel += fuel
		return
	}

	transfer := ByHandCargoTransfer{
		SourceFleetNum:  fleet.Num,
		Cargo:           cargo,
		Fuel:            fuel,
		MapObjectTarget: target,
	}

	// add a new immediate cargo transfer target
	cargoTransfers[key] = append(cargoTransfers[key], transfer)
}

// moveByHandTransfers moves all byHand transfers from one fleet to another (in case a fleet is deleted)
func (cargoTransfers CargoTransfers) moveByHandTransfers(source *Fleet, dest *Fleet) {
	key := source.Position.String()
	transfers, ok := cargoTransfers[key]
	if !ok {
		// no transfers
		return
	}

	for i := range transfers {
		transfer := &transfers[i]
		if transfer.SourceFleetNum == source.Num {
			transfer.SourceFleetNum = dest.Num
		}
		if transfer.Targeting(source.MapObject) {
			transfer.TargetNum = dest.Num
		}
	}
}

// mergeByHandTransfers gives the by hand transfers of fleets merging into fleet to fleet
func (cargoTransfers CargoTransfers) mergeByHandTransfers(fleet *Fleet, mergingFleets []*Fleet) {
	for _, mergingFleet := range mergingFleets {
		if mergingFleet != fleet {
			cargoTransfers.moveByHandTransfers(mergingFleet, fleet)
		}
	}
}

// By hand transfers are applied as soon as the player makes them, so the player sees the result right
// away. Transfers between a player's own fleets, planets and mineral packets are final. Nothing else
// touches those before turn generation, so they are never replayed.
//
// Transfers with anything the player doesn't own (salvage, jettisoned cargo, other players' planets,
// fleets and packets) only changed the player's intel, so they are settled at turn generation. Each
// location is a shared cargo bucket for the player: all of their fleets there, and their planet if they
// own it. Settling moves the player's net exchange with each target in or out of the bucket. If another
// player got to a target first and it has less than the player took, the shortfall comes back out of
// the bucket. If a target can't accept what the player gave, it goes back into the bucket.
//
// All players' loads are settled before any unloads.

// byHandSettlement is a player's net by hand exchange with one target at a location
type byHandSettlement struct {
	player   *Player
	position Vector
	target   MapObjectTarget
	// cargo and fuel given to the target. Negative amounts were taken from it
	cargo Cargo
	fuel  int
	// the player's fleets that transferred with the target, in order
	fleetNums []int
}

// byHandResult is a by hand transfer that didn't go through as the player made it
type byHandResult struct {
	fleet       *Fleet
	target      MapObjectTarget
	cargoType   CargoType
	transferred int
	wanted      int
	status      CargoTransferStatus
}

// byHandBucket is everything a player owns at a location that can hold by hand cargo
type byHandBucket struct {
	// the fleets that transferred with the target, then the player's other fleets here
	fleets []*Fleet
	// how many of fleets transferred with the target
	numTransferred int
	planet         *Planet
}

// byHandSettlements gets a player's by hand transfers that need settling, one per location and target
func (t *cargoTransferer) byHandSettlements(player *Player) []*byHandSettlement {
	settlements := []*byHandSettlement{}

	// settle locations in a consistent order
	keys := make([]string, 0, len(player.CargoTransfers))
	for key := range player.CargoTransfers {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		transfers := player.CargoTransfers[key]
		if len(transfers) == 0 {
			continue
		}

		position, ok := t.byHandPosition(player, transfers)
		if !ok {
			t.log.Error("unable to find location of by hand transfers",
				slog.Int("Player", player.Num),
				slog.String("Location", key))
			continue
		}

		settlementsByTarget := map[MapObjectTarget]*byHandSettlement{}
		for _, transfer := range transfers {
			if t.byHandTransferIsFinal(player, transfer.MapObjectTarget) {
				// earlier versions recorded transfers between owned objects. They're already done
				continue
			}

			target := transfer.MapObjectTarget
			// jettisons are all the same target at a location
			if target.TargetType == MapObjectTypeNone {
				target = MapObjectTarget{TargetPosition: position}
			}

			// don't let a name or position change split a target
			key := MapObjectTarget{TargetType: target.TargetType, TargetNum: target.TargetNum, TargetPlayerNum: target.TargetPlayerNum}
			settlement, ok := settlementsByTarget[key]
			if !ok {
				settlement = &byHandSettlement{player: player, position: position, target: target}
				settlementsByTarget[key] = settlement
				settlements = append(settlements, settlement)
			}
			settlement.cargo = settlement.cargo.Add(transfer.Cargo)
			settlement.fuel += transfer.Fuel
			if !slices.Contains(settlement.fleetNums, transfer.SourceFleetNum) {
				settlement.fleetNums = append(settlement.fleetNums, transfer.SourceFleetNum)
			}
		}
	}

	return settlements
}

// byHandPosition finds where a player's by hand transfers happened
func (t *cargoTransferer) byHandPosition(player *Player, transfers []ByHandCargoTransfer) (Vector, bool) {
	for _, transfer := range transfers {
		if fleet := t.game.getFleet(player.Num, transfer.SourceFleetNum); fleet != nil && !fleet.Delete {
			return fleet.Position, true
		}
	}
	for _, transfer := range transfers {
		if transfer.TargetType != MapObjectTypeNone {
			return transfer.TargetPosition, true
		}
	}
	return Vector{}, false
}

// byHandTransferIsFinal is true for transfers with a player's own planets, fleets and mineral packets.
// These are done when the player makes them and don't need settling
func (t *cargoTransferer) byHandTransferIsFinal(player *Player, target MapObjectTarget) bool {
	switch target.TargetType {
	case MapObjectTypePlanet:
		planet := t.game.getPlanet(target.TargetNum)
		return planet != nil && planet.OwnedBy(player.Num)
	case MapObjectTypeFleet, MapObjectTypeMineralPacket:
		return target.TargetPlayerNum == player.Num
	}
	return false
}

// getByHandBucket gets everything a player owns at a settlement's location that can hold cargo.
// Fleets that transferred with the target come first
func (t *cargoTransferer) getByHandBucket(s *byHandSettlement) byHandBucket {
	bucket := byHandBucket{}
	for _, num := range s.fleetNums {
		if fleet := t.game.getFleet(s.player.Num, num); fleet != nil && !fleet.Delete && fleet.Position == s.position {
			bucket.fleets = append(bucket.fleets, fleet)
		}
	}

	others := []*Fleet{}
	for _, mo := range t.game.getMapObjectsAtPosition(s.position) {
		switch mo := mo.(type) {
		case *Fleet:
			if mo.PlayerNum == s.player.Num && !mo.Delete && !mo.Starbase && !slices.Contains(bucket.fleets, mo) {
				others = append(others, mo)
			}
		case *Planet:
			if mo.OwnedBy(s.player.Num) {
				bucket.planet = mo
			}
		}
	}
	slices.SortFunc(others, func(a, b *Fleet) int { return a.Num - b.Num })
	bucket.numTransferred = len(bucket.fleets)
	bucket.fleets = append(bucket.fleets, others...)

	return bucket
}

// fleet gets the fleet that represents the player in a settlement
func (b byHandBucket) fleet() *Fleet {
	if len(b.fleets) > 0 {
		return b.fleets[0]
	}
	return nil
}

// removeFromTransferred takes cargo out of the fleets that transferred with the target, returning any
// amount they don't have
func (b byHandBucket) removeFromTransferred(cargoType CargoType, amount int) int {
	return removeFromFleets(b.fleets[:b.numTransferred], cargoType, amount)
}

// removeFromOthers takes cargo out of the player's other fleets, the ones with the most first, and
// then their planet. It returns any amount they don't have
func (b byHandBucket) removeFromOthers(cargoType CargoType, amount int) int {
	others := slices.Clone(b.fleets[b.numTransferred:])
	slices.SortStableFunc(others, func(a, b *Fleet) int { return b.Cargo.GetAmount(cargoType) - a.Cargo.GetAmount(cargoType) })
	amount = removeFromFleets(others, cargoType, amount)
	if b.planet != nil {
		removed := min(amount, b.planet.Cargo.GetAmount(cargoType))
		b.planet.Cargo = b.planet.Cargo.SubtractAmount(cargoType, removed)
		amount -= removed
	}
	return amount
}

// addFuel puts fuel back in the bucket's fleets, returning any amount they have no room for
func (b byHandBucket) addFuel(amount int) int {
	for _, fleet := range b.fleets {
		added := min(amount, fleet.availableFuelSpace())
		fleet.Fuel += added
		amount -= added
	}
	return amount
}

// removeFuel takes fuel out of the bucket's fleets
func (b byHandBucket) removeFuel(amount int) {
	for _, fleet := range b.fleets {
		removed := min(amount, fleet.Fuel)
		fleet.Fuel -= removed
		amount -= removed
	}
}

func removeFromFleets(fleets []*Fleet, cargoType CargoType, amount int) int {
	for _, fleet := range fleets {
		removed := min(amount, fleet.Cargo.GetAmount(cargoType))
		fleet.Cargo = fleet.Cargo.SubtractAmount(cargoType, removed)
		amount -= removed
	}
	return amount
}

// add puts cargo back in the bucket's fleets, then its planet. Anything left over is jettisoned
func (t *cargoTransferer) addToBucket(s *byHandSettlement, b byHandBucket, cargoType CargoType, amount int) {
	for _, fleet := range b.fleets {
		added := min(amount, fleet.availableCargoSpace())
		fleet.Cargo = fleet.Cargo.AddAmount(cargoType, added)
		amount -= added
	}
	if amount == 0 {
		return
	}
	if b.planet != nil {
		b.planet.Cargo = b.planet.Cargo.AddAmount(cargoType, amount)
		return
	}
	if cargoType == Colonists {
		// colonists can't survive in space
		t.log.Warn("by hand colonists returned with nowhere to go",
			slog.Int("Player", s.player.Num),
			slog.String("Target", s.target.PrettyString()),
			slog.Int("Colonists", amount*100))
		return
	}
	t.game.getOrCreateSalvage(s.position, s.player.Num, Cargo{}.WithCargo(cargoType, amount))
}

// findByHandTarget finds the cargo holder for a settlement's target
func (t *cargoTransferer) findByHandTarget(s *byHandSettlement) (CargoHolder, bool) {
	dest, ok := t.game.getCargoHolder(s.target.TargetType, s.target.TargetNum, s.target.TargetPlayerNum)
	if !ok || dest.Deleted() {
		return nil, false
	}
	return dest, true
}

// settleByHandLoads takes the cargo a player loaded from a target. If the target is gone or has less
// than the player loaded, the shortfall comes out of the bucket: first the fleets that loaded it, then
// cargo the player is giving away at this location (added to deficit and taken out of the unloads),
// then the player's other fleets and planet. unloads is the cargo the player is giving away here.
func (t *cargoTransferer) settleByHandLoads(s *byHandSettlement, unloads Cargo, deficit *Cargo) (results []byHandResult) {
	toLoad := s.cargo.NegativeOnly().Negative()
	if toLoad == (Cargo{}) && s.fuel >= 0 {
		return nil
	}

	bucket := t.getByHandBucket(s)
	fleet := bucket.fleet()
	if fleet == nil {
		// all the fleets that loaded this cargo are gone, so is the cargo
		t.log.Warn("no fleets left to settle by hand load",
			slog.Int("Player", s.player.Num),
			slog.String("Target", s.target.PrettyString()))
		return nil
	}

	if s.fuel < 0 {
		// fuel can't be taken from other players' fleets. Orders reject it, this guards against bad data
		bucket.removeFuel(-s.fuel)
		results = append(results, byHandResult{fleet: fleet, target: s.target, cargoType: Fuel, wanted: -s.fuel, status: CargoTransferStatusOwned})
	}

	dest, found := t.findByHandTarget(s)
	for _, cargoType := range CargoTypes {
		wanted := toLoad.GetAmount(cargoType)
		if wanted == 0 {
			continue
		}

		loaded := 0
		status := CargoTransferStatusNone
		if !found {
			status = CargoTransferStatusDestCargo
		} else if !dest.CanLoad(fleet) {
			status = CargoTransferStatusOwned
		} else {
			loaded = min(wanted, dest.GetCargo().GetAmount(cargoType))
			dest.SetCargo(dest.GetCargo().SubtractAmount(cargoType, loaded))
		}

		if loaded == wanted {
			continue
		}

		t.log.Debug("by hand load came up short",
			slog.Int("Player", s.player.Num),
			slog.String("Target", s.target.PrettyString()),
			slog.String("CargoType", cargoType.String()),
			slog.Int("Wanted", wanted),
			slog.Int("Loaded", loaded))

		remaining := bucket.removeFromTransferred(cargoType, wanted-loaded)
		lost := min(remaining, unloads.GetAmount(cargoType)-deficit.GetAmount(cargoType))
		*deficit = deficit.AddAmount(cargoType, lost)
		remaining = bucket.removeFromOthers(cargoType, remaining-lost)
		if remaining > 0 {
			// this should never happen. The player can't have used more than they had
			t.log.Error("by hand load shortfall is more than the player has",
				slog.Int("Player", s.player.Num),
				slog.String("Target", s.target.PrettyString()),
				slog.String("CargoType", cargoType.String()),
				slog.Int("Remaining", remaining))
		}
		results = append(results, byHandResult{fleet: fleet, target: s.target, cargoType: cargoType, transferred: loaded, wanted: wanted, status: status})
	}

	if found {
		// delete this salvage or packet if we emptied it
		if salvage, ok := dest.(*Salvage); ok && salvage.Cargo == (Cargo{}) {
			t.game.deleteSalvage(salvage)
		}
		if packet, ok := dest.(*MineralPacket); ok && packet.Cargo == (Cargo{}) {
			t.game.deletePacket(packet)
		}
	}

	return results
}

// settleByHandUnloads gives a target the cargo a player unloaded to it. Any deficit from loads that came
// up short is taken out first. Anything the target can't accept goes back in the bucket.
// Colonists unloaded on another player's planet invade it.
func (t *cargoTransferer) settleByHandUnloads(s *byHandSettlement, deficit *Cargo) (results []byHandResult) {
	toUnload := s.cargo.PositiveOnly()
	if toUnload == (Cargo{}) && s.fuel <= 0 {
		return nil
	}

	bucket := t.getByHandBucket(s)
	fleet := bucket.fleet()

	var dest CargoHolder
	found := false
	if s.target.TargetType == MapObjectTypeNone {
		dest, found = t.game.getOrCreateSalvage(s.position, s.player.Num, Cargo{}), true
	} else {
		dest, found = t.findByHandTarget(s)
	}
	target := s.target
	if found {
		target = dest.GetMapObject().ToTarget()
	}

	for _, cargoType := range CargoTypes {
		wanted := toUnload.GetAmount(cargoType)
		if wanted == 0 {
			continue
		}

		// we can't unload cargo we lost on an earlier load
		lost := min(wanted, deficit.GetAmount(cargoType))
		*deficit = deficit.SubtractAmount(cargoType, lost)
		amount := wanted - lost

		unloaded := 0
		status := CargoTransferStatusNone
		if lost > 0 {
			status = CargoTransferStatusCargo
		}
		if amount > 0 {
			if !found {
				status = CargoTransferStatusDestCargoCapacity
			} else if _, ok := dest.(*Salvage); ok && cargoType == Colonists {
				// colonists can't survive in space. Orders reject this, this guards against old data
				status = CargoTransferStatusDeepSpace
			} else if planet, ok := dest.(*Planet); ok && cargoType == Colonists && !planet.OwnedBy(s.player.Num) {
				unloaded, status = 0, CargoTransferStatusOwned
				if fleet != nil {
					unloaded, status = t.invade(s.player, fleet, planet, amount)
				}
			} else {
				unloaded = amount
				if capacity := dest.GetCargoCapacity(); capacity != Infinite {
					unloaded = min(amount, max(0, capacity-dest.GetCargo().Total()))
				}
				if unloaded < amount {
					status = CargoTransferStatusDestCargoCapacity
				}
				dest.SetCargo(dest.GetCargo().AddAmount(cargoType, unloaded))
			}

			if returned := amount - unloaded; returned > 0 {
				t.addToBucket(s, bucket, cargoType, returned)
			}
		}

		if unloaded == wanted {
			continue
		}

		t.log.Debug("by hand unload came up short",
			slog.Int("Player", s.player.Num),
			slog.String("Target", s.target.PrettyString()),
			slog.String("CargoType", cargoType.String()),
			slog.Int("Wanted", wanted),
			slog.Int("Unloaded", unloaded))

		results = append(results, byHandResult{fleet: fleet, target: target, cargoType: cargoType, transferred: unloaded, wanted: wanted, status: status})
	}

	if s.fuel > 0 {
		// fuel only moves between fleets
		given := 0
		if destFleet, ok := dest.(*Fleet); found && ok {
			given = min(s.fuel, destFleet.availableFuelSpace())
			destFleet.Fuel += given
		}
		if returned := s.fuel - given; returned > 0 {
			if lost := bucket.addFuel(returned); lost > 0 {
				t.log.Debug("by hand fuel returned with nowhere to go",
					slog.Int("Player", s.player.Num),
					slog.String("Target", s.target.PrettyString()),
					slog.Int("Fuel", lost))
			}
			results = append(results, byHandResult{fleet: fleet, target: target, cargoType: Fuel, transferred: given, wanted: s.fuel, status: CargoTransferStatusDestCargoCapacity})
		}
	}

	if found && dest.GetMapObject().Type == MapObjectTypeSalvage && dest.GetCargo() == (Cargo{}) {
		// nothing was jettisoned after all
		t.game.deleteSalvage(dest.(*Salvage))
	}

	return results
}

// load does a fleet's load tasks at a waypoint. Like the original game, cargo loads in a fixed order
// (ironium, boranium, germanium, colonists, fuel) and dunnage loads last, once no other load is waiting.
// It returns whether the fleet should wait at the waypoint for more cargo
func (t *cargoTransferer) load(fleet *Fleet, dest CargoHolder, transportTasks WaypointTransportTasks) (results []cargoTransferResult, wait bool) {
	// a set amount task couldn't get enough from the dest
	waitForCargo := false
	// a wait for percent task isn't filled yet. The fleet waits only while it has room left
	waitForPercent := false

	dunnage := []transportTask{}
	for _, task := range transportTasks.ordered() {
		if task.Action == TransportActionLoadDunnage {
			dunnage = append(dunnage, task)
			continue
		}

		result, taskWait := t.loadTask(fleet, dest, task)
		results = append(results, result)
		if task.Action == TransportActionWaitForPercent && task.cargoType != Fuel {
			waitForPercent = waitForPercent || taskWait
		} else {
			waitForCargo = waitForCargo || taskWait
		}
	}

	wait = waitForCargo || (waitForPercent && fleet.availableCargoSpace() > 0)
	if !wait {
		for _, task := range dunnage {
			result, _ := t.loadTask(fleet, dest, task)
			results = append(results, result)
		}
	}

	// delete this salvage or packet if we emptied it
	if salvage, ok := dest.(*Salvage); ok && salvage.Cargo == (Cargo{}) {
		t.game.deleteSalvage(salvage)
	}
	if packet, ok := dest.(*MineralPacket); ok && packet.Cargo == (Cargo{}) {
		t.game.deletePacket(packet)
	}

	return results, wait
}

// loadTask does a single load task
func (t *cargoTransferer) loadTask(fleet *Fleet, dest CargoHolder, task transportTask) (cargoTransferResult, bool) {
	transferAmount, wanted, wait := t.getCargoLoadAmount(fleet, dest, task.cargoType, task.WaypointTransportTask)
	transferred, status := t.transferCargo(fleet, -transferAmount, task.cargoType, dest)
	return cargoTransferResult{
		fleet:       fleet,
		dest:        dest,
		cargoType:   task.cargoType,
		wanted:      -wanted,
		transferred: transferred,
		status:      status,
	}, wait
}

// unload does a fleet's unload tasks at a waypoint, in the same order as loads
func (t *cargoTransferer) unload(fleet *Fleet, dest CargoHolder, transportTasks WaypointTransportTasks) (results []cargoTransferResult) {
	for _, task := range transportTasks.ordered() {
		transferAmount, wanted := t.getCargoUnloadAmount(fleet, dest, task.cargoType, task.WaypointTransportTask)
		transferred, status := t.transferCargo(fleet, transferAmount, task.cargoType, dest)
		results = append(results, cargoTransferResult{
			fleet:       fleet,
			dest:        dest,
			cargoType:   task.cargoType,
			wanted:      wanted,
			transferred: transferred,
			status:      status,
		})
	}

	return results
}

// transferCargo transfers a single cargo type to/from the fleet to/from the dest. A positive
// transferAmount unloads from the fleet, a negative one loads
func (t *cargoTransferer) transferCargo(fleet *Fleet, transferAmount int, cargoType CargoType, dest CargoHolder) (transferred int, invalid CargoTransferStatus) {
	if transferAmount == 0 {
		return 0, CargoTransferStatusNone
	}

	player := t.game.getPlayer(fleet.PlayerNum)
	if transferAmount > 0 {
		switch dest := dest.(type) {
		case *Salvage:
			if cargoType == Colonists {
				// colonists can't survive in space, they stay aboard
				return 0, CargoTransferStatusDeepSpace
			}
		case *Planet:
			if cargoType == Colonists && !dest.OwnedBy(fleet.PlayerNum) {
				invaders, status := t.invade(player, fleet, dest, transferAmount)
				fleet.Cargo.Colonists -= invaders
				return invaders, status
			}
		case *Fleet:
			if !dest.OwnedBy(fleet.PlayerNum) {
				// like the original game, colonists can't be given to another player's fleet, and a
				// player won't accept anything from an enemy
				if cargoType == Colonists || t.game.getPlayer(dest.PlayerNum).IsEnemy(fleet.PlayerNum) {
					t.log.Debug("fleet cannot unload to another player's fleet",
						slog.Int("Player", fleet.PlayerNum),
						slog.String("Fleet", fleet.Name),
						slog.String("Dest", dest.Name),
						slog.String("cargoType", cargoType.String()))
					return 0, CargoTransferStatusOwned
				}
			}
		}
	}

	if status := t.transferToDest(fleet, dest, cargoType, transferAmount); status != CargoTransferStatusNone {
		return 0, status
	}

	return transferAmount, CargoTransferStatusNone
}

// invade queues an invasion of another player's planet with colonists from a fleet. It returns how many
// colonists (in kT) invade
func (t *cargoTransferer) invade(player *Player, fleet *Fleet, planet *Planet, colonists int) (int, CargoTransferStatus) {
	status := CargoTransferStatusNone
	switch {
	case !planet.Owned():
		// don't beam colonists to their death
		status = CargoTransferStatusDestUnowned
	case planet.Spec.HasStarbase:
		status = CargoTransferStatusDestStarbase
	case player.Race.Spec.LivesOnStarbases:
		status = CargoTransferStatusOwned
	}
	if status != CargoTransferStatusNone {
		t.log.Debug("fleet cannot unload colonists on planet",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Planet", planet.Name),
			slog.String("Status", status.String()))
		return 0, status
	}

	t.invader.addInvasion(invasion{
		planet:    planet,
		attacker:  player,
		defender:  t.game.getPlayer(planet.PlayerNum),
		attackers: colonists * 100,
		fleets:    []*Fleet{fleet},
	})

	return colonists, CargoTransferStatusNone
}

// getCargoLoadAmount gets the amount of cargo to transfer for loading a cargo type from a cargoholder.
// waitAtWaypoint is true when a set amount task can't get enough from the dest, or a wait for percent
// task didn't fill the hold to its percent
func (t *cargoTransferer) getCargoLoadAmount(fleet *Fleet, dest CargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, wantToTransfer int, waitAtWaypoint bool) {
	availableCapacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total()
	availableToLoad := dest.GetCargo().GetAmount(cargoType)
	currentAmount := fleet.Cargo.GetAmount(cargoType)
	totalCapacity := fleet.Spec.CargoCapacity

	// fuel transfers use different tanks
	if cargoType == Fuel {
		availableCapacity = fleet.Spec.FuelCapacity - fleet.Fuel
		// planets with starbases have Infinite fuel, but fuel only moves between fleets
		availableToLoad = max(0, dest.GetFuel())
		currentAmount = fleet.Fuel
		totalCapacity = fleet.Spec.FuelCapacity
	}

	switch task.Action {
	case TransportActionLoadOptimal:
		// fuel only
		// we set our fuel to whatever it takes to finish our waypoints and transfer the rest to the ICargoHolder target.
		// If the target is a planet or starbase (and has infinite fuel capacity), we skip this and don't give them our fuel
		if cargoType == Fuel && dest.GetFuelCapacity() != Infinite {
			fuelRequiredForWaypoints := 0
			for i := 1; i < len(fleet.Waypoints); i++ {
				fuelRequiredForWaypoints += fleet.Waypoints[i].EstFuelUsage
			}
			leftoverFuel := fleet.Fuel - fuelRequiredForWaypoints
			wantToTransfer = leftoverFuel
			fuelCapacityAvailable := dest.GetFuelCapacity() - dest.GetFuel()
			if leftoverFuel > 0 && fuelCapacityAvailable > 0 {
				// transfer the lowest of how much fuel capacity they have available or how much we can give
				// this is a bit weird because we are doing a "Load", but it's actually an unload of fuel
				// from us to a dest fleet, so make the transferAmount negative.
				transferAmount = max(-leftoverFuel, -(dest.GetFuelCapacity() - dest.GetFuel()))
			}
		}
	case TransportActionLoadAll:
		// load all available, based on our constraints
		wantToTransfer = availableToLoad
		transferAmount = min(availableToLoad, availableCapacity)
	case TransportActionLoadAmount:
		wantToTransfer = task.Amount
		transferAmount = min(min(availableToLoad, task.Amount), availableCapacity)
	case TransportActionWaitForPercent, TransportActionFillPercent:
		// we want a percent of our hold to be filled with some amount, figure out how
		// much that is in kT, i.e. 50% of 100kT would be 50kT of this mineral
		var taskAmountkT = int(float64(task.Amount) / 100 * float64(totalCapacity))
		wantToTransfer = taskAmountkT

		if currentAmount >= taskAmountkT {
			// no need to transfer any, move on
			return 0, wantToTransfer, false
		} else {

			// transfer up to our percent specified
			// wait here if we haven't loaded the amount we want. The caller moves on anyway if the hold is
			// full (in case the user suffers from innumeracy and said they wanted 50% 50% 50%)
			transferAmount = min(min(availableToLoad, taskAmountkT-currentAmount), availableCapacity)
			if (transferAmount+currentAmount) < taskAmountkT && task.Action == TransportActionWaitForPercent {
				waitAtWaypoint = true
			}
		}
	case TransportActionSetAmountTo:
		// only transfer the min of what we have, vs what we need, vs the capacity
		wantToTransfer = max(0, task.Amount-currentAmount)
		transferAmount = max(0, min(min(availableToLoad, task.Amount-currentAmount), availableCapacity))
		// like the original game, wait for the dest to have enough, not for room in our hold
		if availableToLoad < task.Amount-currentAmount {
			waitAtWaypoint = true
		}
	case TransportActionSetWaypointTo:
		// Check how much the destination has of what we want
		// if we SetWaypointTo 100kT germanium and they have 120kT, we load 20kT if we can fit it
		if availableToLoad <= task.Amount {
			// they are below the amount, we won't load (this TransportAction will possibly be used to unload later)
			break
		} else {
			wantToTransfer = availableToLoad - task.Amount
			// only transfer down to what we set
			transferAmount = min(min(availableToLoad, availableToLoad-task.Amount), availableCapacity)
		}

	case TransportActionLoadDunnage:
		// (minerals and colonists only) This command waits until all other loads and unloads are
		// complete, then loads as many colonists or amount of a mineral as will fit in the remaining
		// space. For example, setting Load All Germanium, Load Dunnage Ironium, will load all the
		// Germanium that is available, then as much Ironium as possible. If more than one dunnage cargo
		// is specified, they are loaded in the order of Ironium, Boranium, Germanium, and Colonists.
		wantToTransfer = availableToLoad
		transferAmount = min(availableToLoad, availableCapacity)
	}

	// let the caller know how much of this cargo we load
	return transferAmount, wantToTransfer, waitAtWaypoint
}

// getCargoUnloadAmount gets the amount of cargo to transfer for unloading a cargo type from a cargoholder
func (t *cargoTransferer) getCargoUnloadAmount(fleet *Fleet, dest CargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, wantToTransfer int) {

	capacity := dest.GetCargoCapacity()
	if capacity != Infinite {
		capacity = capacity - dest.GetCargo().Total()
	}
	currentAmount := fleet.Cargo.GetAmount(cargoType)

	var availableToUnload int
	if cargoType == Fuel {
		availableToUnload = fleet.Fuel
		capacity = max(0, dest.GetFuelCapacity()-dest.GetFuel())
		currentAmount = fleet.Fuel
	} else {
		availableToUnload = fleet.Cargo.GetAmount(cargoType)
	}
	switch task.Action {
	case TransportActionUnloadAll:
		// unload all available, based on our constraints
		wantToTransfer = availableToUnload
		if capacity == Infinite {
			transferAmount = availableToUnload
		} else {
			transferAmount = min(availableToUnload, capacity)
		}
	case TransportActionUnloadAmount:
		// don't unload more than the task says
		wantToTransfer = task.Amount
		if capacity == Infinite {
			transferAmount = min(availableToUnload, task.Amount)
		} else {
			transferAmount = min(min(availableToUnload, task.Amount), capacity)
		}
	case TransportActionSetAmountTo:
		// set the amount in our hold to amount, or do nothing if we have under that amount
		// unload what the dest has room for
		wantToTransfer = max(0, currentAmount-task.Amount)
		transferAmount = max(0, min(availableToUnload, currentAmount-task.Amount))
		if capacity != Infinite {
			transferAmount = min(transferAmount, capacity)
		}
	case TransportActionSetWaypointTo:
		// Make sure the waypoint has at least whatever we specified
		var currentAmount = dest.GetCargo().GetAmount(cargoType)

		if currentAmount >= task.Amount {
			// no need to transfer any, move on
			break
		} else {
			// only transfer the min of what we have, vs what we need, vs the capacity
			wantToTransfer = task.Amount - currentAmount
			if capacity == Infinite {
				transferAmount = min(availableToUnload, task.Amount-currentAmount)
			} else {
				transferAmount = min(min(availableToUnload, task.Amount-currentAmount), capacity)
			}
		}
	}
	return transferAmount, wantToTransfer
}

// transferToDest performs a transfer of a single cargo type to/from a destination. A positive
// transferAmount unloads from the fleet, a negative one loads.
// returns a status other than None if the transfer fails for some reason
func (t *cargoTransferer) transferToDest(fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) CargoTransferStatus {
	status := t.checkTransferToDest(fleet, dest, cargoType, transferAmount)
	if status != CargoTransferStatusNone {
		t.log.Debug("fleet cannot transfer cargo",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Dest", dest.GetMapObject().Name),
			slog.String("cargoType", cargoType.String()),
			slog.Int("TransferAmount", transferAmount),
			slog.String("Status", status.String()))
		return status
	}

	if cargoType == Fuel {
		fleet.Fuel -= transferAmount
		dest.(*Fleet).Fuel += transferAmount
		return CargoTransferStatusNone
	}

	fleet.Cargo = fleet.Cargo.SubtractAmount(cargoType, transferAmount)
	dest.SetCargo(dest.GetCargo().AddAmount(cargoType, transferAmount))
	return CargoTransferStatusNone
}

// checkTransferToDest checks whether a fleet can transfer a single cargo type to/from a destination
func (t *cargoTransferer) checkTransferToDest(fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) CargoTransferStatus {
	if transferAmount < 0 {
		if !dest.CanLoad(fleet) {
			return CargoTransferStatusOwned
		}
		// like the original game (with its bug fixed), a thief can't take colonists or fuel
		if (cargoType == Colonists || cargoType == Fuel) && !dest.GetMapObject().OwnedBy(fleet.PlayerNum) {
			return CargoTransferStatusOwned
		}
	}

	if cargoType == Fuel {
		_, destIsFleet := dest.(*Fleet)
		switch {
		case transferAmount > 0 && fleet.Fuel < transferAmount:
			return CargoTransferStatusCargo
		case !destIsFleet:
			// like the original game, fuel only moves between fleets
			if transferAmount > 0 {
				return CargoTransferStatusDestCargoCapacity
			}
			return CargoTransferStatusDestCargo
		case transferAmount > 0 && dest.GetFuelCapacity()-dest.GetFuel() < transferAmount:
			return CargoTransferStatusDestCargoCapacity
		case transferAmount < 0 && fleet.availableFuelSpace() < -transferAmount:
			return CargoTransferStatusCargoCapacity
		case transferAmount < 0 && dest.GetFuel() < -transferAmount:
			return CargoTransferStatusDestCargo
		}
		return CargoTransferStatusNone
	}

	destCargo := dest.GetCargo()
	switch {
	case transferAmount > 0 && !fleet.Cargo.CanTransferAmount(cargoType, transferAmount):
		return CargoTransferStatusCargo
	case transferAmount > 0 && dest.GetCargoCapacity() != Infinite && (dest.GetCargoCapacity()-destCargo.Total()) < transferAmount:
		return CargoTransferStatusDestCargoCapacity
	case transferAmount < 0 && fleet.availableCargoSpace() < -transferAmount:
		return CargoTransferStatusCargoCapacity
	case transferAmount < 0 && !destCargo.CanTransferAmount(cargoType, -transferAmount):
		return CargoTransferStatusDestCargo
	}
	return CargoTransferStatusNone
}
