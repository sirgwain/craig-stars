<script lang="ts">
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import { MapObjectTargetSchema, type ProductionQueueItem } from '$lib/types/cs-proto';
	import type {
		ClearProductionQueueProps,
		ShowProductionQueueDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { getMapObjectName } from '$lib/types/MapObject';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { emptyVector } from '$lib/types/Vector';
	import CommandTile from './CommandTile.svelte';
	import { create } from '@bufbuild/protobuf';

	const { settings, universe } = getGameContext();

	type Props = {
		planet: CommandedPlanet;
	} & ClearProductionQueueProps &
		ShowProductionQueueDialogProps;

	let { planet, onShowProductionQueueDialog, onClearProductionQueue }: Props = $props();
	let queueItems: ProductionQueueItem[] | undefined = $derived(planet.planetOrders.productionQueue);

	const clear = async () => {
		if (confirm('Are you sure you want to clear the planet production queue?')) {
			planet.planetOrders.productionQueue = [];
			onClearProductionQueue?.({ planet });
		}
	};
</script>

<CommandTile title="Production">
	<div class="bg-base-100 h-20 overflow-y-auto">
		{#if planet.planetOrders.productionQueue}
			<ul class="w-full h-full">
				{#if queueItems}
					{#each queueItems as queueItem, index (index)}
						<li class="pl-1 cursor-default">
							<ProductionQueueItemLine item={queueItem} {index} shortName={true} />
						</li>
					{/each}
				{/if}
			</ul>
		{/if}
	</div>
	{#if planet.planetOrders.routeTargetNum}
		{@const routeDest = $universe.getMapObject(
			create(MapObjectTargetSchema, {
				targetPosition: emptyVector(),
				targetType: planet.planetOrders.routeTargetType,
				targetNum: planet.planetOrders.routeTargetNum,
				targetPlayerNum: planet.planetOrders.routeTargetPlayerNum
			})
		)}
		<div class="flex justify-between mt-1">
			<span>Route to</span>
			<span>{getMapObjectName(routeDest)}</span>
		</div>
	{/if}
	<div class="flex justify-between">
		<button
			onclick={() => onShowProductionQueueDialog?.({ planet })}
			class="btn btn-outline btn-sm normal-case btn-secondary">Change</button
		>
		<button onclick={clear} class="btn btn-outline btn-sm normal-case btn-secondary">Clear</button>
		<button
			class="btn btn-outline btn-sm normal-case btn-secondary"
			onclick={() => ($settings.setRouteDest = !$settings.setRouteDest)}
			class:btn-accent={$settings.setRouteDest}
			type="button">Route</button
		>
	</div>
</CommandTile>
