import techjson from '$lib/ssr/techs.json';
import { describe, expect, it } from 'vitest';
import { Player } from './Player';
import { LRT } from './Race';
import {
	HullSlotType,
	TerraformHabTypes,
	canFillSlot,
	getBestTerraform,
	type TechStore
} from './Tech';

describe('tech engine test', () => {
	it('engine can fill engine slot', () => {
		expect(canFillSlot(HullSlotType.Engine, HullSlotType.Engine)).toBeTruthy();
	});
	it('engine cannot fill general slot', () => {
		expect(canFillSlot(HullSlotType.Engine, HullSlotType.General)).toBeFalsy();
	});
	it('armor cannot fill engine slot', () => {
		expect(canFillSlot(HullSlotType.Armor, HullSlotType.Engine)).toBeFalsy();
	});
});

describe('tech terraform test', () => {
	const techStore = techjson as TechStore;
	it('get best terraform - undefined', () => {
		const player = new Player();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Gravity)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Temperature)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Radiation)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.All)).toBeUndefined();
	});
	it('get best terraform - gravity', () => {
		const player = new Player();
		player.techLevels.propulsion = 1;
		player.techLevels.biotechnology = 1;
		const gravityTerraform3 = techStore.terraforms.find((t) => t.name === 'Gravity Terraform ±3');
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Gravity)).toBe(gravityTerraform3);
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Temperature)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Radiation)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.All)).toBeUndefined();
	});
	it('get best terraform - temp 11', () => {
		const player = new Player();
		player.techLevels.energy = 10;
		player.techLevels.biotechnology = 3;
		const tempTerraform11 = techStore.terraforms.find((t) => t.name === 'Temp Terraform ±11');
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Gravity)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Temperature)).toBe(
			tempTerraform11
		);
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Radiation)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.All)).toBeUndefined();
	});
	it('get best terraform - total 5', () => {
		const player = new Player();
		player.techLevels.biotechnology = 3;
		player.race.lrts |= LRT.TT;
		const totalTerraform5 = techStore.terraforms.find((t) => t.name === 'Total Terraform ±5');
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Gravity)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Temperature)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.Radiation)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypes.All)).toBe(totalTerraform5);
	});
});
