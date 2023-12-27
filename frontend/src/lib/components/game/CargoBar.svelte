<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { emptyCargo, totalCargo, type Cargo } from '$lib/types/Cargo';
	import type { CargoTransferDialogEvent } from '../../../routes/(user)/games/(game)/[id]/dialogs/cargo/CargoTranfserDialog.svelte';

	const dispatch = createEventDispatcher<CargoTransferDialogEvent>();

	export let value: Cargo = {
		ironium: 0,
		boranium: 0,
		germanium: 0,
		colonists: 0
	};

	export let capacity = 0;
	export let canTransferCargo = false;

	let percent: Cargo = emptyCargo();

	$: if (capacity > 0) {
		percent = {
			ironium: ((value.ironium ?? 0) / capacity) * 100,
			boranium: ((value.boranium ?? 0) / capacity) * 100,
			germanium: ((value.germanium ?? 0) / capacity) * 100,
			colonists: ((value.colonists ?? 0) / capacity) * 100
		};
	} else {
		capacity = 0;
		percent = emptyCargo();
	}
</script>

<div
	on:pointerdown={() => canTransferCargo && dispatch('cargo-transfer-dialog')}
	class="border border-secondary h-[1rem] text-[0rem] relative bg-base-200 select-none"
	class:cursor-pointer={canTransferCargo}
>
	<div
		class="font-semibold text-sm text-center align-middle text-secondary w-full bg-blend-difference absolute"
	>
		{totalCargo(value)} of {capacity}kT
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
