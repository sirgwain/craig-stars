package cs

type CargoTransfers map[string][]ImmediateCargoTransfer

type ImmediateCargoTransfer struct {
	MapObjectTarget
	SourceFleetNum int   `json:"sourceFleetNum,omitempty"`
	Cargo          Cargo `json:"cargo"`
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
