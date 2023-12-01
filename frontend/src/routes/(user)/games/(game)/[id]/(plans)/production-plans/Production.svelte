<script lang="ts">
	import { quantityModifier } from '$lib/quantityModifier';
	import type { DesignFinder } from '$lib/services/Universe';
	import { fromQueueItemType, getQueueItemShortName } from '$lib/types/Planet';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { QueueItemTypes, isAuto } from '$lib/types/QueueItemType';
	import { createEventDispatcher } from 'svelte';
	import ProductionPlanItemsButtons from './ProductionItemsButtons.svelte';

	const dispatch = createEventDispatcher();

	export let designFinder: DesignFinder;
	// default to auto tasks
	export let availableItems: ProductionQueueItem[] = [
		fromQueueItemType(QueueItemTypes.AutoFactories),
		fromQueueItemType(QueueItemTypes.AutoMines),
		fromQueueItemType(QueueItemTypes.AutoDefenses),
		fromQueueItemType(QueueItemTypes.AutoMineralAlchemy),
		fromQueueItemType(QueueItemTypes.AutoMaxTerraform),
		fromQueueItemType(QueueItemTypes.AutoMinTerraform)
	];
	export let queueItems: ProductionQueueItem[] = [];

	export let queueItemDescription = getQueueItemShortName;

	let selectedAvailableItem: ProductionQueueItem | undefined;

	let selectedQueueItemIndex = -1;
	let selectedQueueItem: ProductionQueueItem | undefined;

	const availableItemSelected = (item: ProductionQueueItem) => {
		selectedAvailableItem = item;
		dispatch('available-item-selected', selectedAvailableItem);
	};

	const queueItemClicked = (index: number, item?: ProductionQueueItem) => {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
		dispatch('queue-item-selected', selectedQueueItem);
	};

	const addAvailableItem = (e: MouseEvent, item?: ProductionQueueItem) => {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const quantity = quantityModifier(e);
		if (selectedQueueItem) {
			if (selectedQueueItem.type == item?.type && selectedQueueItem.designNum == item?.designNum) {
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
			// prepend a new queue item
			queueItems = [
				{ type: item.type, designNum: item.designNum, quantity, allocated: {} },
				...queueItems
			];
			selectedQueueItemIndex++;
			selectedQueueItem = queueItems[selectedQueueItemIndex];
		}

		// trigger reaction
		queueItems = queueItems;
	};

	const removeItem = (e: MouseEvent) => {
		if (queueItems && selectedQueueItem) {
			selectedQueueItem.quantity -= quantityModifier(e);
			queueItems = queueItems;
			if (selectedQueueItem.quantity <= 0) {
				// select the item up in the list
				queueItems = queueItems?.filter((item) => item != selectedQueueItem);
				selectedQueueItem =
					queueItems[selectedQueueItemIndex > -1 ? selectedQueueItemIndex - 1 : 0];
				selectedQueueItemIndex--;
			}
		}
	};

	const itemUp = () => {
		if (queueItems && selectedQueueItem && selectedQueueItemIndex > 0) {
			const swap = queueItems[selectedQueueItemIndex - 1];
			queueItems[selectedQueueItemIndex - 1] = selectedQueueItem;
			queueItems[selectedQueueItemIndex] = swap;
			selectedQueueItemIndex--;
			queueItems = queueItems;
		}
	};

	const itemDown = () => {
		if (queueItems && selectedQueueItem && selectedQueueItemIndex < queueItems.length - 1) {
			const swap = queueItems[selectedQueueItemIndex + 1];
			queueItems[selectedQueueItemIndex + 1] = selectedQueueItem;
			queueItems[selectedQueueItemIndex] = swap;
			selectedQueueItemIndex++;
			queueItems = queueItems;
		}
	};

	const clear = () => {
		queueItems = [];
		selectedQueueItem = undefined;
		selectedQueueItemIndex = -1;
	};
</script>

<div class="flex flex-row">
	<div class="grow">
		<ul class="h-full overflow-y-auto bg-base-300 px-1 pb-2">
			{#each availableItems as item}
				<li>
					<button
						type="button"
						on:click={() => availableItemSelected(item)}
						on:dblclick={(e) => addAvailableItem(e, item)}
						class="w-full text-left cursor-default select-none hover:text-secondary-focus {item ==
						selectedAvailableItem
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
		<ProductionPlanItemsButtons
			on:add-item={(e) => addAvailableItem(e.detail)}
			on:remove-item={(e) => removeItem(e.detail)}
			on:item-up={() => itemUp()}
			on:item-down={() => itemDown()}
			on:clear={() => clear()}
		/>
	</div>

	<div class="grow">
		<ul class="h-full bg-base-300 overflow-y-auto px-1 pb-2">
			<li>
				<button
					type="button"
					on:click={() => queueItemClicked(-1)}
					class="w-full italic pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ==
					-1
						? 'bg-primary'
						: ''}"
				>
					Top of the Queue
				</button>
			</li>
			{#if queueItems}
				{#each queueItems as queueItem, index}
					<li>
						<button
							type="button"
							on:click={() => queueItemClicked(index, queueItem)}
							class="w-full text-left pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ==
							index
								? 'bg-primary'
								: ''} {isAuto(queueItem.type) ? 'italic' : ''}"
						>
							<div class="flex justify-between ">
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
