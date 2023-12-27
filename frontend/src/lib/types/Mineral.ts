export interface Mineral {
	ironium?: number;
	boranium?: number;
	germanium?: number;
}

export const MineralTypes = {
	Ironium: 0,
	Boranium: 1,
	Germanium: 2
} as const;

export const totalMinerals = (c: Mineral | undefined) =>
	c ? (c.ironium ?? 0) + (c.boranium ?? 0) + (c.germanium ?? 0) : 0;

export const addInt = (m: Mineral, i: number): Mineral => ({
	ironium: (m.ironium ?? 0) + i,
	boranium: (m.boranium ?? 0) + i,
	germanium: (m.germanium ?? 0) + i
});
