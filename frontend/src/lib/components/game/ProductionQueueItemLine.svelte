<script lang="ts">
	import {
		QueueItemCompletionEstimateSchema,
		type ProductionQueueItem
	} from '#lib/protogen/craig_stars/v1/planet_pb.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { Infinite } from '#lib/types/Consts.js';
	import {
		getAutoAlchemyDescription,
		getFullName,
		getShortName,
		getQuantityText,
		hasQuantity,
		isAuto
	} from '#lib/types/QueueItemType.js';
	import { create } from '@bufbuild/protobuf';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip';

	const { universe } = getGameContext();

	type Props = {
		index: number;
		item: ProductionQueueItem;
		hasFollowingItem?: boolean;
		selected?: boolean;
		shortName?: boolean;
		onQueueItemClicked?: (index: number, queueItem: ProductionQueueItem) => void;
		onQueueItemDoubleClicked?: (index: number, queueItem: ProductionQueueItem) => void;
	};

	let {
		index,
		item,
		hasFollowingItem = false,
		selected = false,
		shortName = false,
		onQueueItemClicked,
		onQueueItemDoubleClicked
	}: Props = $props();
	let autoAlchemy = $derived(!hasQuantity(item.type));

	let estimate = $derived(
		item.queueItemCompletionEstimate ?? create(QueueItemCompletionEstimateSchema)
	);

	// for auto items this is the first year they build their full quantity, so an auto item
	// that only builds some of its quantity this year is colored like any other partial build
	let yearsToBuildAll = $derived(estimate.yearsToBuildAll);
	let skipped = $derived(
		!autoAlchemy &&
			isAuto(item.type) &&
			estimate.yearsToBuildOne == Infinite &&
			estimate.yearsToBuildAll == Infinite
	);
</script>

<button
	type="button"
	title={autoAlchemy ? getAutoAlchemyDescription(hasFollowingItem) : undefined}
	onclick={() => onQueueItemClicked?.(index, item)}
	ondblclick={() => onQueueItemDoubleClicked?.(index, item)}
	oncontextmenu={(e) => onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
	class:italic={isAuto(item.type)}
	class:text-queue-item-this-year={!estimate.skipped &&
		estimate.yearsToBuildOne <= 1 &&
		estimate.yearsToBuildOne != Infinite}
	class:text-queue-item-next-year={!autoAlchemy &&
		!estimate.skipped &&
		(yearsToBuildAll > 1 || yearsToBuildAll === Infinite) &&
		estimate.yearsToBuildOne <= 1 &&
		estimate.yearsToBuildOne != Infinite}
	class:text-queue-item-skipped={skipped}
	class:text-queue-item-never={!autoAlchemy && estimate.yearsToBuildOne == Infinite && !skipped}
	class:bg-primary={selected}
	class="w-full text-left px-1 select-none hover:text-secondary-focus"
>
	<div class="flex justify-between">
		<div class="min-w-0 truncate">
			{shortName ? getShortName(item, $universe) : getFullName(item, $universe)}
		</div>
		{#if hasQuantity(item.type)}
			<div class="shrink-0 pl-1">{getQuantityText(item.type, item.quantity)}</div>
		{/if}
	</div>
	{#if item.note}
		<div
			class="truncate pl-2 text-sm not-italic opacity-70"
			title={item.note}
			data-type="queue-item-note"
		>
			{item.note}
		</div>
	{/if}
</button>
