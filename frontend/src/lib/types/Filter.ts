import type { DesignFinder } from '#lib/services/Universe.js';
import {
	MapObjectType,
	ShipClass,
	WaypointTask,
	type Fleet,
	type Minefield,
	type PlayerNote,
	type ShipDesign
} from './cs-proto';
import type { CommandedPlayer } from './Player';

export type FilterOptions = {
	showIdleFleetsOnly?: boolean;
	filterMyDesigns?: boolean;
	filterAllyDesigns?: boolean;
	filterEnemyDesigns?: boolean;
	// the nums of our designs to show when filtering our designs
	filterDesigns?: number[];
	// the classes of other players' ships to show when filtering them
	filterAllyShipClasses?: ShipClass[];
	filterEnemyShipClasses?: ShipClass[];
	showMinefields?: boolean;
	// whose minefields to show when minefields are shown
	showMinefieldsOf?: MinefieldOwner[];
	showNotes?: boolean;
	// the kinds of note to hide while notes are shown
	hiddenNoteKinds?: string[];
};

// labels placed on the map, and notes without tags, are filtered as kinds of their own
export const NoteKindLabel = 'label';
export const NoteKindUntagged = 'untagged';

// the kind a tag is filtered as. Tags that differ only by case are the same tag.
export function noteKindOfTag(tag: string): string {
	return `tag:${tag.toLowerCase()}`;
}

// the kinds a note belongs to, one for each of its tags
export function noteKinds(note: Pick<PlayerNote, 'targetType' | 'tags'>): string[] {
	if (note.targetType === MapObjectType.UNSPECIFIED) {
		return [NoteKindLabel];
	}
	return note.tags.length ? note.tags.map(noteKindOfTag) : [NoteKindUntagged];
}

// true if this note should be drawn on the map with these filters. A note with several tags is shown
// while any of them is.
export function filterNote(
	note: Pick<PlayerNote, 'targetType' | 'tags'>,
	options: FilterOptions
): boolean {
	if (!(options.showNotes ?? true)) {
		return false;
	}
	const hidden = options.hiddenNoteKinds ?? [];
	return hidden.length === 0 || noteKinds(note).some((kind) => !hidden.includes(kind));
}

// whose minefield it is, from our point of view
export type MinefieldOwner = 'mine' | 'friends' | 'neutrals' | 'enemies';

// the minefield owners that can be shown, in the order they are listed
export const MinefieldOwners: MinefieldOwner[] = ['mine', 'friends', 'neutrals', 'enemies'];

export function getMinefieldOwner(player: CommandedPlayer, minefield: Minefield): MinefieldOwner {
	const playerNum = minefield.mapObject?.playerNum ?? 0;
	if (playerNum === player.num) {
		return 'mine';
	}
	if (player.isFriend(playerNum)) {
		return 'friends';
	}
	return player.isEnemy(playerNum) ? 'enemies' : 'neutrals';
}

// the minefield owners these options show, none if minefields are turned off
export function getShownMinefieldOwners(options: FilterOptions): MinefieldOwner[] {
	return (options.showMinefields ?? true) ? (options.showMinefieldsOf ?? MinefieldOwners) : [];
}

// true if this minefield should be shown with these filters
export function filterMinefield(
	player: CommandedPlayer,
	minefield: Minefield,
	options: FilterOptions
): boolean {
	return getShownMinefieldOwners(options).includes(getMinefieldOwner(player, minefield));
}

// the ship classes that can be filtered, in the order they are listed
export const ShipClasses = Object.values(ShipClass).filter(
	(c): c is ShipClass => typeof c !== 'string' && c !== ShipClass.UNSPECIFIED
);

// the class of a design's hull. Designs from before hulls had classes count as utility ships
export function getShipClass(design: ShipDesign | undefined): ShipClass {
	return design?.spec?.shipClass || ShipClass.UTILITY;
}

// true if this fleet should be shown with these filters
export function filterFleet(
	player: CommandedPlayer,
	fleet: Fleet,
	options: FilterOptions,
	designs: DesignFinder
): boolean {
	if (!filterIdleFleet(player, fleet, options.showIdleFleetsOnly ?? false)) {
		return false;
	}
	return (
		!isDesignFiltered(player, fleet, options) ||
		getVisibleShipCount(player, fleet, options, designs) > 0
	);
}

// true if one of the design filters applies to this fleet's owner
function isDesignFiltered(player: CommandedPlayer, fleet: Fleet, options: FilterOptions): boolean {
	const playerNum = fleet.mapObject?.playerNum;
	if (playerNum === player.num) {
		return options.filterMyDesigns ?? false;
	}
	if (player.isFriend(playerNum)) {
		return options.filterAllyDesigns ?? false;
	}
	return options.filterEnemyDesigns ?? false;
}

// the number of ships in this fleet that the design filters show. Our fleets are filtered
// by design, everyone else's by the class of each ship's hull
export function getVisibleShipCount(
	player: CommandedPlayer,
	fleet: Fleet,
	options: FilterOptions,
	designs: DesignFinder
): number {
	const playerNum = fleet.mapObject?.playerNum;
	let tokens = fleet.tokens;
	if (isDesignFiltered(player, fleet, options)) {
		if (playerNum === player.num) {
			const showDesigns = options.filterDesigns ?? [];
			tokens = tokens.filter((t) => showDesigns.includes(t.designNum));
		} else {
			const showShipClasses =
				(player.isFriend(playerNum)
					? options.filterAllyShipClasses
					: options.filterEnemyShipClasses) ?? [];
			tokens = tokens.filter((t) =>
				showShipClasses.includes(getShipClass(designs.getDesign(playerNum, t.designNum)))
			);
		}
	}
	return tokens.reduce((count, t) => count + t.quantity, 0);
}

// This shows only your fleets that have no movement orders, and any active enemy ships (so you can match one with the other, if you wish).
export function filterIdleFleet(player: CommandedPlayer, fleet: Fleet, enabled: boolean): boolean {
	if (!enabled) {
		// no filter, show all fleets
		return true;
	}

	// show our fleets that are idle
	if (
		fleet.fleetOrders?.waypoints &&
		fleet.fleetOrders?.waypoints.length == 1 &&
		fleet.fleetOrders?.waypoints[0].task == WaypointTask.UNSPECIFIED
	) {
		return true;
	}

	// enemy fleet that is moving, show it so players can match idle fleets to moving fleets
	if (player.num !== fleet.mapObject?.playerNum && fleet.warpSpeed) {
		return true;
	}

	// don't show this fleet if we got here, it's our fleet and moving, or an enemy fleet and idle
	return false;
}
