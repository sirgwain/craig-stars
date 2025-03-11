import { roundTo100 } from '$lib/services/Math';
import { getMinTerraformAmount, getTerraformAmount } from '$lib/services/Terraformer';
import type { AnyPlanet } from '$lib/services/Universe';
import type { CS } from '$lib/wasm';
import { cloneDeep, sortBy } from 'lodash-es';
import { population } from './Cargo';
import type {
	Fleet,
	Planet,
	PlanetSpec,
	ProductionQueueItem,
	Rules,
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
	type QueueItemType,
	type Race,
	type TechStore
} from './cs';
import { absSum } from './Hab';
import { totalMinerals } from './Mineral';
import type { CommandedPlayer } from './Player';
import { fromQueueItemType } from './QueueItemType';
import { getGameContext } from '$lib/services/GameContext';

const { cs } = getGameContext();

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
		return Math.floor((this.cargo.colonists ?? 0) / 100);
	}

	public set population(value: number) {
		this.cargo.colonists = Math.floor(value / 100);
	}

	// get the max popluation this planet will support for a player
	public getMaxPopulation(rules: Rules, player: CommandedPlayer, habitability: number): number {
		const maxPopulationFactor = 1 + (player.race.spec?.maxPopulationOffset ?? 0);
		let maxPossiblePop = rules.maxPopulation ?? 1_000_000;
		const minMaxPop = (maxPossiblePop * maxPopulationFactor * (rules.minHabFloor ?? 5)) / 100.0;

		if (player.race.spec?.livesOnStarbases && this.playerNum === player.num) {
			maxPossiblePop = this.starbase?.spec?.maxPopulation ?? 0;
		}

		return roundTo100(
			Math.max(minMaxPop, (maxPossiblePop * maxPopulationFactor * habitability) / 100.0)
		);
	}

	public getInnateMines(race: Race, population: number): number {
		if (race.spec?.innateMining) {
			return Math.floor(Math.sqrt(population) * (race.spec.innateScannerFactor ?? 0));
		}
		return 0;
	}

	public getMaxMines(race: Race, maxPopulation: number): number {
		if (!race.spec?.innateMining) {
			return Math.floor((maxPopulation * race.numMines) / 10000);
		}
		return 0;
	}

	public getMaxFactories(race: Race, maxPopulation: number): number {
		if (!race.spec?.innateResources) {
			return Math.floor((maxPopulation * race.numFactories) / 10000);
		}
		return 0;
	}

	// get the amount of a given item in the queue
	public getAmountInQueue(
		type: QueueItemType,
		queueItems: ProductionQueueItem[] | undefined = undefined
	): number {
		queueItems = queueItems ?? this.productionQueue;
		return queueItems.reduce((count, i) => count + (i.type === type ? i.quantity : 0), 0);
	}

	public getMaxBuildable(
		techStore: TechStore,
		player: CommandedPlayer,
		maxPopulation: number,
		type: QueueItemType,
		amountInQueue = 0
	): number {
		const productivePop = cs.productivePopulation(this);
		const race = player.race;

		switch (type) {
			case QueueItemTypeAutoDefenses:
			case QueueItemTypeDefenses:
				return Math.max(0, 100 - (this.defenses + amountInQueue));
			case QueueItemTypeAutoMines:
				return Math.max(0, this.getMaxMines(race, productivePop ?? 0) - (this.mines + amountInQueue));
			case QueueItemTypeMine:
				return Math.max(0, this.getMaxMines(race, maxPopulation ?? 0) - (this.mines + amountInQueue));
			case QueueItemTypeAutoFactories:
				return Math.max(
					0,
					this.getMaxFactories(race, productivePop ?? 0) - (this.factories + amountInQueue)
				);
			case QueueItemTypeFactory:
				return Math.max(
					0,
					this.getMaxFactories(race, maxPopulation) - (this.factories + amountInQueue)
				);
			case QueueItemTypeAutoMinTerraform:
				return (
					absSum(getMinTerraformAmount(techStore, this.hab, this.baseHab, player)) - amountInQueue
				);
			case QueueItemTypeAutoMaxTerraform:
			case QueueItemTypeTerraformEnvironment:
				return (
					absSum(getTerraformAmount(techStore, this.hab, this.baseHab, player)) - amountInQueue
				);
			case QueueItemTypeAutoMineralPacket:
			case QueueItemTypeIroniumMineralPacket:
			case QueueItemTypeBoraniumMineralPacket:
			case QueueItemTypeGermaniumMineralPacket:
			case QueueItemTypeMixedMineralPacket:
			case QueueItemTypeAutoMineralAlchemy:
			case QueueItemTypeMineralAlchemy:
				return Number.MAX_SAFE_INTEGER - amountInQueue;
			case QueueItemTypePlanetaryScanner:
				// only one scanner per planet, assuming the race can build scanners...
				return Math.max(0, (this.scanner || race.spec?.innateScanner ? 0 : 1) - amountInQueue);
			case QueueItemTypeGenesisDevice:
				return 1;
			case QueueItemTypeShipToken:
				return Number.MAX_SAFE_INTEGER - amountInQueue;
			case QueueItemTypeStarbase:
				return Math.max(0, 1 - amountInQueue);
			default:
				console.error(`unknown QueueItemType ${type}`);
				return 0;
		}
	}

	// update the production queue estimates for the planet's production queue
	public updateProductionQueueEstimates(cs: CS): ProductionQueueItem[] {
		const planetWithEstimates = cs.estimateProduction(this);
		if (planetWithEstimates?.productionQueue?.length !== this.productionQueue.length) {
			throw Error("failed to estimate production queue. items don't match up");
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

export function getMineralOutput(planet: AnyPlanet, numMines: number, mineOutput: number): Mineral {
	return {
		ironium:
			((((planet.mineralConcentration?.ironium ?? 0)) * numMines) / 1000.0) * mineOutput,
		boranium:
			((((planet.mineralConcentration?.boranium ?? 0)) * numMines) / 1000.0) * mineOutput,
		germanium:
			((((planet.mineralConcentration?.germanium ?? 0)) * numMines) / 1000.0) * mineOutput
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
			return (a, b) => (a.spec.growthAmount ?? 0) - (b.spec.growthAmount ?? 0);
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
