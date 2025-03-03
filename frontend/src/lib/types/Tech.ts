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
 * @param hcType - The type of the hull slot to check.
 * @param slotType - The type of the slot to check.
 * @returns `true` if the slot can be filled with a weapon, `false` otherwise.
 */
export function canFillSlot(hcType: HullSlotType, type: HullSlotType): boolean {
	return (hcType & type) > 0;
}

/**
 * Checks if a given {@linkcode TechHull} is allowed to use a given {@linkcode TechHullComponent}.
 * @param hull - The {@linkcode TechHull|hull} being checked.
 * @param hc - The {@linkcode TechHullComponent|hull component} to be used.
 * @returns `true` if the hull is allowed to use the component, `false` otherwise.
 */
export function hullAllowed(hull: TechHull, hc: TechHullComponent): boolean {
	// nullish coaclescing makes this work ("undefined == -1" is always false)
	const hullAllowed = hc.requirements.hullsAllowed?.indexOf(hull.name) != -1
	const hullDenied = hc.requirements.hullsDenied?.indexOf(hull.name) == -1
	const armedWithUnarmedPart = (hc.cloakUnarmedOnly ?? false) &&
		hull.slots.some(slot => canFillSlot(slot.type, HullSlotTypeWeapon));
	return hullAllowed && !hullDenied && !armedWithUnarmedPart;
}

export function getDefenseCoverage(defense: TechDefense, defenses: number): number {
	return 1.0 - Math.pow(1 - defense.defenseCoverage / 100, clamp(defenses, 0, 100));
}

export function getSmartDefenseCoverage(
	defense: TechDefense,
	defenses: number,
	smartDefenseCoverageFactor?: number
): number {
	smartDefenseCoverageFactor ??= 0.5;
	return (
		1.0 -
		Math.pow(
			1 - (defense.defenseCoverage / 100) * smartDefenseCoverageFactor,
			clamp(defenses, 0, 100)
		)
	);
}

export function getCloakPercentForCloakUnits(cloakUnits: number): number {
	if (cloakUnits <= 100) {
		return cloakUnits / 2;
	} else {
		cloakUnits = cloakUnits - 100;
		if (cloakUnits <= 200) {
			return 50 + cloakUnits / 8;
		} else {
			cloakUnits = cloakUnits - 200;
			if (cloakUnits < 312) {
				return 75 + cloakUnits / 24;
			} else {
				cloakUnits = cloakUnits - 312;
				if (cloakUnits <= 512) {
					return 88 + cloakUnits / 64;
				} else if (cloakUnits < 768) {
					return 96;
				} else if (cloakUnits < 1000) {
					return 97;
				} else {
					return 99;
				}
			}
		}
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
