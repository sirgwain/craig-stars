import { describe, expect, it } from 'vitest';
import { andCommaList } from './andCommaList';

describe('andCommaList', () => {
	it('combines strings', () => {
		expect(andCommaList(['word1'])).toBe('word1');
		expect(andCommaList(['word1', 'word2'])).toBe('word1 and word2');
		expect(andCommaList(['word1', 'word2', 'word3'])).toBe('word1, word2 and word3');
		expect(andCommaList(['word1', '999', 'word3'])).toBe('word1, 999 and word3');
	});
	it('skips falsy strings and undefined arrays', () => {
		expect(andCommaList([''])).toBe('');
		expect(andCommaList(undefined)).toBe('');
		expect(andCommaList(['', 'apple yay'])).toBe('apple yay');
		expect(andCommaList(['aeeee', '', 'apple yay'])).toBe('aeeee and apple yay');
	});
	it('supports custom operators', () => {
		expect(andCommaList(['', '', ''], 'default text')).toBe('nothing');
		expect(andCommaList(undefined, 'default text', 'or')).toBe('default text');
		expect(
			andCommaList(['granny smith', 'honeycrisp', 'fuji', 'gala', 'pink lady'], '', 'or')
		).toBe('granny smith, honeycrisp, fuji, gala or pink lady');
		expect(andCommaList(['aeeee', '', 'apple yay'], '', 'and/or maaaaybe')).toBe(
			'aeeee and/or maaaaybe apple yay'
		);
	});
});
