import type { DesignFinder } from '$lib/services/Universe';
import { sortBy, startCase } from 'lodash-es';
import type { Cargo } from './Cargo';
import type { Cost } from './Cost';
import type { Hab } from './Hab';
import { MapObjectType, None, type MapObject } from './MapObject';
import { totalMinerals, type Mineral } from './Mineral';
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
	mines?: number;
	factories?: number;
	defenses?: number;
	homeworld?: boolean;
	scanner?: boolean;
	reportAge: number;

	spec: PlanetSpec;
} & MapObject &
	PlanetOrders;

export type PlanetOrders = {
	contributesOnlyLeftoverToResearch?: boolean;
	routeTargetType?: MapObjectType;
	routeTargetNum?: number;
	routeTargetPlayerNum?: number;
	packetSpeed?: number;
	packetTargetNum?: number;
	productionQueue?: ProductionQueueItem[];
};

export interface ProductionQueueItem {
	type: QueueItemType;
	quantity: number;
	designNum?: number;
	allocated: Cost;
	costOfOne: Cost;
	skipped?: boolean;
	yearsToBuildOne?: number;
	yearsToBuildAll?: number;
	percentComplete?: number;
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
	homeworld = false;
	scanner = false;
	reportAge = 0;
	position: Vector = { x: 0, y: 0 };
	name = '';
	num = 0;
	playerNum = 0;

	// orders
	contributesOnlyLeftoverToResearch = false;
	productionQueue: ProductionQueueItem[] = [];
	routeTargetType = MapObjectType.None;
	routeTargetNum = None;
	routeTargetPlayerNum = None;
	packetSpeed = 0;
	packetTargetNum = None;

	spec: PlanetSpec = {
		habitability: 0,
		terraformedHabitability: 0,
		maxMines: 0,
		maxPossibleMines: 0,
		maxFactories: 0,
		maxPossibleFactories: 0,
		maxDefenses: 0,
		populationDensity: 0,
		population: 0,
		maxPopulation: 0,
		growthAmount: 0,
		miningOutput: { ironium: 0, boranium: 0, germanium: 0 },
		resourcesPerYear: 0,
		resourcesPerYearAvailable: 0,
		resourcesPerYearResearch: 0,
		resourcesPerYearResearchEstimatedLeftover: 0,
		defense: '',
		defenseCoverage: 0,
		defenseCoverageSmart: 0,
		scanner: '',
		scanRange: 0,
		scanRangePen: 0,
		canTerraform: false,
		terraformAmount: { grav: 0, temp: 0, rad: 0 },
		minTerraformAmount: { grav: 0, temp: 0, rad: 0 },
		hasMassDriver: false,
		hasStarbase: false,
		dockCapacity: 0,
		massDriver: '',
		basePacketSpeed: 0,
		safePacketSpeed: 0,
		hasStargate: false
	};

	/**
	 * get a list of available ProductionQueueItems for ship designs a planet can build
	 * @param planet the planet to get items for
	 * @param designs the designs to load by num to get items for
	 * @returns a list of items for a planet
	 */
	public getAvailableProductionQueueShipDesigns(
		planet: Planet,
		designs: ShipDesign[]
	): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (planet.spec.dockCapacity == UnlimitedSpaceDock || planet.spec.dockCapacity > 0) {
			sortBy(
				designs
					.filter(
						(d) =>
							planet.spec.dockCapacity == UnlimitedSpaceDock ||
							(d.spec.mass ?? 0) <= planet.spec.dockCapacity
					)
					.filter((d) => !d.spec.starbase),
				(d) => d.name
			).forEach((d) => {
				items.push({
					quantity: 0,
					type: QueueItemType.ShipToken,
					designNum: d.num,
					costOfOne: d.spec.cost ?? {},
					allocated: {}
				});
			});
		}

