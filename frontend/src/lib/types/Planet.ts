import type { Cargo } from './Cargo';
import type { Fleet } from './Fleet';
import type { Hab } from './Hab';
import { MapObjectType, type MapObject } from './MapObject';
import type { Mineral } from './Mineral';
import type { Player } from './Player';
import type { ShipDesign } from './ShipDesign';
import { UnlimitedSpaceDock } from './Tech';
import type { Vector } from './Vector';

export const Unexplored = -1;

export type Planet = {
	hab?: Hab;
	baseHab?: Hab;
	terraformedAmount?: Hab;
	mineralConcentration?: Mineral;
	mineYears?: Mineral;
	cargo?: Cargo;
	population?: number;
	mines?: number;
	factories?: number;
	defenses?: number;
	contributesOnlyLeftoverToResearch?: boolean;
	homeworld?: boolean;
	scanner?: boolean;
	reportAge: number;
	productionQueue?: ProductionQueueItem[];
	spec: PlanetSpec;
} & MapObject;

export interface ProductionQueueItem {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	type: QueueItemType;
	quantity: number;
	designName?: string;
}

/**
 * A planet that can be commanded and updated by the player
 */
export class CommandedPlanet implements Planet {
	id = 0;
	gameId = 0;
	createdAt?: string | undefined;
	updatedAt?: string | undefined;
	readonly type = MapObjectType.Planet;

	hab: Hab = { grav: 0, temp: 0, rad: 0 };
	baseHab: Hab = { grav: 0, temp: 0, rad: 0 };
	terraformedAmount = { grav: 0, temp: 0, rad: 0 };
	mineralConcentration: Mineral = { ironium: 0, boranium: 0, germanium: 0 };
	mineYears: Mineral = { ironium: 0, boranium: 0, germanium: 0 };
	cargo: Cargo = { ironium: 0, boranium: 0, germanium: 0, colonists: 0 };
	mines = 0;
	factories = 0;
	defenses = 0;
	contributesOnlyLeftoverToResearch = false;
	homeworld = false;
	scanner = false;
	reportAge = 0;
	productionQueue: ProductionQueueItem[] = [];
	position: Vector = { x: 0, y: 0 };
	name = '';
	num = 0;
	playerNum = 0;
	spec: PlanetSpec = {
		habitability: 0,
		habitabilityTerraformed: 0,
		maxMines: 0,
		maxPossibleMines: 0,
		maxFactories: 0,
		maxPossibleFactories: 0,
		maxDefenses: 0,
		populationDensity: 0,
		maxPopulation: 0,
		growthAmount: 0,
		mineralOutput: { ironium: 0, boranium: 0, germanium: 0 },
		resourcesPerYear: 0,
		resourcesPerYearAvailable: 0,
		resourcesPerYearResearch: 0,
		defense: '',
		defenseCoverage: 0,
		defenseCoverageSmart: 0,
		scanner: '',
		scanRange: 0,
		scanRangePen: 0,
		canTerraform: false,
		terraformAmount: { grav: 0, temp: 0, rad: 0 },
		hasMassDriver: false,
		hasStarbase: false,
		dockCapacity: 0
	};

	/**
	 * get a list of available ProductionQueueItems a planet can build
	 * @param planet the planet to get items for
	 * @param player the player to add items for
	 * @returns a list of items for a planet
	 */
	public getAvailableProductionQueueItems(
		planet: Planet,
		designs: ShipDesign[]
	): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (planet.spec.dockCapacity == UnlimitedSpaceDock || planet.spec.dockCapacity > 0) {
			designs
				.filter(
					(d) =>
						planet.spec.dockCapacity == UnlimitedSpaceDock ||
						(d.spec.mass ?? 0) <= planet.spec.dockCapacity
				)
				.filter((d) => !d.spec.starbase || planet.spec.starbaseDesignNum !== d.num)
				.forEach((d) => {
					items.push({
						quantity: 0,
						type: QueueItemType.ShipToken,
						designName: d.name
					});
				});
		}

		if (planet.spec.hasMassDriver) {
			items.push(
				fromQueueItemType(QueueItemType.IroniumMineralPacket),
				fromQueueItemType(QueueItemType.BoraniumMineralPacket),
				fromQueueItemType(QueueItemType.GermaniumMineralPacket),
				fromQueueItemType(QueueItemType.MixedMineralPacket)
			);
		}

