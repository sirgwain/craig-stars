<script lang="ts">
	import { cargoPercent, emptyCargo, totalCargo } from '$lib/types/Cargo';
	import type { Cargo } from '$lib/types/cs-proto';

	type Props = {
		value: Cargo | undefined;
		capacity: number | undefined;
		canTransferCargo?: boolean;
		onPointerDown?: (e: PointerEvent) => void | undefined;
	};

	let { value, capacity, canTransferCargo = false, onPointerDown: onPointerDown }: Props = $props();

	let percent: Cargo = $derived(cargoPercent(value ?? emptyCargo(), capacity));
</script>

<div
	onpointerdown={(e) => (canTransferCargo && onPointerDown ? onPointerDown(e) : undefined)}
	class="border border-secondary h-[1rem] text-[0rem] relative bg-gauge select-none"
	class:cursor-pointer={canTransferCargo}
>
	<div
		class="font-semibold text-sm text-center align-middle text-white mix-blend-difference w-full bg-blend-difference absolute"
	>
		{totalCargo(value)} of {capacity ?? 0}kT
	</div>
	<div
		style={`left: 0%; width: ${percent.ironium.toFixed()}%`}
		class="ironium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.boranium.toFixed()}%`}
		class="boranium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.germanium.toFixed()}%`}
		class="germanium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.colonists.toFixed()}%`}
		class="colonists-bar h-full inline-block"
	></div>
</div>
