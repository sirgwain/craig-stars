<script lang="ts">
	import CargoTransferer from '$lib/components/game/cargotransfer/CargoTransferer.svelte';
	import type { OnCancel, OnOk, TransferCargoEvent } from '$lib/services/Events';
	import { newCargoTransferRequest } from '$lib/types/CargoTransferRequest';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import type { Planet } from '$lib/types/Planet';
	import type { Salvage } from '$lib/types/Salvage';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';

	type Props = {
		src: CommandedFleet;
		dest: Fleet | Planet | Salvage | undefined;
		onOk: OnOk<TransferCargoEvent>;
		onCancel: OnCancel;
	};

	let { src, dest, onOk, onCancel }: Props = $props();

	let transferAmount = $state(newCargoTransferRequest());

	function reset() {
		transferAmount = newCargoTransferRequest();
		src = src;
	}

	function ok() {
		onOk({ src, dest, transferAmount });
		reset();
	}

	function cancel() {
		reset();
		onCancel();
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
					<button onclick={ok} class="btn btn-primary">Ok</button>
					<button onclick={cancel} class="btn btn-secondary">Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}
