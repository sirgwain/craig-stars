<script lang="ts" module>
	export type ProductionQueueEvent = {
		next: void;
		prev: void;
		ok: void;
		cancel: void;
	};
</script>

<script lang="ts">
	import { asyncToVoidWrapper } from '$lib/asyncToVoid';
	import CostComponent from '$lib/components/game/Cost.svelte';
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import { onAllocatedTooltip } from '$lib/components/game/tooltips/AllocatedTooltip.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import QuantityModifierButtons from '$lib/components/QuantityModifierButtons.svelte';
	import type { CostJson as Cost } from '$lib/types/cs-proto';
	import {
		ProductionQueueItemSchema,
		QueueItemCompletionEstimateSchema,
		type ProductionQueueItem
	} from '$lib/types/cs-proto';
	import { type ProductionPlan } from '$lib/types/cs-proto';
	import { addError, CSError } from '$lib/services/Errors';
	import type { OnCancel, OnOk } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import { divide, multiply } from '$lib/types/Cost';
	import { CommandedPlanet } from '$lib/types/Planet';
	import { getFullName, isAuto } from '$lib/types/QueueItemType';
	import { clone, create } from '@bufbuild/protobuf';
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
	import { onMount } from 'svelte';
	import type { ChangeEventHandler } from 'svelte/elements';
	import { Infinite } from '$lib/types/Consts';

	// used to load the Genesis Device tech
	const GenesisDevice = 'Genesis Device';

	const { cs, player, universe } = getGameContext();

	type Props = {
		planet: CommandedPlanet;
		onOk?: OnOk<CommandedPlanet>;
		onCancel?: OnCancel;
		onNext?: () => Promise<void>;
		onPrev?: () => Promise<void>;
	};

	let { planet, onOk, onCancel, onNext, onPrev }: Props = $props();

	let updatedPlanet = $derived(new CommandedPlanet(planet));

	let availableItems: ProductionQueueItem[] = $state([]);
	let availableShipDesigns: ProductionQueueItem[] = $state([]);
	let availableStarbaseDesigns: ProductionQueueItem[] = $state([]);
	let queueItems: ProductionQueueItem[] = $state([]);
	let contributesOnlyLeftoverToResearch = $state(false);

	let selectedAvailableItem: ProductionQueueItem | undefined = $state();
	let selectedAvailableItemCost: Cost | undefined = $state();

	let selectedQueueItemIndex = $state(-1);
	let selectedQueueItem: ProductionQueueItem | undefined = $state();
	let selectedQueueItemCost: Cost | undefined = $state();

	// keep track of the quantity modifier
	let quantityModifer = $state(1);

	let selectedQueueItemPercentComplete = $state(0);
	$effect(() => {
		if (!selectedQueueItem) {
			return;
		}
		getPercentComplete(selectedQueueItem).then(
			(result) => (selectedQueueItemPercentComplete = result)
		);
	});

	async function availableItemSelected(type: ProductionQueueItem) {
		selectedAvailableItem = type;
		selectedAvailableItemCost = await $player.getItemCost(
			cs,
			selectedAvailableItem,
			$universe,
			planet
		);
	}

	async function onQueueItemClicked(index: number, item?: ProductionQueueItem) {
		selectedQueueItemIndex = index;
		selectedQueueItem = item;
		selectedQueueItemCost = await $player.getItemCost(
			cs,
			selectedQueueItem,
			$universe,
			planet,
			selectedQueueItem?.quantity
		);
	}

	const contributesOnlyLeftoverToResearchChecked: ChangeEventHandler<HTMLInputElement> = (e) => {
		contributesOnlyLeftoverToResearch = e.currentTarget.checked;
		updateQueueEstimates();
	};

	async function updateQueueEstimates() {
		// get updated production queue estimates
		updatedPlanet.planetOrders.productionQueue = [...queueItems];
		updatedPlanet.planetOrders.contributesOnlyLeftoverToResearch =
			contributesOnlyLeftoverToResearch;
		const { planet: planetWithEstimates } = await cs.wasmService.estimateProduction({
			planet: updatedPlanet
		});
		if (!planetWithEstimates?.planetOrders?.productionQueue) {
			addError(new CSError(undefined, 'unable to estimate production', 0));
			return;
		}

		for (let i = 0; i < queueItems.length; i++) {
			const estimate = planetWithEstimates.planetOrders.productionQueue[i];
			queueItems[i].queueItemCompletionEstimate = estimate.queueItemCompletionEstimate;
		}

		// update the reactive variable so the UI updates
		queueItems = updatedPlanet.planetOrders.productionQueue;

		selectedQueueItem = queueItems[selectedQueueItemIndex];
		const itemCost = await $player.getItemCost(cs, selectedQueueItem, $universe, planet);
		selectedQueueItemCost = multiply(itemCost, selectedQueueItem?.quantity);

		for (let i = 0; i < availableItems.length; i++) {
			const item = availableItems[i];
			if (!item.queueItemCompletionEstimate?.yearsToBuildOne) {
				item.queueItemCompletionEstimate = create(QueueItemCompletionEstimateSchema, {
					yearsToBuildOne: await updatedPlanet.getYearsToBuildOne(item, cs)
				});
			}
			if (selectedAvailableItem == item) {
				selectedAvailableItem = item;
			}
		}
		availableItems = [...availableItems];

		for (let i = 0; i < availableShipDesigns.length; i++) {
			const item = availableShipDesigns[i];
			if (!item.queueItemCompletionEstimate?.yearsToBuildOne) {
				item.queueItemCompletionEstimate = create(QueueItemCompletionEstimateSchema, {
					yearsToBuildOne: await updatedPlanet.getYearsToBuildOne(item, cs)
				});
			}
			if (selectedAvailableItem == item) {
				selectedAvailableItem = item;
			}
		}
		availableShipDesigns = [...availableShipDesigns];

		for (let i = 0; i < availableStarbaseDesigns.length; i++) {
			const item = availableStarbaseDesigns[i];
			if (!item.queueItemCompletionEstimate?.yearsToBuildOne) {
				item.queueItemCompletionEstimate = create(QueueItemCompletionEstimateSchema, {
					yearsToBuildOne: await updatedPlanet.getYearsToBuildOne(item, cs)
				});
			}

			if (selectedAvailableItem == item) {
				selectedAvailableItem = item;
			}
		}
		availableStarbaseDesigns = [...availableStarbaseDesigns];
	}

	async function getPercentComplete(item: ProductionQueueItem): Promise<number> {
		if ((item.allocated?.resources ?? 0) === 0) {
			return 0;
		}

		const costOfOne = await $player.getItemCost(cs, item, $universe, planet, 1);
		const percent = divide(item.allocated ?? {}, costOfOne);

		// if we are mineral or resource constrained, report the percent complete based on the lowest.
		return percent;
	}

	async function addAvailableItem(item?: ProductionQueueItem) {
		item = item ?? selectedAvailableItem;
		if (!queueItems || !item) {
			return;
		}

		const amountInQueue = planet.getAmountInQueue(item.type, queueItems);
		const { result: maxBuildable } =
			(await cs.wasmService.getMaxBuildable({ planet: planet, itemType: item.type })) ??
			0 - amountInQueue;
		const quantity = clamp(quantityModifer, 0, maxBuildable);
		if (quantity == 0) {
			// don't add something we can't build any more of
			return;
		}
		if (selectedQueueItem) {
			if (selectedQueueItem.type == item?.type && selectedQueueItem.designNum == item?.designNum) {
				selectedQueueItem.quantity += quantity;
			} else {
				// insert a new item

				queueItems.splice(
					selectedQueueItemIndex + 1,
					0,
					create(ProductionQueueItemSchema, {
						type: item.type,
						quantity,
						designNum: item.designNum
					})
				);
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
				selectedQueueItemCost = await $player.getItemCost(
					cs,
					selectedQueueItem,
					$universe,
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
				selectedQueueItemCost = await $player.getItemCost(
					cs,
					selectedQueueItem,
					$universe,
					planet,
					selectedQueueItem?.quantity
				);
			} else {
				// prepend a new queue item
				queueItems = [
					create(ProductionQueueItemSchema, {
						type: item.type,
						designNum: item.designNum,
						quantity
					}),
					...queueItems
				];
				selectedQueueItemIndex++;
				selectedQueueItem = queueItems[selectedQueueItemIndex];
				selectedQueueItemCost = await $player.getItemCost(
					cs,
					selectedQueueItem,
					$universe,
					planet,
					selectedQueueItem?.quantity
				);
			}
		}

		updateQueueEstimates();
	}

	async function removeItem() {
		if (queueItems && selectedQueueItem) {
			selectedQueueItem.quantity -= quantityModifer;
			selectedQueueItem.quantity = Math.max(0, selectedQueueItem.quantity);
			queueItems = queueItems;
			if (selectedQueueItem.quantity <= 0) {
				// select the item up in the list
				queueItems = queueItems?.filter((item) => item != selectedQueueItem);
				selectedQueueItem =
					queueItems[selectedQueueItemIndex > -1 ? selectedQueueItemIndex - 1 : 0];
				selectedQueueItemCost = await $player.getItemCost(
					cs,
					selectedQueueItem,
					$universe,
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
			updateQueueEstimates();
		}
	}

	function itemDown() {
		if (queueItems && selectedQueueItem && selectedQueueItemIndex < queueItems.length - 1) {
			const swap = queueItems[selectedQueueItemIndex + 1];
			queueItems[selectedQueueItemIndex + 1] = selectedQueueItem;
			queueItems[selectedQueueItemIndex] = swap;
			selectedQueueItemIndex++;
			queueItems = queueItems;
			updateQueueEstimates();
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
			queueItems = [
				...concreteItems,
				...plan.items.map((item) =>
					create(ProductionQueueItemSchema, {
						type: item.type,
						quantity: item.quantity,
						designNum: item.designNum
					})
				)
			];
			contributesOnlyLeftoverToResearch = plan.contributesOnlyLeftoverToResearch ?? false;
			updateQueueEstimates();
		}
	}

	async function next() {
		planet.planetOrders.productionQueue = queueItems ?? [];
		planet.planetOrders.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		await onNext?.();
		await resetQueue();
	}

	async function prev() {
		planet.planetOrders.productionQueue = queueItems ?? [];
		planet.planetOrders.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		await onPrev?.();
		await resetQueue();
	}

	function ok() {
		planet.planetOrders.productionQueue = queueItems ?? [];
		planet.planetOrders.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
		onOk?.(planet);
	}
	function cancel() {
		if (planet) {
			resetQueue();
			onCancel?.();
		}
	}

	function getCompletionDescription(item: ProductionQueueItem) {
		const skipped =
			isAuto(item.type) &&
			item.queueItemCompletionEstimate?.yearsToBuildOne === Infinite &&
			item.queueItemCompletionEstimate?.yearsToBuildAll === Infinite;
		if (skipped) {
			return 'Skipped';
		}

		const yearsToBuildOne = item.queueItemCompletionEstimate?.yearsToBuildOne ?? 1;
		const yearsToBuildAll = isAuto(item.type)
			? item.queueItemCompletionEstimate?.yearsToSkipAuto
			: item.queueItemCompletionEstimate?.yearsToBuildAll;
		if (yearsToBuildOne === yearsToBuildAll) {
			if (yearsToBuildAll == 1) {
				return '1 year';
			}
			if (yearsToBuildAll === Infinite) {
				return 'never';
			}
			return `${yearsToBuildAll} years`;
		}
		if (yearsToBuildAll && yearsToBuildOne != yearsToBuildAll) {
			if (yearsToBuildAll === Infinite) {
				return `${yearsToBuildOne} to ???`;
			}
			return `${yearsToBuildOne} to ${yearsToBuildAll} years`;
		}

		if (yearsToBuildOne == 1) {
			return '1 year';
		}
		if (yearsToBuildOne === Infinite) {
			return 'never';
		}
		return `${yearsToBuildOne} years`;
	}

	async function resetQueue() {
		contributesOnlyLeftoverToResearch = planet.planetOrders.contributesOnlyLeftoverToResearch;
		queueItems = [
			...planet.planetOrders.productionQueue.map((item) => clone(ProductionQueueItemSchema, item))
		];
		const genesisDevice = $techs.getTech(GenesisDevice);
		availableItems = planet.getAvailableProductionQueueItems(
			$player.race.spec.innateMining,
			$player.race.spec.innateResources,
			$player.race.spec.livesOnStarbases,
			genesisDevice && $player.hasTech(genesisDevice)
		);
		availableShipDesigns = planet.getAvailableProductionQueueShipDesigns($universe.designs);
		availableStarbaseDesigns = planet.getAvailableProductionQueueStarbaseDesigns($universe.designs);
		if (availableShipDesigns.length > 0) {
			selectedAvailableItem = availableShipDesigns[0];
		} else if (availableStarbaseDesigns.length > 0) {
			selectedAvailableItem = availableStarbaseDesigns[0];
		} else if (availableItems.length > 0) {
			selectedAvailableItem = availableItems[0];
		}
		selectedAvailableItemCost = await $player.getItemCost(
			cs,
			selectedAvailableItem,
			$universe,
			planet
		);
		contributesOnlyLeftoverToResearch =
			planet.planetOrders.contributesOnlyLeftoverToResearch ?? false;
		await updateQueueEstimates();
	}

	onMount(() => {
		const originalScope = hotkeys.getScope();
		const scope = 'production';
		const syncNext = asyncToVoidWrapper(next);
		const syncPrev = asyncToVoidWrapper(prev);
		hotkeys('Esc', scope, cancel);
		hotkeys('Enter', scope, ok);
		hotkeys('n', scope, syncNext);
		hotkeys('p', scope, syncPrev);
		hotkeys.setScope(scope);

		resetQueue();

		return () => {
			hotkeys.unbind('Esc', scope, cancel);
			hotkeys.unbind('Enter', scope, ok);
			hotkeys.unbind('n', scope, syncNext);
			hotkeys.unbind('p', scope, syncPrev);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

<div class="flex flex-col h-full bg-base-200 shadow rounded-sm border-2 border-base-300 text-base">
	<div class="text-center"><h2 class="text-lg">{planet.mapObject.name}</h2></div>
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
								{#each availableShipDesigns as item (item.designNum)}
									{@const estimate =
										item.queueItemCompletionEstimate ?? create(QueueItemCompletionEstimateSchema)}
									<li>
										<button
											type="button"
											onclick={() => availableItemSelected(item)}
											ondblclick={() => addAvailableItem(item)}
											oncontextmenu={(e) =>
												onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
											class:italic={isAuto(item.type)}
											class:bg-primary={item === selectedAvailableItem}
											class:text-queue-item-this-year={estimate.yearsToBuildOne == 1}
											class:text-queue-item-next-year={estimate.yearsToBuildOne == 2}
											class:text-queue-item-never={estimate.yearsToBuildOne == Infinite}
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
								{#each availableStarbaseDesigns as item (item.designNum)}
									{@const estimate =
										item.queueItemCompletionEstimate ?? create(QueueItemCompletionEstimateSchema)}
									<li>
										<button
											type="button"
											onclick={() => availableItemSelected(item)}
											ondblclick={() => addAvailableItem(item)}
											oncontextmenu={(e) =>
												onShipDesignTooltip(e, $universe.getMyDesign(item.designNum))}
											class:italic={isAuto(item.type)}
											class:bg-primary={item === selectedAvailableItem}
											class:text-queue-item-this-year={estimate.yearsToBuildOne == 1}
											class:text-queue-item-next-year={estimate.yearsToBuildOne == 2}
											class:text-queue-item-never={estimate.yearsToBuildOne == Infinite}
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
							{#each availableItems as item (item.type)}
								<li>
									<button
										type="button"
										onclick={() => availableItemSelected(item)}
										ondblclick={() => addAvailableItem(item)}
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
						<div class="divider"></div>
						<div class="h-32">
							{#if selectedAvailableItem && selectedAvailableItemCost}
								<h3>
									{#if selectedAvailableItem.designNum}
										<button
											type="button"
											onpointerdown={(e) =>
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
								{#if selectedAvailableItem.queueItemCompletionEstimate?.yearsToBuildOne}
									Completion {getCompletionDescription(selectedAvailableItem)}
								{/if}
							{/if}
						</div>
					</div>
				</div>
				<div class="flex-none h-full mx-0.5 md:w-34 px-1">
					<div class="flex-row flex-none gap-y-2">
						<button
							onclick={() => addAvailableItem()}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Add </span><Icon
								src={ArrowLongRight}
								size="16"
								class="hover:stroke-accent inline"
							/></button
						>
						<button
							onclick={removeItem}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /><span
								class="hidden sm:inline"
							>
								Remove</span
							>
						</button>
						<button
							onclick={itemUp}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Item Up </span><Icon
								src={ArrowLongUp}
								size="16"
								class="hover:stroke-accent inline"
							/>
						</button>
						<button
							onclick={itemDown}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Item Down </span><Icon
								src={ArrowLongDown}
								size="16"
								class="hover:stroke-accent inline"
							/>
						</button>
						<button
							onclick={clear}
							class="btn btn-outline btn-sm normal-case btn-secondary block w-full"
							><span class="hidden sm:inline">Clear </span><Icon
								src={XCircle}
								size="16"
								class="hover:stroke-accent inline"
							/>
						</button>
						<select
							class="select select-outline select-sm select-secondary w-12 sm:w-full text-secondary"
							onchange={(e) => {
								e.preventDefault();
								applyPlan(
									$player.playerPlans.productionPlans.find(
										(p) => p.num == parseInt(e.currentTarget.value)
									)
								);
								e.currentTarget.value = '0';
							}}
						>
							<option value={0}>Apply Plan</option>
							{#each $player.playerPlans.productionPlans as plan (plan.num)}
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
									onclick={() => onQueueItemClicked(-1)}
									class:bg-primary={selectedQueueItemIndex === -1}
									class="w-full pl-1 select-none cursor-default hover:text-secondary-focus"
								>
									-- Top of the Queue --
								</button>
							</li>
							{#if queueItems}
								{#each queueItems as queueItem, index (index)}
									<li class="cursor-default">
										<ProductionQueueItemLine
											item={queueItem}
											{index}
											{onQueueItemClicked}
											selected={queueItem === selectedQueueItem}
										/>
									</li>
								{/each}
							{/if}
						</ul>
						<div class="divider"></div>
						<div class="h-32">
							{#if selectedQueueItem}
								<h3>
									{#if selectedQueueItem.designNum}
										<button
											type="button"
											onpointerdown={(e) =>
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
										<button
											type="button"
											onpointerdown={(e) => onAllocatedTooltip(e, selectedQueueItem?.allocated)}
											>{(selectedQueueItemPercentComplete * 100)?.toFixed()}%<Icon
												src={QuestionMarkCircle}
												size="16"
												class="cursor-help inline-block ml-1"
											/> Done,</button
										>
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
							checked={contributesOnlyLeftoverToResearch}
							onchange={contributesOnlyLeftoverToResearchChecked}
							class="checkbox checkbox-xs"
							type="checkbox"
						/> Contributes Only Leftover to Research
					</label>
				</div>
				<div class="w-1/2 flex flex-row flex-wrap justify-between sm:justify-end">
					<div class="w-1/2 md:w-auto md:grow">
						<button class="btn btn-sm btn-outline btn-secondary w-full" onclick={prev}>Prev</button>
					</div>
					<div class="w-1/2 md:w-auto md:grow">
						<button class="btn btn-sm btn-outline btn-secondary w-full" onclick={next}>Next</button>
					</div>
					<div class="w-1/2 md:w-auto md:grow">
						<button onclick={cancel} class="btn btn-sm btn-outline btn-secondary w-full"
							>Cancel</button
						>
					</div>
					<div class="w-1/2 md:w-auto md:grow">
						<button onclick={ok} class="btn btn-sm btn-primary w-full">Ok</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
