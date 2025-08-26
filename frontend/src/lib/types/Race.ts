import {
	Prt,
	RaceSchema,
	ResearchCostLevel,
	SpendLeftoverPointsOn,
	type Race
} from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';
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

export function getHabChance(race: Race): number {
	const habWidth = getHabWidth(race);
	// do a straight calc of hab width, so if we have a hab with widths of 50, 50% of planets will be habitable
	// so we get (.5 * .5 * .5) = .125, or 1 in 8 planets
	const gravChance = race.immuneGrav ? 1.0 : habWidth.grav / 100.0;
	const tempChance = race.immuneTemp ? 1.0 : habWidth.temp / 100.0;
	const radChance = race.immuneRad ? 1.0 : habWidth.rad / 100.0;
	return gravChance * tempChance * radChance;
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
