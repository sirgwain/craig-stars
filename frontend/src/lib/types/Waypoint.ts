import type { Universe } from '#lib/services/Universe.js';
import { MapObjectType, type Vector, type Waypoint } from './cs-proto.js';

// Only use the player's known wormhole links. Unknown or missing exits leave estimates
// starting at the entrance, just as the scanner path does.
export function getWormholeExit(
	universe: Pick<Universe, 'getWormhole'>,
	waypoint: Waypoint
): Vector | undefined {
	const target = waypoint.mapObjectTarget;
	if (target?.targetType !== MapObjectType.WORMHOLE || !target.targetNum) {
		return undefined;
	}
	const destinationNum = universe.getWormhole(target.targetNum)?.destinationNum;
	return destinationNum ? universe.getWormhole(destinationNum)?.mapObject?.position : undefined;
}
