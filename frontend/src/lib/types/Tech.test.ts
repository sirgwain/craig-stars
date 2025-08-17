import { describe, expect, it } from 'vitest';
import { canFillSlot } from './Tech';
import { HullSlotTypeArmor, HullSlotTypeEngine, HullSlotTypeGeneral } from './Consts';

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
