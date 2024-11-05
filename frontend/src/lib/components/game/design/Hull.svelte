<script lang="ts">
	import { preventDefault } from 'svelte/legacy';

	import type { ShipDesignSlot } from '$lib/types/ShipDesign';
	import type { HullSlot, TechHull } from '$lib/types/Tech';
	import { createEventDispatcher } from 'svelte';
	import CargoComponent from '../../tech/hull/CargoComponent.svelte';
	import HullComponent from '../../tech/hull/HullComponent.svelte';
	import SpaceDockComponent from '../../tech/hull/SpaceDockComponent.svelte';
	import { onTechTooltip } from '../tooltips/TechTooltip.svelte';
	import { techs } from '$lib/services/Stores';

	const dispatch = createEventDispatcher();

	const componentSize = 64; // each component block is 64px
	const containerWidth = componentSize * 5;
	const containerHeight = componentSize * 5;

	interface Props {
		hull: TechHull;
		shipDesignSlots?: ShipDesignSlot[];
		highlightedSlots?: HullSlot[];
		highlightedClass?: string;
		cargoCapacity?: any;
		showTooltips?: boolean;
	}

	let {
		hull,
		shipDesignSlots = $bindable([]),
		highlightedSlots = [],
		highlightedClass = '',
		cargoCapacity = hull.cargoCapacity ?? 0,
		showTooltips = true
	}: Props = $props();
</script>

<div
	class="relative m-2 bg-base-200 dark:bg-base-300"
	style={`width: ${containerWidth}px; height: ${containerHeight}px`}
>
	{#each hull.slots as slot, index}
		{@const shipDesignSlot = shipDesignSlots.find((s) => s.hullSlotIndex === index + 1)}
		{#if index === 1 && cargoCapacity > 0}
			<div
				class="absolute"
				style={`left: ${
					(hull.cargoSlotPosition?.x ?? 0) * componentSize +
					(containerWidth / 2 - componentSize / 2)
				}px; top: ${
					(hull.cargoSlotPosition?.y ?? 0) * componentSize +
					(containerHeight / 2 - componentSize / 2)
				}px; width: ${(hull.cargoSlotSize?.x ?? 0) * componentSize}px; height: ${
					(hull.cargoSlotSize?.y ?? 0) * componentSize
				}px;`}
			>
				<CargoComponent capacity={cargoCapacity} />
			</div>
		{/if}
		{#if index === 1 && hull.spaceDock && hull.spaceDock !== 0}
			<div
				class="absolute"
				style={`left: ${
					(hull.spaceDockSlotPosition?.x ?? 0) * componentSize +
					(containerWidth / 2 - componentSize / 2)
				}px; top: ${
					(hull.spaceDockSlotPosition?.y ?? 0) * componentSize +
					(containerHeight / 2 - componentSize / 2)
				}px; width: ${(hull.spaceDockSlotSize?.x ?? 0) * componentSize + 1}px; height: ${
					(hull.spaceDockSlotSize?.y ?? 0) * componentSize + 1
				}px;`}
			>
				<SpaceDockComponent spaceDock={hull.spaceDock} rounded={hull.spaceDockSlotCircle} />
			</div>
		{/if}
		<div
			class="absolute"
			style={`left: ${
				slot.position.x * componentSize + (containerWidth / 2 - componentSize / 2)
			}px; top: ${slot.position.y * componentSize + (containerHeight / 2 - componentSize / 2)}px;`}
			role="link"
			tabindex="-1"
			oncontextmenu={preventDefault(
				(e) =>
					shipDesignSlot && onTechTooltip(e, $techs.getHullComponent(shipDesignSlot?.hullComponent))
			)}
		>
			<HullComponent
				{shipDesignSlot}
				type={slot.type}
				capacity={slot.capacity}
				required={slot.required}
				highlighted={highlightedSlots.findIndex((s) => s === slot) !== -1}
				{highlightedClass}
				{showTooltips}
				on:clicked={(e) => {
					dispatch('slot-clicked', { index, slot, shipDesignSlot });
				}}
				on:deleted={() => {
					shipDesignSlots = shipDesignSlots.filter((s) => s !== shipDesignSlot);
				}}
				on:updated={() => {
					shipDesignSlots = shipDesignSlots;
				}}
			/>
		</div>
	{/each}
</div>
