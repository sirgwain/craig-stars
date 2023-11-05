export interface Cost {
	ironium?: number;
	boranium?: number;
	germanium?: number;
	resources?: number;
}

// return this cost with a minimum of zero for each value
export function minZero(c: Cost): Cost {
	return {
		ironium: Math.max(c.ironium ?? 0, 0),
		boranium: Math.max(c.boranium ?? 0, 0),
		germanium: Math.max(c.germanium ?? 0, 0),
		resources: Math.max(c.resources ?? 0, 0)
	};
}

// divide two costs
export function divide(a: Cost, b: Cost): number {
	const newIronium = !b.ironium ? Infinity : (a.ironium ?? 0) / (b.ironium ?? 0);
	const newBoranium = !b.boranium ? Infinity : (a.boranium ?? 0) / (b.boranium ?? 0);
	const newGermanium = !b.germanium ? Infinity : (a.germanium ?? 0) / (b.germanium ?? 0);
	const newResources = !b.resources ? Infinity : (a.resources ?? 0) / (b.resources ?? 0);

	return Math.min(newResources, Math.min(newIronium, Math.min(newBoranium, newGermanium)));
}

export function minus(a: Cost, b: Cost): Cost {
	return {
		ironium: (a.ironium ?? 0) - (b.ironium ?? 0),
		boranium: (a.boranium ?? 0) - (b.boranium ?? 0),
		germanium: (a.germanium ?? 0) - (b.germanium ?? 0),
		resources: (a.resources ?? 0) - (b.resources ?? 0)
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
