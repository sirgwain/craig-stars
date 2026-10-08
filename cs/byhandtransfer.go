package cs

import (
	"fmt"
	"log/slog"
	"slices"
)

// ByHandCargoTransfers are cargo transfers the player made in the UI, in order. Cargo and fuel are what
// the source fleet gave the target; negative amounts were taken from it. Splits and merges are recorded
// as transfers between the fleets, so cargo can be followed if a load comes up short.
type ByHandCargoTransfer struct {
	MapObjectTarget
	SourceFleetNum int   `json:"sourceFleetNum,omitempty"`
	Cargo          Cargo `json:"cargo"`
	Fuel           int   `json:"fuel,omitempty"`
}

// CargoTransfers are a player's ByHandCargoTransfers, keyed by the location they were made. See
// fleetByHandTransfers for how they're settled.
type CargoTransfers map[string][]ByHandCargoTransfer

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

// By hand transfers are applied as soon as the player makes them, so the player sees the result right
// away. Every transfer is recorded in order, including transfers between a player's own fleets,
// planets and mineral packets, and splits and merges. Transfers with the player's own things are final,
// nothing else touches them before turn generation.
//
// Transfers with anything the player doesn't own (salvage, jettisoned cargo, other players' planets,
// fleets and packets) only changed the player's intel, so they are settled at turn generation. Each
// player's net exchange with each target is settled, all players' loads before any unloads. If another
// player got to a target first and it has less than the player loaded, the player's transfers at that
// location are replayed in order, like the original game, with each one limited to what was actually
// there. Only the fleets and planet the missing cargo was passed through come up short. If a target
// can't accept what the player gave, it goes back to the fleets that gave it.

// byHandHolder identifies something a player transferred cargo with by hand
type byHandHolder struct {
	Type      MapObjectType
	Num       int
	PlayerNum int
}

// byHandLocation is a player's by hand transfers at one location
type byHandLocation struct {
	player    *Player
	position  Vector
	transfers []ByHandCargoTransfer
	// cargo in the player's own fleets, planet and packets here, as the player saw it before anything is
	// settled
	seen map[byHandHolder]Cargo
	// settlements with things the player doesn't own
	settlements []*byHandSettlement
	// cargo the player's loads from each target came up short
	shortfalls map[byHandHolder]Cargo
	// cargo each target had, and the cargo the player loaded from it, when its loads were settled
	available map[byHandHolder]Cargo
	loaded    map[byHandHolder]Cargo
	// cargo the player gave to each target but didn't have, because of a shortfall
	lost map[byHandHolder]Cargo
}

