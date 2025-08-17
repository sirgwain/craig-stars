import { MineralSchema, type CargoJson, type Mineral, type MineralJson } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

export const totalMinerals = (c: MineralJson | CargoJson | undefined) =>
	c ? (c.ironium ?? 0) + (c.boranium ?? 0) + (c.germanium ?? 0) : 0;

export const addToAll = (m: Mineral, i: number): Mineral =>
	create(MineralSchema, {
		ironium: (m.ironium ?? 0) + i,
		boranium: (m.boranium ?? 0) + i,
		germanium: (m.germanium ?? 0) + i
	});
