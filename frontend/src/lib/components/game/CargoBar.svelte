<script lang="ts">
	import { cargoPercent, emptyCargo, totalCargo, type Cargo } from '$lib/types/Cargo';
	import { createEventDispatcher } from 'svelte';
	import type { CargoTransferDialogEvent } from '../../../routes/(user)/games/(game)/[id]/dialogs/cargo/CargoTranfserDialog.svelte';

	const dispatch = createEventDispatcher<CargoTransferDialogEvent>();

	export let value: Cargo = {
		ironium: 0,
		boranium: 0,
		germanium: 0,
		colonists: 0
	};

	export let capacity: number | undefined = 0;
	export let canTransferCargo = false;

	let percent: Cargo = emptyCargo();

	$: percent = cargoPercent(value, capacity);
</script>

<div
	on:pointerdown={() => canTransferCargo && dispatch('cargo-transfer-dialog')}
	class="border border-secondary h-[1rem] text-[0rem] relative bg-gauge select-none"
	class:cursor-pointer={canTransferCargo}
>
	<div
		class="font-semibold text-sm text-center align-middle text-white mix-blend-difference w-full bg-blend-difference absolute"
	>
		{totalCargo(value)} of {capacity ?? 0}kT
	</div>
	<div
		style={`left: 0%; width: ${percent.ironium?.toFixed()}%`}
		class="ironium-bar h-full inline-block"
	/>
	<div style={`width: ${percent.boranium?.toFixed()}%`} class="boranium-bar h-full inline-block" />
	<div
		style={`width: ${percent.germanium?.toFixed()}%`}
		class="germanium-bar h-full inline-block"
	/>
	<div
		style={`width: ${percent.colonists?.toFixed()}%`}
		class="colonists-bar h-full inline-block"
	/>
</div>
