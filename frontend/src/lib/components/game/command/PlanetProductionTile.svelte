<script lang="ts">
	import { getCommandedPlanet } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import { isAuto, QueueItemType } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	const planet = getCommandedPlanet();
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
			planet.productionQueue = [];
			await planetService.updatePlanet(planet);
		}
	};
</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-16 overflow-y-scroll">
		{#if planet.productionQueue}
			<table class="w-full h-full">
				<tbody>
					{#each planet.productionQueue as queueItem}
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
	<div class="justify-center">
		<button class="btn btn-sm normal-case">Change</button>
		<button on:click={clear} class="btn btn-sm normal-case">Clear</button>
		<button class="btn btn-sm normal-case">Route</button>
	</div>
</CommandTile>
