import { describe, expect, it } from 'vitest';
import { divide } from './Cost';

describe('cost test', () => {
	it('divicdes two costs', () => {
		expect(divide({}, {})).toBe(Infinity);

		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 0, resources: 0 },
				{ ironium: 0, boranium: 0, germanium: 0, resources: 0 }
			)
		).toBe(Infinity);
		expect(
			divide(
				{ ironium: 1, boranium: 0, germanium: 0, resources: 0 },
				{ ironium: 1, boranium: 0, germanium: 0, resources: 0 }
			)
		).toBe(1);
		expect(
			divide(
				{ ironium: 0, boranium: 1, germanium: 0, resources: 0 },
				{ ironium: 0, boranium: 1, germanium: 0, resources: 0 }
			)
		).toBe(1);
		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 1, resources: 0 },
				{ ironium: 0, boranium: 0, germanium: 1, resources: 0 }
			)
		).toBe(1);
		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 0, resources: 1 },
				{ ironium: 0, boranium: 0, germanium: 0, resources: 1 }
			)
		).toBe(1);
		expect(
			divide(
				{ ironium: 2, boranium: 0, germanium: 0, resources: 0 },
				{ ironium: 1, boranium: 0, germanium: 0, resources: 0 }
			)
		).toBe(2);
		expect(
			divide(
				{ ironium: 0, boranium: 2, germanium: 0, resources: 0 },
				{ ironium: 0, boranium: 1, germanium: 0, resources: 0 }
			)
		).toBe(2);
		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 2, resources: 0 },
				{ ironium: 0, boranium: 0, germanium: 1, resources: 0 }
			)
		).toBe(2);
		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 0, resources: 2 },
				{ ironium: 0, boranium: 0, germanium: 0, resources: 1 }
			)
		).toBe(2);
		expect(
			divide(
				{ ironium: 2, boranium: 2, germanium: 2, resources: 2 },
				{ ironium: 1, boranium: 1, germanium: 1, resources: 1 }
			)
		).toBe(2);
		expect(
			divide(
				{ ironium: 1, boranium: 0, germanium: 0, resources: 0 },
				{ ironium: 2, boranium: 0, germanium: 0, resources: 0 }
			)
		).toBe(0.5);
		expect(
			divide(
				{ ironium: 0, boranium: 1, germanium: 0, resources: 0 },
				{ ironium: 0, boranium: 2, germanium: 0, resources: 0 }
			)
		).toBe(0.5);
		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 1, resources: 0 },
				{ ironium: 0, boranium: 0, germanium: 2, resources: 0 }
			)
		).toBe(0.5);
		expect(
			divide(
				{ ironium: 0, boranium: 0, germanium: 0, resources: 1 },
				{ ironium: 0, boranium: 0, germanium: 0, resources: 2 }
			)
		).toBe(0.5);
		expect(
			divide(
				{ ironium: 1, boranium: 1, germanium: 1, resources: 1 },
				{ ironium: 2, boranium: 2, germanium: 2, resources: 2 }
			)
		).toBe(0.5);
	});
});
