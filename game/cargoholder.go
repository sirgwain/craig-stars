package game

const Unlimited = -1

type CargoHolder interface {
	GetMapObject() MapObject
	GetCargo() *Cargo
	GetCargpCapacity() int
	GetFuel() int
	GetFuelCapacity() int
}

func (ch *Planet) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *Planet) GetCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Planet) GetCargpCapacity() int {
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
	return 0
}

func (ch *Fleet) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *Fleet) GetCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Fleet) GetFuel() int {
	return ch.Fuel
}

func (ch *Fleet) GetCargpCapacity() int {
	return ch.Spec.CargoCapacity
}

func (ch *Fleet) GetFuelCapacity() int {
	return ch.Spec.FuelCapacity
}

func (ch *Salvage) GetMapObject() MapObject {
	return ch.MapObject
}

func (ch *Salvage) GetCargo() *Cargo {
	return &ch.Cargo
}

func (ch *Salvage) GetCargpCapacity() int {
	return Unlimited
}

func (ch *Salvage) GetFuel() int {
	return 0
}

func (ch *Salvage) GetFuelCapacity() int {
	return 0
}
