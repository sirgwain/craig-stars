<!--
  @component This component displays a line of text for a ProductionQueueItem.

  The first 2 props (`index` and `item`) are used to render the text and handle button click events,
  while the next 4 props (`selected`, `shortName`, `availableItem` and `maxBuildable`)
  modify how the item is rendered in both text and coloration.

  `maxBuildable` is only used if `availableItem` is set to `true` and should be ommitted otherwise.
-->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { Infinite } from '$lib/types/cs';
	import type { ProductionQueueItem } from '$lib/types/cs';
	import { getFullName, getShortName, isAuto, isPlanetary } from '$lib/types/QueueItemType';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip.svelte';

	const { universe } = getGameContext();

	type Props = {
		index: number;
		item: ProductionQueueItem;
		selected?: boolean;
		shortName?: boolean;
		availableItem?: boolean;
		maxBuildable?: number;
		onQueueItemClicked?: (index: number, queueItem: ProductionQueueItem) => void;
		onQueueItemDoubleClicked?: () => void;
	};

	let {
		index,
		item,
		selected,
		shortName,
		availableItem = false,
		maxBuildable = 0,
		onQueueItemClicked,
		onQueueItemDoubleClicked
	}: Props = $props();

	// Note: A lot of this logic relies heavily on the fact that none of
	// QueueItemCompletionEstimate's fields are naturally set to 0 due to `omitempty`.

	let yearsToBuildAll = $derived(
		isAuto(item.type) ? item.yearsToSkipOrCancel : item.yearsToBuildAll
	);
	let skipOrCancel = $derived((item.yearsToSkipOrCancel ?? 0) > 0)
	// true if this is a planetary structure not in the queue; these have most formatting disabled
	let unbuiltStructure = $derived(availableItem && isPlanetary(item.type));
	let builtFirstYear = $derived(skipOrCancel && item.yearsToBuildOne == 1);
	let skipped = $derived(
		// grey out option to add queue items if we can't add any more
		// This mostly applies to concrete installations (but also other stuff if we happen to have 5K of them queued up)
		// Normal items in queue use the textbook definition of skipped (can't build more auto)
		availableItem ? maxBuildable == 0 : isAuto(item.type) && skipOrCancel
	);
	let canceled = $derived(!isAuto(item.type) && skipOrCancel);
</script>

<button
	type="button"
	onclick={() => onQueueItemClicked?.(index, item)}
	ondblclick={onQueueItemDoubleClicked}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:text-queue-item-auto={isAuto(item.type)}
	class:text-queue-item-this-year={!unbuiltStructure && builtFirstYear}
	class:text-queue-item-next-year={!unbuiltStructure && // started this year but not finished yet
		builtFirstYear &&
		(yearsToBuildAll ?? 0) > 1 ||
		yearsToBuildAll === Infinite}
	class:text-queue-item-never={!unbuiltStructure &&
		!canceled &&
		item.yearsToBuildOne == Infinite}
	class:text-queue-item-canceled={canceled}
	class:text-queue-item-skipped={skipped}
	class:bg-primary={selected}
	class="w-full text-left {availableItem
		? 'pl-0.5'
		: 'px-1'} select-none hover:text-secondary-focus"
>
	<div class="flex justify-between">
		<div>
			{shortName ? getShortName(item, $universe) : getFullName(item, $universe)}
		</div>
		{#if !availableItem}
			<div>
				{item.quantity}
			</div>
		{/if}
	</div>
</button>
