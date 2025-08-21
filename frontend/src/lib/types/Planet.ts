import { roundTo100 } from '$lib/services/Math';
import type { DesignFinder } from '$lib/services/Universe';
import type { ProductionPlanItem, ShipDesign } from '$lib/types/cs-proto';
import {
	CargoSchema,
	HabSchema,
	MapObjectSchema,
	MapObjectType,
	MineralSchema,
	PlanetOrdersSchema,
	PlanetSchema,
	PlanetSpecSchema,
	ProductionQueueItemSchema,
	QueueItemType,
	type Fleet,
	type Mineral,
	type Planet,
	type PlanetSpec,
	type ProductionQueueItem
} from '$lib/types/cs-proto';
import type { CS } from '$lib/wasm';
import { clone, create, merge, type UnknownField } from '@bufbuild/protobuf';
import { sortBy } from 'lodash-es';
import { population } from './Cargo';
import { Infinite, None, UnlimitedSpaceDock } from './Consts';
import { enumToString } from './Enums';
import { totalMinerals } from './Mineral';

/**
 * A planet that can be commanded and updated by the player
 */
export class CommandedPlanet implements Planet {
	$typeName: 'craig_stars.v1.Planet';
	$unknown?: UnknownField[] | undefined;
	readonly type = MapObjectType.PLANET;

	mapObject = create(MapObjectSchema);
	hab = create(HabSchema);
	baseHab = create(HabSchema);
	terraformedAmount = create(HabSchema);
	mineralConcentration = create(MineralSchema);
	mineYears = create(MineralSchema);
	cargo = create(CargoSchema);
	partialPopulation = 0;
	mines = 0;
	factories = 0;
	defenses = 0;
	homeworld = false;
	scanner = false;
	reportAge = 0;
	starbase: Fleet | undefined = undefined;
	planetOrders = create(PlanetOrdersSchema);
	spec: PlanetSpec = create(PlanetSpecSchema);

	constructor(data?: Planet) {
		this.$typeName = 'craig_stars.v1.Planet';

		if (data) {
			merge(PlanetSchema, this, data);
		}
	}

	/**
	 * Get the amount of a given item in a planet's production queue
	 * @param type the {@linkcode QueueItemType} of the item being checked
	 * @param queueItems an array of queue items to check
	 * @returns the number of items in the queue
	 */
	public getAmountInQueue(
		type: QueueItemType,
		queueItems: ProductionQueueItem[] = this.planetOrders.productionQueue
	): number {
		return queueItems.reduce((count, i) => count + (i.type === type ? (i.quantity ?? 0) : 0), 0);
	}

	// update the production queue estimates for the planet's production queue
	public async updateProductionQueueEstimates(cs: CS): Promise<ProductionQueueItem[]> {
		const { planet: planetWithEstimates } = await cs.wasmService.estimateProduction({
			planet: this
		});
		if (
			(planetWithEstimates?.planetOrders?.productionQueue?.length ?? 0) !==
			(this.planetOrders.productionQueue.length ?? 0)
		) {
			// something went wrong
			console.error("failed to estimate production queue. items don't match up");
			return this.planetOrders.productionQueue;
		}

		const estimatesPQ: ProductionQueueItem[] =
			planetWithEstimates?.planetOrders?.productionQueue ?? [];
		for (let i = 0; i < this.planetOrders.productionQueue.length; i++) {
			const estimate = estimatesPQ[i];
			if (estimate) {
				Object.assign(this.planetOrders.productionQueue[i], {
					yearsToBuildOne: estimate.queueItemCompletionEstimate?.yearsToBuildOne,
					yearsToBuildAll: estimate.queueItemCompletionEstimate?.yearsToBuildAll,
					yearsToSkipAuto: estimate.queueItemCompletionEstimate?.yearsToSkipAuto
				});
			}
		}
		return this.planetOrders.productionQueue;
	}

