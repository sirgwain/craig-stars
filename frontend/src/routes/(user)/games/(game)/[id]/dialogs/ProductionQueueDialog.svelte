<script lang="ts">
	import CostComponent from '$lib/components/game/Cost.svelte';
	import { getQuantityModifier } from '$lib/quantityModifier';
	import { commandMapObject } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import type { Cost } from '$lib/types/Cost';
	import type { CommandedPlanet, ProductionQueueItem } from '$lib/types/Planet';
	import { isAuto, QueueItemType } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import {
		ArrowLongDown,
		ArrowLongLeft,
		ArrowLongRight,
		ArrowLongUp,
		XCircle
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher } from 'svelte';

	export let player: Player;
	export let planet: CommandedPlanet;

	const getFullName = (item: ProductionQueueItem) => {
		switch (item.type) {
			case QueueItemType.Starbase:
			case QueueItemType.ShipToken:
				return item.designName ?? '';
			case QueueItemType.AutoMineralAlchemy:
				return 'Alchemy (Auto Build)';
			case QueueItemType.MineralAlchemy:
				return 'Alchemy';
			case QueueItemType.AutoMines:
				return 'Mine (Auto Build)';
			case QueueItemType.AutoFactories:
				return 'Factory (Auto Build)';
			case QueueItemType.AutoDefenses:
				return 'Defense (Auto Build)';
			case QueueItemType.AutoMinTerraform:
				return 'Minimum Terraform';
			case QueueItemType.AutoMaxTerraform:
				return 'Maximum Terraform';
			case QueueItemType.IroniumMineralPacket:
				return 'Mineral Packet (Ironium)';
			case QueueItemType.BoraniumMineralPacket:
				return 'Mineral Packet (Boranium)';
			case QueueItemType.GermaniumMineralPacket:
				return 'Mineral Packet (Germanium)';
			case QueueItemType.TerraformEnvironment:
				return 'Terraform Environment';
			case QueueItemType.MixedMineralPacket:
				return 'Mixed Mineral Packet';
			case QueueItemType.AutoMineralPacket:
				return 'Mixed Mineral Packet (Auto)';
			default:
				return item.type.toString();
		}
	};

	const availableItemSelected = (type: ProductionQueueItem) => {
		selectedAvailableItem = type;
		selectedAvailableItemCost = getAvailableItemCost();
	};

	const queueItemClicked = (index: number, item?: ProductionQueueItem) => {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
		selectedQueueItemCost = getSelectedItemCost();
	};

	const addAvailableItem = (item?: ProductionQueueItem) => {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const quantity = getQuantityModifier();
		if (selectedQueueItem) {
			if (selectedQueueItem.type == item?.type) {
				selectedQueueItem.quantity += quantity;
			} else {
				// insert a new item
				queueItems.splice(selectedQueueItemIndex + 1, 0, { type: item.type, quantity });
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
			}
		} else {
			// prepend a new queue item
			queueItems = [{ type: item.type, quantity }, ...queueItems];
			selectedQueueItemIndex++;
			selectedQueueItem = queueItems[selectedQueueItemIndex];
		}

		// trigger reaction
		queueItems = queueItems;
	};

	const removeItem = () => {
		if (queueItems && selectedQueueItem) {
			selectedQueueItem.quantity -= getQuantityModifier();
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

	const ok = async () => {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		const result = await planetService.updatePlanet(planet);
		commandMapObject(result);
		dispatch('ok');
	};

	const cancel = () => {
		if (planet) {
			queueItems = planet.productionQueue?.map((item) => ({ ...item } as ProductionQueueItem));
			contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
			dispatch('cancel');
		}
	};

	const getSelectedItemCost = (): Cost | undefined => {
		if (selectedQueueItem) {
			const typeCost = player.race.spec?.costs[selectedQueueItem.type];
			if (typeCost) {
				return {
					ironium: (typeCost.ironium ?? 0) * selectedQueueItem.quantity,
					boranium: (typeCost.boranium ?? 0) * selectedQueueItem.quantity,
					germanium: (typeCost.germanium ?? 0) * selectedQueueItem.quantity,
					resources: (typeCost.resources ?? 0) * selectedQueueItem.quantity
				};
			}
		}
		return;
	};

	const getAvailableItemCost = (): Cost | undefined => {
		if (player && selectedAvailableItem) {
			return player.race.spec?.costs[selectedAvailableItem.type];
		}
		return;
	};

	const dispatch = createEventDispatcher();

	hotkeys('Esc', () => cancel());
	hotkeys('Enter', () => {
		ok();
	});

	const planetService = new PlanetService();

	let availableItems: ProductionQueueItem[] = [];
	let queueItems: ProductionQueueItem[] = [];
	let contributesOnlyLeftoverToResearch = false;

	let selectedAvailableItem: ProductionQueueItem | undefined;
	let selectedAvailableItemCost = getAvailableItemCost();

	let selectedQueueItemIndex = -1;
	let selectedQueueItem: ProductionQueueItem | undefined;
	let selectedQueueItemCost = getSelectedItemCost();

	const resetQueue = () => {
		queueItems = planet.productionQueue?.map((item) => ({ ...item } as ProductionQueueItem));
		availableItems = planetService.getAvailableProductionQueueItems(planet, player);

		selectedAvailableItem = availableItems.length > 0 ? availableItems[0] : selectedAvailableItem;
		contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
	};

	// clone the production queue whenever the planet is updated
	$: planetNum = planet.num && resetQueue();
</script>

<div
	class="flex h-full bg-base-200 shadow-xl max-h-fit min-h-fit rounded-sm border-2 border-base-300"
>
	<div class="flex-col h-full w-full">
		<div class="flex flex-col h-full w-full">
			<div class="flex flex-row h-full w-full grid-cols-3">
				<div class="flex-1 h-full bg-base-100 py-1 px-1">
					<div class="flex flex-col h-full">
						<ul class="grow h-20 overflow-y-auto">
							{#each availableItems as item}
								<li
									on:click={() => availableItemSelected(item)}
									on:dblclick={() => addAvailableItem(item)}
									class="cursor-default select-none hover:text-secondary-focus {item ==
									selectedAvailableItem
										? ' bg-primary'
										: ''}
									{isAuto(item.type) ? ' italic' : ''}"
								>
									{getFullName(item)}
								</li>
							{/each}
						</ul>
						<div class="divider" />
						<div class="h-32">
							{#if selectedAvailableItem}
								<h3>Cost of one {getFullName(selectedAvailableItem)}</h3>
								<CostComponent cost={selectedAvailableItemCost} />
							{/if}
						</div>
					</div>
				</div>
				<div class="flex-none h-full mx-0.5 w-32 px-1">
					<div class="flex-row gap-y-2">
						<button
							on:click={() => addAvailableItem()}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							>Add <Icon
								src={ArrowLongRight}
								size="16"
								class="hover:stroke-accent inline"
							/></button
						>
						<button
							on:click={removeItem}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /> Remove
						</button>
						<button
							on:click={itemUp}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							>Item Up <Icon src={ArrowLongUp} size="16" class="hover:stroke-accent inline" />
						</button>
						<button
							on:click={itemDown}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							>Item Down <Icon src={ArrowLongDown} size="16" class="hover:stroke-accent inline" />
						</button>
						<button
							on:click={clear}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							>Clear <Icon src={XCircle} size="16" class="hover:stroke-accent inline" />
						</button>
					</div>
				</div>
				<div class="flex-1 h-full bg-base-100 py-1 px-1">
					<div class="flex flex-col h-full">
						<ul class="grow h-20 overflow-y-auto">
							<li
								on:click={() => queueItemClicked(-1)}
								class="pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ==
								-1
									? 'bg-primary'
									: ''}"
							>
								-- Top of the Queue --
							</li>
							{#if queueItems}
								{#each queueItems as queueItem, index}
									<li
										on:click={() => queueItemClicked(index, queueItem)}
										class="pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ==
										index
											? 'bg-primary'
											: ''} {isAuto(queueItem.type) ? 'italic' : ''}"
									>
										<div class="flex justify-between ">
											<div>
												{getFullName(queueItem)}
											</div>
											<div>
												{queueItem.quantity}
											</div>
										</div>
									</li>
								{/each}
							{/if}
						</ul>
						<div class="divider" />
						<div class="h-32">
							{#if selectedQueueItem}
								<h3>
									Cost of {getFullName(selectedQueueItem)} x {selectedQueueItem.quantity}
								</h3>
								<CostComponent cost={selectedQueueItemCost} />
							{/if}
						</div>
					</div>
				</div>
			</div>
			<div class="flex justify-end pt-2">
				<div class="grow">
					<label>
						<input
							bind:checked={contributesOnlyLeftoverToResearch}
							class="checkbox-xs"
							type="checkbox"
						/> Contributes Only Leftover to Research
					</label>
				</div>
				<div>
					<button class="btn">Prev</button>
					<button class="btn">Next</button>
					<button on:click={cancel} class="btn">Cancel</button>
					<button on:click={ok} class="btn btn-primary">Ok</button>
				</div>
			</div>
		</div>
	</div>
</div>
