<script lang="ts">
	import FormError from '$lib/components/FormError.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import Hull from '$lib/components/game/design/Hull.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import { canLearnTech } from '$lib/types/Player';
	import type { ShipDesign, ShipDesignSlot, Spec } from '$lib/types/ShipDesign';
	import {
		HullSlotType,
		canFillSlot,
		hullAllowed,
		type HullSlot,
		type TechHull,
		type TechHullComponent
	} from '$lib/types/Tech';
	import { hasRequiredLevels } from '$lib/types/TechLevel';
	import { ChevronLeft, ChevronRight, QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';
	import Cost from '../Cost.svelte';
	import CostMini from '../CostMini.svelte';
	import DesignStats from '../DesignStats.svelte';
	import { onTechTooltip } from '../tooltips/TechTooltip.svelte';
	import { shipDesignerContext } from './ShipDesignerContext';

	const { cs, player } = getGameContext();

	type Props = {
		hull: TechHull;
		design: ShipDesign;
		error?: string;
		numHullSets?: number;
		onSave?: () => void;
	};

	let { hull, design = $bindable(), error = '', numHullSets = 4, onSave }: Props = $props();

	let designSpec: Spec = $derived(cs.computeShipDesignSpec(design) ?? ({} as Spec));
	let highlightedSlots: number[] = $state([]);

	// only show hull components that actually fit on this hull
	let validHullSlotTypes = hull.slots.reduce((type, slot) => type | +slot.type, HullSlotType.None);

	let selectedComponent = $derived(
		$shipDesignerContext.selectedHullComponent ??
			($shipDesignerContext.selectedShipDesignSlot?.hullComponent
				? $techs.getHullComponent($shipDesignerContext.selectedShipDesignSlot?.hullComponent)
				: undefined)
	);

	onMount(() => {
		design.hull = hull.name;
		if (!design.name) {
			design.name = hull.name;
		}
		shipDesignerContext.update(() => ({
			selectedSlotIndex: undefined,
			selectedSlot: undefined,
			selectedShipDesignSlot: undefined,
			selectedHullComponent: undefined
		}));
	});

	function updateHullSetNumber(num: number) {
		if (num < 0) {
			design.hullSetNumber = numHullSets - 1;
		} else if (num >= numHullSets) {
			design.hullSetNumber = 0;
		} else {
			design.hullSetNumber = num;
		}
	}

	// when a tech is selected from the tech tree
	function techHullComponentClicked(hc: TechHullComponent) {
		if ($shipDesignerContext.selectedSlot && $shipDesignerContext.selectedSlotIndex !== undefined) {
			// clear out the selected hull component
			$shipDesignerContext.selectedHullComponent = undefined;
			addHullComponent(
				hc,
				$shipDesignerContext.selectedSlot,
				$shipDesignerContext.selectedSlotIndex
			);
		} else {
			if ($shipDesignerContext.selectedHullComponent === hc) {
				$shipDesignerContext.selectedHullComponent = undefined;
			} else {
				$shipDesignerContext.selectedHullComponent = hc;
			}

			// we need a filtered list of slot indices that this selected component
			// can go into so we can highlight it
			// first build a map of slot/index, then filter it, then return an array
			// of just the indices
			highlightedSlots =
				hull?.slots
					.map((slot, index) => ({ slot, index }))
					.filter(
						(s) =>
							$shipDesignerContext.selectedHullComponent &&
							canFillSlot($shipDesignerContext.selectedHullComponent.hullSlotType, s.slot.type)
					)
					.map((s) => s.index) ?? [];
		}
	}

	// when a slot is clicked on the hull
	function slotClicked(index: number, slot: HullSlot, shipDesignSlot: ShipDesignSlot | undefined) {
		if (
			$shipDesignerContext.selectedHullComponent &&
			canFillSlot($shipDesignerContext.selectedHullComponent.hullSlotType, slot.type)
		) {
			addHullComponent($shipDesignerContext.selectedHullComponent, slot, index);
		} else {
			$shipDesignerContext.selectedHullComponent = undefined;
			if (highlightedSlots.length == 1 && highlightedSlots[0] == index) {
				highlightedSlots = [];
				$shipDesignerContext.selectedSlotIndex = undefined;
				$shipDesignerContext.selectedSlot = undefined;
				$shipDesignerContext.selectedShipDesignSlot = undefined;
			} else {
				highlightedSlots = [index];
				$shipDesignerContext.selectedSlotIndex = index;
				$shipDesignerContext.selectedSlot = slot;
				$shipDesignerContext.selectedShipDesignSlot = shipDesignSlot;
			}
		}
	}

	function addHullComponent(hc: TechHullComponent, slot: HullSlot, index: number) {
		const existingShipDesignSlot = design.slots.find((s) => s.hullSlotIndex === index + 1);

		if (existingShipDesignSlot) {
			existingShipDesignSlot.hullComponent = hc.name;
			existingShipDesignSlot.quantity = slot.capacity;
			design.slots = design.slots;
		} else {
			design.slots = [
				...design.slots,
				{
					hullSlotIndex: index + 1,
					hullComponent: hc.name,
					quantity: slot.capacity
				}
			];
		}
	}
</script>

<form
	onsubmit={(e) => {
		e.preventDefault();
		onSave?.();
	}}
>
	<FormError {error} />

	<div class="flex flex-col md:flex-row-reverse justify-center">
		<div class="flex flex-col w-full mx-1">
			<div class="flex flex-row justify-between">
				<div class="flex flex-col">
					<div class="mx-auto border border-secondary bg-black p-1">
						<TechAvatar tech={hull} hullSetNumber={design.hullSetNumber} />
					</div>
					<div class="flex flex-row justify-between">
						<div>
							<button
								type="button"
								onclick={() => updateHullSetNumber(design.hullSetNumber - 1)}
								class="btn btn-outline btn-xs normal-case btn-secondary"
							>
								<Icon src={ChevronLeft} size="16" class="hover:stroke-accent" />
							</button>
						</div>
						<div>
							<button
								type="button"
								onclick={() => updateHullSetNumber(design.hullSetNumber + 1)}
								class="btn btn-outline btn-xs normal-case btn-secondary"
							>
								<Icon src={ChevronRight} size="16" class="hover:stroke-accent" />
							</button>
						</div>
					</div>
				</div>

				<div class="grow">
					<TextInput
						name="name"
						bind:value={design.name}
						required
						titleClass="label-text w-16 text-right"
					/>
				</div>
			</div>
			<div class="flex flex-row justify-center">
				<Hull
					bind:shipDesignSlots={design.slots}
					{hull}
					cargoCapacity={designSpec.cargoCapacity}
					{highlightedSlots}
					highlightedClass={'border-accent'}
					showTooltips={false}
					onslotclicked={slotClicked}
				/>
			</div>
			<div class="flex flex-row justify-between pl-2">
				<div class="flex flex-col">
					<div>Cost of one {design.name}</div>
					<div class="pl-2 hidden sm:block">
						<Cost cost={designSpec?.cost} />
					</div>
					<div class="pl-2 sm:hidden flex justify-between">
						<CostMini cost={designSpec?.cost} />
						<!-- <div class="ml-2"><button type="button" class="btn btn-sm btn-outline btn-secondary">Stats</button></div> -->
					</div>
				</div>
				<div class="hidden sm:block">
					<DesignStats spec={designSpec} />
				</div>
			</div>
		</div>
		<div>
			<div class="font-bold text-2xl">Hull Components</div>
			<ul class="w-full h-[400px] border-b sm:w-[16rem] px-1 p-1 overflow-y-auto">
				{#each $techs.hullComponents.filter((hc) => hullAllowed(hull.name, hc)) as hc}
					{#if canLearnTech($player, hc) && hasRequiredLevels($player.techLevels, hc.requirements) && (!$shipDesignerContext.selectedSlot || canFillSlot(hc.hullSlotType, $shipDesignerContext.selectedSlot.type)) && canFillSlot(hc.hullSlotType, validHullSlotTypes)}
						<li>
							<div
								class={`flex ${
									$shipDesignerContext.selectedHullComponent === hc ? 'border border-accent' : ''
								}`}
							>
								<button
									type="button"
									class="w-full h-full"
									onclick={() => techHullComponentClicked(hc)}
								>
									<div class="flex flex-row place-items-center">
										<div class="mr-2 pt-1 pl-1">
											<TechAvatar tech={hc} />
										</div>
										<div>
											{hc.name}
										</div>
									</div>
								</button>
							</div>
						</li>
					{/if}
				{/each}
			</ul>
			<div class="flex flex-col mt-3">
				{#if selectedComponent}
					<div>
						Cost of one {selectedComponent.name}
						<span class="inline-block" onpointerdown={(e) => onTechTooltip(e, selectedComponent)}
							><Icon
								src={QuestionMarkCircle}
								size="16"
								class=" cursor-help hover:stroke-accent"
							/></span
						>
					</div>
					<div class="pl-2">
						<Cost cost={$player.getTechCost(selectedComponent)} />
					</div>
				{/if}
			</div>
		</div>
	</div>
</form>
