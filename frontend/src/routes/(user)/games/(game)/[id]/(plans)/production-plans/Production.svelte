<script lang="ts">
	import QuantityModifierButtons from '$lib/components/QuantityModifierButtons.svelte';
	import type { DesignFinder } from '$lib/services/Universe';
	import { fromQueueItemType, getQueueItemShortName } from '$lib/types/Planet';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { QueueItemTypes, isAuto } from '$lib/types/QueueItemType';
	import ProductionItemsButtons from './ProductionItemsButtons.svelte';

	type Props = {
		designFinder: DesignFinder;
		// default to auto tasks
		availableItems?: ProductionQueueItem[];
		queueItems?: ProductionQueueItem[];
		queueItemDescription?: any;
		onAvailableItemSelected?: (item: ProductionQueueItem) => void;
		onQueueItemSelected?: (item: ProductionQueueItem | undefined) => void;
	};

	let {
		designFinder,
		availableItems = [
			fromQueueItemType(QueueItemTypes.AutoFactories),
			fromQueueItemType(QueueItemTypes.AutoMines),
			fromQueueItemType(QueueItemTypes.AutoDefenses),
			fromQueueItemType(QueueItemTypes.AutoMineralAlchemy),
			fromQueueItemType(QueueItemTypes.AutoMaxTerraform),
			fromQueueItemType(QueueItemTypes.AutoMinTerraform)
		],
		queueItems = $bindable([]),
		queueItemDescription = getQueueItemShortName,
		onAvailableItemSelected,
		onQueueItemSelected
	}: Props = $props();

	let quantityModifier = $state(1);

	let selectedAvailableItem: ProductionQueueItem | undefined = $state();
	let selectedAvailableItemIndex = $state(-1);

	let selectedQueueItemIndex = $state(-1);
	let selectedQueueItem: ProductionQueueItem | undefined;

	function availableItemSelected(index: number, item: ProductionQueueItem) {
		selectedAvailableItemIndex = index;
		selectedAvailableItem = item;
		onAvailableItemSelected?.(selectedAvailableItem);
	}

	function queueItemClicked(index: number, item?: ProductionQueueItem) {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
		onQueueItemSelected?.(selectedQueueItem);
	}

	function addAvailableItem(item?: ProductionQueueItem) {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const quantity = quantityModifier;
		if (selectedQueueItem) {
			if (
				selectedQueueItem.type === item?.type &&
				selectedQueueItem.designNum === item?.designNum
			) {
				selectedQueueItem.quantity += quantity;
			} else {
				// insert a new item
				queueItems.splice(selectedQueueItemIndex + 1, 0, {
					type: item.type,
					quantity,
					designNum: item.designNum,
					allocated: {}
				});
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
			}
		} else {
			let nextItem = queueItems.length ? queueItems[0] : undefined;
			if (nextItem && nextItem.type === item?.type && nextItem.designNum == item.designNum) {
				nextItem.quantity++;
				selectedQueueItemIndex = 0;
				selectedQueueItem = nextItem;
			} else {
				// prepend a new queue item
				queueItems = [
					{ type: item.type, designNum: item.designNum, quantity, allocated: {} },
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
		if (queueItems && selectedQueueItem) {
			selectedQueueItem.quantity -= quantityModifier;
			queueItems = queueItems;
			if (selectedQueueItem.quantity <= 0) {
				// select the item up in the list
				queueItems = queueItems?.filter((item) => item != selectedQueueItem);
				selectedQueueItem =
					queueItems[selectedQueueItemIndex > -1 ? selectedQueueItemIndex - 1 : 0];
				selectedQueueItemIndex--;
			}
		}
	}

	function itemUp() {
		if (queueItems && selectedQueueItem && selectedQueueItemIndex > 0) {
			const swap = queueItems[selectedQueueItemIndex - 1];
			queueItems[selectedQueueItemIndex - 1] = selectedQueueItem;
			queueItems[selectedQueueItemIndex] = swap;
			selectedQueueItemIndex--;
			queueItems = queueItems;
		}
	}

	function itemDown() {
		if (queueItems && selectedQueueItem && selectedQueueItemIndex < queueItems.length - 1) {
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
