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
	export let queueItem: ProductionQueueItem;
	export let selected = false;
	export let shortName = false;

	$: yearsToBuildAll = isAuto(queueItem.type)
		? queueItem.yearsToSkipAuto
		: queueItem.yearsToBuildAll;
</script>

<button
	type="button"
	on:click={() => dispatch('queue-item-clicked', { index, queueItem })}
	on:contextmenu|preventDefault={(e) =>
		onShipDesignTooltip(e, $universe.getMyDesign(queueItem.designNum))}
	class:italic={isAuto(queueItem.type)}
	class:text-queue-item-this-year={!queueItem.skipped &&
		(queueItem.yearsToBuildOne ?? 0) <= 1 &&
		queueItem.yearsToBuildOne != NeverBuilt}
	class:text-queue-item-next-year={!queueItem.skipped &&
		((yearsToBuildAll ?? 0) > 1 || yearsToBuildAll === NeverBuilt) &&
		(queueItem.yearsToBuildOne ?? 0) <= 1 &&
		queueItem.yearsToBuildOne != NeverBuilt}
	class:text-queue-item-skipped={queueItem.yearsToSkipAuto === 1}
	class:text-queue-item-never={queueItem.yearsToBuildOne == NeverBuilt &&
		queueItem.yearsToSkipAuto !== 1}
	class:bg-primary={selected}
	class="w-full text-left px-1 select-none cursor-default hover:text-secondary-focus"
>
	<div class="flex justify-between ">
		<div>
			{shortName ? getShortName(queueItem, $universe) : getFullName(queueItem, $universe)}
		</div>
		<div>
			{queueItem.quantity}
		</div>
	</div>
</button>
