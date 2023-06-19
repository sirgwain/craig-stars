import { describe, expect, it } from 'vitest';
import { HullSlotType, canFillSlot } from './Tech';

describe('tech test', () => {
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
