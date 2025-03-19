<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { Infinite } from '$lib/types/cs';
	import type { ProductionQueueItem } from '$lib/types/cs';
	import {
		getFullName,
		getShortName,
		isAuto,
		isFullySkipped,
		skippedFirstYear
	} from '$lib/types/QueueItemType';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip.svelte';

	const { universe } = getGameContext();

	type Props = {
		index: number;
		item: ProductionQueueItem;
		selected?: boolean;
		shortName?: boolean;
		onQueueItemClicked?: (index: number, queueItem: ProductionQueueItem) => void;
		onQueueItemDoubleClicked?: () => void;
	};

	let { index, item, selected = false, shortName = false, onQueueItemClicked, onQueueItemDoubleClicked }: Props = $props();

	let yearsToBuildAll = $derived(isAuto(item.type) ? item.yearsToSkipAuto : item.yearsToBuildAll);
	let fullySkipped = $derived(isFullySkipped(item));
</script>

<button
	type="button"
	onclick={() => onQueueItemClicked?.(index, item)}
	ondblclick={() => onQueueItemDoubleClicked}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:italic={isAuto(item.type)}
	class:strikethrough={item.canceled}
	class:text-queue-item-this-year={!skippedFirstYear(item) &&
		(item.yearsToBuildOne ?? 0) <= 1 &&
		item.yearsToBuildOne != Infinite}
	class:text-queue-item-next-year={!skippedFirstYear(item) &&
		((yearsToBuildAll ?? 0) > 1 || yearsToBuildAll === Infinite) &&
		(item.yearsToBuildOne ?? 0) <= 1 &&
		item.yearsToBuildOne != Infinite}
	class:text-queue-item-skipped={fullySkipped}
	class:text-queue-item-never={item.yearsToBuildOne == Infinite && !fullySkipped}
	class:bg-primary={selected}
	class="w-full text-left px-1 select-none hover:text-secondary-focus"
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
