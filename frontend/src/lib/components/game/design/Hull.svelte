<script lang="ts">
	import { techs } from '$lib/services/Stores';
	import type { ShipDesignSlot } from '$lib/types/ShipDesign';
	import type { HullSlot, TechHull } from '$lib/types/Tech';
	import CargoComponent from '../../tech/hull/CargoComponent.svelte';
	import HullComponent from '../../tech/hull/HullComponent.svelte';
	import SpaceDockComponent from '../../tech/hull/SpaceDockComponent.svelte';
	import { onTechTooltip } from '../tooltips/TechTooltip.svelte';

	const componentSize = 64; // each component block is 64px
	const containerWidth = componentSize * 5;
	const containerHeight = componentSize * 5;

	type Props = {
		hull: TechHull;
		shipDesignSlots?: ShipDesignSlot[];
		highlightedSlots?: number[];
		highlightedClass?: string;
		cargoCapacity?: any;
		showTooltips?: boolean;
		onslotclicked?: (
			index: number,
			hullSlot: HullSlot,
			shipDesignSlot: ShipDesignSlot | undefined
		) => void;
	};

	let {
		hull,
		shipDesignSlots = $bindable([]),
		highlightedSlots = [],
		highlightedClass = '',
		cargoCapacity = hull.cargoCapacity ?? 0,
		showTooltips = true,
		onslotclicked
	}: Props = $props();
</script>

<div
	class="relative m-2 bg-base-200 dark:bg-base-300"
	style={`width: ${containerWidth}px; height: ${containerHeight}px`}
>
	{#each hull.slots as slot, index (index)}
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
			oncontextmenu={(e) =>
				shipDesignSlot && onTechTooltip(e, $techs.getHullComponent(shipDesignSlot?.hullComponent))}
		>
			<HullComponent
				{shipDesignSlot}
				type={slot.type}
				capacity={slot.capacity}
				required={slot.required}
				highlighted={highlightedSlots.findIndex((s) => s === index) !== -1}
				{highlightedClass}
				{showTooltips}
				onclick={() => {
					onslotclicked?.(index, slot, shipDesignSlot);
				}}
				ondelete={() => {
					shipDesignSlots = shipDesignSlots.filter((s) => s != shipDesignSlot);
				}}
				onupdate={() => {
					shipDesignSlots = shipDesignSlots;
				}}
			/>
		</div>
	{/each}
</div>
