import type { CostJson as Cost } from '$lib/types/cs-proto';

export const emptyCost: Readonly<Cost> = {
	ironium: 0,
	boranium: 0,
	germanium: 0,
	resources: 0
};

// divide two costs, returning a number of how many times b goes into a
// this is used for calculating how many of an item we can build
export function divide(a: Cost, b: Cost): number {
	const newIronium = !b.ironium ? Infinity : (a.ironium ?? 0) / (b.ironium ?? 0);
	const newBoranium = !b.boranium ? Infinity : (a.boranium ?? 0) / (b.boranium ?? 0);
	const newGermanium = !b.germanium ? Infinity : (a.germanium ?? 0) / (b.germanium ?? 0);
	const newResources = !b.resources ? Infinity : (a.resources ?? 0) / (b.resources ?? 0);

	return Math.min(newResources, Math.min(newIronium, Math.min(newBoranium, newGermanium)));
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