// byHandSettlement is a player's net by hand exchange with one target at a location
type byHandSettlement struct {
	location *byHandLocation
	player   *Player
	position Vector
	holder   byHandHolder
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

// byHandBucket is where cargo that didn't go through goes back to: the fleets that transferred with a
// target, then the player's planet
type byHandBucket struct {
	fleets []*Fleet
	planet *Planet
}

// newByHandHolder gets the holder for a transfer's target. Jettisons are all the same target at a location
func newByHandHolder(target MapObjectTarget) byHandHolder {
	if target.TargetType == MapObjectTypeNone {
		return byHandHolder{}
	}
	return byHandHolder{Type: target.TargetType, Num: target.TargetNum, PlayerNum: target.TargetPlayerNum}
}

// fleetHolder gets the holder for one of a player's fleets
func fleetHolder(player *Player, num int) byHandHolder {
	return byHandHolder{Type: MapObjectTypeFleet, Num: num, PlayerNum: player.Num}
}

// byHandLocations gets a player's by hand transfers, one per location, with the settlements they need
func (t *cargoTransferer) byHandLocations(player *Player) []*byHandLocation {
	locations := []*byHandLocation{}

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

		location := &byHandLocation{
			player:     player,
			position:   position,
			transfers:  transfers,
			seen:       map[byHandHolder]Cargo{},
			shortfalls: map[byHandHolder]Cargo{},
			available:  map[byHandHolder]Cargo{},
			loaded:     map[byHandHolder]Cargo{},
			lost:       map[byHandHolder]Cargo{},
		}
		locations = append(locations, location)

		settlementsByHolder := map[byHandHolder]*byHandSettlement{}
		for _, transfer := range transfers {
			source := fleetHolder(player, transfer.SourceFleetNum)
			location.seen[source] = t.ownedCargo(player, source)

			holder := newByHandHolder(transfer.MapObjectTarget)
			if t.byHandTransferIsFinal(player, holder) {
				location.seen[holder] = t.ownedCargo(player, holder)
				continue
			}

			settlement, ok := settlementsByHolder[holder]
			if !ok {
				target := transfer.MapObjectTarget
				if holder == (byHandHolder{}) {
					target = MapObjectTarget{TargetPosition: position}
				}
				settlement = &byHandSettlement{location: location, player: player, position: position, holder: holder, target: target}
				settlementsByHolder[holder] = settlement
				location.settlements = append(location.settlements, settlement)
			}
			settlement.cargo = settlement.cargo.Add(transfer.Cargo)
			settlement.fuel += transfer.Fuel
			if !slices.Contains(settlement.fleetNums, transfer.SourceFleetNum) {
				settlement.fleetNums = append(settlement.fleetNums, transfer.SourceFleetNum)
			}
		}
	}

	return locations
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

// byHandTransferIsFinal is true for a player's own planets, fleets and mineral packets. Transfers with
// these are done when the player makes them and don't need settling
func (t *cargoTransferer) byHandTransferIsFinal(player *Player, holder byHandHolder) bool {
	switch holder.Type {
	case MapObjectTypePlanet:
		planet := t.game.getPlanet(holder.Num)
		return planet != nil && planet.OwnedBy(player.Num)
	case MapObjectTypeFleet, MapObjectTypeMineralPacket:
		return holder.PlayerNum == player.Num
	}
	return false
}

// ownedCargoHolder gets one of a player's own fleets, planets or packets, if it still exists
func (t *cargoTransferer) ownedCargoHolder(player *Player, holder byHandHolder) CargoHolder {
	switch holder.Type {
	case MapObjectTypeFleet:
		if fleet := t.game.getFleet(player.Num, holder.Num); fleet != nil && !fleet.Delete {
			return fleet
		}
	case MapObjectTypePlanet:
		if planet := t.game.getPlanet(holder.Num); planet != nil && planet.OwnedBy(player.Num) {
			return planet
		}
	case MapObjectTypeMineralPacket:
		if packet := t.game.getMineralPacket(player.Num, holder.Num); packet != nil && !packet.Delete {
			return packet
		}
	}
	return nil
}

// ownedCargo gets the cargo in one of a player's own fleets, planets or packets
func (t *cargoTransferer) ownedCargo(player *Player, holder byHandHolder) Cargo {
	if ch := t.ownedCargoHolder(player, holder); ch != nil {
		return ch.GetCargo()
	}
	return Cargo{}
}

// byHandFleet gets the fleet that holds what a fleet transferred by hand. A fleet that merged into another
// fleet (or split away entirely) gave everything to it, so follow it there
func (t *cargoTransferer) byHandFleet(location *byHandLocation, num int) *Fleet {
	player := location.player
	// each step follows a transfer, so this can't take more steps than there are transfers (fleet nums
	// can be reused, which could otherwise loop)
	for range location.transfers {
		if fleet := t.game.getFleet(player.Num, num); fleet != nil && !fleet.Delete {
			return fleet
		}
		next := None
		for _, transfer := range location.transfers {
			if transfer.SourceFleetNum == num && transfer.TargetType == MapObjectTypeFleet && transfer.TargetPlayerNum == player.Num {
				next = transfer.TargetNum
			}
		}
		if next == None {
			return nil
		}
		num = next
	}
	return nil
}

// getByHandBucket gets the fleets that transferred with a settlement's target, and the player's planet
func (t *cargoTransferer) getByHandBucket(s *byHandSettlement) byHandBucket {
	bucket := byHandBucket{}
	for _, num := range s.fleetNums {
		if fleet := t.byHandFleet(s.location, num); fleet != nil && fleet.Position == s.position && !slices.Contains(bucket.fleets, fleet) {
			bucket.fleets = append(bucket.fleets, fleet)
		}
	}

	for _, mo := range t.game.getMapObjectsAtPosition(s.position) {
		if planet, ok := mo.(*Planet); ok && planet.OwnedBy(s.player.Num) {
			bucket.planet = planet
		}
	}

	return bucket
}

// fleet gets the fleet that represents the player in a settlement
func (b byHandBucket) fleet() *Fleet {
	if len(b.fleets) > 0 {
		return b.fleets[0]
	}
	return nil
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

// addToBucket puts cargo back in the bucket's fleets, then its planet. Anything left over is jettisoned
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

// settleByHandLoads takes the cargo a player loaded from a target. If the target is gone or has less than
// the player loaded, the shortfall is recorded on the location for replayByHandShortfalls
func (t *cargoTransferer) settleByHandLoads(s *byHandSettlement) (results []byHandResult) {
	toLoad := s.cargo.NegativeOnly().Negative()
	if toLoad == (Cargo{}) && s.fuel >= 0 {
		return nil
	}

	bucket := t.getByHandBucket(s)
	fleet := bucket.fleet()
	if fleet == nil {
		// the fleets that loaded this cargo are gone, and so is the cargo
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
	if found && dest.CanLoad(fleet) {
		s.location.available[s.holder] = dest.GetCargo()
	}
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
		s.location.loaded[s.holder] = s.location.loaded[s.holder].AddAmount(cargoType, loaded)

		if loaded == wanted {
			continue
		}

		t.log.Debug("by hand load came up short",
			slog.Int("Player", s.player.Num),
			slog.String("Target", s.target.PrettyString()),
			slog.String("CargoType", cargoType.String()),
			slog.Int("Wanted", wanted),
			slog.Int("Loaded", loaded))

		s.location.shortfalls[s.holder] = s.location.shortfalls[s.holder].AddAmount(cargoType, wanted-loaded)
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

// replayByHandShortfalls handles loads that came up short at a location. Like the original game, it
// replays the player's transfers there in order from what they started with, with each transfer limited
// to what was actually there. Each of the player's fleets, planets and packets loses whatever it ends up
// short compared to what the player saw, and each target is corrected to what it actually got. Cargo the
// player gave away but didn't have is recorded on the location for settleByHandUnloads
func (t *cargoTransferer) replayByHandShortfalls(location *byHandLocation) (results []byHandResult) {
	player := location.player
	for _, cargoType := range CargoTypes {
		short := false
		for _, shortfall := range location.shortfalls {
			short = short || shortfall.GetAmount(cargoType) > 0
		}
		if !short {
			continue
		}

		// work backwards from what the player saw to what they started with
		cargo := map[byHandHolder]int{}
		for holder, seen := range location.seen {
			cargo[holder] = seen.GetAmount(cargoType)
		}
		for i := len(location.transfers) - 1; i >= 0; i-- {
			transfer := location.transfers[i]
			amount := transfer.Cargo.GetAmount(cargoType)
			cargo[fleetHolder(player, transfer.SourceFleetNum)] += amount
			if holder := newByHandHolder(transfer.MapObjectTarget); t.byHandTransferIsFinal(player, holder) {
				cargo[holder] -= amount
			}
		}
		for holder, amount := range cargo {
			cargo[holder] = max(0, amount)
		}

		// targets we don't own start with what they actually had for us
		settlements := map[byHandHolder]*byHandSettlement{}
		start := map[byHandHolder]int{}
		for _, s := range location.settlements {
			settlements[s.holder] = s
			switch {
			case s.cargo.GetAmount(cargoType) < 0:
				start[s.holder] = location.available[s.holder].GetAmount(cargoType)
			case s.holder == byHandHolder{}:
				// we can only load what we jettisoned
			default:
				if dest, found := t.findByHandTarget(s); found {
					start[s.holder] = dest.GetCargo().GetAmount(cargoType)
				}
			}
			cargo[s.holder] = start[s.holder]
		}

		// replay the transfers in order, limiting each one to what's there
		for _, transfer := range location.transfers {
			amount := transfer.Cargo.GetAmount(cargoType)
			if amount == 0 {
				continue
			}
			source := fleetHolder(player, transfer.SourceFleetNum)
			holder := newByHandHolder(transfer.MapObjectTarget)

			// given to the target, or taken from it
			transferred := min(amount, cargo[source])
			status := CargoTransferStatusCargo
			if amount < 0 {
				transferred = max(amount, -cargo[holder])
				status = CargoTransferStatusDestCargo
			}
			cargo[source] -= transferred
			cargo[holder] += transferred

			if transferred == amount {
				continue
			}
			// settling loads and unloads reports the shortfalls with targets we don't own
			if s := settlements[holder]; s != nil && (s.cargo.GetAmount(cargoType) < 0) == (amount < 0) {
				continue
			}
			if fleet := t.game.getFleet(player.Num, transfer.SourceFleetNum); fleet != nil && !fleet.Delete {
				results = append(results, byHandResult{fleet: fleet, target: transfer.MapObjectTarget, cargoType: cargoType, transferred: Abs(transferred), wanted: Abs(amount), status: status})
			}
		}

		// take away whatever each of our fleets, planets and packets came up short
		for holder, seen := range location.seen {
			short := seen.GetAmount(cargoType) - cargo[holder]
			ch := t.ownedCargoHolder(player, holder)
			if short <= 0 || ch == nil {
				continue
			}
			removed := min(short, ch.GetCargo().GetAmount(cargoType))
			ch.SetCargo(ch.GetCargo().SubtractAmount(cargoType, removed))
			if removed < short {
				t.log.Error("by hand shortfall is more than the player has",
					slog.Int("Player", player.Num),
					slog.String("Holder", fmt.Sprintf("%v", holder)),
					slog.String("CargoType", cargoType.String()),
					slog.Int("Short", short-removed))
			}
		}

		// correct each target to what it actually got from us, or gave us
		for _, s := range location.settlements {
			got := cargo[s.holder] - start[s.holder]
			net := s.cargo.GetAmount(cargoType)
			if net > 0 {
				// settleByHandUnloads gives the target what we actually had to give
				location.lost[s.holder] = location.lost[s.holder].AddAmount(cargoType, net-max(0, got))
				got = min(0, got)
			} else {
				// settleByHandLoads already took what we loaded
				got += location.loaded[s.holder].GetAmount(cargoType)
			}
			if got != 0 {
				t.correctByHandTarget(s, cargoType, got)
			}
		}
	}
	return results
}

// correctByHandTarget gives a target cargo, or takes it if amount is negative
func (t *cargoTransferer) correctByHandTarget(s *byHandSettlement, cargoType CargoType, amount int) {
	var dest CargoHolder
	if s.holder == (byHandHolder{}) {
		dest = t.game.getOrCreateSalvage(s.position, s.player.Num, Cargo{})
	} else if found, ok := t.findByHandTarget(s); ok {
		dest = found
	}
	if dest == nil {
		return
	}
	dest.SetCargo(dest.GetCargo().AddAmount(cargoType, max(amount, -dest.GetCargo().GetAmount(cargoType))))
	if salvage, ok := dest.(*Salvage); ok && salvage.Cargo == (Cargo{}) {
		t.game.deleteSalvage(salvage)
	}
}

// settleByHandUnloads gives a target the cargo a player unloaded to it, minus anything the player turned
// out not to have. Anything the target can't accept goes back to the fleets that gave it. Colonists
// unloaded on another player's planet invade it.
func (t *cargoTransferer) settleByHandUnloads(s *byHandSettlement) (results []byHandResult) {
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

		// we can't unload cargo we didn't have
		lost := min(wanted, s.location.lost[s.holder].GetAmount(cargoType))
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