	/**
	 * get a list of available ProductionQueueItems for ship designs a planet can build
	 * @param planet the planet to get items for
	 * @param designs the designs to load by num to get items for
	 * @returns a list of items for a planet
	 */
	public getAvailableProductionQueueShipDesigns(designs: ShipDesign[]): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (
			this.spec?.planetStarbaseSpec?.dockCapacity == UnlimitedSpaceDock ||
			(this.spec?.planetStarbaseSpec?.dockCapacity ?? 0) > 0
		) {
			sortBy(
				designs
					.filter(
						(d) =>
							this.spec?.planetStarbaseSpec?.dockCapacity == UnlimitedSpaceDock ||
							(d.spec?.mass ?? 0) <= (this.spec?.planetStarbaseSpec?.dockCapacity ?? 0)
					)
					.filter((d) => !d.spec?.starbase)
					.filter((d) => d.originalPlayerNum == None),
				(d) => d.name
			).forEach((d) => {
				items.push(
					create(ProductionQueueItemSchema, {
						quantity: 1,
						type: QueueItemType.SHIP_TOKEN,
						designNum: d.num
					})
				);
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
			designs.filter(
				(d) => d.spec?.starbase && this.spec?.planetStarbaseSpec?.starbaseDesignNum !== d.num
			),
			(d) => d.name
		).map<ProductionQueueItem>(
			(d: ShipDesign): ProductionQueueItem =>
				create(ProductionQueueItemSchema, {
					quantity: 1,
					type: QueueItemType.STARBASE,
					designNum: d.num
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
			items.push(fromQueueItemType(QueueItemType.FACTORY));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemType.MINE));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemType.DEFENSES));
		}

		items.push(fromQueueItemType(QueueItemType.MINERAL_ALCHEMY));

		if (!this.scanner) {
			items.push(fromQueueItemType(QueueItemType.PLANETARY_SCANNER));
		}
		if (genesisDevice) {
			items.push(fromQueueItemType(QueueItemType.GENESIS_DEVICE));
		}

		if (this.spec?.canTerraform) {
			items.push(fromQueueItemType(QueueItemType.TERRAFORM_ENVIRONMENT));
		}

		if (this.spec?.planetStarbaseSpec?.hasMassDriver) {
			items.push(
				fromQueueItemType(QueueItemType.IRONIUM_MINERAL_PACKET),
				fromQueueItemType(QueueItemType.BORANIUM_MINERAL_PACKET),
				fromQueueItemType(QueueItemType.GERMANIUM_MINERAL_PACKET),
				fromQueueItemType(QueueItemType.MIXED_MINERAL_PACKET)
			);
		}

		// add auto items
		if (!innateResources) {
			items.push(fromQueueItemType(QueueItemType.AUTO_FACTORIES));
		}
		if (!innateMining) {
			items.push(fromQueueItemType(QueueItemType.AUTO_MINES));
		}
		if (!livesOnStarbases) {
			items.push(fromQueueItemType(QueueItemType.AUTO_DEFENSES));
		}

		items.push(
			fromQueueItemType(QueueItemType.AUTO_MINERAL_ALCHEMY),
			fromQueueItemType(QueueItemType.AUTO_MAX_TERRAFORM),
			fromQueueItemType(QueueItemType.AUTO_MIN_TERRAFORM)
		);

		if (this.spec?.planetStarbaseSpec?.hasMassDriver) {
			items.push(fromQueueItemType(QueueItemType.AUTO_MINERAL_PACKET));
		}

		return items;
	}

	// get the estimated years to build one item
	public async getYearsToBuildOne(item: ProductionQueueItem, cs: CS): Promise<number> {
		const planetCopy = clone(PlanetSchema, this);
		planetCopy.planetOrders = create(PlanetOrdersSchema, planetCopy.planetOrders);
		planetCopy.planetOrders.productionQueue = [item];

		const { planet: planetWithEstimates } = await cs.wasmService.estimateProduction({
			planet: planetCopy
		});
		return (
			planetWithEstimates?.planetOrders?.productionQueue[0].queueItemCompletionEstimate
				?.yearsToBuildOne ?? Infinite
		);
	}
}

export const fromQueueItemType = (type: QueueItemType): ProductionQueueItem =>
	create(ProductionQueueItemSchema, {
		type,
		quantity: 1
	});

export const getQueueItemShortName = (
	item: ProductionQueueItem | ProductionPlanItem,
	designFinder: DesignFinder
): string => {
	switch (item.type) {
		case QueueItemType.STARBASE:
		case QueueItemType.SHIP_TOKEN:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemType.TERRAFORM_ENVIRONMENT:
			return 'Terraform Environment';
		case QueueItemType.AUTO_MINES:
			return 'Mine (Auto)';
		case QueueItemType.AUTO_FACTORIES:
			return 'Factory (Auto)';
		case QueueItemType.AUTO_DEFENSES:
			return 'Defenses (Auto)';
		case QueueItemType.AUTO_MINERAL_ALCHEMY:
			return 'Alchemy (Auto)';
		case QueueItemType.AUTO_MAX_TERRAFORM:
			return 'Max Terraform (Auto)';
		case QueueItemType.AUTO_MIN_TERRAFORM:
			return 'Min Terraform (Auto)';
		default:
			return `${enumToString(QueueItemType, item.type)}`;
	}
};

