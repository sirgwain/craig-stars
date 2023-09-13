import { describe, it, expect } from 'vitest';
import { getGravString, getTempString } from './Hab';

describe('hab test', () => {
	it('returns a temp string', () => {
		expect(getTempString(100)).toBe('200°C');
	});

	it('returns a grav string', () => {
		expect(getGravString(0)).toBe('0.25g');
		expect(getGravString(25)).toBe('0.50g');
		expect(getGravString(50)).toBe('1.00g');
		expect(getGravString(75)).toBe('2.00g');
		expect(getGravString(100)).toBe('4.00g');
	});
});
