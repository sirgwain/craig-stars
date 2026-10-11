import { MapObjectType } from '#lib/types/cs-proto.js';
import type { MapObjectLike } from '#lib/types/MapObject.js';
import { describe, expect, it } from 'vitest';
import { pickSummaryMapObject, summaryKey } from './CommandPaneSummary';

const at = (type: MapObjectType, num: number, x: number) =>
	({ mapObject: { type, num, playerNum: 1, position: { x, y: 0 } } }) as MapObjectLike;
const planet = at(MapObjectType.PLANET, 5, 100);
const scout = at(MapObjectType.FLEET, 6, 100);
const elsewhere = at(MapObjectType.PLANET, 9, 300);

// On a phone, tapping a planet with one of our fleets in orbit cycles between the two. The planet
// stays selected on the map the whole time, so the summary goes by what was picked last.
describe('pickSummaryMapObject', () => {
	it('shows the fleet the player just cycled to', () => {
		expect(pickSummaryMapObject(planet, scout, false, 'commanded')).toBe(scout);
	});

	it('shows the planet the player cycled back to, even with our fleet in orbit', () => {
		expect(pickSummaryMapObject(planet, scout, false, 'selected')).toBe(planet);
	});

	it('shows a selected fleet, and anything the player selects somewhere else', () => {
		expect(pickSummaryMapObject(scout, scout, false, 'selected')).toBe(scout);
		expect(pickSummaryMapObject(elsewhere, scout, false, 'selected')).toBe(elsewhere);
	});

	// commanding a fleet with a route moves the selection to the end of the route
	it('shows a fleet that was just commanded, though the selection moved to its route end', () => {
		expect(pickSummaryMapObject(elsewhere, scout, false, 'commanded')).toBe(scout);
	});

	it('shows the selection when a planet is commanded', () => {
		expect(pickSummaryMapObject(elsewhere, planet, false, 'commanded')).toBe(elsewhere);
	});

	it('keeps showing the commanded fleet while adding waypoints', () => {
		expect(pickSummaryMapObject(elsewhere, scout, true, 'selected')).toBe(scout);
	});

	it('shows what is commanded when nothing is selected', () => {
		expect(pickSummaryMapObject(undefined, scout, false, 'selected')).toBe(scout);
	});
});

it('gives a refreshed copy of a mapobject the same key', () => {
	expect(summaryKey(at(MapObjectType.FLEET, 6, 100))).toBe(summaryKey(scout));
	expect(summaryKey(planet)).not.toBe(summaryKey(scout));
	expect(summaryKey(undefined)).toBe('');
});
