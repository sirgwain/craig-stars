<script lang="ts">
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import type {
		ClearProductionQueueProps,
		ShowProductionQueueDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import type { ProductionQueueItem } from '$lib/types/cs';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const { universe } = getGameContext();

	type Props = {
		planet: CommandedPlanet;
	} & ClearProductionQueueProps &
		ShowProductionQueueDialogProps;

	let { planet, onShowProductionQueueDialog, onClearProductionQueue }: Props = $props();
	let queueItems: ProductionQueueItem[] | undefined = $derived(planet.productionQueue);

	const clear = async () => {
		if (planet && confirm('Are you sure you want to clear the planet production queue?')) {
			planet.productionQueue = [];
			onClearProductionQueue?.({ planet });
		}
	};
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
		<!-- TODO: Make this a button that updates routing dest stuff-->
		<span
			>{$universe.getPlanet(planet.routeTargetNum)?.name
				? 'Routing to ' + ($universe.getPlanet(planet.routeTargetNum)?.name ?? 'somewhere??')
				: 'Not routing'}</span
		>
	</div>
	<div class="flex justify-between">
		<button
			onclick={() => onShowProductionQueueDialog?.({ planet })}
			class="btn btn-outline btn-sm normal-case btn-secondary">Change</button
		>
		<button onclick={clear} class="btn btn-outline btn-sm normal-case btn-secondary">Clear</button>

		<button class="btn btn-outline btn-sm normal-case btn-secondary">Route</button>
	</div>
</CommandTile>
