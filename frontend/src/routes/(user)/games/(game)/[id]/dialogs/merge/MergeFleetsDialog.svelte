<script lang="ts">
	import type {
		MergeFleetsDialogEvent,
		MergeFleetsEvent,
		OnCancel,
		OnOk
	} from '#lib/services/Events.js';
	import MergeFleets from './MergeFleets.svelte';

	type Props = {
		show?: boolean;
		props: MergeFleetsDialogEvent | undefined;
		onOk?: OnOk<MergeFleetsEvent>;
		onCancel?: OnCancel;
	};

	let { show = false, props, onOk, onCancel }: Props = $props();
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max size-full md:max-w-lg md:max-h-128">
		{#if props && show}
			<MergeFleets fleet={props.fleet} otherFleetsHere={props.otherFleetsHere} {onOk} {onCancel} />
		{/if}
	</div>
</div>
