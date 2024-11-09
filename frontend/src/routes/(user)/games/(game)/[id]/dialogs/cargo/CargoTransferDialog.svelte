<script lang="ts">
	import type {
		CargoTransferDialogEvent,
		OnCancel,
		OnOk,
		TransferCargoEvent
	} from '$lib/services/Events';
	import CargoTransfer from './CargoTransfer.svelte';

	interface Props {
		show?: boolean;
		props: CargoTransferDialogEvent | undefined;
		onOk: OnOk<TransferCargoEvent>;
		onCancel: OnCancel;
	}

	let { show, props, onOk, onCancel }: Props = $props();
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
		{#if props && show}
			<CargoTransfer src={props.src} dest={props.dest} {onOk} {onCancel} />
		{/if}
	</div>
</div>
