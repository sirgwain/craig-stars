import { describe, it, expect } from 'vitest';
import { getTempString } from './Hab';

describe('hab test', () => {
	it('returns a temp string', () => {
		expect(getTempString(100)).toBe('200°C');
	});
});
