import type { Fleet } from './Fleet';
import type { Planet } from './Planet';
import type { Race } from './Race';
import type { ShipDesign } from './ShipDesign';
import type { TechDefense, TechPlanetaryScanner } from './Tech';

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
