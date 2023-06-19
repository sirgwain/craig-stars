<script lang="ts">
	import CostComponent from '$lib/components/game/Cost.svelte';
	import { getQuantityModifier } from '$lib/quantityModifier';
	import { commandMapObject } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import type { Cost } from '$lib/types/Cost';
	import type { Fleet } from '$lib/types/Fleet';
	import type { Game } from '$lib/types/Game';
	import type { CommandedPlanet, ProductionQueueItem } from '$lib/types/Planet';
	import { QueueItemType, isAuto } from '$lib/types/Planet';
	import type { Player, PlayerResponse } from '$lib/types/Player';
	import type { ShipDesign } from '$lib/types/ShipDesign';
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

	export let game: Game;
	export let player: Player;
	export let designs: ShipDesign[];
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
	};

	const queueItemClicked = (index: number, item?: ProductionQueueItem) => {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
	};

	const addAvailableItem = (item?: ProductionQueueItem) => {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const quantity = getQuantityModifier();
		if (selectedQueueItem) {
			if (
				selectedQueueItem.type == item?.type &&
				selectedQueueItem.designName == item?.designName
			) {
				selectedQueueItem.quantity += quantity;
			} else {
				// insert a new item
				queueItems.splice(selectedQueueItemIndex + 1, 0, {
					type: item.type,
					quantity,
					designName: item.designName
				});
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
			}
		} else {
			// prepend a new queue item
			queueItems = [{ type: item.type, designName: item.designName, quantity }, ...queueItems];
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
		const result = await PlanetService.update(game.id, planet);
		Object.assign(planet, result);
		commandMapObject(planet);
		dispatch('ok');
	};
	const cancel = () => {
		if (planet) {
			queueItems = planet.productionQueue?.map((item) => ({ ...item } as ProductionQueueItem));
			contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
			dispatch('cancel');
		}
	};

	/**
	 * Get the cost of a ProductionQueueItem
	 * @param item the item to get cost for
	 * @param quantity the quantity of items to multiply by cost, defaults to 1
	 */
	const getItemCost = (item: ProductionQueueItem | undefined, quantity = 1): Cost | undefined => {
		let cost: Cost | undefined;
		switch (item?.type) {
			case QueueItemType.ShipToken:
			case QueueItemType.Starbase:
				cost = designs.find((d) => d.name === item?.designName)?.spec.cost;
				break;
			default:
				if (item && player.race) {
					cost = player.race.spec?.costs[item.type];
				}
		}
		return cost
			? {
					ironium: (cost.ironium ?? 0) * quantity,
					boranium: (cost.boranium ?? 0) * quantity,
					germanium: (cost.germanium ?? 0) * quantity,
					resources: (cost.resources ?? 0) * quantity
			  }
			: undefined;
	};

	const dispatch = createEventDispatcher();

	hotkeys('Esc', () => cancel());
	hotkeys('Enter', () => {
		ok();
	});

	let availableItems: ProductionQueueItem[] = [];
	let queueItems: ProductionQueueItem[] = [];
	let contributesOnlyLeftoverToResearch = false;

	let selectedAvailableItem: ProductionQueueItem | undefined;

	let selectedQueueItemIndex = -1;
	let selectedQueueItem: ProductionQueueItem | undefined;

	const resetQueue = () => {
		queueItems = planet.productionQueue?.map((item) => ({ ...item } as ProductionQueueItem));
		availableItems = planet.getAvailableProductionQueueItems(planet, designs ?? []);
		selectedAvailableItem = availableItems.length > 0 ? availableItems[0] : selectedAvailableItem;
		contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
	};

	// clone the production queue whenever the planet is updated
	$: planet && designs && resetQueue();
</script>

<div
	class="flex h-full bg-base-200 shadow max-h-fit min-h-fit rounded-sm border-2 border-base-300"
>
	<div class="flex-col h-full w-full">
		<div class="flex flex-col h-full w-full">
			<div class="flex flex-row h-full w-full grid-cols-3">
				<div class="flex-1 h-full bg-base-100 py-1 px-1">
					<div class="flex flex-col h-full">
						<ul class="grow h-20 overflow-y-auto">
							{#each availableItems as item}
								<li>
									<button
										type="button"
										on:click={() => availableItemSelected(item)}
										on:dblclick={() => addAvailableItem(item)}
										class="w-full text-left cursor-default select-none hover:text-secondary-focus {item ==
										selectedAvailableItem
											? ' bg-primary'
											: ''}
									{isAuto(item.type) ? ' italic' : ''}"
									>
										{getFullName(item)}
									</button>
								</li>
							{/each}
						</ul>
						<div class="divider" />
						<div class="h-32">
							{#if selectedAvailableItem}
								<h3>Cost of one {getFullName(selectedAvailableItem)}</h3>
								<CostComponent cost={getItemCost(selectedAvailableItem)} />
							{/if}
						</div>
					</div>
				</div>
				<div class="flex-none h-full mx-0.5 md:w-32 px-1">
					<div class="flex-row flex-none gap-y-2">
						<button
							on:click={() => addAvailableItem()}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Add </span><Icon
								src={ArrowLongRight}
								size="16"
								class="hover:stroke-accent inline"
							/></button
						>
						<button
							on:click={removeItem}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /><span
								class="hidden sm:inline"
							>
								Remove</span
							>
						</button>
						<button
							on:click={itemUp}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Item Up </span><Icon
								src={ArrowLongUp}
								size="16"
								class="hover:stroke-accent inline"
							/>
						</button>
						<button
							on:click={itemDown}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Item Down </span><Icon
								src={ArrowLongDown}
								size="16"
								class="hover:stroke-accent inline"
							/>
						</button>
						<button
							on:click={clear}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Clear </span><Icon
								src={XCircle}
								size="16"
								class="hover:stroke-accent inline"
							/>
						</button>
					</div>
				</div>
				<div class="flex-1 h-full bg-base-100 py-1 px-1">
					<div class="flex flex-col h-full">
						<ul class="grow h-20 overflow-y-auto">
							<li>
								<button
									type="button"
									on:click={() => queueItemClicked(-1)}
									class="w-full pl-1 select-none cursor-default hover:text-secondary-focus {selectedQueueItemIndex ==
									-1
										? 'bg-primary'
										: ''}"
								>
									-- Top of the Queue --
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
													{getFullName(queueItem)}
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
						<div class="divider" />
						<div class="h-32">
							{#if selectedQueueItem}
								<h3>
									Cost of {getFullName(selectedQueueItem)} x {selectedQueueItem.quantity}
								</h3>
								<CostComponent cost={getItemCost(selectedQueueItem, selectedQueueItem?.quantity)} />
							{/if}
						</div>
					</div>
				</div>
			</div>
			<div class="flex justify-between pt-2">
				<div class="w-1/2 mr-14">
					<label>
						<input
							bind:checked={contributesOnlyLeftoverToResearch}
							class="checkbox checkbox-xs"
							type="checkbox"
						/> Contributes Only Leftover to Research
					</label>
				</div>
				<div class="w-1/2 flex flex-row flex-wrap justify-between sm:justify-end">
					<div class="grow">
						<button class="btn btn-sm btn-outline btn-secondary w-full">Prev</button>
					</div>
					<div class="grow">
						<button class="btn btn-sm btn-outline btn-secondary w-full">Next</button>
					</div>
					<div class="grow">
						<button on:click={cancel} class="btn btn-sm btn-outline btn-secondary w-full"
							>Cancel</button
						>
					</div>
					<div class="grow">
						<button on:click={ok} class="btn btn-sm btn-primary w-full">Ok</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