		items.push(
			fromQueueItemType(QueueItemType.Factory),
			fromQueueItemType(QueueItemType.Mine),
			fromQueueItemType(QueueItemType.Defenses),
			fromQueueItemType(QueueItemType.MineralAlchemy)
		);

		if (planet.spec.canTerraform) {
			items.push(fromQueueItemType(QueueItemType.TerraformEnvironment));
		}

		// add auto items
		items.push(
			fromQueueItemType(QueueItemType.AutoFactories),
			fromQueueItemType(QueueItemType.AutoMines),
			fromQueueItemType(QueueItemType.AutoDefenses),
			fromQueueItemType(QueueItemType.AutoMineralAlchemy),
			fromQueueItemType(QueueItemType.AutoMaxTerraform),
			fromQueueItemType(QueueItemType.AutoMinTerraform)
		);

		if (planet.spec.hasMassDriver) {
			items.push(fromQueueItemType(QueueItemType.AutoMineralPacket));
		}

		return items;
	}
}

export const fromQueueItemType = (type: QueueItemType): ProductionQueueItem => ({
	type,
	quantity: 0
});

export enum QueueItemType {
	AutoMines = 'AutoMines',
	AutoFactories = 'AutoFactories',
	AutoDefenses = 'AutoDefenses',
	AutoMineralAlchemy = 'AutoMineralAlchemy',
	AutoMinTerraform = 'AutoMinTerraform',
	AutoMaxTerraform = 'AutoMaxTerraform',
	AutoMineralPacket = 'AutoMineralPacket',
	Factory = 'Factory',
	Mine = 'Mine',
	Defenses = 'Defenses',
	MineralAlchemy = 'MineralAlchemy',
	TerraformEnvironment = 'TerraformEnvironment',
	IroniumMineralPacket = 'IroniumMineralPacket',
	BoraniumMineralPacket = 'BoraniumMineralPacket',
	GermaniumMineralPacket = 'GermaniumMineralPacket',
	MixedMineralPacket = 'MixedMineralPacket',
	ShipToken = 'ShipToken',
	Starbase = 'Starbase'
}

export const getQueueItemShortName = (item: ProductionQueueItem): string => {
	switch (item.type) {
		case QueueItemType.Starbase:
		case QueueItemType.ShipToken:
			return item.designName ?? '';
		case QueueItemType.TerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemType.AutoMines:
			return 'Mine (Auto)';
		case QueueItemType.AutoFactories:
			return 'Factory (Auto)';
		case QueueItemType.AutoDefenses:
			return 'Defenses (Auto)';
		case QueueItemType.AutoMineralAlchemy:
			return 'Alchemy (Auto)';
		case QueueItemType.AutoMaxTerraform:
			return 'Max Terraform (Auto)';
		case QueueItemType.AutoMinTerraform:
			return 'Min Terraform (Auto)';
		default:
			return `${item.type}`;
	}
};

export const stringToQueueItemType = (value: string): QueueItemType => {
	return QueueItemType[value as keyof typeof QueueItemType];
};

/**
 * Determine if a ProductionQueueItem is an auto item
 * @param type The type to check
 * @returns
 */
export const isAuto = (type: QueueItemType): boolean => {
	switch (type) {
		case QueueItemType.AutoMines:
		case QueueItemType.AutoFactories:
		case QueueItemType.AutoDefenses:
		case QueueItemType.AutoMineralAlchemy:
		case QueueItemType.AutoMinTerraform:
		case QueueItemType.AutoMaxTerraform:
		case QueueItemType.AutoMineralPacket:
			return true;
		default:
			return false;
	}
};

export interface PlanetSpec {
	habitability?: number;
	habitabilityTerraformed?: number;
	maxMines: number;
	maxPossibleMines: number;
	maxFactories: number;
	maxPossibleFactories: number;
	maxDefenses: number;
	populationDensity: number;
	maxPopulation: number;
	growthAmount: number;
	mineralOutput: Mineral;
	resourcesPerYear: number;
	resourcesPerYearAvailable: number;
	resourcesPerYearResearch: number;
	defense: string;
	defenseCoverage: number;
	defenseCoverageSmart: number;
	scanner: string;
	scanRange: number;
	scanRangePen: number;
	canTerraform: boolean;
	terraformAmount?: Hab;
	hasMassDriver: boolean;
	hasStarbase: boolean;
	starbaseDesignNum?: number;
	starbaseDesignName?: string;
	dockCapacity: number;
}
