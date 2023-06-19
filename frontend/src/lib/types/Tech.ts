import type { Cost } from './Cost';
import type { TechLevel as TechLevel } from './Player';

export interface TechStore {
	engines: TechEngine[];
	planetaryScanners: TechPlanetaryScanner[];
	defenses: TechDefense[];
	hullComponents: TechHullComponent[];
	hulls: TechHull[];
}

export interface Tech {
	id: number;
	createdAt: string;
	updatedat: string;
	deletedAt: null;
	techStoreId: number;
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

export interface TechDefense extends Tech {
	defenseCoverage: number;
}

export interface TechHullComponent extends Tech {
	hullSlotType: HullSlotType;
	scanRange: number;
	scanRangePen: number;
	safeHullMass: number;
	safeRange: number;
	maxHullMass: number;
	maxRange: number;
	ranking?: number;
	packetSpeed?: number;
	mass?: number;
	miningRate?: number;
	cloakUnits?: number;
	terraformRate?: number;
	killRate?: number;
	minKillRate?: number;
	structureDestroyRate?: number;
	unterraformRate?: number;
	smart?: boolean;
	canStealFleetCargo?: boolean;
	canStealPlanetCargo?: boolean;
	armor?: number;
	shield?: number;
	cloakUnarmedOnly?: boolean;
	torpedoBonus?: number;
	initiativeBonus?: number;
	torpedoJamming?: number;
	beamBonus?: number;
	reduceMovement?: number;
	reduceCloaking?: boolean;
	fuelBonus?: number;
	fuelRegenerationRate?: number;
	mineFieldType?: MineFieldType;
	mineLayingRate?: number;
	colonizationModule?: boolean;
	orbitalConstructionModule?: boolean;
	cargoBonus?: number;
	movementBonus?: number;
	beamDefense?: number;
	power?: number;
	range?: number;
	initiative?: number;
	gattling?: boolean;
	hitsAllTargets?: boolean;
	damageShieldsOnly?: boolean;
	accuracy?: number;
	capitalShipMissile?: boolean;
}

export interface TechHull extends Tech {
	mass?: number;
	armor: number;
	fuelCapacity?: number;
	cargoCapacity?: number;
	slots: HullSlot[];
	builtInScanner?: boolean;
	initiative?: number;
	repairBonus?: number;
	mineLayingFactor?: number;
	immuneToOwnDetonation?: boolean;
	rangeBonus?: number;
	starbase?: boolean;
	orbitalConstructionHull?: boolean;
	spaceDock?: number;
	innateScanRangePenFactor?: number;
}

export interface HullSlot {
	type: HullSlotType;
	capacity: number;
	required?: boolean;
}

export enum HullSlotType {
	None = 0,
	Engine = 1 << 0,
	Scanner = 1 << 1,
	Mechanical = 1 << 2,
	Bomb = 1 << 3,
	Mining = 1 << 4,
	Electrical = 1 << 5,
	Shield = 1 << 6,
	Armor = 1 << 7,
	Cargo = 1 << 8,
	SpaceDock = 1 << 9,
	Weapon = 1 << 10,
	Orbital = 1 << 11,
	MineLayer = 1 << 12,
	OrbitalElectrical = HullSlotType.Orbital | HullSlotType.Electrical,
	ShieldElectricalMechanical = HullSlotType.Shield |
		HullSlotType.Electrical |
		HullSlotType.Mechanical,
	ScannerElectricalMechanical = HullSlotType.Scanner |
		HullSlotType.Electrical |
		HullSlotType.Mechanical,
	ArmorScannerElectricalMechanical = HullSlotType.Armor |
		HullSlotType.Scanner |
		HullSlotType.Electrical |
		HullSlotType.Mechanical,
	MineElectricalMechanical = HullSlotType.MineLayer |
		HullSlotType.Electrical |
		HullSlotType.Mechanical,
	ShieldArmor = HullSlotType.Shield | HullSlotType.Armor,
	WeaponShield = HullSlotType.Shield | HullSlotType.Weapon,
	General = HullSlotType.Scanner |
		HullSlotType.Mechanical |
		HullSlotType.Electrical |
		HullSlotType.Shield |
		HullSlotType.Armor |
		HullSlotType.Weapon |
		HullSlotType.MineLayer
}

export enum MineFieldType {
	Heavy = 'Heavy',
	SpeedBump = 'SpeedBump',
	Standard = 'Standard'
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
