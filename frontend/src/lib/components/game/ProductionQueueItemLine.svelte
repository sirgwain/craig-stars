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
	import {
		getFullName,
		getShortName,
		isAuto,
		isFullySkipped,
		isPlanetary,
		skippedFirstYear
	} from '$lib/types/QueueItemType';
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

	let yearsToBuildAll = $derived(isAuto(item.type) ? item.yearsToSkipAuto : item.yearsToBuildAll);
	// true if this is a planetary structure not in the queue; these have most formatting disabled
	let unbuiltStructure = $derived(availableItem && isPlanetary(item.type));
	let skipped = $derived(
		// grey out option to add more queue items if we can't add any more
		// This mostly applies to concrete installations (but also other stuff if we happen to have 5K of them queued up)
		availableItem ? maxBuildable == 0 : isFullySkipped(item)
	);
	let builtFirstYear = $derived(
		!skippedFirstYear(item) && (item.yearsToBuildOne ?? 0) <= 1 && item.yearsToBuildOne != Infinite
	);
</script>

<!-- Due to CSS precedence rules, later coloring rules will override prior ones-->
<button
	type="button"
	onclick={() => onQueueItemClicked?.(index, item)}
	ondblclick={onQueueItemDoubleClicked}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:text-queue-item-auto={isAuto(item.type)}
	class:text-queue-item-this-year={!unbuiltStructure && builtFirstYear}
	class:text-queue-item-next-year={(!unbuiltStructure &&
		builtFirstYear &&
		(yearsToBuildAll ?? 0) > 1) ||
		yearsToBuildAll === Infinite}
	class:text-queue-item-never={!unbuiltStructure &&
		!item.canceled &&
		item.yearsToBuildOne == Infinite}
	class:text-queue-item-canceled={!unbuiltStructure && item.canceled}
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
