<script lang="ts">
	import { EventManager } from '$lib/EventManager';

	import { commandedPlanet } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import { isAuto, QueueItemType } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	const planetService = new PlanetService();

	const getShortName = (type: QueueItemType): string => {
		switch (type) {
			case QueueItemType.Starbase:
			case QueueItemType.ShipToken:
				return ''; //Design.Name;
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
				return `${type}`;
		}
	};

	const clear = async () => {
		if (confirm('Are you sure you want to clear the planet production queue?')) {
			$commandedPlanet.productionQueue = [];
			$commandedPlanet = await planetService.updatePlanet($commandedPlanet);
		}
	};

	const change = () => {
		EventManager.publishProductionQueueDialogRequestedEvent($commandedPlanet);
	};
</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-20 overflow-y-auto">
		{#if $commandedPlanet.productionQueue}
			<table class="w-full h-full">
				<tbody>
					{#each $commandedPlanet.productionQueue as queueItem}
						<tr>
							<td class="pl-1 {isAuto(queueItem.type) ? 'italic' : ''}"
								>{getShortName(queueItem.type)}</td
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
