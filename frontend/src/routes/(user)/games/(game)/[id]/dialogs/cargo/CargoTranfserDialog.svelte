<script lang="ts" context="module">
	import type { Cargo } from '$lib/types/Cargo';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import type { Planet } from '$lib/types/Planet';

	export type CargoTransferDialogEventDetails = {
		src: CommandedFleet;
		dest?: Fleet | Planet | Salvage | undefined;
	};
	export type TransferCargoEventDetails = {
		src: CommandedFleet;
		dest?: Fleet | Planet | Salvage;
		transferAmount: Cargo;
	};
	export type CargoTransferEvent = {
		'cargo-transfer-dialog': CargoTransferDialogEventDetails;
		'transfer-cargo': TransferCargoEventDetails;
		cancel: void;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { subtract } from '$lib/types/Cargo';
	import CargoTransfer from './CargoTransfer.svelte';
	import { newSalvage, type Salvage } from '$lib/types/Salvage';

	const { game, player, universe } = getGameContext();

	export let show = false;
	export let props: CargoTransferDialogEventDetails | undefined;

	const onTransferCargo = async (detail: TransferCargoEventDetails) => {
		if (detail) {
			if (!detail.dest) {
				detail.dest = newSalvage();
			}
			await $game.transferCargo(detail.src, detail.dest, detail.transferAmount);

			if (detail?.dest?.cargo) {
				detail.dest.cargo = subtract(detail.dest.cargo, detail.transferAmount);
			}
		}

		// close the dialog
		show = false;
	};
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
		{#if props}
			<CargoTransfer
				src={props.src}
				dest={props.dest}
				on:transfer-cargo={(e) => onTransferCargo(e.detail)}
				on:cancel={() => (show = false)}
			/>
		{/if}
	</div>
</div>
