import {
	CargoSchema,
	ResourceType,
	type Cargo,
	type CargoJson,
	type MineralJson as Mineral
} from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

export type CargoType = ResourceType;

export const totalCargo = (c: Cargo | CargoJson | undefined) =>
	c ? (c.ironium ?? 0) + (c.boranium ?? 0) + (c.germanium ?? 0) + (c.colonists ?? 0) : 0;

export const emptyCargo = (): Cargo => {
	return create(CargoSchema, emptyCargoJson());
};

export const emptyCargoJson = (): CargoJson => {
	return {
		ironium: 0,
		boranium: 0,
		germanium: 0,
		colonists: 0
	};
};

// return this cargo with all fields negated
export const negativeCargo = (c: Cargo): Cargo => {
	return create(CargoSchema, {
		ironium: -(c.ironium ?? 0),
		boranium: -(c.boranium ?? 0),
		germanium: -(c.germanium ?? 0),
		colonists: -(c.colonists ?? 0)
	});
};

export const add = (c1: Cargo, c2: Cargo | undefined): Cargo => {
	return create(CargoSchema, {
		ironium: (c1.ironium ?? 0) + (c2?.ironium ?? 0),
		boranium: (c1.boranium ?? 0) + (c2?.boranium ?? 0),
		germanium: (c1.germanium ?? 0) + (c2?.germanium ?? 0),
		colonists: (c1.colonists ?? 0) + (c2?.colonists ?? 0)
	});
};

export const subtract = (c1: Cargo, c2: Cargo): Cargo => {
	return create(CargoSchema, {
		ironium: (c1.ironium ?? 0) - (c2.ironium ?? 0),
		boranium: (c1.boranium ?? 0) - (c2.boranium ?? 0),
		germanium: (c1.germanium ?? 0) - (c2.germanium ?? 0),
		colonists: (c1.colonists ?? 0) - (c2.colonists ?? 0)
	});
};

export const addMineral = (c1: Cargo, m1: Mineral): Cargo => {
	return create(CargoSchema, {
		ironium: (c1.ironium ?? 0) + (m1.ironium ?? 0),
		boranium: (c1.boranium ?? 0) + (m1.boranium ?? 0),
		germanium: (c1.germanium ?? 0) + (m1.germanium ?? 0),
		colonists: c1.colonists ?? 0
	});
};

export function toMineral(cargo: Cargo): Mineral {
	return {
		ironium: cargo.ironium,
		boranium: cargo.boranium,
		germanium: cargo.germanium
	};
}

export function population(cargo: Cargo | undefined): number {
	return (cargo?.colonists ?? 0) * 100;
}

// if we are displaying cargo as a percent, don't let it be more than 100%
export function cargoPercent(cargo: Cargo, capacity: number | undefined): Cargo {
	if (capacity == 0 || capacity == undefined) {
		return emptyCargo();
	}

	const percent: Cargo = create(CargoSchema, {
		ironium: Math.round(((cargo.ironium ?? 0) / capacity) * 100)
	});
	let total = percent.ironium ?? 0;
	percent.boranium = Math.min(100 - total, Math.round(((cargo.boranium ?? 0) / capacity) * 100));
	total += percent.boranium ?? 0;
	percent.germanium = Math.min(100 - total, Math.round(((cargo.germanium ?? 0) / capacity) * 100));
	total += percent.germanium ?? 0;
	percent.colonists = Math.min(100 - total, Math.round(((cargo.colonists ?? 0) / capacity) * 100));

	return percent;
}

export function cargoDescription(cargoType: CargoType, amount: number | undefined): string {
	switch (cargoType) {
		case ResourceType.FUEL:
			return `${(amount ?? 0).toLocaleString()}mg`;
		case ResourceType.COLONISTS:
			return `${((amount ?? 0) * 100).toLocaleString()}`;
		default:
			return `${(amount ?? 0).toLocaleString()}kt`;
	}
}
