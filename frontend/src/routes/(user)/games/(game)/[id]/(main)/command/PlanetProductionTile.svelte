<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { Infinite } from '$lib/types/MapObject';
	import { CommandedPlanet, getQueueItemShortName, isAuto } from '$lib/types/Planet';
	import { createEventDispatcher } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher();
	const { game, universe } = getGameContext();

	export let planet: CommandedPlanet;

	const clear = async () => {
		if (planet && confirm('Are you sure you want to clear the planet production queue?')) {
			planet.productionQueue = [];
			$game.updatePlanetOrders(planet);
		}
	};

</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-20 overflow-y-auto">
		{#if planet.productionQueue}
			<ul class="w-full h-full">
				{#each planet.productionQueue as queueItem}
					<li class="pl-1">
						<div
							class="flex flex-row justify-between"
							class:italic={isAuto(queueItem.type)}
							class:text-queue-item-this-year={!queueItem.skipped &&
								(queueItem.yearsToBuildOne ?? 0) <= 1}
							class:text-queue-item-next-year={!queueItem.skipped &&
								((queueItem.yearsToBuildAll ?? 0) > 1 || queueItem.yearsToBuildAll === Infinite) &&
								(queueItem.yearsToBuildOne ?? 0) <= 1}
							class:text-queue-item-skipped={queueItem.skipped}
						>
							<div>
								{getQueueItemShortName(queueItem, $universe)}
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
		<button on:click={() => dispatch("change-production", planet)} class="btn btn-outline btn-sm normal-case btn-secondary"
			>Change</button
		>
		<button on:click={clear} class="btn btn-outline btn-sm normal-case btn-secondary">Clear</button>
		<button class="btn btn-outline btn-sm normal-case btn-secondary">Route</button>
	</div>
</CommandTile>
