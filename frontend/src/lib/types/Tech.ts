import type { Cost } from './Cost';
import type { TechLevels as TechLevel } from './Player';

export interface TechStore {
	engines: TechEngine[];
	planetaryScanners: TechPlanetaryScanner[];
}

export interface Tech {
	name: string;
	cost: Cost;
	requirements: TechRequirements;
	ranking?: number;
	category: TechCategory;
}

export interface TechPlanetaryScanner extends Tech {
	scanRange: number;
	scanRangePen?: number;
}

export interface TechHullComponent extends Tech {
	hullSlotType?: string;
	mass?: number;
	scanRange: number;
	scanRangePen?: number;
	safeHullMass?: number;
	safeRange?: number;
	maxHullMass?: number;
	maxRange?: number;
	radiating?: boolean;
}

export interface TechEngine extends TechHullComponent {
	idealSpeed?: number;
	freeSpeed?: number;
	fuelUsage?: number[];
}

export enum TechCategory {
	Armor = 'Armor',
	BeamWeapon = 'BeamWeapon',
	Bomb = 'Bomb',
	Electrical = 'Electrical',
	Engine = 'Engine',
	Mechanical = 'Mechanical',
	MineLayer = 'MineLayer',
	MineRobot = 'MineRobot',
	Orbital = 'Orbital',
	PlanetaryScanner = 'PlanetaryScanner',
	PlanetaryDefense = 'PlanetaryDefense',
	Scanner = 'Scanner',
	Shield = 'Shield',
	ShipHull = 'ShipHull',
	StarbaseHull = 'StarbaseHull',
	Terraforming = 'Terraforming',
	Torpedo = 'Torpedo'
}

export interface TechRequirements extends TechLevel {
	lrtsRequired: number[] | null;
	lrtsDenied: number[] | null;
	prtRequired?: string;
	prtDenied?: string;
}
