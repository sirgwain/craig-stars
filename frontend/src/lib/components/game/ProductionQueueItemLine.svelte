<script lang="ts" context="module">
	export type QueueItemClickedEventDetails = {
		index: number;
		queueItem: ProductionQueueItem;
	};

	export type QueueItemClickedEvent = {
		'queue-item-clicked': QueueItemClickedEventDetails;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { NeverBuilt } from '$lib/services/Producer';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { getFullName, getShortName, isAuto } from '$lib/types/QueueItemType';
	import { createEventDispatcher } from 'svelte';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip.svelte';

	const dispatch = createEventDispatcher<QueueItemClickedEvent>();
	const { universe } = getGameContext();

	export let index: number;
	export let item: ProductionQueueItem;
	export let selected = false;
	export let shortName = false;

	$: yearsToBuildAll = isAuto(item.type) ? item.yearsToSkipAuto : item.yearsToBuildAll;
	$: skipped =
		isAuto(item.type) && item.yearsToBuildOne == NeverBuilt && item.yearsToBuildAll == NeverBuilt;
</script>

<button
	type="button"
	on:click={() => dispatch('queue-item-clicked', { index, queueItem: item })}
	on:contextmenu|preventDefault={(e) =>
		onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:italic={isAuto(item.type)}
	class:text-queue-item-this-year={!item.skipped &&
		(item.yearsToBuildOne ?? 0) <= 1 &&
		item.yearsToBuildOne != NeverBuilt}
	class:text-queue-item-next-year={!item.skipped &&
		((yearsToBuildAll ?? 0) > 1 || yearsToBuildAll === NeverBuilt) &&
		(item.yearsToBuildOne ?? 0) <= 1 &&
		item.yearsToBuildOne != NeverBuilt}
	class:text-queue-item-skipped={skipped}
	class:text-queue-item-never={item.yearsToBuildOne == NeverBuilt && !skipped}
	class:bg-primary={selected}
	class="w-full text-left px-1 select-none cursor-default hover:text-secondary-focus"
>
	<div class="flex justify-between ">
		<div>
			{shortName ? getShortName(item, $universe) : getFullName(item, $universe)}
		</div>
		<div>
			{item.quantity}
		</div>
	</div>
</button>