/**
 * Return the amount this {@linkcode Planet} or {@linkcode Planet} will grow next year,
 * truncated to the nearest multiple of 100.
 * @param planet The planet to check
 * @returns The planet's growth next year if `planet` is a {@linkcode Planet}, or 0 for a {@linkcode Planet}
 */
export function getGrowth(planet: Planet): number {
	const pPop = planet.partialPopulation;
	return roundTo100(planet.spec?.growthAmount ?? 0 + pPop, Math.trunc);
}

export function getMineralOutput(planet: Planet, numMines: number, mineOutput: number): Mineral {
	return create(MineralSchema, {
		ironium: (((planet.mineralConcentration?.ironium ?? 0) * numMines) / 1000.0) * mineOutput,
		boranium: (((planet.mineralConcentration?.boranium ?? 0) * numMines) / 1000.0) * mineOutput,
		germanium: (((planet.mineralConcentration?.germanium ?? 0) * numMines) / 1000.0) * mineOutput
	});
}

// planetsSortBy returns a sortBy function for planets by key. This is used by the planets report page
// and sorting when cycling through Planets
export function planetsSortBy(key: string): ((a: Planet, b: Planet) => number) | undefined {
	switch (key) {
		case 'name':
			return (a, b) => (a.mapObject?.name ?? '').localeCompare(b.mapObject?.name ?? '');
		case 'production':
			return (a, b) => {
				if (!('planetOrders' in a && 'planetOrders' in b)) {
					return 0;
				}
				const aItem =
					a.planetOrders?.productionQueue && (a.planetOrders.productionQueue?.length ?? 0) > 0
						? `${JSON.stringify({
								type: a.planetOrders.productionQueue[0].type,
								design: a.planetOrders.productionQueue[0].designNum,
								quantity: a.planetOrders.productionQueue[0].quantity
							})}`
						: '';
				const bItem =
					b.planetOrders?.productionQueue && (b.planetOrders.productionQueue?.length ?? 0) > 0
						? `${JSON.stringify({
								type: b.planetOrders.productionQueue[0].type,
								design: b.planetOrders.productionQueue[0].designNum,
								quantity: b.planetOrders.productionQueue[0].quantity
							})}`
						: '';
				return aItem.localeCompare(bItem);
			};

		case 'starbase':
			return (a, b) =>
				(a.spec?.planetStarbaseSpec?.starbaseDesignName ?? '').localeCompare(
					b.spec?.planetStarbaseSpec?.starbaseDesignName ?? ''
				);
		case 'population':
			return (a, b) => (population(a.cargo) ?? 0) - (population(b.cargo) ?? 0);
		case 'populationDensity':
			return (a, b) => (a.spec?.populationDensity ?? 0) - (b.spec?.populationDensity ?? 0);
		case 'populationGrowth':
			return (a, b) => getGrowth(a) - getGrowth(b);
		case 'habitability':
			return (a, b) => (a.spec?.habitability ?? 0) - (b.spec?.habitability ?? 0);
		case 'mines':
			return (a, b) => ('mines' in a && 'mines' in b ? (a.mines ?? 0) - (b.mines ?? 0) : 0);
		case 'factories':
			return (a, b) =>
				'factories' in a && 'factories' in b ? (a.factories ?? 0) - (b.factories ?? 0) : 0;
		case 'defense':
			return (a, b) => (a.spec?.defenseCoverage ?? 0) - (b.spec?.defenseCoverage ?? 0);
		case 'minerals':
			return (a, b) => totalMinerals(a.cargo) - totalMinerals(b.cargo);
		case 'miningRate':
			return (a, b) => totalMinerals(a.spec?.miningOutput) - totalMinerals(b.spec?.miningOutput);
		case 'mineralConcentration':
			return (a, b) =>
				totalMinerals(a.mineralConcentration) - totalMinerals(b.mineralConcentration);
		case 'resources':
			return (a, b) =>
				(a.spec?.resourcesPerYearAvailable ?? 0) - (b.spec?.resourcesPerYearAvailable ?? 0);
		case 'contributesOnlyLeftoverToResearch':
			return (a, b) =>
				'contributesOnlyLeftoverToResearch' in a && 'contributesOnlyLeftoverToResearch' in b
					? ((a.contributesOnlyLeftoverToResearch ?? false) ? 1 : 0) -
						((b.contributesOnlyLeftoverToResearch ?? false) ? 1 : 0)
					: 0;
		default:
			return (a, b) => (a.mapObject?.num ?? 0) - (b.mapObject?.num ?? 0);
	}
}
