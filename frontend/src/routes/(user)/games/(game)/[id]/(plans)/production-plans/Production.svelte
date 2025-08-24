<script lang="ts">
	import QuantityModifierButtons from '$lib/components/QuantityModifierButtons.svelte';
	import {
		ProductionPlanItemSchema,
		QueueItemType,
		type ProductionPlanItem
	} from '$lib/types/cs-proto';
	import type { DesignFinder } from '$lib/services/Universe';
	import { getQueueItemShortName } from '$lib/types/Planet';
	import { isAuto } from '$lib/types/QueueItemType';
	import { create } from '@bufbuild/protobuf';
	import ProductionItemsButtons from './ProductionItemsButtons.svelte';
	import { planItemFromQueueItemType } from '$lib/types/Player';

	type Props = {
		designFinder: DesignFinder;
		// default to auto tasks
		availableItems?: ProductionPlanItem[];
		queueItems?: ProductionPlanItem[];
		queueItemDescription?: (item: ProductionPlanItem, designFinder: DesignFinder) => string;
		onAvailableItemSelected?: (item: ProductionPlanItem) => void;
		onQueueItemSelected?: (item: ProductionPlanItem | undefined) => void;
	};

	let {
		designFinder,
		availableItems = [
			planItemFromQueueItemType(QueueItemType.AUTO_FACTORIES),
			planItemFromQueueItemType(QueueItemType.AUTO_MINES),
			planItemFromQueueItemType(QueueItemType.AUTO_DEFENSES),
			planItemFromQueueItemType(QueueItemType.AUTO_MINERAL_ALCHEMY),
			planItemFromQueueItemType(QueueItemType.AUTO_MAX_TERRAFORM),
			planItemFromQueueItemType(QueueItemType.AUTO_MIN_TERRAFORM)
		],
		queueItems = $bindable([]),
		queueItemDescription = getQueueItemShortName,
		onAvailableItemSelected,
		onQueueItemSelected
	}: Props = $props();

	let quantityModifier = $state(1);

	let selectedAvailableItem: ProductionPlanItem | undefined = $state();
	let selectedAvailableItemIndex = $state(-1);

	let selectedQueueItemIndex = $state(-1);
	let selectedQueueItem: ProductionPlanItem | undefined;

	function availableItemSelected(index: number, item: ProductionPlanItem) {
		selectedAvailableItemIndex = index;
		selectedAvailableItem = item;
		onAvailableItemSelected?.(selectedAvailableItem);
	}

	function queueItemClicked(index: number, item?: ProductionPlanItem) {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
		onQueueItemSelected?.(selectedQueueItem);
	}

	function addAvailableItem(item?: ProductionPlanItem) {
		item = item ?? selectedAvailableItem;
		if (!item) {
			return;
		}

		const quantity = quantityModifier;
		if (selectedQueueItem) {
			if (selectedQueueItem.type === item.type && selectedQueueItem.designNum === item.designNum) {
				selectedQueueItem.quantity += quantity;
			} else {
				// insert a new item
				queueItems.splice(
					selectedQueueItemIndex + 1,
					0,
					create(ProductionPlanItemSchema, {
						type: item.type,
						quantity,
						designNum: item.designNum
					})
				);
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
			}
		} else {
			let nextItem = queueItems.length ? queueItems[0] : undefined;
			if (nextItem && nextItem.type === item.type && nextItem.designNum == item.designNum) {
				nextItem.quantity++;
				selectedQueueItemIndex = 0;
				selectedQueueItem = nextItem;
			} else {
				// prepend a new queue item
				queueItems = [
					create(ProductionPlanItemSchema, {
						type: item.type,
						designNum: item.designNum,
						quantity
					}),
					...queueItems
				];
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
			}
		}

		// trigger reaction
		queueItems = queueItems;
	}

	function removeItem() {
		if (selectedQueueItem) {
			selectedQueueItem.quantity -= quantityModifier;
			queueItems = queueItems;
			if (selectedQueueItem.quantity <= 0) {
				// select the item up in the list
				queueItems = queueItems.filter((item) => item != selectedQueueItem);
				selectedQueueItem =
					queueItems[selectedQueueItemIndex > -1 ? selectedQueueItemIndex - 1 : 0];
				selectedQueueItemIndex--;
			}
		}
	}

	function itemUp() {
		if (selectedQueueItem && selectedQueueItemIndex > 0) {
			const swap = queueItems[selectedQueueItemIndex - 1];
			queueItems[selectedQueueItemIndex - 1] = selectedQueueItem;
			queueItems[selectedQueueItemIndex] = swap;
			selectedQueueItemIndex--;
			queueItems = queueItems;
		}
	}

	function itemDown() {
		if (selectedQueueItem && selectedQueueItemIndex < queueItems.length - 1) {
			const swap = queueItems[selectedQueueItemIndex + 1];
			queueItems[selectedQueueItemIndex + 1] = selectedQueueItem;
			queueItems[selectedQueueItemIndex] = swap;
			selectedQueueItemIndex++;
			queueItems = queueItems;
		}
	}

	function clear() {
		queueItems = [];
		selectedQueueItem = undefined;
		selectedQueueItemIndex = -1;
	}
</script>

<div class="flex flex-row">
	<div class="grow">
		<ul class="h-full overflow-y-auto bg-base-300 px-1 pb-2">
			{#each availableItems as item, index (index)}
				<li>
					<button
						type="button"
						onclick={() => availableItemSelected(index, item)}
						ondblclick={() => addAvailableItem(item)}
						class="w-full text-left cursor-default select-none hover:text-secondary-focus {index ===
						selectedAvailableItemIndex
							? ' bg-primary'
							: ''}
				{isAuto(item.type) ? ' italic' : ''}"
					>
						{queueItemDescription(item, designFinder)}
					</button>
				</li>
			{/each}
		</ul>
	</div>

	<div>
		<ProductionItemsButtons
			onAddItem={() => addAvailableItem()}
			onRemoveItem={() => removeItem()}
			onItemUp={() => itemUp()}
			onItemDown={() => itemDown()}
			onClear={() => clear()}
		/>
		<div class="flex flex-col sm:flex-row justify-between mt-2 gap-1 mx-1">
			<QuantityModifierButtons bind:modifier={quantityModifier} />
		</div>
	</div>

	<div class="grow">
		<ul class="h-full bg-base-300 overflow-y-auto px-1 pb-2">
			<li>
				<button
					type="button"
					onclick={() => queueItemClicked(-1)}
					class="w-full italic pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ===
					-1
						? 'bg-primary'
						: ''}"
				>
					Top of the Queue
				</button>
			</li>
			{#if queueItems}
				{#each queueItems as queueItem, index (index)}
					<li>
						<button
							type="button"
							onclick={() => queueItemClicked(index, queueItem)}
							class="w-full text-left pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ===
							index
								? 'bg-primary'
								: ''} {isAuto(queueItem.type) ? 'italic' : ''}"
						>
							<div class="flex justify-between">
								<div>
									{queueItemDescription(queueItem, designFinder)}
								</div>
								<div>
									{queueItem.quantity}
								</div>
							</div>
						</button>
					</li>
				{/each}
			{/if}
		</ul>
	</div>
</div>
