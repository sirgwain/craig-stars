import { roundTo100 } from '$lib/services/Math';
import type { AnyPlanet, DesignFinder } from '$lib/services/Universe';
import type { CS } from '$lib/wasm';
import { cloneDeep, sortBy, startCase } from 'lodash-es';
import { population } from './Cargo';
import type {
	Fleet,
	Planet,
	PlanetSpec,
	ProductionQueueItem,
	ShipDesign,
	Tags,
	Vector
} from './cs';
import {
	Infinite,
	MapObjectTypeNone,
	MapObjectTypePlanet,
	None,
	QueueItemTypeAutoDefenses,
	QueueItemTypeAutoFactories,
	QueueItemTypeAutoMaxTerraform,
	QueueItemTypeAutoMineralAlchemy,
	QueueItemTypeAutoMineralPacket,
	QueueItemTypeAutoMines,
	QueueItemTypeAutoMinTerraform,
	QueueItemTypeBoraniumMineralPacket,
	QueueItemTypeDefenses,
	QueueItemTypeFactory,
	QueueItemTypeGenesisDevice,
	QueueItemTypeGermaniumMineralPacket,
	QueueItemTypeIroniumMineralPacket,
	QueueItemTypeMine,
	QueueItemTypeMineralAlchemy,
	QueueItemTypeMixedMineralPacket,
	QueueItemTypePlanetaryScanner,
	QueueItemTypeShipToken,
	QueueItemTypeStarbase,
	QueueItemTypeTerraformEnvironment,
	UnlimitedSpaceDock,
	type Cargo,
	type Hab,
	type Mineral,
	type QueueItemType
} from './cs';
import { totalMinerals } from './Mineral';

/**
 * A planet that can be commanded and updated by the player
 */
export class CommandedPlanet implements Planet {
	readonly type = MapObjectTypePlanet;
	tags: Tags = {};

	hab: Hab = { grav: 0, temp: 0, rad: 0 };
	baseHab: Hab = { grav: 0, temp: 0, rad: 0 };
	terraformedAmount = { grav: 0, temp: 0, rad: 0 };
	mineralConcentration: Mineral = { ironium: 0, boranium: 0, germanium: 0 };
	mineYears: Mineral = { ironium: 0, boranium: 0, germanium: 0 };
	cargo: Cargo = { ironium: 0, boranium: 0, germanium: 0, colonists: 0 };
	partialPopulation = 0;
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
	starbase: Fleet | undefined = undefined;

	// orders
	contributesOnlyLeftoverToResearch = false;
	productionQueue: ProductionQueueItem[] = [];
	routeTargetType = MapObjectTypeNone;
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

	public get population(): number {
		return (this.cargo.colonists ?? 0) * 100;
	}

	public set population(value: number) {
		this.cargo.colonists = Math.trunc(value / 100);
	}

	/**
	 * Get the amount of a given item in a planet's production queue
	 * @param type the {@linkcode QueueItemType} of the item being checked
	 * @param queueItems an array of queue items to check
	 * @returns the number of items in the queue
	 */
	public getAmountInQueue(
		type: QueueItemType,
		queueItems: ProductionQueueItem[] = this.productionQueue
	): number {
		return queueItems.reduce((count, i) => count + (i.type === type ? i.quantity : 0), 0);
	}

	// update the production queue estimates for the planet's production queue
	public updateProductionQueueEstimates(cs: CS): ProductionQueueItem[] {
		const planetWithEstimates = cs.estimateProduction(this);
		if (planetWithEstimates?.productionQueue?.length !== this.productionQueue.length) {
			// something went wrong
			// flag everything as never
			console.error("failed to estimate production queue. items don't match up");
			return this.productionQueue;
		}

		for (let i = 0; i < this.productionQueue.length; i++) {
			const estimate = planetWithEstimates.productionQueue[i];
			Object.assign(this.productionQueue[i], {
				yearsToBuildOne: estimate.yearsToBuildOne,
				yearsToBuildAll: estimate.yearsToBuildAll,
				yearsToSkipAuto: estimate.yearsToSkipAuto
			});
		}
		return this.productionQueue;
	}

