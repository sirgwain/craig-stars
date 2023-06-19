import {
	CommandedPlanet,
	fromQueueItemType,
	QueueItemType,
	type Planet,
	type ProductionQueueItem
} from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { Service } from './Service';

export class PlanetService extends Service {
	async updatePlanet(planet: CommandedPlanet): Promise<CommandedPlanet> {
		const updated = await Service.update<Planet>(planet, `/api/planets/${planet.id}`);
		return Object.assign(planet, updated)
	}

	/**
	 * get a list of available ProductionQueueItems a planet can build
	 * @param planet the planet to get items for
	 * @param player the player to add items for
	 * @returns a list of items for a planet
	 */
	getAvailableProductionQueueItems(planet: Planet, player: Player): ProductionQueueItem[] {
		const items: ProductionQueueItem[] = [];

		if (planet.spec) {
			if (planet.spec.dockCapacity > 0) {
				// todo: add designs
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
		}

		return items;
	}
}
