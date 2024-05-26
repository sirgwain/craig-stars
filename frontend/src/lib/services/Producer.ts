import {
	add,
	divide,
	emptyCost,
	minZero,
	minus,
	multiply,
	numBuildable,
	total,
	type Cost
} from '$lib/types/Cost';
import type { CommandedPlanet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import type { ProductionQueueItem } from '$lib/types/Production';
import { QueueItemTypes, concreteType, isAuto } from '$lib/types/QueueItemType';
import { getPlanetHabitability } from '$lib/types/Race';
import type { Rules } from '$lib/types/Rules';
import type { TechStore } from '$lib/types/Tech';
import { cloneDeep } from 'lodash-es';
import type { DesignFinder } from './Universe';

// the number of years for something to be considered not built
export const NeverBuilt = 100;

export function getProductionEstimates(
	rules: Rules,
	techStore: TechStore,
	player: Player,
	planet: CommandedPlanet,
	designFinder: DesignFinder
): ProductionQueueItem[] {
	const planetCopy = cloneDeep(planet);
	const race = player.race;

	planetCopy.productionQueue.forEach((i) => {
		i.yearsToBuildAll = NeverBuilt;
		i.yearsToBuildOne = NeverBuilt;
		i.yearsToSkipAuto = NeverBuilt;
	});

	// copy the items. We will process up to 100 turns and update these items as we go
	const items = [...planetCopy.productionQueue];

	for (let year = 1; year <= NeverBuilt; year++) {
		// mine for minerals each year
		planetCopy.mine(rules, race);

		// build things on the planet
		const completedQueue = produce(rules, techStore, planetCopy, player, designFinder, year);

		if (completedQueue) {
			// all done, no need to keep going
			break;
		}

		// grow pop
		planetCopy.grow(rules, player);

		// colonists died off, no more production
		if (planetCopy.population < 0) {
			break;
		}
	}

	return items;
}

// produce things on this planet this year, recording the year built per item
export function produce(
	rules: Rules,
	techStore: TechStore,
	planet: CommandedPlanet,
	player: Player,
	designFinder: DesignFinder,
	year: number
): boolean {
	let completedQueue = false;

	let items = [...planet.productionQueue];
	let availableToSpend: Cost = {
		ironium: planet.cargo.ironium ?? 0,
		boranium: planet.cargo.boranium ?? 0,
		germanium: planet.cargo.germanium ?? 0,
		resources:
			planet.getResourcesAvailable(player) *
			(!planet.contributesOnlyLeftoverToResearch ? (100 - player.researchAmount) / 100 : 1)
	};

	const habitability = getPlanetHabitability(player.race, planet.hab);
	const maxPopulation = planet.getMaxPopulation(rules, player, habitability);
	for (let i = 0; i < items.length; i++) {
		const item = items[i];
		const maxBuildable = planet.getMaxBuildable(techStore, player, maxPopulation, item.type);
		const cost = player.getItemCost(item, designFinder, techStore, planet);

		// check for auto items we should skip
		// we skip auto items if we can't build any more because we don't have the minerals
		if (
			isAuto(item.type) &&
			(maxBuildable <= 0 || divide(availableToSpend, { ...cost, resources: 0 }) < 1)
		) {
			if (item.yearsToSkipAuto == NeverBuilt) {
				item.yearsToSkipAuto = year;
				if (year == 1) {
					item.skipped = true;
				}
			}
			continue;
		}

		// add in any previously allocated resources to this item into our pot
		availableToSpend = add(availableToSpend, item.allocated);
		item.allocated = {};

		// determine how many we can build
		const numBuiltResult = getNumBuilt(item, cost, availableToSpend, maxBuildable);

		// deduct what was built from available
		availableToSpend = minus(availableToSpend, numBuiltResult.spent);

		// add mines and factories to the planet, terraform
		addPlanetaryInstallations(techStore, player, planet, item, numBuiltResult.numBuilt);

		// record what year this item was first and last built
		if (numBuiltResult.numBuilt > 0) {
			// we built at least one
			if (item.yearsToBuildOne === NeverBuilt) {
				item.yearsToBuildOne = year;
			}

			// if we built all we requested or can, we're
			if (numBuiltResult.numBuilt >= item.quantity) {
				// we built them all, record it
				if (item.yearsToBuildAll == NeverBuilt) {
					item.yearsToBuildAll = year;
				}
			}
			if (numBuiltResult.numBuilt >= maxBuildable) {
				// we built them all, record it
				if (item.yearsToSkipAuto == NeverBuilt) {
					item.yearsToSkipAuto = year;
				}
			}
		}

		// if we built all of the last item in the queue, we're all done
		if (
			i === items.length - 1 &&
			(numBuiltResult.numBuilt >= item.quantity || numBuiltResult.numBuilt >= maxBuildable)
		) {
			completedQueue = true;
			break;
		}

		if (isAuto(item.type)) {
			// if we are an auto item and we have resources left, move on to the next
			// auto items don't block the queue
			if ((availableToSpend.resources ?? 0) > 0) {
				// if we still have auto items to build, and
				// if we have enough minerals to complete an auto item, add a concrete one to the queue
				if (
					!(numBuiltResult.numBuilt >= item.quantity || numBuiltResult.numBuilt >= maxBuildable) &&
					divide(availableToSpend, { ...cost, resources: 0 }) >= 1
				) {
					items = [
						{
							type: concreteType(item.type),
							quantity: 1,
							allocated: { resources: availableToSpend.resources }
						},
						...items
					];
					break;
				}
				// we have resources left, but we can't build anymore of this auto item, continue to the next
				continue;
			} else {
				// no resources left, we're all done
				break;
			}
		} else {
			item.quantity -= numBuiltResult.numBuilt;

			if (item.quantity == 0 || numBuiltResult.numBuilt >= maxBuildable) {
				i--;
				items.splice(0, 1);
			} else {
				// allocate any remaining resources to the item
				const resourcesToAllocate = Math.min(availableToSpend.resources ?? 0, cost.resources ?? 0);
				availableToSpend.resources = (availableToSpend.resources ?? 0) - resourcesToAllocate;
				item.allocated.resources = resourcesToAllocate;
				break;
			}
		}
	}

	// update the items
	planet.productionQueue = items;
	planet.cargo = {
		ironium: availableToSpend.ironium,
		boranium: availableToSpend.boranium,
		germanium: availableToSpend.germanium,
		colonists: planet.cargo.colonists
	};

	return completedQueue;
}

export class ProcessQueueItemResult {
	numBuilt = 0;
	spent: Cost = {};
	partialItem?: ProductionQueueItem;
}

// for a single item in the production queue, determine how many are built
// and how much is allocated for the unbuilt items
export function getNumBuilt(
	item: ProductionQueueItem,
	cost: Cost,
	availableToSpend: Cost,
	maxBuildable: number
): ProcessQueueItemResult {
	const result = new ProcessQueueItemResult();

	// if it costs nothing, we return as many are as buildable
	// this will be the case with starbase upgrades until the upgrade logic is implemented
	if (total(minZero(cost)) == 0) {
		result.numBuilt = maxBuildable;
		return result;
	}

	// add in anything allocated in previous turns
	availableToSpend = add(availableToSpend, item.allocated);
	item.allocated = {};

	if (total(cost) > 0) {
		// figure out how many we can build
		// and make sure we only build up to the quantity, and we don't build more than the planet supports
		result.numBuilt = Math.max(
			0,
			Math.min(item.quantity, maxBuildable, numBuildable(availableToSpend, cost))
		);

		result.spent = multiply(cost, result.numBuilt);
	}

	return result;
}

function addPlanetaryInstallations(
	techStore: TechStore,
	player: Player,
	planet: CommandedPlanet,
	item: ProductionQueueItem,
	numBuilt: number
) {
	switch (item.type) {
		case QueueItemTypes.AutoMines:
		case QueueItemTypes.Mine:
			planet.mines += numBuilt;
			break;
		case QueueItemTypes.AutoFactories:
		case QueueItemTypes.Factory:
			planet.factories += numBuilt;
			break;
		case QueueItemTypes.AutoDefenses:
		case QueueItemTypes.Defenses:
			planet.defenses += numBuilt;
			break;
		case QueueItemTypes.PlanetaryScanner:
			planet.scanner = true;
			break;

		case QueueItemTypes.AutoMinTerraform:
		case QueueItemTypes.AutoMaxTerraform:
		case QueueItemTypes.TerraformEnvironment:
			planet.terraformOneStep(techStore, player);
			break;
		case QueueItemTypes.AutoMineralAlchemy:
		case QueueItemTypes.MineralAlchemy:
			// todo
			break;
		case QueueItemTypes.Starbase:
			// todo
			break;
		case QueueItemTypes.AutoMineralPacket:
		case QueueItemTypes.IroniumMineralPacket:
		case QueueItemTypes.BoraniumMineralPacket:
		case QueueItemTypes.GermaniumMineralPacket:
		case QueueItemTypes.MixedMineralPacket:
		case QueueItemTypes.ShipToken:
	}
}
