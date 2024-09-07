import { describe, expect, it } from 'vitest';
import { andCommaList } from './andCommandList';

describe('andCommandList test', () => {
	it('combines strings', () => {
		expect(andCommaList(['word1'])).toBe('word1');
		expect(andCommaList(['word1', 'word2'])).toBe('word1 and word2');
		expect(andCommaList(['word1', 'word2', 'word3'])).toBe('word1, word2 and word3');
		expect(andCommaList(['word1', '', 'word3'])).toBe('word1 and word3');
	});
});
