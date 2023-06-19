import type { Fleet } from './Fleet';
import type { Planet } from './Planet';
import type { Race } from './Race';
import type { ShipDesign } from './ShipDesign';
import type {
	Tech,
	TechDefense,
	TechPlanetaryScanner
} from './Tech';

export type Player = {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	gameId: number;
	userId?: number;
	num: number;
	color: string;
	name?: string;
	race: Race;
	ready?: boolean;
	aIControlled?: boolean;
	submittedTurn?: boolean;
	techLevels: TechLevel;
	techLevelsSpent: TechLevel;
	messages: Message[];
	designs: ShipDesign[];
	planets: Planet[];
	fleets: Fleet[];
	starbases: Fleet[];
	planetIntels: Planet[];
	fleetIntels?: Fleet[];
	playerIntels: PlayerIntel[];
	researchSpentLastYear?: number;
	spec: PlayerSpec;
} & PlayerOrders;

export type PlayerOrders = {
	researching: TechField;
	nextResearchField: NextResearchField;
	researchAmount: number;
};

export type PlayerSpec = {
	planetaryScanner?: TechPlanetaryScanner;
	defense?: TechDefense;
	resourcesPerYear?: number;
	resourcesPerYearResearch?: number;
	currentResearchCost?: number;
};

export type PlayerIntel = {
	name: string;
	num: number;
	color: string;
	seen: boolean;
	raceName?: string;
	racePluralName?: string;
};

export enum NextResearchField {
	SameField = 'SameField',
	Energy = 'Energy',
	Weapons = 'Weapons',
	Propulsion = 'Propulsion',
	Construction = 'Construction',
	Electronics = 'Electronics',
	Biotechnology = 'Biotechnology',
	LowestField = 'LowestField'
}

export enum TechField {
	Energy = 'Energy',
	Weapons = 'Weapons',
	Propulsion = 'Propulsion',
	Construction = 'Construction',
	Electronics = 'Electronics',
	Biotechnology = 'Biotechnology'
}

export type TechLevel = {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
};

export interface Message {
	type: string;
	text: string;
	targetType?: MessageTargetType;
	targetNum?: number;
	targetPlayerNum?: number;
}

export enum MessageTargetType {
	None = '',
	Planet = 'Planet',
	Fleet = 'Fleet',
	Wormhole = 'Wormhole',
	MineField = 'MineField',
	MysteryTrader = 'MysteryTrader',
	Battle = 'Battle'
}

export function hasRequiredLevels(tl: TechLevel, required: TechLevel): boolean {
	return (
		(tl.energy ?? 0) >= (required.energy ?? 0) &&
		(tl.weapons ?? 0) >= (required.weapons ?? 0) &&
		(tl.propulsion ?? 0) >= (required.propulsion ?? 0) &&
		(tl.construction ?? 0) >= (required.construction ?? 0) &&
		(tl.electronics ?? 0) >= (required.electronics ?? 0) &&
		(tl.biotechnology ?? 0) >= (required.biotechnology ?? 0)
	);
}

export function canLearnTech(player: Player, tech: Tech): boolean {
	const requirements = tech.requirements;
	if (requirements.prtRequired && requirements.prtRequired !== player.race.prt) {
		return false;
	}
	if (requirements.prtDenied && player.race.prt === requirements.prtDenied) {
		return false;
	}
	if (
		requirements.lrtsRequired &&
		(player.race.lrts & (1 << requirements.lrtsRequired)) !== 1 << requirements.lrtsRequired
	) {
		return false;
	}
	if (requirements.lrtsDenied && (player.race.lrts & (1 << requirements.lrtsDenied)) !== 0) {
		return false;
	}
	return true;
}
