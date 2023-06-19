<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { commandMapObject } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import {
		CommandedPlanet,
		isAuto,
		QueueItemType,
		type ProductionQueueItem
	} from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	export let planet: CommandedPlanet;

	const planetService = new PlanetService();

	const getShortName = (item: ProductionQueueItem): string => {
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

	const clear = async () => {
		if (planet && confirm('Are you sure you want to clear the planet production queue?')) {
			planet.productionQueue = [];
			const updated = await PlanetService.update(planet.gameId, planet);
			Object.assign(planet, updated);
			commandMapObject(updated);
		}
	};

	const change = () => {
		EventManager.publishProductionQueueDialogRequestedEvent(planet);
	};
</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-20 overflow-y-auto">
		{#if planet.productionQueue}
			<table class="w-full h-full">
				<tbody>
					{#each planet.productionQueue as queueItem}
						<tr>
							<td class="pl-1 {isAuto(queueItem.type) ? 'italic' : ''}"
								>{getShortName(queueItem)}</td
							>
							<td class="pr-1 text-right">{queueItem.quantity}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
	<div class="flex justify-between mt-1">
		<span>Route to</span>
		<span>{''}</span>
	</div>
	<div class="flex justify-between">
		<button on:click={change} class="btn btn-outline btn-sm normal-case btn-secondary"
			>Change</button
		>
		<button on:click={clear} class="btn btn-outline btn-sm normal-case btn-secondary">Clear</button>
		<button class="btn btn-outline btn-sm normal-case btn-secondary">Route</button>
	</div>
</CommandTile>
