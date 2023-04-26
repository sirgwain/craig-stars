<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import Hull from '$lib/components/game/design/Hull.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { player, techs } from '$lib/services/Context';
	import { canLearnTech, hasRequiredLevels } from '$lib/types/Player';
	import type { ShipDesign, ShipDesignSlot, Spec } from '$lib/types/ShipDesign';
	import { canFillSlot, type HullSlot, type TechHullComponent } from '$lib/types/Tech';
	import { createEventDispatcher, onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { shipDesignerContext } from './ShipDesignerContext';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ExclamationTriangle } from '@steeze-ui/heroicons';
	import { DesignService } from '$lib/services/DesignService';
	import Cost from '../Cost.svelte';
	import DesignStats from '../DesignStats.svelte';
	import HullComponent from '$lib/components/tech/hull/HullComponent.svelte';

	const dispatch = createEventDispatcher();

	export let gameId: number | string;
	export let hullName: string;
	export let design: ShipDesign;
	export let error: string = '';

	let designSpec: Spec = design?.spec || {};
	let highlightedSlots: HullSlot[] = [];
	let highlightedClass: string;

	$: hull = $techs.getHull(hullName);
	$: design.hull = hull?.name ?? '';

	$: design && DesignService.computeSpec(gameId, design).then((s) => (designSpec = s));

	onMount(() => {
		shipDesignerContext.update(() => ({
			selectedSlotIndex: undefined,
			selectedSlot: undefined,
			selectedShipDesignSlot: undefined,
			selectedHullComponent: undefined
		}));
	});

	// when a tech is selected from the tech tree
	const onTechHullComponentClicked = (hc: TechHullComponent) => {
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

			highlightedSlots =
				hull?.slots.filter(
					(slot) =>
						$shipDesignerContext.selectedHullComponent &&
						canFillSlot($shipDesignerContext.selectedHullComponent.hullSlotType, slot.type)
				) ?? [];
		}
	};

	// when a slot is clicked on the hull
	const onSlotClicked = (index: number, slot: HullSlot, shipDesignSlot: ShipDesignSlot) => {
		if (
			$shipDesignerContext.selectedHullComponent &&
			canFillSlot($shipDesignerContext.selectedHullComponent.hullSlotType, slot.type)
		) {
			addHullComponent($shipDesignerContext.selectedHullComponent, slot, index);
		} else {
			$shipDesignerContext.selectedHullComponent = undefined;
			if (highlightedSlots.length == 1 && highlightedSlots[0] == slot) {
				highlightedSlots = [];
				$shipDesignerContext.selectedSlotIndex = undefined;
				$shipDesignerContext.selectedSlot = undefined;
				$shipDesignerContext.selectedShipDesignSlot = undefined;
			} else {
				highlightedSlots = [slot];
				highlightedClass = 'border-accent';
				$shipDesignerContext.selectedSlotIndex = index;
				$shipDesignerContext.selectedSlot = slot;
				$shipDesignerContext.selectedShipDesignSlot = shipDesignSlot;
			}
		}
	};

	const addHullComponent = (hc: TechHullComponent, slot: HullSlot, index: number) => {
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
	};

	const onSubmit = async () => {
		dispatch('save');
	};
</script>

<div>
	<form on:submit|preventDefault={onSubmit}>
		<div class="flex justify-end gap-2">
			<button class="btn btn-success" type="submit">Save</button>
		</div>

		<ItemTitle>Design: {design.name}</ItemTitle>

		{#if error !== ''}
			<div
				class="alert alert-error shadow-lg w-1/2 mx-auto"
				in:fade
				out:fade={{ delay: 5000 }}
				on:introend={(e) => (error = '')}
			>
				<div>
					<Icon src={ExclamationTriangle} size="24" class="hover:stroke-accent" />
					<span>{error}</span>
				</div>
			</div>
		{/if}

		<div class="flex flex-col md:flex-row-reverse justify-center">
			<div class="grow flex flex-col">
				<div>
					<TextInput name="name" bind:value={design.name} required />
				</div>
				{#if hull}
					<div class="flex flex-row justify-center">
						<Hull
							bind:shipDesignSlots={design.slots}
							{hull}
							{highlightedSlots}
							highlightedClass={'border-accent'}
							on:slot-clicked={(e) =>
								onSlotClicked(e.detail.index, e.detail.slot, e.detail.shipDesignSlot)}
						/>
					</div>
					<div class="flex flex-row justify-between pl-2">
						<Cost cost={designSpec?.cost} />
						<DesignStats spec={designSpec} />
					</div>
				{/if}
			</div>
			<div class="">
				<div class="font-bold text-2xl">Hull Components</div>
				<ul class="h-[600px] w-full sm:w-[16rem] px-1 border p-1 overflow-y-auto">
					{#each $techs.hullComponents as hc}
						{#if $player && canLearnTech($player, hc) && hasRequiredLevels($player.techLevels, hc.requirements) && (!$shipDesignerContext.selectedSlot || canFillSlot(hc.hullSlotType, $shipDesignerContext.selectedSlot.type))}
							<li>
								<div
									class={`flex ${
										$shipDesignerContext.selectedHullComponent === hc ? 'border border-accent' : ''
									}`}
								>
									<button
										type="button"
										class="w-full h-full"
										on:click={(e) => onTechHullComponentClicked(hc)}
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
			</div>
		</div>
	</form>
</div>