		return items;
	}

	/**
	 * get a list of available ProductionQueueItems for ship designs a planet can build
	 * @param planet the planet to get items for
	 * @param designs the designs to load by num to get items for
	 * @returns a list of items for a planet
	 */
	public getAvailableProductionQueueStarbaseDesigns(
		planet: Planet,
		designs: ShipDesign[]
	): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		// filter starbase designs
		return sortBy(
			designs.filter((d) => d.spec.starbase && planet.spec.starbaseDesignNum !== d.num),
			(d) => d.name
		).map<ProductionQueueItem>(
			(d: ShipDesign): ProductionQueueItem => ({
				quantity: 0,
				type: QueueItemType.Starbase,
				designNum: d.num,
				costOfOne: d.spec.cost ?? {},
				allocated: {}
			})
		);

		return items;
	}

	/**
	 * get a list of available ProductionQueueItems a planet can build
	 */
	public getAvailableProductionQueueItems(
		planet: Planet,
		innateMining: boolean | undefined,
		innateResources: boolean | undefined,
		livesOnStarbases: boolean | undefined
	): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (!innateResources) {
			items.push(fromQueueItemType(QueueItemType.Factory));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemType.Mine));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemType.Defenses));
		}

		items.push(fromQueueItemType(QueueItemType.MineralAlchemy));

		if (!planet.scanner) {
			items.push(fromQueueItemType(QueueItemType.PlanetaryScanner));
		}

		if (planet.spec.canTerraform) {
			items.push(fromQueueItemType(QueueItemType.TerraformEnvironment));
		}

		if (planet.spec.hasMassDriver) {
			items.push(
				fromQueueItemType(QueueItemType.IroniumMineralPacket),
				fromQueueItemType(QueueItemType.BoraniumMineralPacket),
				fromQueueItemType(QueueItemType.GermaniumMineralPacket),
				fromQueueItemType(QueueItemType.MixedMineralPacket)
			);
		}

		// add auto items
		if (!innateResources) {
			items.push(fromQueueItemType(QueueItemType.AutoFactories));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemType.AutoMines));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemType.AutoDefenses));
		}

		items.push(
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
	quantity: 0,
	costOfOne: {},
	allocated: {}
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
	Starbase = 'Starbase',
	PlanetaryScanner = 'PlanetaryScanner'
}

