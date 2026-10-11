import {
	CargoSchema,
	ResourceType,
	type Cargo,
	type CargoJson,
	type MineralJson as Mineral
} from '#lib/types/cs-proto.js';
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

// The share of a hold each kind of cargo takes up, as exact percents. An overfull hold is capped at
// 100% in total. The shares aren't rounded, so a full hold fills the whole bar and a small share
// still shows.
export type CargoPercent = {
	ironium: number;
	boranium: number;
	germanium: number;
	colonists: number;
};

export function cargoPercent(cargo: Cargo, capacity: number | undefined): CargoPercent {
	const percent: CargoPercent = { ironium: 0, boranium: 0, germanium: 0, colonists: 0 };
	if (capacity == 0 || capacity == undefined) {
		return percent;
	}

	let remaining = 100;
	for (const type of ['ironium', 'boranium', 'germanium', 'colonists'] as const) {
		percent[type] = Math.min(remaining, Math.max(0, (cargo[type] / capacity) * 100));
		remaining -= percent[type];
	}
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
