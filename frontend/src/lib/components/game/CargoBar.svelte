<script lang="ts">
	import { run } from 'svelte/legacy';

	import { cargoPercent, emptyCargo, totalCargo, type Cargo } from '$lib/types/Cargo';
	import { createEventDispatcher } from 'svelte';
	import type { CargoTransferDialogEvent } from '../../../routes/(user)/games/(game)/[id]/dialogs/cargo/CargoTranfserDialog.svelte';

	const dispatch = createEventDispatcher<CargoTransferDialogEvent>();

	interface Props {
		value?: Cargo;
		capacity?: number | undefined;
		canTransferCargo?: boolean;
	}

	let {
		value = {
			ironium: 0,
			boranium: 0,
			germanium: 0,
			colonists: 0
		},
		capacity = 0,
		canTransferCargo = false
	}: Props = $props();

	let percent: Cargo = $state(emptyCargo());

	run(() => {
		percent = cargoPercent(value, capacity);
	});
</script>

<div
	onpointerdown={() => canTransferCargo && dispatch('cargo-transfer-dialog')}
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
	></div>
	<div
		style={`width: ${percent.boranium?.toFixed()}%`}
		class="boranium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.germanium?.toFixed()}%`}
		class="germanium-bar h-full inline-block"
	></div>
	<div
		style={`width: ${percent.colonists?.toFixed()}%`}
		class="colonists-bar h-full inline-block"
	></div>
</div>
