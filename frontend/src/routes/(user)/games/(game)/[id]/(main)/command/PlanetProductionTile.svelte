<script lang="ts">
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { createEventDispatcher } from 'svelte';
	import type { ProductionQueueDialogEvent } from '../../dialogs/production/ProductionQueueDialog.svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher<ProductionQueueDialogEvent>();
	const { game, player, universe, updatePlanetOrders } = getGameContext();

	export let planet: CommandedPlanet;

	const clear = async () => {
		if (planet && confirm('Are you sure you want to clear the planet production queue?')) {
			planet.productionQueue = [];
			updatePlanetOrders(planet);
		}
	};

	$: queueItems = planet.updateProductionQueueEstimates($game.rules, $techs, $player, $universe);
</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-20 overflow-y-auto">
		{#if planet.productionQueue}
			<ul class="w-full h-full">
				{#each queueItems as queueItem, index}
					<li class="pl-1">
						<ProductionQueueItemLine item={queueItem} {index} shortName={true} />
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
		<button
			on:click={() => dispatch('change-production', planet)}
			class="btn btn-outline btn-sm normal-case btn-secondary">Change</button
		>
		<button on:click={clear} class="btn btn-outline btn-sm normal-case btn-secondary">Clear</button>
		<button class="btn btn-outline btn-sm normal-case btn-secondary">Route</button>
	</div>
</CommandTile>
