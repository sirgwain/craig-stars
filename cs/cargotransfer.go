package cs

import "slices"

type CargoTransfers map[string][]ImmediateCargoTransfer

type ImmediateCargoTransfer struct {
	MapObjectTarget `tstype:",extends"`
	SourceFleetNum  int   `json:"sourceFleetNum,omitempty"`
	Cargo           Cargo `json:"cargo"`
}

func (cargoTransfers CargoTransfers) getTransfers(position Vector) []ImmediateCargoTransfer {
	return cargoTransfers[position.String()]
}

// getJettison sums all jettison cargo calls for this position to determine the total amount of jettisoned cargo here
func (cargoTransfers CargoTransfers) getJettison(position Vector) Cargo {
	cargo := Cargo{}
	if cargoTransfers == nil {
		return cargo
	}
	key := position.String()

	for _, transfer := range cargoTransfers[key] {
		if transfer.TargetType != MapObjectTypeNone {
			continue
		}
		cargo = cargo.Add(transfer.Cargo)
	}
	return cargo
}

// jettisonCargo jettison's cargo
func (cargoTransfers CargoTransfers) jettisonCargo(fleet *Fleet, jettison Cargo) {
	transfer := ImmediateCargoTransfer{
		SourceFleetNum: fleet.Num,
		Cargo:          jettison,
		MapObjectTarget: MapObjectTarget{
			TargetPosition: fleet.Position,
		},
	}

	// add the new cargo transfer to the player
	key := fleet.Position.String()
	transfers := cargoTransfers[key]

	// if the last transfer is the same fleet/target, just update it instead of creating another one
	if transfers != nil && transfers[len(transfers)-1].SourceFleetNum == fleet.Num && transfers[len(transfers)-1].MapObjectTarget == transfer.MapObjectTarget {
		transfers[len(transfers)-1].Cargo = transfers[len(transfers)-1].Cargo.Add(jettison)
		return
	}

	// add a new immediate cargo transfer for this jettison
	cargoTransfers[key] = append(cargoTransfers[key], transfer)
}

// splitFleetCargoTransfers splits the ImmediateCargoTransfers for a source fleet into two
// ImmediateCargoTransfers, based on capacity of each fleet
func (cargoTransfers CargoTransfers) splitFleetCargoTransfers(source *Fleet, dest *Fleet) error {
	key := source.Position.String()
	transfers, ok := cargoTransfers[key]
	if !ok {
		// no transfers
		return nil
	}

	updatedTransfers := make([]ImmediateCargoTransfer, 0, len(transfers))
	for _, transfer := range transfers {
		if transfer.SourceFleetNum != source.Num {
			updatedTransfers = append(updatedTransfers, transfer)
			continue
		}
		// split this transfer
		sourceCargo := transfer.Cargo.ToArray()
		cargo1, cargo2, err := splitValues(
			source.Spec.CargoCapacity+dest.Spec.CargoCapacity,
			source.Spec.CargoCapacity,
			dest.Spec.CargoCapacity,
			sourceCargo[:]...)
		if err != nil {
			return err
		}
		// copy this transfer into two transfers
		transfer1 := transfer
		transfer2 := transfer

		transfer1.Cargo = NewCargoFromArray([4]int(cargo1))
		if transfer1.Cargo.absSum() > 0 {
			updatedTransfers = append(updatedTransfers, transfer1)
		}

		transfer2.Cargo = NewCargoFromArray([4]int(cargo2))
		if transfer2.Cargo.absSum() > 0 {
			transfer2.SourceFleetNum = dest.Num
			updatedTransfers = append(updatedTransfers, transfer2)
		}
	}

	cargoTransfers[key] = updatedTransfers
	return nil
}

// mergeFleetCargoTransfers merges cargo transfers from merging fleets into a source fleet
func (cargoTransfers CargoTransfers) mergeFleetCargoTransfers(fleet *Fleet, mergingFleets []*Fleet) {
	key := fleet.Position.String()
	transfers, ok := cargoTransfers[key]
	if !ok {
		// no transfers
		return
	}

	updatedTransfers := make([]ImmediateCargoTransfer, 0, len(transfers))
	for i, transfer := range transfers {
		if !slices.ContainsFunc(mergingFleets, func(f *Fleet) bool { return transfer.SourceFleetNum == f.Num }) {
			// this transfer isn't about a merging fleet, continue
			updatedTransfers = append(updatedTransfers, transfer)
			continue
		}

		var prevTransfer *ImmediateCargoTransfer
		if i > 0 {
			prevTransfer = &updatedTransfers[i-1]
		}

		// if the previous transfer is the fleet we're merging into AND the target is the same, just merge the request
		// into the previous transfer
		if prevTransfer != nil && prevTransfer.SourceFleetNum == fleet.Num && prevTransfer.MapObjectTarget == transfer.MapObjectTarget {
			prevTransfer.Cargo = prevTransfer.Cargo.Add(transfer.Cargo)
			continue
		}

		// give this transfer to the merged fleet
		transfer.SourceFleetNum = fleet.Num
		updatedTransfers = append(updatedTransfers, transfer)
	}

	cargoTransfers[key] = updatedTransfers
}

// getUnloadTasks creates a WaypointTransportTask for each positive cargo value in the ImmediateCargoTransfer
func (o ImmediateCargoTransfer) getUnloadTasks() WaypointTransportTasks {
	tt := WaypointTransportTasks{}
	if o.Cargo.Ironium > 0 {
		tt.Ironium.Action = TransportActionUnloadAmount
		tt.Ironium.Amount = o.Cargo.Ironium
	}
	if o.Cargo.Boranium > 0 {
		tt.Boranium.Action = TransportActionUnloadAmount
		tt.Boranium.Amount = o.Cargo.Boranium
	}
	if o.Cargo.Germanium > 0 {
		tt.Germanium.Action = TransportActionUnloadAmount
		tt.Germanium.Amount = o.Cargo.Germanium
	}
	if o.Cargo.Colonists > 0 {
		tt.Colonists.Action = TransportActionUnloadAmount
		tt.Colonists.Amount = o.Cargo.Colonists
	}

	return tt
}

// getLoadTasks creates a WaypointTransportTask for each negative cargo value in the ImmediateCargoTransfer
func (o ImmediateCargoTransfer) getLoadTasks() WaypointTransportTasks {
	tt := WaypointTransportTasks{}
	if o.Cargo.Ironium < 0 {
		tt.Ironium.Action = TransportActionUnloadAmount
		tt.Ironium.Amount = o.Cargo.Ironium
	}
	if o.Cargo.Boranium < 0 {
		tt.Boranium.Action = TransportActionUnloadAmount
		tt.Boranium.Amount = o.Cargo.Boranium
	}
	if o.Cargo.Germanium < 0 {
		tt.Germanium.Action = TransportActionUnloadAmount
		tt.Germanium.Amount = o.Cargo.Germanium
	}
	if o.Cargo.Colonists < 0 {
		tt.Colonists.Action = TransportActionUnloadAmount
		tt.Colonists.Amount = o.Cargo.Colonists
	}

	return tt
}
