package cs

const Unlimited = -1

type cargoHolder interface {
	getMapObject() MapObject
	getCargo() *Cargo
	getCargpCapacity() int
	getFuel() int
	getFuelCapacity() int
}

func (ch *Planet) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *Planet) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Planet) getCargpCapacity() int {
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

func (ch *Fleet) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *Fleet) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Fleet) getFuel() int {
	return ch.Fuel
}

func (ch *Fleet) getCargpCapacity() int {
	return ch.Spec.CargoCapacity
}

func (ch *Fleet) getFuelCapacity() int {
	return ch.Spec.FuelCapacity
}

func (ch *Salvage) getMapObject() MapObject {
	return ch.MapObject
}

func (ch *Salvage) getCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Salvage) getCargpCapacity() int {
	return Unlimited
}

func (ch *Salvage) getFuel() int {
	return 0
}

func (ch *Salvage) getFuelCapacity() int {
	return 0
}
