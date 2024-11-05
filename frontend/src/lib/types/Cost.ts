import type { Mineral } from './Mineral';

export type Cost = {
	resources?: number;
} & Mineral;

export const emptyCost: Readonly<Cost> = {
	ironium: 0,
	boranium: 0,
	germanium: 0,
	resources: 0
};

// return this cost with a minimum of zero for each value
export function minZero(c: Cost): Cost {
	return {
		ironium: Math.max(c.ironium ?? 0, 0),
		boranium: Math.max(c.boranium ?? 0, 0),
		germanium: Math.max(c.germanium ?? 0, 0),
		resources: Math.max(c.resources ?? 0, 0)
	};
}

// divide two costs, returning a number of how many times b goes into a
// this is used for calculating how many of an item we can build
export function divide(a: Cost, b: Cost): number {
	const newIronium = !b.ironium ? Infinity : (a.ironium ?? 0) / (b.ironium ?? 0);
	const newBoranium = !b.boranium ? Infinity : (a.boranium ?? 0) / (b.boranium ?? 0);
	const newGermanium = !b.germanium ? Infinity : (a.germanium ?? 0) / (b.germanium ?? 0);
	const newResources = !b.resources ? Infinity : (a.resources ?? 0) / (b.resources ?? 0);

	return Math.min(newResources, Math.min(newIronium, Math.min(newBoranium, newGermanium)));
}

export function add(a: Cost, b: Cost | Mineral | undefined): Cost {
	return {
		ironium: (a.ironium ?? 0) + (b?.ironium ?? 0),
		boranium: (a.boranium ?? 0) + (b?.boranium ?? 0),
		germanium: (a.germanium ?? 0) + (b?.germanium ?? 0),
		resources: (a.resources ?? 0) + (b && 'resources' in b ? (b.resources ?? 0) : 0)
	};
}

export function subtract(a: Cost, b: Cost | Mineral): Cost {
	return {
		ironium: (a.ironium ?? 0) - (b.ironium ?? 0),
		boranium: (a.boranium ?? 0) - (b.boranium ?? 0),
		germanium: (a.germanium ?? 0) - (b.germanium ?? 0),
		resources: (a.resources ?? 0) - ('resources' in b ? (b.resources ?? 0) : 0)
	};
}

export function multiply(cost: Cost, scalar: number): Cost {
	return {
		ironium: Math.floor((cost.ironium ?? 0) * scalar),
		boranium: Math.floor((cost.boranium ?? 0) * scalar),
		germanium: Math.floor((cost.germanium ?? 0) * scalar),
		resources: Math.floor((cost.resources ?? 0) * scalar)
	};
}

export function total(cost?: Cost): number {
	return cost
		? (cost.ironium ?? 0) + (cost.boranium ?? 0) + (cost.germanium ?? 0) + (cost.resources ?? 0)
		: 0;
}

export function totalMinerals(cost?: Cost): number {
	return cost ? (cost.ironium ?? 0) + (cost.boranium ?? 0) + (cost.germanium ?? 0) : 0;
}

export function numBuildable(available: Cost, cost: Cost): number {
	const buildable: Cost = {
		ironium: Number.MAX_SAFE_INTEGER,
		boranium: Number.MAX_SAFE_INTEGER,
		germanium: Number.MAX_SAFE_INTEGER,
		resources: Number.MAX_SAFE_INTEGER
	};

	if (cost.ironium && cost.ironium > 0) {
		buildable.ironium = Math.floor((available.ironium ?? 0) / cost.ironium);
	}
	if (cost.boranium && cost.boranium > 0) {
		buildable.boranium = Math.floor((available.boranium ?? 0) / cost.boranium);
	}
	if (cost.germanium && cost.germanium > 0) {
		buildable.germanium = Math.floor((available.germanium ?? 0) / cost.germanium);
	}
	if (cost.resources && cost.resources > 0) {
		buildable.resources = Math.floor((available.resources ?? 0) / cost.resources);
	}

	return Math.min(
		buildable.ironium ?? 0,
		buildable.boranium ?? 0,
		buildable.germanium ?? 0,
		buildable.resources ?? 0
	);
}
