<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import Hull from '$lib/components/game/design/Hull.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { player, techs } from '$lib/services/Context';
	import { canLearnTech, hasRequiredLevels } from '$lib/types/Player';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { canFillSlot, type TechHullComponent } from '$lib/types/Tech';
	import { createEventDispatcher, onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { shipDesignerContext } from './ShipDesignerContext';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ExclamationTriangle } from '@steeze-ui/heroicons';

	const dispatch = createEventDispatcher();

	export let hullName: string;
	export let design: ShipDesign;
	export let error: string = '';

	onMount(() => {
		shipDesignerContext.update(() => ({
			selectedSlotIndex: undefined,
			selectedSlot: undefined,
			selectedShipDesignSlot: undefined
		}));
	});

	$: hull = $techs.getHull(hullName);
	$: design.hull = hull?.name ?? '';

	// callback when a tech is selected from the tech tree
	const onTechHullComponentClicked = (hc: TechHullComponent) => {
		if ($shipDesignerContext.selectedSlot && $shipDesignerContext.selectedSlotIndex !== undefined) {
			const existingShipDesignSlot = design.slots.find(
				(s) =>
					$shipDesignerContext.selectedSlotIndex &&
					s.hullSlotIndex === $shipDesignerContext.selectedSlotIndex + 1
			);

			if (existingShipDesignSlot) {
				existingShipDesignSlot.hullComponent = hc.name;
				existingShipDesignSlot.quantity = $shipDesignerContext.selectedSlot.capacity;
				design.slots = design.slots;
			} else {
				design.slots = [
					...design.slots,
					{
						hullSlotIndex: $shipDesignerContext.selectedSlotIndex + 1,
						hullComponent: hc.name,
						quantity: $shipDesignerContext.selectedSlot.capacity
					}
				];
			}
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
						<Hull bind:shipDesignSlots={design.slots} {hull} />
					</div>
					<!-- <Cost cost={design.spec?.cost} />
			<DesignStats {design} /> -->
				{/if}
			</div>
			<div class="">
				<div class="font-bold text-2xl">Hull Components</div>
				<ul class="h-[400px] w-full sm:w-[16rem] px-1 border p-1 overflow-y-auto">
					{#each $techs.hullComponents as hc}
						{#if $player && canLearnTech($player, hc) && hasRequiredLevels($player.techLevels, hc.requirements) && (!$shipDesignerContext.selectedSlot || canFillSlot(hc.hullSlotType, $shipDesignerContext.selectedSlot.type))}
							<li>
								<button
									type="button"
									class="flex flex-row place-items-center"
									on:click={(e) => onTechHullComponentClicked(hc)}
								>
									<div class="mr-2 mb-2">
										<TechAvatar tech={hc} />
									</div>
									<div>
										{hc.name}
									</div>
								</button>
							</li>
						{/if}
					{/each}
				</ul>
			</div>
		</div>
	</form>
</div>
