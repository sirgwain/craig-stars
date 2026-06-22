<script lang="ts">
	import {
		QueueItemCompletionEstimateSchema,
		type ProductionQueueItem
	} from '$lib/protogen/craig_stars/v1/planet_pb';
	import { getGameContext } from '$lib/services/GameContext';
	import { Infinite } from '$lib/types/Consts';
	import { getFullName, getShortName, isAuto } from '$lib/types/QueueItemType';
	import { create } from '@bufbuild/protobuf';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip';

	const { universe } = getGameContext();

	type Props = {
		index: number;
		item: ProductionQueueItem;
		selected?: boolean;
		shortName?: boolean;
		onQueueItemClicked?: (index: number, queueItem: ProductionQueueItem) => void;
	};

	let { index, item, selected = false, shortName = false, onQueueItemClicked }: Props = $props();

	let estimate = $derived(
		item.queueItemCompletionEstimate ?? create(QueueItemCompletionEstimateSchema)
	);

	let yearsToBuildAll = $derived(
		isAuto(item.type) ? estimate.yearsToSkipAuto : estimate.yearsToBuildAll
	);
	let skipped = $derived(
		isAuto(item.type) &&
			estimate.yearsToBuildOne == Infinite &&
			estimate.yearsToBuildAll == Infinite
	);
</script>

<button
	type="button"
	onclick={() => onQueueItemClicked?.(index, item)}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:italic={isAuto(item.type)}
	class:text-queue-item-this-year={!estimate.skipped &&
		estimate.yearsToBuildOne <= 1 &&
		estimate.yearsToBuildOne != Infinite}
	class:text-queue-item-next-year={!estimate.skipped &&
		(yearsToBuildAll > 1 || yearsToBuildAll === Infinite) &&
		estimate.yearsToBuildOne <= 1 &&
		estimate.yearsToBuildOne != Infinite}
	class:text-queue-item-skipped={skipped}
	class:text-queue-item-never={estimate.yearsToBuildOne == Infinite && !skipped}
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
