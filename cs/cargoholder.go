package cs

const Unlimited = -1

type cargoHolder interface {
	getMapObject() MapObject
	getCargo() *Cargo
	getCargoCapacity() int
	getFuel() int
	getFuelCapacity() int
	canLoad(playerNum int) bool
	MarkDirty()
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
	return 0
}

// players can load from unowned planets or planets they own
func (ch *Planet) canLoad(playerNum int) bool {
	return !ch.owned() || ch.OwnedBy(playerNum)
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
