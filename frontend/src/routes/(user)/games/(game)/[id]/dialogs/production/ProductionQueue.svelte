<script lang="ts">
	import CostComponent from '$lib/components/game/Cost.svelte';
	import { getQuantityModifier } from '$lib/quantityModifier';
	import { getGameContext } from '$lib/services/Contexts';
	import { PlanetService } from '$lib/services/PlanetService';
	import type { Cost } from '$lib/types/Cost';
	import { Infinite } from '$lib/types/MapObject';
	import type { CommandedPlanet, ProductionQueueItem } from '$lib/types/Planet';
	import { QueueItemType, isAuto } from '$lib/types/Planet';
	import type { ProductionPlan } from '$lib/types/Player';
	import {
		ArrowLongDown,
		ArrowLongLeft,
		ArrowLongRight,
		ArrowLongUp,
		XCircle
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import hotkeys from 'hotkeys-js';
	import { clamp } from 'lodash-es';
	import { createEventDispatcher, onMount } from 'svelte';

	const { game, player, universe, designs } = getGameContext();

	export let planet: CommandedPlanet;

	let availableItems: ProductionQueueItem[] = [];
	let queueItems: ProductionQueueItem[] = [];
	let contributesOnlyLeftoverToResearch = false;

	let selectedAvailableItem: ProductionQueueItem | undefined;

	let selectedQueueItemIndex = -1;
	let selectedQueueItem: ProductionQueueItem | undefined;

	$: updatedPlanet = planet;

	const getFullName = (item: ProductionQueueItem): string => {
		switch (item.type) {
			case QueueItemType.Starbase:
			case QueueItemType.ShipToken:
				return $universe.getMyDesign(item.designNum)?.name ?? '';
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
			case QueueItemType.PlanetaryScanner:
				return 'Planetary Scanner';
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

	const updateQueueEstimates = async () => {
		// update with estimates from the server
		updatedPlanet.productionQueue = queueItems;
		updatedPlanet = await PlanetService.getPlanetProductionEstimates(updatedPlanet);
		queueItems = updatedPlanet.productionQueue;
		selectedQueueItem = queueItems[selectedQueueItemIndex];
	};

	const addAvailableItem = async (item?: ProductionQueueItem) => {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const max = getMaxBuildable(item.type);
		const quantity = clamp(getQuantityModifier(), 0, max);
		if (quantity == 0) {
			// don't add something we can't build any more of
			return;
		}
		const cost = getItemCost(item) ?? {};
		if (selectedQueueItem) {
			if (selectedQueueItem.type == item?.type && selectedQueueItem.designNum == item?.designNum) {
				selectedQueueItem.quantity += quantity;
			} else {
				// insert a new item

				queueItems.splice(selectedQueueItemIndex + 1, 0, {
					type: item.type,
					quantity,
					designNum: item.designNum,
					costOfOne: cost,
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
					{
						type: item.type,
						designNum: item.designNum,
						costOfOne: cost,
						allocated: {},
						quantity
					},
					...queueItems
				];
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
			}
		}

		updateQueueEstimates();
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
			updateQueueEstimates();
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

	function applyPlan(plan: ProductionPlan | undefined) {
		if (plan) {
			if (queueItems[0].percentComplete) {
				queueItems = [queueItems[0], ...plan.items];
			} else {
				queueItems = plan.items;
			}
			queueItems.forEach((item) => (item.costOfOne = getItemCost(item) ?? {}));
			contributesOnlyLeftoverToResearch = plan.contributesOnlyLeftoverToResearch ?? false;
			updateQueueEstimates();
		}
	}

	const next = () => {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		dispatch('next');
	};

	const prev = () => {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		dispatch('prev');
	};

	const ok = () => {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		dispatch('ok');
	};
	const cancel = () => {
		if (planet) {
			queueItems = planet.productionQueue?.map((item) => ({ ...item } as ProductionQueueItem));
			contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
			dispatch('cancel');
		}
	};

	const getAmountInQueue = (type: QueueItemType): number => {
		return queueItems.reduce((count, i) => count + (i.type === type ? i.quantity : 0), 0);
	};

	const getMaxBuildable = (type: QueueItemType): number => {
		const amountInQueue = getAmountInQueue(type);
		switch (type) {
			case QueueItemType.AutoMines:
				return 1000;
			case QueueItemType.Mine:
				return Math.max(0, (planet.spec.maxPossibleMines ?? 0) - (planet.mines + amountInQueue));
			case QueueItemType.AutoFactories:
				return 1000;
			case QueueItemType.Factory:
				return Math.max(
					0,
					(planet.spec.maxPossibleFactories ?? 0) - (planet.factories + amountInQueue)
				);
			case QueueItemType.AutoDefenses:
				return 100;
			case QueueItemType.Defenses:
				return 100 - planet.defenses + amountInQueue;
			case QueueItemType.AutoMineralAlchemy:
			case QueueItemType.MineralAlchemy:
				return 1000;
			case QueueItemType.AutoMinTerraform:
			case QueueItemType.AutoMaxTerraform:
				return 100;
			case QueueItemType.AutoMineralPacket:
				return 1000;
			case QueueItemType.TerraformEnvironment:
				return (
					Math.abs(planet.spec.terraformAmount?.grav ?? 0) +
					Math.abs(planet.spec.terraformAmount?.temp ?? 0) +
					Math.abs(planet.spec.terraformAmount?.rad ?? 0) -
					amountInQueue
				);
			case QueueItemType.IroniumMineralPacket:
			case QueueItemType.BoraniumMineralPacket:
			case QueueItemType.GermaniumMineralPacket:
			case QueueItemType.MixedMineralPacket:
				return 1000;
			case QueueItemType.ShipToken:
				return 1000;
			case QueueItemType.Starbase:
				return clamp(1 - amountInQueue, 0, 1);

			case QueueItemType.PlanetaryScanner:
				return planet.scanner || amountInQueue > 0 ? 0 : 1;
		}
		return 0;
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
				if (item.designNum) {
					const design = $universe.getMyDesign(item.designNum);
					cost = design?.spec.cost;
				}
				break;
			default:
				if (item && $player.race.spec?.costs) {
					cost = $player.race.spec.costs[item.type];
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

	const getCompletionDescription = (item: ProductionQueueItem) => {
		if (item.skipped) {
			return 'Skipped';
		}

		const yearsToBuildOne = item.yearsToBuildOne ?? 1;
		const yearsToBuildAll = item.yearsToBuildAll ?? 1;
		if (yearsToBuildOne === yearsToBuildAll) {
			if (yearsToBuildAll == 1) {
				return '1 year';
			}
			if (yearsToBuildAll === Infinite) {
				return 'never';
			}
			return `${yearsToBuildAll} years`;
		}
		if (yearsToBuildOne != yearsToBuildAll) {
			if (yearsToBuildAll === Infinite) {
				return `${yearsToBuildOne} to ???`;
			}
			return `${yearsToBuildOne} to ${yearsToBuildAll} years`;
		}
	};

	const dispatch = createEventDispatcher();

	onMount(() => {
		hotkeys('Esc', cancel);
		hotkeys('Enter', ok);
		hotkeys('n', 'production', next);
		hotkeys('p', 'production', prev);
		hotkeys.setScope('production');

		return () => {
			hotkeys.unbind('Esc', cancel);
			hotkeys.unbind('Enter', ok);
			hotkeys.unbind('n', 'production', next);
			hotkeys.unbind('p', 'production', prev);
			hotkeys.deleteScope('production');
		};
	});

	const resetQueue = () => {
		queueItems = planet.productionQueue?.map((item) => ({ ...item } as ProductionQueueItem));
		availableItems = planet.getAvailableProductionQueueItems(planet, $designs);
		selectedAvailableItem = availableItems.length > 0 ? availableItems[0] : selectedAvailableItem;
		contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
	};

	// clone the production queue whenever the planet is updated
	$: planet && resetQueue();
</script>

<div
	class="flex flex-col h-full bg-base-200 shadow max-h-fit min-h-fit rounded-sm border-2 border-base-300 text-base"
>
	<div class="text-center"><h2 class="text-lg">{planet.name}</h2></div>
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
										class:italic={isAuto(item.type)}
										class:bg-primary={item === selectedAvailableItem}
										class="w-full text-left cursor-default select-none hover:text-secondary-focus }
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
						<select
							class="select select-outline select-sm select-secondary w-12 sm:w-full text-secondary"
							on:change|preventDefault={(e) => {
								applyPlan(
									$player.productionPlans.find((p) => p.num == parseInt(e.currentTarget.value))
								);
								e.currentTarget.value = '0';
							}}
						>
							<option value={0}>Apply Plan</option>
							{#each $player.productionPlans as plan}
								<option value={plan.num}>{plan.name}</option>
							{/each}
						</select>
					</div>
				</div>
				<div class="flex-1 h-full bg-base-100 py-1 px-1">
					<div class="flex flex-col h-full">
						<ul class="grow h-20 overflow-y-auto">
							<li>
								<button
									type="button"
									on:click={() => queueItemClicked(-1)}
									class:bg-primary={selectedQueueItemIndex === -1}
									class="w-full pl-1 select-none cursor-default hover:text-secondary-focus"
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
											class:italic={isAuto(queueItem.type)}
											class:text-queue-item-this-year={!queueItem.skipped &&
												(queueItem.yearsToBuildOne ?? 0) <= 1}
											class:text-queue-item-next-year={!queueItem.skipped &&
												((queueItem.yearsToBuildAll ?? 0) > 1 ||
													queueItem.yearsToBuildAll === Infinite) &&
												(queueItem.yearsToBuildOne ?? 0) <= 1}
											class:text-queue-item-skipped={queueItem.skipped}
											class:bg-primary={queueItem === selectedQueueItem}
											class="w-full text-left pl-1 select-none cursor-default hover:text-secondary-focus"
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
								<div class="mt-1 text-base">
									{((selectedQueueItem.percentComplete ?? 0) * 100)?.toFixed()}% Done, Completion {getCompletionDescription(
										selectedQueueItem
									)}
								</div>
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
						<button class="btn btn-sm btn-outline btn-secondary w-full" on:click={prev}>Prev</button
						>
					</div>
					<div class="grow">
						<button class="btn btn-sm btn-outline btn-secondary w-full" on:click={next}>Next</button
						>
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
