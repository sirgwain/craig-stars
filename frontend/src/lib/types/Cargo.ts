export interface Cargo {
	ironium?: number;
	boranium?: number;
	germanium?: number;
	colonists?: number;
}

export const totalCargo = (c: Cargo) =>
	(c.ironium ?? 0) + (c.boranium ?? 0) + (c.germanium ?? 0) + (c.colonists ?? 0);

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
