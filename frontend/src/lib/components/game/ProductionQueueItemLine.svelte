<!--
  @component This component displays a line of text for a ProductionQueueItem.

  The first 2 props (`index` and `item`) are used to render the text and handle button click events,
  while the next 4 props (`selected`, `shortName`, `notInQueue` and `maxBuildable`)
  dictate how to render the item in both text and coloration.

  `maxBuildable` is only used if `notInQueue` is set to `true` and should be ommitted otherwise.
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
		notInQueue?: boolean;
		maxBuildable?: number;
		onQueueItemClicked?: (index: number, queueItem: ProductionQueueItem) => void;
		onQueueItemDoubleClicked?: () => void;
	};

	let {
		index,
		item,
		selected,
		shortName,
		notInQueue = false,
		maxBuildable = 0,
		onQueueItemClicked,
		onQueueItemDoubleClicked
	}: Props = $props();

	let yearsToBuildAll = $derived(isAuto(item.type) ? item.yearsToSkipAuto : item.yearsToBuildAll);
	let skipped = $derived(
		notInQueue ? isPlanetary(item.type) &&
		maxBuildable == 0 // grey out option to add concrete queue items if we can't make more
		: isFullySkipped(item));
	let builtFirstYear = $derived(
		!skippedFirstYear(item) && (item.yearsToBuildOne ?? 0) <= 1 && item.yearsToBuildOne != Infinite
	);
</script>

<!-- Due to CSS precedence rules, later coloring rules will override prior ones-->
<button
	type="button"
	onclick={() => onQueueItemClicked?.(index, item)}
	ondblclick={() => onQueueItemDoubleClicked}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:text-queue-item-auto={isAuto(item.type)}
	class:text-queue-item-this-year={builtFirstYear}
	class:text-queue-item-next-year={(builtFirstYear && (yearsToBuildAll ?? 0) > 1) ||
		yearsToBuildAll === Infinite}
	class:text-queue-item-canceled={!notInQueue && item.canceled}
	class:text-queue-item-never={!notInQueue && item.yearsToBuildOne == Infinite}
	class:text-queue-item-skipped={skipped}
	class:bg-primary={selected}
	class="w-full text-left {notInQueue ? 'pl-0.5' : 'px-1'} select-none hover:text-secondary-focus"
>
	<div class="flex justify-between">
		<div>
			{shortName ? getShortName(item, $universe) : getFullName(item, $universe)}
		</div>
		<div>
			{item.quantity}
		</div>
	</div>
</button>
