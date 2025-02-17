import techjson from '$lib/ssr/techs.json';
import { describe, expect, it } from 'vitest';
import { CommandedPlayer } from './Player';
import { canFillSlot, getBestTerraform } from './Tech';
import {
	HullSlotTypeArmor,
	HullSlotTypeEngine,
	HullSlotTypeGeneral,
	TerraformHabTypeAll,
	TerraformHabTypeGrav,
	TerraformHabTypeRad,
	TerraformHabTypeTemp,
	TT,
	type TechStore
} from './cs';

describe('tech engine test', () => {
	it('engine can fill engine slot', () => {
		expect(canFillSlot(HullSlotTypeEngine, HullSlotTypeEngine)).toBeTruthy();
	});
	it('engine cannot fill general slot', () => {
		expect(canFillSlot(HullSlotTypeEngine, HullSlotTypeGeneral)).toBeFalsy();
	});
	it('armor cannot fill engine slot', () => {
		expect(canFillSlot(HullSlotTypeArmor, HullSlotTypeEngine)).toBeFalsy();
	});
});

describe('tech terraform test', () => {
	const techStore = techjson as unknown as TechStore;
	it('get best terraform - undefined', () => {
		const player = new CommandedPlayer();
		expect(getBestTerraform(techStore, player, TerraformHabTypeGrav)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeTemp)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeRad)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeAll)).toBeUndefined();
	});
	it('get best terraform - gravity', () => {
		const player = new CommandedPlayer();
		player.techLevels.propulsion = 1;
		player.techLevels.biotechnology = 1;
		const gravityTerraform3 = techStore.terraforms.find((t) => t.name === 'Gravity Terraform ±3');
		expect(getBestTerraform(techStore, player, TerraformHabTypeGrav)).toBe(gravityTerraform3);
		expect(getBestTerraform(techStore, player, TerraformHabTypeTemp)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeRad)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeAll)).toBeUndefined();
	});
	it('get best terraform - temp 11', () => {
		const player = new CommandedPlayer();
		player.techLevels.energy = 10;
		player.techLevels.biotechnology = 3;
		const tempTerraform11 = techStore.terraforms.find((t) => t.name === 'Temp Terraform ±11');
		expect(getBestTerraform(techStore, player, TerraformHabTypeGrav)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeTemp)).toBe(tempTerraform11);
		expect(getBestTerraform(techStore, player, TerraformHabTypeRad)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeAll)).toBeUndefined();
	});
	it('get best terraform - total 5', () => {
		const player = new CommandedPlayer();
		player.techLevels.biotechnology = 3;
		player.race.lrts |= TT;
		const totalTerraform5 = techStore.terraforms.find((t) => t.name === 'Total Terraform ±5');
		expect(getBestTerraform(techStore, player, TerraformHabTypeGrav)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeTemp)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeRad)).toBeUndefined();
		expect(getBestTerraform(techStore, player, TerraformHabTypeAll)).toBe(totalTerraform5);
	});
});