export const getQueueItemShortName = (
	item: ProductionQueueItem,
	designFinder: DesignFinder
): string => {
	switch (item.type) {
		case QueueItemType.Starbase:
		case QueueItemType.ShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
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
			return `${startCase(item.type)}`;
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

/**
 * Get the concrete type for a queue item type,
 * @param type The QueueItemType
 * @returns Factory for AuotFactories, Mine for AutoMines, etc
 */
export const concreteType = (type: QueueItemType): QueueItemType => {
	switch (type) {
		case QueueItemType.AutoMines:
			return QueueItemType.Mine;
		case QueueItemType.AutoFactories:
			return QueueItemType.Factory;
		case QueueItemType.AutoDefenses:
			return QueueItemType.Defenses;
		case QueueItemType.AutoMineralAlchemy:
			return QueueItemType.MineralAlchemy;
		case QueueItemType.AutoMinTerraform:
		case QueueItemType.AutoMaxTerraform:
			return QueueItemType.TerraformEnvironment;
		case QueueItemType.AutoMineralPacket:
			return QueueItemType.MixedMineralPacket;
		default:
			return type;
	}
};

export interface PlanetSpec {
	habitability?: number;
	terraformedHabitability?: number;
	maxMines: number;
	maxPossibleMines: number;
	maxFactories: number;
	maxPossibleFactories: number;
	maxDefenses: number;
	population?: number;
	populationDensity: number;
	maxPopulation?: number;
	growthAmount: number;
	miningOutput: Mineral;
	resourcesPerYear: number;
	resourcesPerYearAvailable: number;
	resourcesPerYearResearch: number;
	resourcesPerYearResearchEstimatedLeftover: number;
	defense: string;
	defenseCoverage: number;
	defenseCoverageSmart: number;
	scanner: string;
	scanRange: number;
	scanRangePen: number;
	canTerraform: boolean;
	terraformAmount?: Hab;
	minTerraformAmount?: Hab;
	hasStarbase: boolean;
	starbaseDesignNum?: number;
	starbaseDesignName?: string;
	dockCapacity: number;

	hasMassDriver: boolean;
	massDriver: string;
	basePacketSpeed?: number;
	safePacketSpeed?: number;

	hasStargate: boolean;
	stargate?: string;
	safeHullMass?: number;
	safeRange?: number;
	maxHullMass?: number;
	maxRange?: number;
}

export function getMineralOutput(planet: Planet, numMines: number, mineOutput: number): Mineral {
	return {
		ironium:
			((((planet.mineralConcentration?.ironium ?? 0) / 100.0) * numMines) / 10.0) * mineOutput,
		boranium:
			((((planet.mineralConcentration?.boranium ?? 0) / 100.0) * numMines) / 10.0) * mineOutput,
		germanium:
			((((planet.mineralConcentration?.germanium ?? 0) / 100.0) * numMines) / 10.0) * mineOutput
	};
}

// planetsSortBy returns a sortBy function for planets by key. This is used by the planets report page
// and sorting when cycling through Planets
export function planetsSortBy(key: string): ((a: Planet, b: Planet) => number) | undefined {
	switch (key) {
		case 'name':
			return (a, b) => a.name.localeCompare(b.name);
		case 'production':
			return (a, b) => {
				const aItem =
					a.productionQueue && (a.productionQueue?.length ?? 0) > 0
						? `${JSON.stringify({
								type: a.productionQueue[0].type,
								design: a.productionQueue[0].designNum,
								quantity: a.productionQueue[0].quantity
						  })}`
						: '';
				const bItem =
					b.productionQueue && (b.productionQueue?.length ?? 0) > 0
						? `${JSON.stringify({
								type: b.productionQueue[0].type,
								design: b.productionQueue[0].designNum,
								quantity: b.productionQueue[0].quantity
						  })}`
						: '';
				return aItem.localeCompare(bItem);
			};
		case 'starbase':
			return (a, b) =>
				(a.spec.starbaseDesignName ?? '').localeCompare(b.spec.starbaseDesignName ?? '');
		case 'population':
			return (a, b) => (a.cargo?.colonists ?? 0) - (b.cargo?.colonists ?? 0);
		case 'populationDensity':
			return (a, b) => (a.spec.populationDensity ?? 0) - (b.spec.populationDensity ?? 0);
		case 'habitability':
			return (a, b) => (a.spec.habitability ?? 0) - (b.spec.habitability ?? 0);
		case 'mines':
			return (a, b) => (a.mines ?? 0) - (b.mines ?? 0);
		case 'factories':
			return (a, b) => (a.factories ?? 0) - (b.factories ?? 0);
		case 'defense':
			return (a, b) => (a.spec.defenseCoverage ?? 0) - (b.spec.defenseCoverage ?? 0);
		case 'minerals':
			return (a, b) => totalMinerals(a.cargo) - totalMinerals(b.cargo);
		case 'miningRate':
			return (a, b) => totalMinerals(a.spec.miningOutput) - totalMinerals(b.spec.miningOutput);
		case 'mineralConcentration':
			return (a, b) =>
				totalMinerals(a.mineralConcentration) - totalMinerals(b.mineralConcentration);
		case 'resources':
			return (a, b) =>
				(a.spec.resourcesPerYearAvailable ?? 0) - (b.spec.resourcesPerYearAvailable ?? 0);
		case 'contributesOnlyLeftoverToResearch':
			return (a, b) =>
				(a.contributesOnlyLeftoverToResearch ?? false ? 1 : 0) -
				(b.contributesOnlyLeftoverToResearch ?? false ? 1 : 0);
		default:
			return (a, b) => a.num - b.num;
	}
}
