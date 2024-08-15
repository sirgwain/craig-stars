package cs

const Unlimited = -1

// The cargoHolder is an interface implemented by any map object that can hold cargo. It's used for handling
// cargo transfers between different types of map objects.
type cargoHolder interface {
	getMapObject() MapObject
	getCargo() *Cargo
	getCargoCapacity() int
	getFuel() int
	getFuelCapacity() int
	canLoad(playerNum int) bool
	canTransfer(transferAmount CargoTransferRequest) bool
}

func (ch *Planet) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *Planet) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Planet) getCargoCapacity() int {
	return Unlimited
}

func (ch *Planet) getFuel() int {
	if ch.Spec.HasStarbase {
		return Unlimited
	} else {
		return 0
	}
}

func (ch *Planet) getFuelCapacity() int {
	return Unlimited
}

// players can load from unowned planets or planets they own
func (ch *Planet) canLoad(playerNum int) bool {
	return !ch.Owned() || ch.OwnedBy(playerNum)
}

// planets can't transfer fuel
func (ch *Planet) canTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel > 0 {
		return false
	}
	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}

func (ch *Fleet) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *Fleet) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Fleet) getFuel() int {
	return ch.Fuel
}

func (ch *Fleet) getCargoCapacity() int {
	return ch.Spec.CargoCapacity
}

func (ch *Fleet) getFuelCapacity() int {
	return ch.Spec.FuelCapacity
}

// players can load from fleets they own
func (ch *Fleet) canLoad(playerNum int) bool {
	return ch.OwnedBy(playerNum)
}

// planets can't transfer fuel
func (ch *Fleet) canTransfer(transferAmount CargoTransferRequest) bool {
	return ch.Fuel >= transferAmount.Fuel && ch.Cargo.CanTransfer(transferAmount.Cargo)
}

func (ch *Salvage) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *Salvage) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Salvage) getCargoCapacity() int {
	return Unlimited
}

func (ch *Salvage) getFuel() int {
	return 0
}

func (ch *Salvage) getFuelCapacity() int {
	return 0
}

// salvage can't transfer fuel
func (ch *Salvage) canTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel > 0 {
		return false
	}
	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}

// players can load from all salvages
func (ch *Salvage) canLoad(playerNum int) bool {
	return true
}

func (ch *MineralPacket) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *MineralPacket) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *MineralPacket) getCargoCapacity() int {
	return Unlimited
}

func (ch *MineralPacket) getFuel() int {
	return 0
}

func (ch *MineralPacket) getFuelCapacity() int {
	return 0
}

// players can load from all mineralPackets
func (ch *MineralPacket) canLoad(playerNum int) bool {
	return true
}

// mineral packets can't transfer fuel
func (ch *MineralPacket) canTransfer(transferAmount CargoTransferRequest) bool {
	if transferAmount.Fuel > 0 {
		return false
	}
	return ch.Cargo.CanTransfer(transferAmount.Cargo)
}
