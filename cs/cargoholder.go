package cs

const Unlimited = -1

// The CargoHolder is an interface implemented by any map object that can hold cargo. It's used for handling
// cargo transfers between different types of map objects.
type CargoHolder interface {
	GetMapObject() MapObject
	GetCargo() Cargo
	SetCargo(cargo Cargo)
	GetCargoCapacity() int
	GetFuel() int
	GetFuelCapacity() int
	CanLoad(fleet *Fleet) bool
	CanTransfer(transferAmount CargoTransferRequest) bool
	Deleted() bool
}

func (ch *Planet) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *Planet) GetCargo() Cargo {
	return ch.Cargo
}

func (ch *Planet) SetCargo(cargo Cargo) {
	ch.Cargo = cargo
}

func (ch *Planet) GetCargoCapacity() int {
	return Unlimited
}

func (ch *Planet) GetFuel() int {
	if ch.Spec.HasStarbase {
		return Unlimited
	} else {
		return 0
	}
}

func (ch *Planet) GetFuelCapacity() int {
	return Unlimited
}

// players can load from unowned planets or planets they own
func (ch *Planet) CanLoad(fleet *Fleet) bool {
	return !ch.Owned() || ch.OwnedBy(fleet.PlayerNum) || fleet.Spec.CanStealPlanetCargo
}

// planets can't transfer fuel
func (ch *Planet) CanTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel > 0 {
		return false
	}
	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}
func (ch *Planet) Deleted() bool {
	return ch.Delete
}

func (ch *PlanetIntel) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *PlanetIntel) GetCargo() Cargo {
	return ch.Cargo
}

func (ch *PlanetIntel) SetCargo(cargo Cargo) {
	ch.Cargo = cargo
}

func (ch *PlanetIntel) GetCargoCapacity() int {
	return Unlimited
}

func (ch *PlanetIntel) GetFuel() int {
	if ch.Spec.HasStarbase {
		return Unlimited
	} else {
		return 0
	}
}

func (ch *PlanetIntel) GetFuelCapacity() int {
	return Unlimited
}

// players can load from unowned planets or planets they own
func (ch *PlanetIntel) CanLoad(fleet *Fleet) bool {
	return !ch.Owned() || ch.OwnedBy(fleet.PlayerNum) || fleet.Spec.CanStealPlanetCargo
}

// planets can't transfer fuel
func (ch *PlanetIntel) CanTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel > 0 {
		return false
	}
	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}
func (ch *PlanetIntel) Deleted() bool {
	return false
}

func (ch *Fleet) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *Fleet) GetCargo() Cargo {
	return ch.Cargo
}

func (ch *Fleet) SetCargo(cargo Cargo) {
	ch.Cargo = cargo
}

func (ch *Fleet) GetFuel() int {
	return ch.Fuel
}

func (ch *Fleet) GetCargoCapacity() int {
	return ch.Spec.CargoCapacity
}

func (ch *Fleet) GetFuelCapacity() int {
	return ch.Spec.FuelCapacity
}

// players can load from fleets they own
func (ch *Fleet) CanLoad(fleet *Fleet) bool {
	return ch.OwnedBy(fleet.PlayerNum) || fleet.Spec.CanStealFleetCargo
}

// planets can't transfer fuel
func (ch *Fleet) CanTransfer(transferAmount CargoTransferRequest) bool {
	return ch.Fuel >= transferAmount.Fuel && ch.Cargo.CanTransfer(transferAmount.Cargo)
}
func (ch *Fleet) Deleted() bool {
	return ch.Delete
}

func (ch *Salvage) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *Salvage) GetCargo() Cargo {
	return ch.Cargo
}

func (ch *Salvage) SetCargo(cargo Cargo) {
	ch.Cargo = cargo
}

func (ch *Salvage) GetCargoCapacity() int {
	return Unlimited
}

func (ch *Salvage) GetFuel() int {
	return 0
}

func (ch *Salvage) GetFuelCapacity() int {
	return 0
}

// salvage can't transfer fuel
func (ch *Salvage) CanTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel != 0 {
		return false
	}
	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}

// players can load from all salvages
func (ch *Salvage) CanLoad(fleet *Fleet) bool {
	return true
}

func (ch *Salvage) Deleted() bool {
	return ch.Delete
}

func (ch *MineralPacket) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *MineralPacket) GetCargo() Cargo {
	return ch.Cargo
}

func (ch *MineralPacket) SetCargo(cargo Cargo) {
	ch.Cargo = cargo
}

func (ch *MineralPacket) GetCargoCapacity() int {
	// can't add to it, only take away
	return ch.Cargo.Total()
}

func (ch *MineralPacket) GetFuel() int {
	return 0
}

func (ch *MineralPacket) GetFuelCapacity() int {
	return 0
}

// players can load from all mineralPackets
func (ch *MineralPacket) CanLoad(fleet *Fleet) bool {
	return true
}

// mineral packets can't transfer fuel
func (ch *MineralPacket) CanTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel != 0 || transferAmount.Colonists != 0 {
		return false
	}

	// can't receive cargo, only give it away
	if transferAmount.HasNegative() {
		return false
	}

	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}

func (ch *MineralPacket) Deleted() bool {
	return ch.Delete
}
