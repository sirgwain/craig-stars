<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { NeverBuilt } from '$lib/types/Constants';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { getFullName, getShortName, isAuto } from '$lib/types/QueueItemType';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip.svelte';

	const { universe } = getGameContext();

	type Props = {
		index: number;
		item: ProductionQueueItem;
		selected?: boolean;
		shortName?: boolean;
		onQueueItemClicked?: (index: number, queueItem: ProductionQueueItem) => void;
	};

	let { index, item, selected = false, shortName = false, onQueueItemClicked }: Props = $props();

	let yearsToBuildAll = $derived(isAuto(item.type) ? item.yearsToSkipAuto : item.yearsToBuildAll);
	let skipped = $derived(
		isAuto(item.type) && item.yearsToBuildOne == NeverBuilt && item.yearsToBuildAll == NeverBuilt
	);
</script>

<button
	type="button"
	onclick={() => onQueueItemClicked?.(index, item)}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
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
