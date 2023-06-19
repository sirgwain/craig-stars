<script lang="ts">
	import { techs } from '$lib/services/Context';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import type { HullSlot, TechHull, TechHullComponent } from '$lib/types/Tech';
	import CargoComponent from './CargoComponent.svelte';
	import HullComponent from './HullComponent.svelte';
	import SpaceDockComponent from './SpaceDockComponent.svelte';

	export let hull: TechHull;
	export let design: ShipDesign | undefined = undefined;
	const componentSize = 64; // each component block is 64px
	const containerWidth = componentSize * 5;
	const containerHeight = componentSize * 5;

	// get a hull component for a slot
	function getHullComponent(slot: HullSlot, slotIndex: number): TechHullComponent | undefined {
		if (design) {
			const slot = design.slots.find((s) => s.hullSlotIndex === slotIndex + 1);
			if (slot) {
				return $techs.getHullComponent(slot.hullComponent);
			}
		}
	}

	function getSlotQuantity(slot: HullSlot, slotIndex: number): number | undefined {
		if (design) {
			const slot = design.slots.find((s) => s.hullSlotIndex === slotIndex + 1);
			if (slot) {
				return slot.quantity;
			}
		}
	}

</script>

<div class="relative m-2" style={`width: ${containerWidth}px; height: ${containerHeight}px`}>
	{#each hull.slots as slot, index}
		{#if index == 1 && hull.cargoCapacity && hull.cargoCapacity > 0}
			<div
				class="absolute"
				style={`left: ${
					(hull.cargoSlotPosition?.x ?? 0) * componentSize +
					(containerWidth / 2 - componentSize / 2)
				}px; top: ${
					(hull.cargoSlotPosition?.y ?? 0) * componentSize +
					(containerHeight / 2 - componentSize / 2)
				}px; width: ${(hull.cargoSlotSize?.x ?? 0) * componentSize + 1}px; height: ${
					(hull.cargoSlotSize?.y ?? 0) * componentSize + 1
				}px;`}
			>
				<CargoComponent capacity={hull.cargoCapacity} />
			</div>
		{/if}
		{#if index == 1 && hull.spaceDock && hull.spaceDock != 0}
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
		>
			<HullComponent
				component={getHullComponent(slot, index)}
				quantity={getSlotQuantity(slot, index)}
				type={slot.type}
				capacity={slot.capacity}
				required={slot.required}
			/>
		</div>
	{/each}
</div>