	// get the mineral output of a planet based on mineOutput (10 for remote mining)
	public getMineralOutput(numMines: number, mineOutput: number): Mineral {
		return {
			ironium: Math.floor(
				((((this.mineralConcentration.ironium ?? 0) / 100) * numMines) / 10) * mineOutput
			),
			boranium: Math.floor(
				((((this.mineralConcentration.boranium ?? 0) / 100) * numMines) / 10) * mineOutput
			),
			germanium: Math.floor(
				((((this.mineralConcentration.germanium ?? 0) / 100) * numMines) / 10) * mineOutput
			)
		};
	}

	/**
	 * get a list of available ProductionQueueItems for ship designs a planet can build
	 * @param planet the planet to get items for
	 * @param designs the designs to load by num to get items for
	 * @returns a list of items for a planet
	 */
	public getAvailableProductionQueueShipDesigns(designs: ShipDesign[]): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (this.spec.dockCapacity == UnlimitedSpaceDock || (this.spec.dockCapacity ?? 0) > 0) {
			sortBy(
				designs
					.filter(
						(d) =>
							this.spec.dockCapacity == UnlimitedSpaceDock ||
							(d.spec.mass ?? 0) <= (this.spec.dockCapacity ?? 0)
					)
					.filter((d) => !d.spec.starbase)
					.filter((d) => d.originalPlayerNum == None),
				(d) => d.name
			).forEach((d) => {
				items.push({
					quantity: 1,
					type: QueueItemTypeShipToken,
					designNum: d.num,
					tags: {},
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
	public getAvailableProductionQueueStarbaseDesigns(designs: ShipDesign[]): ProductionQueueItem[] {
		// filter starbase designs
		const items = sortBy(
			designs.filter((d) => d.spec.starbase && this.spec.starbaseDesignNum !== d.num),
			(d) => d.name
		).map<ProductionQueueItem>(
			(d: ShipDesign): ProductionQueueItem => ({
				quantity: 1,
				type: QueueItemTypeStarbase,
				designNum: d.num,
				allocated: {},
				tags: {},
				yearsToBuildAll: 0
			})
		);

		return items;
	}

	/**
	 * get a list of available ProductionQueueItems a planet can build
	 */
	public getAvailableProductionQueueItems(
		innateMining: boolean | undefined,
		innateResources: boolean | undefined,
		livesOnStarbases: boolean | undefined,
		genesisDevice: boolean | undefined
	): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (!innateResources) {
			items.push(fromQueueItemType(QueueItemTypeFactory));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemTypeMine));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemTypeDefenses));
		}

		items.push(fromQueueItemType(QueueItemTypeMineralAlchemy));

		if (!this.scanner) {
			items.push(fromQueueItemType(QueueItemTypePlanetaryScanner));
		}
		if (genesisDevice) {
			items.push(fromQueueItemType(QueueItemTypeGenesisDevice));
		}

		if (this.spec.canTerraform) {
			items.push(fromQueueItemType(QueueItemTypeTerraformEnvironment));
		}

		if (this.spec.hasMassDriver) {
			items.push(
				fromQueueItemType(QueueItemTypeIroniumMineralPacket),
				fromQueueItemType(QueueItemTypeBoraniumMineralPacket),
				fromQueueItemType(QueueItemTypeGermaniumMineralPacket),
				fromQueueItemType(QueueItemTypeMixedMineralPacket)
			);
		}

		// add auto items
		if (!innateResources) {
			items.push(fromQueueItemType(QueueItemTypeAutoFactories));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemTypeAutoMines));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemTypeAutoDefenses));
		}

		items.push(
			fromQueueItemType(QueueItemTypeAutoMineralAlchemy),
			fromQueueItemType(QueueItemTypeAutoMaxTerraform),
			fromQueueItemType(QueueItemTypeAutoMinTerraform)
		);

		if (this.spec.hasMassDriver) {
			items.push(fromQueueItemType(QueueItemTypeAutoMineralPacket));
		}

		return items;
	}

	// get the estimated years to build one item
	public getYearsToBuildOne(item: ProductionQueueItem, cs: CS): number {
		const planetCopy = cloneDeep(this);
		planetCopy.productionQueue = [item];
		const planetWithEstimates = cs.estimateProduction(planetCopy);
		return planetWithEstimates?.productionQueue?.length == 1
			? (planetWithEstimates.productionQueue[0].yearsToBuildOne ?? Infinite)
			: Infinite;
	}
}

