import type { DesignFinder } from '$lib/services/Universe';
import { sortBy, startCase } from 'lodash-es';
import { addMineral, type Cargo } from './Cargo';
import { divide, minus, minZero, type Cost } from './Cost';
import { absSum, getHabValue, getLargest, withHabValue, type Hab, add } from './Hab';
import { Infinite, MapObjectType, None, type MapObject } from './MapObject';
import { totalMinerals, type Mineral, addInt } from './Mineral';
import type { ShipDesign } from './ShipDesign';
import { UnlimitedSpaceDock, type TechStore, type Tech } from './Tech';
import type { Vector } from './Vector';
import { getPlanetHabitability, type Race } from './Race';
import { roundToNearest100 } from '$lib/services/Math';
import type { Rules } from './Rules';
import type { Player } from './Player';
import type { Fleet } from './Fleet';
import { type QueueItemType, QueueItemTypes } from './QueueItemType';
import type { ProductionQueueItem } from './Production';
import { getTerraformAmount } from '$lib/services/Terraformer';
import { getProductionEstimates } from '$lib/services/Producer';

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
	starbase: Fleet | undefined = undefined;

	// orders
	contributesOnlyLeftoverToResearch = true;
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

	// get the population from a planet's cargo
	public get population() {
		return (this.cargo.colonists ?? 0) * 100;
	}

	public set population(value: number) {
		this.cargo.colonists = Math.floor(value / 100);
	}

	// get the max popluation this planet will support for a player
	public getMaxPopulation(rules: Rules, player: Player, habitability: number): number {
		const maxPopulationFactor = 1 + (player.race.spec?.maxPopulationOffset ?? 0);
		let maxPossiblePop = rules.maxPopulation;
		const minMaxPop = maxPossiblePop * maxPopulationFactor * rules.minMaxPopulationPercent;

		if (player.race.spec?.livesOnStarbases && this.playerNum === player.num) {
			maxPossiblePop = this.starbase?.spec?.maxPopulation ?? 0;
		}

		return roundToNearest100(
			Math.max(minMaxPop, (maxPossiblePop * maxPopulationFactor * habitability) / 100.0)
		);
	}

	public getGrowthAmount(
		race: Race,
		maxPopulation: number,
		populationOvercrowdDieoffRate: number,
		populationOvercrowdDieoffRateMax: number
	): number {
		const growthFactor = race.spec?.growthFactor ?? 0;
		const capacity = this.population / maxPopulation;
		const habValue = getPlanetHabitability(race, this.hab);

		if (habValue > 0) {
			let popGrowth =
				((this.population * race.growthRate * growthFactor) / 100.0) * (habValue / 100.0) + 0.5;

			if (capacity > 1) {
				// Overpopulation calculations
				const dieoffPercent = Math.max(
					Math.min((1 - capacity) * populationOvercrowdDieoffRate, 0),
					-populationOvercrowdDieoffRateMax
				);
				popGrowth = this.population * dieoffPercent;
			} else if (capacity > 0.25) {
				const crowdingFactor = (16 / 9) * (1 - capacity) * (1 - capacity);
				popGrowth *= crowdingFactor;
			}

			// Round to the nearest 100 colonists
			return roundToNearest100(popGrowth);
		} else {
			// Kill off (habValue / 10)% colonists every year
			const deathAmount = this.population * (habValue / 1000);
			return roundToNearest100(Math.max(deathAmount, -100));
		}
	}

	public getProductivePopulation(maxPop: number): number {
		return Math.min(this.population, 3 * maxPop);
	}

	public getInnateMines(race: Race, population: number): number {
		if (race.spec?.innateMining) {
			return Math.floor(Math.sqrt(population) * (race.spec.innatePopulationFactor ?? 0));
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
		player: Player,
		maxPopulation: number,
		type: QueueItemType,
		amountInQueue = 0
	): number {
		const productivePop = this.getProductivePopulation(maxPopulation);
		const race = player.race;

		switch (type) {
			case QueueItemTypes.AutoDefenses:
			case QueueItemTypes.Defenses:
				return Math.max(0, 100 - (this.defenses + amountInQueue));
			case QueueItemTypes.AutoMines:
				return Math.max(0, this.getMaxMines(race, productivePop) - (this.mines + amountInQueue));
			case QueueItemTypes.Mine:
				return Math.max(0, this.getMaxMines(race, maxPopulation) - (this.mines + amountInQueue));
			case QueueItemTypes.AutoFactories:
				return Math.max(
					0,
					this.getMaxFactories(race, productivePop) - (this.factories + amountInQueue)
				);
			case QueueItemTypes.Factory:
				return Math.max(
					0,
					this.getMaxFactories(race, maxPopulation) - (this.factories + amountInQueue)
				);
			case QueueItemTypes.AutoMinTerraform:
			case QueueItemTypes.AutoMaxTerraform:
			case QueueItemTypes.TerraformEnvironment:
				return (
					absSum(getTerraformAmount(techStore, this.hab, this.baseHab, player)) - amountInQueue
				);

			case QueueItemTypes.AutoMineralPacket:
			case QueueItemTypes.IroniumMineralPacket:
			case QueueItemTypes.BoraniumMineralPacket:
			case QueueItemTypes.GermaniumMineralPacket:
			case QueueItemTypes.MixedMineralPacket:
			case QueueItemTypes.AutoMineralAlchemy:
			case QueueItemTypes.MineralAlchemy:
				return Number.MAX_SAFE_INTEGER - amountInQueue;
			case QueueItemTypes.PlanetaryScanner:
				// only one scanner per planet, assuming the race can build scanners...
				return Math.max(0, (this.scanner || race.spec?.innateScanner ? 0 : 1) - amountInQueue);
			case QueueItemTypes.ShipToken:
				return Number.MAX_SAFE_INTEGER - amountInQueue;
			case QueueItemTypes.Starbase:
				return Math.max(0, 1 - amountInQueue);
		}
	}

	// update the production queue estimates for the planet's production queue
	public updateProductionQueueEstimates(
		rules: Rules,
		techStore: TechStore,
		player: Player,
		designFinder: DesignFinder
	): ProductionQueueItem[] {
		const itemEstimates = getProductionEstimates(rules, techStore, player, this, designFinder);

		for (let i = 0; i < this.productionQueue.length; i++) {
			const estimate = itemEstimates[i];
			Object.assign(this.productionQueue[i], {
				yearsToBuildOne: estimate.yearsToBuildOne,
				yearsToBuildAll: estimate.yearsToBuildAll,
				yearsToSkipAuto: estimate.yearsToSkipAuto
			});
		}

		return this.productionQueue;
	}

	// grow pop on this planet. This is used when estimating production queues
	public grow(rules: Rules, player: Player) {
		const habitability = getPlanetHabitability(player.race, this.hab);
		const maxPopulation = this.getMaxPopulation(rules, player, habitability);
		const growthAmount = this.getGrowthAmount(
			player.race,
			maxPopulation,
			rules.populationOvercrowdDieoffRate,
			rules.populationOvercrowdDieoffRateMax
		);
		this.population = this.population + growthAmount;

		if (player.race.spec?.innateMining) {
			const productivePop = this.getProductivePopulation(maxPopulation);
			this.mines = this.getInnateMines(player.race, productivePop);
		}
	}

	public mine(rules: Rules, race: Race) {
		this.cargo = addMineral(this.cargo, this.getMineralOutput(this.mines, race.mineOutput));
		this.mineYears = addInt(this.mineYears, this.mines);
		this.reduceMineralConcentration(rules);
	}

	reduceMineralConcentration(rules: Rules) {
		const mineralDecayFactor = rules.mineralDecayFactor;
		let minMineralConcentration = rules.minMineralConcentration;
		if (this.homeworld) {
			minMineralConcentration = rules.minHomeworldMineralConcentration;
		}

		const planetMineYears = [
			this.mineYears.ironium ?? 0,
			this.mineYears.boranium ?? 0,
			this.mineYears.germanium ?? 0
		];
		const planetMineralConcentration = [
			this.mineralConcentration.ironium ?? 0,
			this.mineralConcentration.boranium ?? 0,
			this.mineralConcentration.germanium ?? 0
		];

		for (let i = 0; i < 3; i++) {
			let conc = planetMineralConcentration[i];

			if (conc < minMineralConcentration) {
				// Ensure the concentration is at least the minimum value
				conc = minMineralConcentration;
				planetMineralConcentration[i] = conc;
			}

			const minesPer = Math.floor(mineralDecayFactor / conc / conc);
			let mineYears = planetMineYears[i];

			if (mineYears > minesPer) {
				conc -= Math.floor(mineYears / minesPer);

				if (conc < minMineralConcentration) {
					conc = minMineralConcentration;
				}

				mineYears %= minesPer;

				planetMineYears[i] = mineYears;
				planetMineralConcentration[i] = conc;
			}
		}

		this.mineYears = {
			ironium: planetMineYears[0],
			boranium: planetMineYears[1],
			germanium: planetMineYears[2]
		};
		this.mineralConcentration = {
			ironium: planetMineralConcentration[0],
			boranium: planetMineralConcentration[1],
			germanium: planetMineralConcentration[2]
		};
	}

	// terraform this planet one step
	public terraformOneStep(techStore: TechStore, player: Player) {
		const terraformAmount = getTerraformAmount(techStore, this.hab, this.baseHab, player);

		if (absSum(terraformAmount) === 0) {
			// no need to terraform, return
			return;
		}

		const habType = getLargest(terraformAmount);
		const terraformPossibleAmount = getHabValue(terraformAmount, habType);
		if (terraformPossibleAmount > 0) {
			this.hab = add(this.hab, withHabValue(habType, 1));
		} else {
			this.hab = add(this.hab, withHabValue(habType, -1));
		}
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

	// get the resources produced by this planet each year
	public getResourcesAvailable(player: Player): number {
		const productivePop = this.getProductivePopulation(this.population);
		const race = player.race;
		if (race.spec?.innateMining) {
			return Math.floor(
				Math.sqrt((productivePop * (player.techLevels.energy ?? 0)) / race.popEfficiency)
			);
		} else {
			// compute resources from population
			const resourcesFromPop = productivePop / (race.popEfficiency * 100);

			// compute resources from factories
			const resourcesFromFactories = (this.factories * race.factoryOutput) / 10;

			return Math.floor(resourcesFromPop + resourcesFromFactories);
		}
	}

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
			const yearlyAvailableToSpend = {
				resources: planet.spec.resourcesPerYearAvailable,
				...planet.spec.miningOutput
			};

			sortBy(
				designs
					.filter(
						(d) =>
							planet.spec.dockCapacity == UnlimitedSpaceDock ||
							(d.spec.mass ?? 0) <= planet.spec.dockCapacity
					)
					.filter((d) => !d.spec.starbase)
					.filter((d) => d.originalPlayerNum == None),
				(d) => d.name
			).forEach((d) => {
				items.push({
					quantity: 0,
					type: QueueItemTypes.ShipToken,
					designNum: d.num,
					allocated: {},
					yearsToBuildOne: this.getYearsToBuildOne(
						d.spec.cost ?? {},
						this.cargo,
						yearlyAvailableToSpend
					)
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
		const yearlyAvailableToSpend = {
			resources: planet.spec.resourcesPerYearAvailable,
			...planet.spec.miningOutput
		};
		// filter starbase designs
		const items = sortBy(
			designs.filter((d) => d.spec.starbase && planet.spec.starbaseDesignNum !== d.num),
			(d) => d.name
		).map<ProductionQueueItem>(
			(d: ShipDesign): ProductionQueueItem => ({
				quantity: 0,
				type: QueueItemTypes.Starbase,
				designNum: d.num,
				allocated: {},
				yearsToBuildOne: this.getYearsToBuildOne(d.spec.cost, this.cargo, yearlyAvailableToSpend),
				yearsToBuildAll: 0
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
			items.push(fromQueueItemType(QueueItemTypes.Factory));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemTypes.Mine));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemTypes.Defenses));
		}

		items.push(fromQueueItemType(QueueItemTypes.MineralAlchemy));

		if (!planet.scanner) {
			items.push(fromQueueItemType(QueueItemTypes.PlanetaryScanner));
		}

		if (planet.spec.canTerraform) {
			items.push(fromQueueItemType(QueueItemTypes.TerraformEnvironment));
		}

		if (planet.spec.hasMassDriver) {
			items.push(
				fromQueueItemType(QueueItemTypes.IroniumMineralPacket),
				fromQueueItemType(QueueItemTypes.BoraniumMineralPacket),
				fromQueueItemType(QueueItemTypes.GermaniumMineralPacket),
				fromQueueItemType(QueueItemTypes.MixedMineralPacket)
			);
		}

		// add auto items
		if (!innateResources) {
			items.push(fromQueueItemType(QueueItemTypes.AutoFactories));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemTypes.AutoMines));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemTypes.AutoDefenses));
		}

		items.push(
			fromQueueItemType(QueueItemTypes.AutoMineralAlchemy),
			fromQueueItemType(QueueItemTypes.AutoMaxTerraform),
			fromQueueItemType(QueueItemTypes.AutoMinTerraform)
		);

		if (planet.spec.hasMassDriver) {
			items.push(fromQueueItemType(QueueItemTypes.AutoMineralPacket));
		}

		return items;
	}

	// get the estimated years to build one item
	public getYearsToBuildOne(
		cost: Cost = {},
		mineralsOnHand: Mineral,
		yearlyAvailableToSpend: Cost
	): number {
		const numBuiltInAYear = divide(yearlyAvailableToSpend, minZero(minus(cost, mineralsOnHand)));

		if (numBuiltInAYear === 0 || isNaN(numBuiltInAYear) || numBuiltInAYear == Infinity) {
			return Infinite;
		}

		return Math.ceil(1 / numBuiltInAYear);
	}
}

export const fromQueueItemType = (type: QueueItemType): ProductionQueueItem => ({
	type,
	quantity: 0,
	allocated: {}
});

export const getQueueItemShortName = (
	item: ProductionQueueItem,
	designFinder: DesignFinder
): string => {
	switch (item.type) {
		case QueueItemTypes.Starbase:
		case QueueItemTypes.ShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemTypes.TerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemTypes.AutoMines:
			return 'Mine (Auto)';
		case QueueItemTypes.AutoFactories:
			return 'Factory (Auto)';
		case QueueItemTypes.AutoDefenses:
			return 'Defenses (Auto)';
		case QueueItemTypes.AutoMineralAlchemy:
			return 'Alchemy (Auto)';
		case QueueItemTypes.AutoMaxTerraform:
			return 'Max Terraform (Auto)';
		case QueueItemTypes.AutoMinTerraform:
			return 'Min Terraform (Auto)';
		default:
			return `${startCase(item.type)}`;
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
