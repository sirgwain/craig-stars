<script lang="ts">
	import type {
		CargoTransferDialogEvent,
		OnCancel,
		OnOk,
		TransferCargoEvent
	} from '#lib/services/Events.js';
	import CargoTransfer from './CargoTransfer.svelte';

	type Props = {
		show?: boolean;
		props: CargoTransferDialogEvent | undefined;
		onOk: OnOk<TransferCargoEvent>;
		onCancel: OnCancel;
	};

	let { show, props, onOk, onCancel }: Props = $props();
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max size-full lg:max-w-160 lg:max-h-192">
		{#if props && show}
			<CargoTransfer src={props.src} dest={props.dest} {onOk} {onCancel} />
		{/if}
	</div>
</div>
