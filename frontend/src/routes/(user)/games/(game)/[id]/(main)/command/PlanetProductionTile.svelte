<script lang="ts">
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import type {
		ClearProductionQueueProps,
		ShowProductionQueueDialogProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { getMapObjectName } from '$lib/types/MapObject';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { emptyVector } from '$lib/types/Vector';
	import { type ProductionQueueItem } from '$lib/types/cs';
	import CommandTile from './CommandTile.svelte';

	const { settings, universe } = getGameContext();

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
					{#each queueItems as queueItem, index (index)}
						<li class="pl-1 cursor-default">
							<ProductionQueueItemLine item={queueItem} {index} shortName={true} />
						</li>
					{/each}
				{/if}
			</ul>
		{/if}
	</div>
	{#if planet.routeTargetNum}
		{@const routeDest = $universe.getMapObject({
			targetPosition: emptyVector,
			targetType: planet.routeTargetType ?? '',
			targetNum: planet.routeTargetNum ?? 0,
			targetPlayerNum: planet.routeTargetPlayerNum ?? 0
		})}
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
