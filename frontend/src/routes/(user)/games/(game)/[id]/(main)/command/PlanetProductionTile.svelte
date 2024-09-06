<script lang="ts">
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { createEventDispatcher } from 'svelte';
	import type { ProductionQueueDialogEvent } from '../../dialogs/production/ProductionQueueDialog.svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher<ProductionQueueDialogEvent>();
	const { cs, game, player, universe, updatePlanetOrders } = getGameContext();

	export let planet: CommandedPlanet;
	let queueItems: ProductionQueueItem[] | undefined = undefined;

	const clear = async () => {
		if (planet && confirm('Are you sure you want to clear the planet production queue?')) {
			planet.productionQueue = [];
			updatePlanetOrders(planet);
		}
	};

	$: queueItems = planet.updateProductionQueueEstimates(cs);
</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-20 overflow-y-auto">
		{#if planet.productionQueue}
			<ul class="w-full h-full">
				{#if queueItems}
					{#each queueItems as queueItem, index}
						<li class="pl-1 cursor-default">
							<ProductionQueueItemLine item={queueItem} {index} shortName={true} />
						</li>
					{/each}
				{/if}
			</ul>
		{/if}
	</div>
	<div class="flex justify-between mt-1">
		<span>Route to</span>
		<span>{''}</span>
	</div>
	<div class="flex justify-between">
		<button
			on:click={() => dispatch('change-production', planet)}
			class="btn btn-outline btn-sm normal-case btn-secondary">Change</button
		>
		<button on:click={clear} class="btn btn-outline btn-sm normal-case btn-secondary">Clear</button>
		<button class="btn btn-outline btn-sm normal-case btn-secondary">Route</button>
	</div>
</CommandTile>
