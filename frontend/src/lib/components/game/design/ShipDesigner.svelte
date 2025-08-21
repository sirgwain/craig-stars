<script lang="ts">
	import FormError from '$lib/components/FormError.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import Hull from '$lib/components/game/design/Hull.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import { HullSlotTypeNone } from '$lib/types/Consts';
	import { canLearnTech } from '$lib/types/Player';
	import { canFillSlot, hullAllowed } from '$lib/types/Tech';
	import { hasRequiredLevels } from '$lib/types/TechLevel';
	import type { TechHull, TechHullComponent, TechHullSlot } from '$lib/types/cs-proto';
	import {
		CostSchema,
		ShipDesignSlotSchema,
		ShipDesignSpecSchema,
		type ShipDesign,
		type ShipDesignSlot,
		type ShipDesignSpec
	} from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
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

	// Make plan-like local state for runes-compliant bindings, then sync back to design
	let { hull, design = $bindable(), error = '', numHullSets = 4, onSave }: Props = $props();

	// Local runes state for frequently bound fields
	let name: string = $state(design.name ?? '');
	let hullSetNumber: number = $state(design.hullSetNumber ?? 0);
	let slots: ShipDesignSlot[] = $state(design.slots ?? []);

	// Derived/spec state from wasm
	let designSpec: ShipDesignSpec = $state(create(ShipDesignSpecSchema));

	// Keep parent prop in sync with local state
	$effect(() => {
		design.name = name;
		design.hullSetNumber = hullSetNumber;
		design.slots = slots;
	});

	// Recompute spec when design changes
	$effect(() => {
		cs.wasmService
			.computeShipDesignSpec({ design })
			.then((resp) => (designSpec = resp.spec ?? create(ShipDesignSpecSchema)));
	});

	let highlightedSlots: number[] = $state([]);

	// only show hull components that actually fit on this hull
	let validHullSlotTypes = hull.slots.reduce((type, slot) => type | +slot.type, HullSlotTypeNone);

	let selectedComponent = $derived(
		$shipDesignerContext.selectedHullComponent ??
			($shipDesignerContext.selectedShipDesignSlot?.hullComponent
				? $techs.getHullComponent($shipDesignerContext.selectedShipDesignSlot?.hullComponent)
				: undefined)
	);

	let selectedComponentCost = $state(create(CostSchema));
	$effect(() => {
		if (!selectedComponent) {
			return;
		}
		cs.wasmService
			.getTechCost({ tech: selectedComponent.tech })
			.then((resp) => (selectedComponentCost = resp.cost ?? selectedComponentCost));
	});

	function updateHullSetNumber(num: number) {
		if (num < 0) {
			hullSetNumber = numHullSets - 1;
		} else if (num >= numHullSets) {
			hullSetNumber = 0;
		} else {
			hullSetNumber = num;
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

			// highlight compatible slots
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
	function slotClicked(
		index: number,
		slot: TechHullSlot,
		shipDesignSlot: ShipDesignSlot | undefined
	) {
		if (
			$shipDesignerContext.selectedHullComponent &&
			canFillSlot($shipDesignerContext.selectedHullComponent.hullSlotType, slot.type)
		) {
			addHullComponent($shipDesignerContext.selectedHullComponent, slot, index);
		} else {
			$shipDesignerContext.selectedHullComponent = undefined;
			if (highlightedSlots.length === 1 && highlightedSlots[0] === index) {
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

	function addHullComponent(hc: TechHullComponent, slot: TechHullSlot, index: number) {
		const existingShipDesignSlot = slots.find((s) => s.hullSlotIndex === index + 1);

		if (existingShipDesignSlot) {
			// mutate then reassign to trigger reactivity
			existingShipDesignSlot.hullComponent = hc.tech?.name ?? '';
			existingShipDesignSlot.quantity = slot.capacity ?? 0;
			slots = [...slots];
		} else {
			slots = [
				...slots,
				create(ShipDesignSlotSchema, {
					hullSlotIndex: index + 1,
					hullComponent: hc.tech?.name,
					quantity: slot.capacity
				})
			];
		}
	}

	onMount(() => {
		design.hull = hull.tech?.name ?? '';
		if (!name) {
			name = hull.tech?.name ?? '';
		}
		shipDesignerContext.update(() => ({
			selectedSlotIndex: undefined,
			selectedSlot: undefined,
			selectedShipDesignSlot: undefined,
			selectedHullComponent: undefined
		}));
	});
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
						<TechAvatar tech={hull} {hullSetNumber} />
					</div>
					<div class="flex flex-row justify-between">
						<div>
							<button
								type="button"
								onclick={() => updateHullSetNumber(hullSetNumber - 1)}
								class="btn btn-outline btn-xs normal-case btn-secondary"
								data-type="prev-hull-set-button"
							>
								<Icon src={ChevronLeft} size="16" class="hover:stroke-accent" />
							</button>
						</div>
						<div>
							<button
								type="button"
								onclick={() => updateHullSetNumber(hullSetNumber + 1)}
								class="btn btn-outline btn-xs normal-case btn-secondary"
								data-type="next-hull-set-button"
							>
								<Icon src={ChevronRight} size="16" class="hover:stroke-accent" />
							</button>
						</div>
					</div>
				</div>

				<div class="grow">
					<TextInput
						name="name"
						bind:value={name}
						required
						titleClass="label-text w-16 text-right"
					/>
				</div>
			</div>
			<div class="flex flex-row justify-center">
				<Hull
					bind:shipDesignSlots={slots}
					{hull}
					cargoCapacity={designSpec.cargoCapacity}
					{highlightedSlots}
					highlightedClass="border-accent"
					showTooltips={false}
					onSlotClicked={slotClicked}
				/>
			</div>
			<div class="flex flex-row justify-between pl-2">
				<div class="flex flex-col">
					<div>Cost of one {name}</div>
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
				{#each $techs.hullComponents.filter((hc) => hullAllowed(hull, hc)) as hc (hc.tech?.name)}
					{#if canLearnTech($player, hc) && hasRequiredLevels($player.techLevels, hc.tech?.requirements?.techLevel) && (!$shipDesignerContext.selectedSlot || canFillSlot(hc.hullSlotType, $shipDesignerContext.selectedSlot.type)) && canFillSlot(hc.hullSlotType, validHullSlotTypes)}
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
											{hc.tech?.name}
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
						Cost of one {selectedComponent.tech?.name}
						<span class="inline-block" onpointerdown={(e) => onTechTooltip(e, selectedComponent)}
							><Icon
								src={QuestionMarkCircle}
								size="16"
								class=" cursor-help hover:stroke-accent"
							/></span
						>
					</div>
					<div class="pl-2">
						<Cost cost={selectedComponentCost} />
					</div>
				{/if}
			</div>
		</div>
	</div>
</form>
