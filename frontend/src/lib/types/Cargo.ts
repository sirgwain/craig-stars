import type { Mineral } from './Mineral';

export type Cargo = {
	colonists?: number;
} & Mineral;

export const CargoTypes = {
	Ironium: 'ironium',
	Boranium: 'boranium',
	Germanium: 'germanium',
	Colonists: 'colonists'
} as const;

export type CargoType = (typeof CargoTypes)[keyof typeof CargoTypes];

export class CargoTransferRequest implements Cargo {
	ironium = 0;
	boranium = 0;
	germanium = 0;
	colonists = 0;
	fuel = 0;

	constructor(cargo?: Cargo, fuel?: number) {
		Object.assign(this, cargo);
		this.fuel = fuel ?? 0;
	}

	// return a new CargoTransferRequest from this one, but negated
	negative(): CargoTransferRequest {
		return new CargoTransferRequest(negativeCargo(this), -this.fuel);
	}

	// return the absolute size of this transfer request
	absoluteSize(): number {
		return (
			Math.abs(this.ironium) +
			Math.abs(this.boranium) +
			Math.abs(this.germanium) +
			Math.abs(this.colonists) +
			Math.abs(this.fuel)
		);
	}

	// return the absolute size of this transfer request
	absoluteCargoSize(): number {
		return (
			Math.abs(this.ironium) +
			Math.abs(this.boranium) +
			Math.abs(this.germanium) +
			Math.abs(this.colonists)
		);
	}

	// add a CargoTransferRequest to another CargoTransferRequest, returning a new one
	add(c: CargoTransferRequest) {
		return new CargoTransferRequest(
			{
				ironium: this.ironium + c.ironium,
				boranium: this.boranium + c.boranium,
				germanium: this.germanium + c.germanium,
				colonists: this.colonists + c.colonists
			},
			this.fuel + c.fuel
		);
	}
}

export const totalCargo = (c: Cargo | undefined) =>
	c ? (c.ironium ?? 0) + (c.boranium ?? 0) + (c.germanium ?? 0) + (c.colonists ?? 0) : 0;

export const emptyCargo = (): Cargo => {
	return {
		ironium: 0,
		boranium: 0,
		germanium: 0,
		colonists: 0
	};
};

// return this cargo with all fields negated
export const negativeCargo = (c: Cargo) => {
	return {
		ironium: -(c.ironium ?? 0),
		boranium: -(c.boranium ?? 0),
		germanium: -(c.germanium ?? 0),
		colonists: -(c.colonists ?? 0)
	};
};

export const add = (c1: Cargo, c2: Cargo) => {
	return {
		ironium: (c1.ironium ?? 0) + (c2.ironium ?? 0),
		boranium: (c1.boranium ?? 0) + (c2.boranium ?? 0),
		germanium: (c1.germanium ?? 0) + (c2.germanium ?? 0),
		colonists: (c1.colonists ?? 0) + (c2.colonists ?? 0)
	};
};

export const subtract = (c1: Cargo, c2: Cargo) => {
	return {
		ironium: (c1.ironium ?? 0) - (c2.ironium ?? 0),
		boranium: (c1.boranium ?? 0) - (c2.boranium ?? 0),
		germanium: (c1.germanium ?? 0) - (c2.germanium ?? 0),
		colonists: (c1.colonists ?? 0) - (c2.colonists ?? 0)
	};
};

export const addMineral = (c1: Cargo, m1: Mineral) => {
	return {
		ironium: (c1.ironium ?? 0) + (m1.ironium ?? 0),
		boranium: (c1.boranium ?? 0) + (m1.boranium ?? 0),
		germanium: (c1.germanium ?? 0) + (m1.germanium ?? 0),
		colonists: c1.colonists ?? 0
	};
};

export function toMineral(cargo: Cargo): Mineral {
	return {
		ironium: cargo.ironium,
		boranium: cargo.boranium,
		germanium: cargo.germanium
	};
}

// if we are displaying cargo as a percent, don't let it be more than 100%
export function cargoPercent(cargo: Cargo, capacity: number | undefined): Cargo {
	if (capacity == 0 || capacity == undefined) {
		return emptyCargo();
	}

	let total = 0;
	const percent: Cargo = {
		ironium: Math.round(((cargo.ironium ?? 0) / capacity) * 100)
	};
	total = percent.ironium ?? 0;
	percent.boranium = Math.min(100 - total, Math.round(((cargo.boranium ?? 0) / capacity) * 100));
	total += percent.boranium ?? 0;
	percent.germanium = Math.min(100 - total, Math.round(((cargo.germanium ?? 0) / capacity) * 100));
	total += percent.germanium ?? 0;
	percent.colonists = Math.min(100 - total, Math.round(((cargo.colonists ?? 0) / capacity) * 100));

	return percent;
}
