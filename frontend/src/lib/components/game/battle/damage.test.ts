import { describe, expect, it } from 'vitest';
import { ARMOR_COLOR, SHIELD_COLOR, damageMarks, formatDamage } from './damage';

describe('battle damage marks', () => {
	it('marks armor, shields and kills in a fixed order', () => {
		const marks = damageMarks({ damageDoneShields: 10, damageDoneArmor: 30, tokensDestroyed: 2 });
		expect(marks.map((m) => m.kind)).toEqual(['armor', 'shield', 'destroyed']);
		expect(marks.map((m) => m.text)).toEqual(['-30', '-10', '2 destroyed']);
	});

	it('tells shield damage from armor damage', () => {
		const [shield] = damageMarks({ damageDoneShields: 50 });
		const [armor] = damageMarks({ damageDoneArmor: 50 });
		expect(shield.kind).toBe('shield');
		expect(armor.kind).toBe('armor');
		expect(shield.color).toBe(SHIELD_COLOR);
		expect(armor.color).toBe(ARMOR_COLOR);
		expect(shield.color).not.toBe(armor.color);
		expect(shield.label).toBe('50 shield damage');
		expect(armor.label).toBe('50 armor damage');
	});

	it('has no marks for a miss or a move', () => {
		expect(damageMarks({})).toEqual([]);
		expect(damageMarks(undefined)).toEqual([]);
	});

	it('formats damage', () => {
		expect(formatDamage('armor', 12)).toBe('-12');
		expect(formatDamage('destroyed', 1)).toBe('destroyed');
	});
});