export const fromQueueItemType = (type: QueueItemType): ProductionQueueItem => ({
	type,
	quantity: 1,
	allocated: {},
	tags: {}
});

export const getQueueItemShortName = (
	item: ProductionQueueItem,
	designFinder: DesignFinder
): string => {
	switch (item.type) {
		case QueueItemTypeStarbase:
		case QueueItemTypeShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemTypeTerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemTypeAutoMines:
			return 'Mine (Auto)';
		case QueueItemTypeAutoFactories:
			return 'Factory (Auto)';
		case QueueItemTypeAutoDefenses:
			return 'Defenses (Auto)';
		case QueueItemTypeAutoMineralAlchemy:
			return 'Alchemy (Auto)';
		case QueueItemTypeAutoMaxTerraform:
			return 'Max Terraform (Auto)';
		case QueueItemTypeAutoMinTerraform:
			return 'Min Terraform (Auto)';
		default:
			return `${startCase(item.type)}`;
	}
};

/**
 * Return the amount this {@linkcode Planet} or {@linkcode PlanetIntel} will grow next year,
 * truncated to the nearest multiple of 100.
 * @param planet The planet to check
 * @returns The planet's growth next year if `planet` is a {@linkcode Planet}, or 0 for a {@linkcode PlanetIntel}
 */
export function getGrowth(planet: AnyPlanet): number {
	// TODO: Change once isIntel is added
	const pPop = 'reportAge' in planet ? 0 : planet.partialPopulation;
	return roundTo100(planet.spec.growthAmount ?? 0 + pPop, Math.trunc);
}

export function getMineralOutput(planet: AnyPlanet, numMines: number, mineOutput: number): Mineral {
	return {
		ironium: (((planet.mineralConcentration?.ironium ?? 0) * numMines) / 1000.0) * mineOutput,
		boranium: (((planet.mineralConcentration?.boranium ?? 0) * numMines) / 1000.0) * mineOutput,
		germanium: (((planet.mineralConcentration?.germanium ?? 0) * numMines) / 1000.0) * mineOutput
	};
}

// planetsSortBy returns a sortBy function for planets by key. This is used by the planets report page
// and sorting when cycling through Planets
export function planetsSortBy(key: string): ((a: AnyPlanet, b: AnyPlanet) => number) | undefined {
	switch (key) {
		case 'name':
			return (a, b) => a.name.localeCompare(b.name);
		case 'production':
			return (a, b) => {
				if (!('productionQueue' in a && 'productionQueue' in b)) {
					return 0;
				}
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
			return (a, b) => (population(a.cargo) ?? 0) - (population(b.cargo) ?? 0);
		case 'populationDensity':
			return (a, b) => (a.spec.populationDensity ?? 0) - (b.spec.populationDensity ?? 0);
		case 'populationGrowth':
			return (a, b) => getGrowth(a) - getGrowth(b);
		case 'habitability':
			return (a, b) => (a.spec.habitability ?? 0) - (b.spec.habitability ?? 0);
		case 'mines':
			return (a, b) => ('mines' in a && 'mines' in b ? (a.mines ?? 0) - (b.mines ?? 0) : 0);
		case 'factories':
			return (a, b) =>
				'factories' in a && 'factories' in b ? (a.factories ?? 0) - (b.factories ?? 0) : 0;
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
				'contributesOnlyLeftoverToResearch' in a && 'contributesOnlyLeftoverToResearch' in b
					? ((a.contributesOnlyLeftoverToResearch ?? false) ? 1 : 0) -
						((b.contributesOnlyLeftoverToResearch ?? false) ? 1 : 0)
					: 0;
		default:
			return (a, b) => a.num - b.num;
	}
}
