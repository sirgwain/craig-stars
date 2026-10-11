<script lang="ts">
	import { onPointerKeyDown, onPointerKeyUp } from '#lib/services/Events.js';
	import { cargoPercent, emptyCargo, totalCargo } from '#lib/types/Cargo.js';
	import type { Cargo } from '#lib/types/cs-proto.js';

	type Props = {
		value: Cargo | undefined;
		capacity: number | undefined;
		canTransferCargo?: boolean;
		onPointerDown?: (e: PointerEvent) => void | undefined;
	};

	let { value, capacity, canTransferCargo = false, onPointerDown }: Props = $props();

	let percent = $derived(cargoPercent(value ?? emptyCargo(), capacity));
</script>

<div
	role="button"
	aria-disabled={!(canTransferCargo && onPointerDown)}
	tabindex={canTransferCargo && onPointerDown ? 0 : undefined}
	aria-label={canTransferCargo ? 'Transfer cargo' : undefined}
	onkeydown={onPointerKeyDown}
	onkeyup={onPointerKeyUp}
	onpointerdown={(e) => (canTransferCargo && onPointerDown ? onPointerDown(e) : undefined)}
	class="border border-secondary h-4 text-[0rem] relative bg-gauge select-none"
	class:cursor-pointer={canTransferCargo}
>
	<div
		class="font-semibold text-sm text-center align-middle text-white mix-blend-difference w-full bg-blend-difference absolute"
	>
		{totalCargo(value)} of {capacity ?? 0}kT
	</div>
	<!-- Each share keeps its exact width. Rounded to whole percents, three thirds left a 1% gap in a
	     full hold, and a small share could vanish. -->
	<div class="flex h-full" data-type="cargo-fill">
		<div style={`width: ${percent.ironium}%`} class="ironium-bar h-full"></div>
		<div style={`width: ${percent.boranium}%`} class="boranium-bar h-full"></div>
		<div style={`width: ${percent.germanium}%`} class="germanium-bar h-full"></div>
		<div style={`width: ${percent.colonists}%`} class="colonists-bar h-full"></div>
	</div>
</div>
