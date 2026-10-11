import { MapObjectType, WaypointSchema } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { expect, it } from 'vitest';
import { waypointClearance } from './draw';

type Args = Parameters<typeof waypointClearance>;

// Paths stop short of whatever is at a waypoint so they don't cover it. There's nothing at a
// waypoint in deep space, so the path runs straight through it without a gap.
it('leaves a gap in the path at targets but not in deep space', () => {
	const frame = { planetsByNum: new Map() } as unknown as Args[0];
	const view = { iconScale: 1 } as unknown as Args[1];

	const deepSpace = create(WaypointSchema, { position: { x: 10, y: 10 } });
	expect(waypointClearance(frame, view, deepSpace)).toBe(0);

	const atFleet = create(WaypointSchema, {
		position: { x: 10, y: 10 },
		mapObjectTarget: { targetType: MapObjectType.FLEET, targetNum: 1, targetPlayerNum: 2 }
	});
	expect(waypointClearance(frame, view, atFleet)).toBeGreaterThan(0);
});
