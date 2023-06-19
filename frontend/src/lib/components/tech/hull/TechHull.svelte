<script lang="ts">
	import type { TechHull } from '$lib/types/Tech';
	import CargoComponent from './CargoComponent.svelte';
	import HullComponent from './HullComponent.svelte';

	export let hull: TechHull;
</script>

<div
	class={`flex flex-row place-items-center w-[${64 * 5}px] h-[${64 * 5}px]`}
	style={`width: ${64 * 5}px; height: ${64 * 5}px`}
>
	{#each hull.slots as slot, index}
		{#if index == 1 && hull.cargoCapacity && hull.cargoCapacity > 0}
			<CargoComponent capacity={hull.cargoCapacity} />
		{/if}
		<HullComponent type={slot.type} capacity={slot.capacity} />
	{/each}
</div>
