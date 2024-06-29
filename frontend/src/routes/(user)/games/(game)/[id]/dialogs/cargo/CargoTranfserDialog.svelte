<script lang="ts" context="module">
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import type { Planet } from '$lib/types/Planet';

	export type CargoTransferDialogEventDetails = {
		src: CommandedFleet;
		dest?: Fleet | Planet | Salvage | undefined;
	};
	export type CargoTransferDialogEvent = {
		'cargo-transfer-dialog'?: CargoTransferDialogEventDetails;
		cancel: void;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { newSalvage, type Salvage } from '$lib/types/Salvage';
	import CargoTransfer, { type TransferCargoEventDetails } from './CargoTransfer.svelte';

	const { transferCargo } = getGameContext();

	export let show = false;
	export let props: CargoTransferDialogEventDetails | undefined;

	const onTransferCargo = async (detail: TransferCargoEventDetails) => {
		if (detail && detail.transferAmount.absoluteSize() > 0) {
			if (!detail.dest) {
				detail.dest = newSalvage();
			}
			await transferCargo(detail.src, detail.dest, detail.transferAmount);
		}

		// close the dialog
		show = false;
	};
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
		{#if props && show}
			<CargoTransfer
				src={props.src}
				dest={props.dest}
				on:transfer-cargo={(e) => onTransferCargo(e.detail)}
				on:cancel={() => (show = false)}
			/>
		{/if}
	</div>
</div>
