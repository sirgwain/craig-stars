<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { commandMapObject } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import {
		CommandedPlanet,
		isAuto,
		QueueItemType,
		type ProductionQueueItem,
		getQueueItemShortName
	} from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	export let planet: CommandedPlanet;

	const planetService = new PlanetService();

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
			<ul class="w-full h-full">
				{#each planet.productionQueue as queueItem}
					<li class="pl-1">
						<div class="flex flex-row justify-between" class:italic={isAuto(queueItem.type)}>
							<div>
								{getQueueItemShortName(queueItem)}
							</div>
							<div>
								{queueItem.quantity}
							</div>
						</div>
					</li>
				{/each}
			</ul>
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
