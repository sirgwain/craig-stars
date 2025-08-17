import { clamp } from '$lib/services/Math';
import {
	TechCategory,
	TerraformHabType,
	type GetTechsResponse,
	type Tech,
	type TechDefense,
	type TechHull,
	type TechHullComponent
} from '$lib/types/cs-proto';
import { type HullSlotType, HullSlotTypeWeapon } from './Consts';

export type TechLike = {
	tech?: Tech;
};

export type TechStore = GetTechsResponse;

/**
 * Return the "long-form" name of a TerraformHabType, given its abbreviated form
 * @param type The TerraformHabType being expanded
 * @returns The full name of the TerraformHabType
 */
export function getLongHabName(type: TerraformHabType): string {
	switch (type) {
		case TerraformHabType.GRAV:
			return 'Gravity';
		case TerraformHabType.TEMP:
			return 'Temperature';
		case TerraformHabType.RAD:
			return 'Radiation';
		case TerraformHabType.ALL:
			return 'All';
		default:
			return 'None';
	}
}

export const TechCategories: TechCategory[] = [
	TechCategory.UNSPECIFIED,
	TechCategory.ARMOR,
	TechCategory.BEAM_WEAPON,
	TechCategory.BOMB,
	TechCategory.ELECTRICAL,
	TechCategory.ENGINE,
	TechCategory.MECHANICAL,
	TechCategory.MINE_LAYER,
	TechCategory.MINE_ROBOT,
	TechCategory.ORBITAL,
	TechCategory.PLANETARY,
	TechCategory.PLANETARY_SCANNER,
	TechCategory.PLANETARY_DEFENSE,
	TechCategory.SCANNER,
	TechCategory.SHIELD,
	TechCategory.SHIP_HULL,
	TechCategory.STARBASE_HULL,
	TechCategory.TERRAFORMING,
	TechCategory.TORPEDO
];

/**
 * Determine if a tech is a hull component
 * @param category The TechCategory to check
 * @returns true if the tech category contains TechHullComponents
 */
export function isHullComponent(category: TechCategory | undefined): boolean {
	switch (category) {
		case TechCategory.ARMOR:
		case TechCategory.BEAM_WEAPON:
		case TechCategory.BOMB:
		case TechCategory.ELECTRICAL:
		case TechCategory.ENGINE:
		case TechCategory.MECHANICAL:
		case TechCategory.MINE_LAYER:
		case TechCategory.MINE_ROBOT:
		case TechCategory.ORBITAL:
		case TechCategory.SCANNER:
		case TechCategory.TORPEDO:
		case TechCategory.SHIELD:
			return true;
		case TechCategory.PLANETARY:
		case TechCategory.PLANETARY_SCANNER:
		case TechCategory.PLANETARY_DEFENSE:
		case TechCategory.SHIP_HULL:
		case TechCategory.STARBASE_HULL:
		case TechCategory.TERRAFORMING:
			return false;
		default:
			return false;
	}
}

/** check if this tech is a hull
 * @param tech The tech to check
 * @returns true if this tech is defined and is a ship hull; false otherwise
 */
export function isHull(tech: TechLike | undefined): boolean {
	if (!tech) {
		return false;
	}
	return [TechCategory.SHIP_HULL, TechCategory.STARBASE_HULL].includes(
		tech.tech?.category ?? TechCategory.UNSPECIFIED
	);
}

/**
 * Checks if the {@linkcode HullSlotType} of a given {@linkcode TechHullSlot}
   can be filled with an item of another HullSlotType.
 * @param hcType - The {@linkcode HullSlotType} to check against.
 * @param slotType - The {@linkcode HullSlotType} of the item being placed.
 * @returns `true` if the two slots are compatible, `false` otherwise.
 */
export function canFillSlot(
	hcType: HullSlotType | undefined,
	type: HullSlotType | undefined
): boolean {
	if (!hcType || !type) {
		return false;
	}
	return (hcType & type) > 0;
}

/**
 * Checks if a given {@linkcode TechHull} is allowed to use a given {@linkcode TechHullComponent}.
 * @param hull - The {@linkcode TechHull} being checked.
 * @param hc - The {@linkcode TechHullComponent} to be used.
 * @returns `true` if the hull can use the component, `false` otherwise.
 */
export function hullAllowed(hull: TechHull, hc: TechHullComponent): boolean {
	/* nullish coaclescing makes this work - undefined is never equal to -1,
	so a null array will be treated as allowing all or denying no hulls
	*/
	const hullAllowed =
		(hc.tech?.requirements?.hullsAllowed?.length ?? 0) === 0 ||
		hc.tech?.requirements?.hullsAllowed?.indexOf(hull.tech?.name ?? '') != -1;
	const hullDenied =
		(hc.tech?.requirements?.hullsDenied?.length ?? 0) > 0 &&
		hc.tech?.requirements?.hullsDenied?.indexOf(hull.tech?.name ?? '') == -1;
	const armedWithUnarmedPart =
		(hc.cloakUnarmedOnly ?? false) &&
		hull.slots.some((slot) => canFillSlot(slot.type, HullSlotTypeWeapon));
	return hullAllowed && !hullDenied && !armedWithUnarmedPart;
}

export function getDefenseCoverage(defense: TechDefense, defenses: number): number {
	return 1.0 - Math.pow(1 - defense.defenseCoverage / 100, clamp(defenses, 0, 100));
}

export function getSmartDefenseCoverage(
	defense: TechDefense,
	numDefenses: number,
	smartDefenseCoverageFactor: number = 0.5
): number {
	return (
		1.0 -
		Math.pow(
			1 - (defense.defenseCoverage / 100) * smartDefenseCoverageFactor,
			clamp(numDefenses, 0, 100)
		)
	);
}

export function getCloakPercentForCloakUnits(cloakUnits: number): number {
	if (cloakUnits <= 100) {
		return cloakUnits / 2;
	}
	cloakUnits -= 100;
	if (cloakUnits <= 200) {
		return 50 + cloakUnits / 8;
	}
	cloakUnits -= 200;
	if (cloakUnits < 312) {
		return 75 + cloakUnits / 24;
	}
	cloakUnits -= 512;
	switch (true) {
		case cloakUnits <= 512:
			return 88 + cloakUnits / 64;
		case cloakUnits < 768:
			return 96;
		case cloakUnits < 1000:
			return 97;
		default:
			return 98;
	}
}
