import { MapObjectType } from '#lib/types/cs-proto.js';
import type { MapObjectLike } from '#lib/types/MapObject.js';

// which of the two the player picked most recently
export type LastPicked = 'commanded' | 'selected';

// identifies a mapobject, so we can tell picking a different one from a refresh of the same one
export function summaryKey(mo: MapObjectLike | undefined): string {
	const o = mo?.mapObject;
	return o ? `${o.type}/${o.playerNum}/${o.num}/${o.position?.x},${o.position?.y}` : '';
}

/**
 * The phone layout has room for one summary, so it has to choose between what's selected on the
 * map and what's commanded.
 *
 * Commanding a fleet doesn't select it. Tapping a planet cycles through the fleets there by
 * commanding each in turn while the planet stays selected, and commanding a fleet with a route
 * moves the selection to the end of the route. So we show a commanded fleet until the player picks
 * something else on the map, like the planet they cycle back to.
 */
export function pickSummaryMapObject(
	selected: MapObjectLike | undefined,
	commanded: MapObjectLike | undefined,
	addWaypoint: boolean,
	lastPicked: LastPicked
): MapObjectLike | undefined {
	if (!selected) {
		return commanded;
	}

	// don't change the summary while adding waypoints
	if (addWaypoint) {
		return commanded;
	}

	// if we select a fleet, put it in the summary
	if (commanded && selected.mapObject?.type === MapObjectType.FLEET) {
		return selected;
	}

	// the player just commanded this fleet
	if (commanded?.mapObject?.type === MapObjectType.FLEET && lastPicked === 'commanded') {
		return commanded;
	}

	return selected;
}
