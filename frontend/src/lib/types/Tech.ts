import { clamp } from '$lib/services/Math';
import {
	HullSlotTypeWeapon,
	TechCategoryArmor,
	TechCategoryBeamWeapon,
	TechCategoryBomb,
	TechCategoryElectrical,
	TechCategoryEngine,
	TechCategoryMechanical,
	TechCategoryMineLayer,
	TechCategoryMineRobot,
	TechCategoryNone,
	TechCategoryOrbital,
	TechCategoryPlanetary,
	TechCategoryPlanetaryDefense,
	TechCategoryPlanetaryScanner,
	TechCategoryScanner,
	TechCategoryShield,
	TechCategoryShipHull,
	TechCategoryStarbaseHull,
	TechCategoryTerraforming,
	TechCategoryTorpedo,
	TerraformHabTypeAll,
	TerraformHabTypeGrav,
	TerraformHabTypeNone,
	TerraformHabTypeRad,
	TerraformHabTypeTemp,
	type HullSlotType,
	type Tech,
	type TechCategory,
	type TechDefense,
	type TechHull,
	type TechHullComponent,
	type TechStore,
	type TechTerraform,
	type TerraformHabType
} from './cs';
import type { CommandedPlayer } from './Player';

export const TerraformHabTypes = [
	TerraformHabTypeNone,
	TerraformHabTypeGrav,
	TerraformHabTypeTemp,
	TerraformHabTypeRad,
	TerraformHabTypeAll
];

/**
 * Return the "long-form" name of a TerraformHabType, given its abbreviated form
 * @param type The TerraformHabType being expanded
 * @returns The full name of the TerraformHabType
 */
export function getLongHabName(type: TerraformHabType): string {
	switch (type) {
		case TerraformHabTypeGrav:
			return 'Gravity';
		case TerraformHabTypeTemp:
			return 'Temperature';
		case TerraformHabTypeRad:
			return 'Radiation';
		case TerraformHabTypeAll:
			return 'All';
		default:
			return 'None';
	}
}

export const TechCategories: TechCategory[] = [
	TechCategoryNone,
	TechCategoryArmor,
	TechCategoryBeamWeapon,
	TechCategoryBomb,
	TechCategoryElectrical,
	TechCategoryEngine,
	TechCategoryMechanical,
	TechCategoryMineLayer,
	TechCategoryMineRobot,
	TechCategoryOrbital,
	TechCategoryPlanetary,
	TechCategoryPlanetaryScanner,
	TechCategoryPlanetaryDefense,
	TechCategoryScanner,
	TechCategoryShield,
	TechCategoryShipHull,
	TechCategoryStarbaseHull,
	TechCategoryTerraforming,
	TechCategoryTorpedo
];

/**
 * Determine if a tech is a hull component
 * @param category The TechCategory to check
 * @returns true if the tech category contains TechHullComponents
 */
export function isHullComponent(category: TechCategory | undefined): boolean {
	switch (category) {
		case TechCategoryArmor:
		case TechCategoryBeamWeapon:
		case TechCategoryBomb:
		case TechCategoryElectrical:
		case TechCategoryEngine:
		case TechCategoryMechanical:
		case TechCategoryMineLayer:
		case TechCategoryMineRobot:
		case TechCategoryOrbital:
		case TechCategoryScanner:
		case TechCategoryTorpedo:
		case TechCategoryShield:
			return true;
		case TechCategoryPlanetary:
		case TechCategoryPlanetaryScanner:
		case TechCategoryPlanetaryDefense:
		case TechCategoryShipHull:
		case TechCategoryStarbaseHull:
		case TechCategoryTerraforming:
			return false;
		default:
			return false;
	}
}

/** check if this tech is a hull
 * @param tech The tech to check
 * @returns true if this tech is defined and is a ship hull; false otherwise
 */
export function isHull(tech: Tech | undefined): boolean {
	if (!tech) {
		return false;
	}
	return [TechCategoryShipHull, TechCategoryStarbaseHull].includes(tech.category);
}

/**
 * Checks if the {@linkcode HullSlotType} of a given {@linkcode TechHullSlot}
   can be filled with an item of another HullSlotType.
 * @param hcType - The {@linkcode HullSlotType} to check against.
 * @param slotType - The {@linkcode HullSlotType} of the item being placed.
 * @returns `true` if the two slots are compatible, `false` otherwise.
 */
export function canFillSlot(hcType: HullSlotType, type: HullSlotType): boolean {
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
	const hullAllowed = hc.requirements.hullsAllowed?.indexOf(hull.name) != -1;
	const hullDenied = hc.requirements.hullsDenied?.indexOf(hull.name) == -1;
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

export function getBestTerraform(
	techStore: TechStore,
	player: CommandedPlayer,
	habType: TerraformHabType
): TechTerraform | undefined {
	// get the best terraform for a given type, sorted largest to smallest ranking
	return techStore.terraforms
		.filter((t) => player.hasTech(t) && t.habType == habType)
		.sort((a, b) => (b.ranking ?? 0) - (a.ranking ?? 0))[0];
}
