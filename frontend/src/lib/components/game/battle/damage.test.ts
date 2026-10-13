import { describe, expect, it } from 'vitest';
import { ARMOR_COLOR, SHIELD_COLOR, damageByTarget, damageMarks, formatDamage } from './damage';

describe('battle damage marks', () => {
	it('marks armor, then shields', () => {
		const marks = damageMarks({ damageDoneShields: 10, damageDoneArmor: 30 });
		expect(marks.map((m) => m.kind)).toEqual(['armor', 'shield']);
		expect(marks.map((m) => m.text)).toEqual(['-30', '-10']);
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
		expect(formatDamage(12)).toBe('-12');
	});

	it('adds up the shots of an attack for each square', () => {
		const targets = damageByTarget([
			{ to: { x: 5, y: 5 }, damageDoneShields: 50 },
			{ to: { x: 5, y: 5 }, damageDoneShields: 10, damageDoneArmor: 30 },
			{ to: { x: 5, y: 6 }, damageDoneArmor: 40 },
			// a miss
			{ to: { x: 5, y: 7 } }
		]);
		expect(targets.map((t) => [t.x, t.y])).toEqual([
			[5, 5],
			[5, 6]
		]);
		expect(targets[0].marks.map((m) => m.label)).toEqual(['30 armor damage', '60 shield damage']);
		expect(targets[1].marks.map((m) => m.label)).toEqual(['40 armor damage']);
	});
});
