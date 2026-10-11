import {
	Prt,
	RaceSchema,
	ResearchCostLevel,
	SpendLeftoverPointsOn,
	type Race
} from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import type { Rules } from '#lib/types/cs-proto.js';
import { defaultRules } from './Rules';
import { Grav, Rad, Temp, type HabType } from './Hab';

export type LRT = number;
/**
 * No LRT; only used for tech requirements
 */
export const LRTNone = 0;
/**
 * Improved Fuel Efficiency
 */
export const IFE: LRT = 1 << (1 - 1);
/**
 * Total Terraforming
 */
export const TT: LRT = 1 << (2 - 1);
/**
 * Advanced Remote Mining
 */
export const ARM: LRT = 1 << (3 - 1);
/**
 * Improved Starbases
 */
export const ISB: LRT = 1 << (4 - 1);
/**
 * Generalized Research
 */
export const GR: LRT = 1 << (5 - 1);
/**
 * Ultimate Recycling
 */
export const UR: LRT = 1 << (6 - 1);
/**
 * No Ramscoop Engines
 */
export const NRSE: LRT = 1 << (7 - 1);
/**
 * Only Basic Remote Mining
 */
export const OBRM: LRT = 1 << (8 - 1);
/**
 * No Advanced Scanners
 */
export const NAS: LRT = 1 << (9 - 1);
/**
 * Low Starting Population
 */
export const LSP: LRT = 1 << (10 - 1);
/**
 * Bleeding Edge Technology
 */
export const BET: LRT = 1 << (11 - 1);
/**
 * Regenerating Shields
 */
export const RS: LRT = 1 << (12 - 1);
/**
 * Mineral Alchemy
 */
export const MA: LRT = 1 << (13 - 1);
/**
 * Cheap Engines
 */
export const CE: LRT = 1 << (14 - 1);

export const lrts = [IFE, TT, ARM, ISB, GR, UR, NRSE, OBRM, NAS, LSP, BET, RS, MA, CE] as const;

export const humanoid = (): Race =>
	create(RaceSchema, {
		name: 'Humanoid',
		pluralName: 'Humanoids',
		spendLeftoverPointsOn: SpendLeftoverPointsOn.SURFACE_MINERALS,
		prt: Prt.JOAT,
		lrts: LRTNone,
		habLow: { grav: 15, temp: 15, rad: 15 },
		habHigh: { grav: 85, temp: 85, rad: 85 },
		growthRate: 15,
		popEfficiency: 10,
		factoryOutput: 10,
		factoryCost: 10,
		numFactories: 10,
		mineOutput: 10,
		mineCost: 5,
		numMines: 10,
		researchCost: {
			energy: ResearchCostLevel.STANDARD,
			weapons: ResearchCostLevel.STANDARD,
			propulsion: ResearchCostLevel.STANDARD,
			construction: ResearchCostLevel.STANDARD,
			electronics: ResearchCostLevel.STANDARD,
			biotechnology: ResearchCostLevel.STANDARD
		}
	});

export const getLabelForPRT = (prt: Prt): string => {
	switch (prt) {
		case Prt.HE:
			return 'Hyper Expansion';
		case Prt.SS:
			return 'Super Stealth';
		case Prt.WM:
			return 'Warmonger';
		case Prt.CA:
			return 'Claim Adjuster';
		case Prt.IS:
			return 'Inner Strength';
		case Prt.SD:
			return 'Space Demolition';
		case Prt.PP:
			return 'Packet Physics';
		case Prt.IT:
			return 'Interstellar Traveler';
		case Prt.AR:
			return 'Alternate Reality';
		case Prt.JOAT:
			return 'Jack of All Trades';
		default:
			return toString();
	}
};

export const getLabelForLRT = (lrt: LRT): string => {
	switch (lrt) {
		case IFE:
			return 'Improved Fuel Efficiency';
		case TT:
			return 'Total Terraforming';
		case ARM:
			return 'Advanced Remote Mining';
		case ISB:
			return 'Improved Starbases';
		case GR:
			return 'Generalized Research';
		case UR:
			return 'Ultimate Recycling';
		case NRSE:
			return 'No Ram Scoop Engines';
		case OBRM:
			return 'Only Basic Remote Mining';
		case NAS:
			return 'No Advanced Scanners';
		case LSP:
			return 'Low Starting Population';
		case BET:
			return 'Bleeding Edge Technology';
		case RS:
			return 'Regenerating Shields';
		case MA:
			return 'Mineral Alchemy';
		case CE:
			return 'Cheap Engines';
		default:
			return toString();
	}
};

export function getHabWidth(race: Race) {
	return {
		grav: (race.habHigh?.grav ?? 0) - (race.habLow?.grav ?? 0),
		temp: (race.habHigh?.temp ?? 0) - (race.habLow?.temp ?? 0),
		rad: (race.habHigh?.rad ?? 0) - (race.habLow?.rad ?? 0)
	};
}

// The number of planets a race can expect to search to find one it can live on, so 8 means "1 in 8 planets".
// Returns 0 if every planet is habitable.
// The rules decide how planet values are spread out. A value with no dropoff range (radiation) is equally
// likely anywhere. One with a dropoff range (gravity, temperature) is rarer near the extremes, so each value
// in the race's range is weighted by how often planets roll it.
export function getHabitablePlanetOdds(race: Race, rules: Rules = defaultRules): number {
	const universe = rules.universeGenerationRules ?? defaultRules.universeGenerationRules;
	const minHab = universe?.minHab ?? 1;
	const maxHab = universe?.maxHab ?? 99;
	// hab values run from one below the lowest planet value to one above the highest
	const scale = maxHab - minHab + 2;
	let chance = 1;
	const add = (immune: boolean, low: number, high: number, dropoff: number) => {
		if (immune || high - low === scale) {
			chance *= scale;
		} else if (dropoff <= 0) {
			chance *= high - low;
		} else {
			// a planet rolls its value in two parts, one of them up to the dropoff range, so there are
			// more ways to roll a middle value than one near either end
			const span = maxHab - minHab - dropoff;
			let weight = 0;
			for (let value = low; value <= high; value++) {
				const roll = value - minHab;
				if (roll >= 0 && roll <= span + dropoff) {
					weight += Math.min(roll, dropoff, span, span + dropoff - roll) + 1;
				}
			}
			chance = Math.trunc((chance * weight * scale) / ((span + 1) * (dropoff + 1)));
		}
	};
	const dropoff = universe?.habDropoffRange;
	add(
		race.immuneGrav ?? false,
		race.habLow?.grav ?? 0,
		race.habHigh?.grav ?? 0,
		dropoff?.grav ?? 0
	);
	add(
		race.immuneTemp ?? false,
		race.habLow?.temp ?? 0,
		race.habHigh?.temp ?? 0,
		dropoff?.temp ?? 0
	);
	add(race.immuneRad ?? false, race.habLow?.rad ?? 0, race.habHigh?.rad ?? 0, dropoff?.rad ?? 0);

	const allPlanets = scale * scale * scale;
	if (chance === allPlanets) {
		return 0;
	}
	chance = Math.max(1, chance);
	return Math.trunc((Math.trunc(chance / 2) + allPlanets) / chance);
}

export function isImmune(race: Race, habType: HabType): boolean {
	switch (habType) {
		case Grav:
			return race.immuneGrav ?? false;
		case Temp:
			return race.immuneTemp ?? false;
		case Rad:
			return race.immuneRad ?? false;
	}

	return false;
}
