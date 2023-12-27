<script lang="ts" context="module">
	import { CargoTransferRequest } from '$lib/types/Cargo';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import type { Planet } from '$lib/types/Planet';

	export type TransferCargoEventDetails = {
		src: CommandedFleet;
		dest?: Fleet | Planet | Salvage;
		transferAmount: CargoTransferRequest;
	};
	export type CargoTransferEvent = {
		'transfer-cargo': TransferCargoEventDetails;
		cancel: void;
	};
</script>

<script lang="ts">
	import CargoTransferer from '$lib/components/game/cargotransfer/CargoTransferer.svelte';
	import type { Salvage } from '$lib/types/Salvage';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher<CargoTransferEvent>();

	export let src: CommandedFleet;
	export let dest: Fleet | Planet | Salvage | undefined;

	let transferAmount = new CargoTransferRequest();

	function reset() {
		transferAmount = new CargoTransferRequest();
		src = src;
	}

	function ok() {
		dispatch('transfer-cargo', { src, dest, transferAmount });
		reset();
	}

	function cancel() {
		reset();
		dispatch('cancel');
	}

	onMount(() => {
		const originalScope = hotkeys.getScope();
		const scope = 'cargoTransfer';
		hotkeys('Esc', scope, cancel);
		hotkeys('Enter', scope, ok);
		hotkeys.setScope(scope);

		return () => {
			hotkeys.unbind('Esc', scope, cancel);
			hotkeys.unbind('Enter', scope, ok);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

{#if src?.spec}
	<div
		class="flex h-full bg-base-200 shadow max-h-fit min-h-fit rounded-sm border-2 border-base-300"
	>
		<div class="flex-col h-full w-full">
			<div class="flex flex-col h-full w-full">
				<CargoTransferer {src} {dest} bind:transferAmount />
				<div class="flex justify-end pt-2">
					<button on:click={ok} class="btn btn-primary">Ok</button>
					<button on:click={cancel} class="btn btn-secondary">Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}
