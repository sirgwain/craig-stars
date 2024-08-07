<script lang="ts" context="module">
	export type ProductionQueueEvent = {
		next: void;
		prev: void;
		ok: void;
		cancel: void;
	};
</script>

<script lang="ts">
	import CostComponent from '$lib/components/game/Cost.svelte';
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import QuantityModifierButtons from '$lib/components/QuantityModifierButtons.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { NeverBuilt, getProductionEstimates } from '$lib/services/Producer';
	import { techs } from '$lib/services/Stores';
	import { divide, multiply, type Cost } from '$lib/types/Cost';
	import { CommandedPlanet } from '$lib/types/Planet';
	import type { ProductionPlan } from '$lib/types/Player';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { getFullName, isAuto } from '$lib/types/QueueItemType';
	import { getPlanetHabitability } from '$lib/types/Race';
	import {
		ArrowLongDown,
		ArrowLongLeft,
		ArrowLongRight,
		ArrowLongUp,
		QuestionMarkCircle,
		XCircle
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import hotkeys from 'hotkeys-js';
	import { clamp } from 'lodash-es';
	import { createEventDispatcher, onMount } from 'svelte';

	const { game, player, universe } = getGameContext();
	const dispatch = createEventDispatcher<ProductionQueueEvent>();

	export let planet: CommandedPlanet;

	let availableItems: ProductionQueueItem[] = [];
	let availableShipDesigns: ProductionQueueItem[] = [];
	let availableStarbaseDesigns: ProductionQueueItem[] = [];
	let queueItems: ProductionQueueItem[] = [];
	let contributesOnlyLeftoverToResearch = false;

	let selectedAvailableItem: ProductionQueueItem | undefined;
	let selectedAvailableItemCost: Cost | undefined;

	let selectedQueueItemIndex = -1;
	let selectedQueueItem: ProductionQueueItem | undefined;
	let selectedQueueItemCost: Cost | undefined;
	$: selectedQueueItemPercentComplete = selectedQueueItem
		? getPercentComplete(selectedQueueItem)
		: 0;

	$: updatedPlanet = Object.assign(new CommandedPlanet(), planet);

	// keep track of the quantity modifier
	let quantityModifer = 1;

	function availableItemSelected(type: ProductionQueueItem) {
		selectedAvailableItem = type;
		selectedAvailableItemCost = $player.getItemCost(
			selectedAvailableItem,
			$universe,
			$techs,
			planet
		);
	}

	function queueItemClicked(index: number, item?: ProductionQueueItem) {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
		selectedQueueItemCost = $player.getItemCost(
			selectedQueueItem,
			$universe,
			$techs,
			planet,
			selectedQueueItem?.quantity
		);
	}

	function updateQueueEstimates() {
		// get updated production queue estimates
		updatedPlanet.productionQueue = [...queueItems];
		const itemEstimates = getProductionEstimates(
			$game.rules,
			$techs,
			$player,
			updatedPlanet,
			$universe
		);

		for (let i = 0; i < queueItems.length; i++) {
			const estimate = itemEstimates[i];
			Object.assign(queueItems[i], {
				yearsToBuildOne: estimate.yearsToBuildOne,
				yearsToBuildAll: estimate.yearsToBuildAll,
				yearsToSkipAuto: estimate.yearsToSkipAuto
			});
		}

		// update the reactive variable so the UI updates
		queueItems = updatedPlanet.productionQueue;

		selectedQueueItem = queueItems[selectedQueueItemIndex];
		selectedQueueItemCost = multiply(
			$player.getItemCost(selectedQueueItem, $universe, $techs, planet),
			selectedQueueItem?.quantity
		);
	}

	function getPercentComplete(item: ProductionQueueItem): number {
		if ((item.allocated?.resources ?? 0) === 0) {
			return 0;
		}

		const cost = $player.getItemCost(item, $universe, $techs, planet, item.quantity);
		const resourcePercent = cost.resources
			? (item.allocated?.resources ?? 0) / cost.resources
			: 1.0;
		const mineralsPercent = divide(planet.cargo, { ...item.allocated, resources: 0 });

		// if we are mineral or resource constrained, report the percent complete based on the lowest.
		return Math.min(resourcePercent, mineralsPercent);
	}

	function addAvailableItem(e: MouseEvent, item?: ProductionQueueItem) {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const maxPopulation = planet.getMaxPopulation(
			$game.rules,
			$player,
			getPlanetHabitability($player.race, planet.hab)
		);
		const amountInQueue = planet.getAmountInQueue(item.type, queueItems);
		// get the max number of items we can build on this planet. For auto items, let them add 5k because it's ok to
		// add more than our auto items will build. This getMaxBuildable function returns the number of usuable mines for auto, but
		// when updating the production queue we don't care about that
		const max = isAuto(item.type)
			? 5000
			: planet.getMaxBuildable($techs, $player, maxPopulation, item.type, amountInQueue);
		const quantity = clamp(quantityModifer, 0, max);
		if (quantity == 0) {
			// don't add something we can't build any more of
			return;
		}
		const cost = $player.getItemCost(item, $universe, $techs, planet) ?? {};
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
				selectedQueueItemCost = $player.getItemCost(
					selectedQueueItem,
					$universe,
					$techs,
					planet,
					selectedQueueItem?.quantity
				);
			}
		} else {
			let nextItem = queueItems.length ? queueItems[0] : undefined;
			if (nextItem && nextItem.type === item?.type && nextItem.designNum == item.designNum) {
				nextItem.quantity++;
				selectedQueueItemIndex = 0;
				selectedQueueItem = nextItem;
				selectedQueueItemCost = $player.getItemCost(
					selectedQueueItem,
					$universe,
					$techs,
					planet,
					selectedQueueItem?.quantity
				);
			} else {
				// prepend a new queue item
				queueItems = [
					{
						type: item.type,
						designNum: item.designNum,
						allocated: {},
						quantity
					},
					...queueItems
				];
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
				selectedQueueItemCost = $player.getItemCost(
					selectedQueueItem,
					$universe,
					$techs,
					planet,
					selectedQueueItem?.quantity
				);
			}
		}

		updateQueueEstimates();
	}

	function removeItem(e: MouseEvent) {
		if (queueItems && selectedQueueItem) {
			selectedQueueItem.quantity -= quantityModifer;
			selectedQueueItem.quantity = Math.max(0, selectedQueueItem.quantity);
			queueItems = queueItems;
			if (selectedQueueItem.quantity <= 0) {
				// select the item up in the list
				queueItems = queueItems?.filter((item) => item != selectedQueueItem);
				selectedQueueItem =
					queueItems[selectedQueueItemIndex > -1 ? selectedQueueItemIndex - 1 : 0];
				selectedQueueItemCost = $player.getItemCost(
					selectedQueueItem,
					$universe,
					$techs,
					planet,
					selectedQueueItem?.quantity
				);

				selectedQueueItemIndex--;
			}
			updateQueueEstimates();
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
		selectedQueueItemCost = {};
	}

	function applyPlan(plan: ProductionPlan | undefined) {
		if (plan) {
			const concreteItems = queueItems.filter((i) => !isAuto(i.type));
			queueItems = [...concreteItems, ...plan.items];
			contributesOnlyLeftoverToResearch = plan.contributesOnlyLeftoverToResearch ?? false;
			updateQueueEstimates();
		}
	}

	function next() {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		dispatch('next');
	}

	function prev() {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		dispatch('prev');
	}

	function ok() {
		planet.productionQueue = queueItems ?? [];
		planet.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		dispatch('ok');
	}
	function cancel() {
		if (planet) {
			resetQueue();
			dispatch('cancel');
		}
	}

	function getCompletionDescription(item: ProductionQueueItem) {
		const skipped =
			isAuto(item.type) && item.yearsToBuildOne == NeverBuilt && item.yearsToBuildAll == NeverBuilt;
		if (skipped) {
			return 'Skipped';
		}

		const yearsToBuildOne = item.yearsToBuildOne ?? 1;
		const yearsToBuildAll = isAuto(item.type) ? item.yearsToSkipAuto : item.yearsToBuildAll;
		if (yearsToBuildOne === yearsToBuildAll) {
			if (yearsToBuildAll == 1) {
				return '1 year';
			}
			if (yearsToBuildAll === NeverBuilt) {
				return 'never';
			}
			return `${yearsToBuildAll} years`;
		}
		if (yearsToBuildAll && yearsToBuildOne != yearsToBuildAll) {
			if (yearsToBuildAll === NeverBuilt) {
				return `${yearsToBuildOne} to ???`;
			}
			return `${yearsToBuildOne} to ${yearsToBuildAll} years`;
		}

		if (yearsToBuildOne == 1) {
			return '1 year';
		}
		if (yearsToBuildOne === NeverBuilt) {
			return 'never';
		}
		return `${yearsToBuildOne} years`;
	}

	onMount(() => {
		const originalScope = hotkeys.getScope();
		const scope = 'production';
		hotkeys('Esc', cancel);
		hotkeys('Enter', ok);
		hotkeys('n', scope, next);
		hotkeys('p', scope, prev);
		hotkeys.setScope(scope);

		return () => {
			hotkeys.unbind('Esc', cancel);
			hotkeys.unbind('Enter', ok);
			hotkeys.unbind('n', scope, next);
			hotkeys.unbind('p', scope, prev);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});

	function resetQueue() {
		queueItems = [...planet.productionQueue?.map((item) => ({ ...item }) as ProductionQueueItem)];
		availableItems = planet.getAvailableProductionQueueItems(
			planet,
			$player.race.spec?.innateMining,
			$player.race.spec?.innateResources,
			$player.race.spec?.livesOnStarbases
		);
		availableShipDesigns = planet.getAvailableProductionQueueShipDesigns(planet, $universe.designs);
		availableStarbaseDesigns = planet.getAvailableProductionQueueStarbaseDesigns(
			planet,
			$universe.designs
		);
		if (availableShipDesigns.length > 0) {
			selectedAvailableItem = availableShipDesigns[0];
		} else if (availableStarbaseDesigns.length > 0) {
			selectedAvailableItem = availableStarbaseDesigns[0];
		} else if (availableItems.length > 0) {
			selectedAvailableItem = availableItems[0];
		}
		selectedAvailableItemCost = $player.getItemCost(
			selectedAvailableItem,
			$universe,
			$techs,
			planet
		);
		contributesOnlyLeftoverToResearch = planet.contributesOnlyLeftoverToResearch ?? false;
		updateQueueEstimates();
	}

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
							{#if availableShipDesigns.length > 0}
								<li class="font-semibold text-secondary text-lg border-b border-b-secondary mb-0.5">
									Ships
								</li>
								{#each availableShipDesigns as item}
									<li>
										<button
											type="button"
											on:click={() => availableItemSelected(item)}
											on:dblclick={(e) => addAvailableItem(e, item)}
											on:contextmenu|preventDefault={(e) =>
												onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
											class:italic={isAuto(item.type)}
											class:bg-primary={item === selectedAvailableItem}
											class:text-queue-item-this-year={(item.yearsToBuildOne ?? 0) == 1}
											class:text-queue-item-next-year={(item.yearsToBuildOne ?? 0) == 2}
											class:text-queue-item-never={(item.yearsToBuildOne ?? 0) == NeverBuilt}
											class="w-full pl-0.5 text-left cursor-default select-none hover:text-secondary-focus }
									{isAuto(item.type) ? ' italic' : ''}"
										>
											{getFullName(item, $universe)}
										</button>
									</li>
								{/each}
							{/if}

							{#if availableStarbaseDesigns.length > 0}
								<li class="font-semibold text-secondary text-lg border-b border-b-secondary my-0.5">
									Starbases
								</li>
								{#each availableStarbaseDesigns as item}
									<li>
										<button
											type="button"
											on:click={() => availableItemSelected(item)}
											on:dblclick={(e) => addAvailableItem(e, item)}
											on:contextmenu|preventDefault={(e) =>
												onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
											class:italic={isAuto(item.type)}
											class:bg-primary={item === selectedAvailableItem}
											class:text-queue-item-this-year={(item.yearsToBuildOne ?? 0) == 1}
											class:text-queue-item-next-year={(item.yearsToBuildOne ?? 0) == 2}
											class:text-queue-item-never={(item.yearsToBuildOne ?? 0) == NeverBuilt}
											class="w-full pl-0.5 text-left cursor-default select-none hover:text-secondary-focus }
									{isAuto(item.type) ? ' italic' : ''}"
										>
											{getFullName(item, $universe)}
										</button>
									</li>
								{/each}
							{/if}
							<li class="font-semibold text-secondary text-lg border-b border-b-secondary mb-0.5">
								Planetary Structures
							</li>
							{#each availableItems as item}
								<li>
									<button
										type="button"
										on:click={() => availableItemSelected(item)}
										on:dblclick={(e) => addAvailableItem(e, item)}
										class:italic={isAuto(item.type)}
										class:bg-primary={item === selectedAvailableItem}
										class="w-full pl-0.5 text-left cursor-default select-none hover:text-secondary-focus }
									{isAuto(item.type) ? ' italic' : ''}"
									>
										{getFullName(item, $universe)}
									</button>
								</li>
							{/each}
						</ul>
						<div class="divider" />
						<div class="h-32">
							{#if selectedAvailableItem && selectedAvailableItemCost}
								<h3>
									{#if selectedAvailableItem.designNum}
										<button
											type="button"
											on:pointerdown={(e) =>
												onShipDesignTooltip(
													e,
													$universe.getMyDesign(selectedAvailableItem?.designNum)
												)}
											>Cost of {getFullName(selectedAvailableItem, $universe)}<Icon
												src={QuestionMarkCircle}
												size="16"
												class="cursor-help inline-block ml-1"
											/></button
										>
									{:else}
										Cost of {getFullName(selectedAvailableItem, $universe)}
									{/if}
								</h3>
								<CostComponent cost={selectedAvailableItemCost} />
								{#if selectedAvailableItem.yearsToBuildOne}
									Completion {getCompletionDescription(selectedAvailableItem)}
								{/if}
							{/if}
						</div>
					</div>
				</div>
				<div class="flex-none h-full mx-0.5 md:w-34 px-1">
					<div class="flex-row flex-none gap-y-2">
						<button
							on:click={(e) => addAvailableItem(e)}
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
						<div class="flex flex-col sm:flex-row justify-between mt-2 gap-1 mx-1">
							<QuantityModifierButtons bind:modifier={quantityModifer} />
						</div>
					</div>
				</div>
				<div class="flex-1 h-full bg-base-100 py-1">
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
									<li class="cursor-default">
										<ProductionQueueItemLine
											item={queueItem}
											{index}
											on:queue-item-clicked={() => queueItemClicked(index, queueItem)}
											selected={queueItem === selectedQueueItem}
										/>
									</li>
								{/each}
							{/if}
						</ul>
						<div class="divider" />
						<div class="h-32">
							{#if selectedQueueItem}
								<h3>
									{#if selectedQueueItem.designNum}
										<button
											type="button"
											on:pointerdown={(e) =>
												onShipDesignTooltip(e, $universe.getMyDesign(selectedQueueItem?.designNum))}
											>Cost of {getFullName(selectedQueueItem, $universe)} x {selectedQueueItem.quantity}<Icon
												src={QuestionMarkCircle}
												size="16"
												class="cursor-help inline-block ml-1"
											/></button
										>
									{:else}
										Cost of {getFullName(selectedQueueItem, $universe)} x {selectedQueueItem.quantity}
									{/if}
								</h3>
								<CostComponent cost={selectedQueueItemCost} />
								<div class="mt-1 text-base">
									{#if selectedQueueItemPercentComplete}
										{(selectedQueueItemPercentComplete * 100)?.toFixed()}% Done,
									{/if}
									Completion {getCompletionDescription(selectedQueueItem)}
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
			<div class="flex justify-between p-1 pt-2">
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
