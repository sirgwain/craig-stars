<script lang="ts">
	import type {
		OnCancel,
		OnOk,
		SplitFleetDialogEvent,
		SplitFleetEvent
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import SplitFleet from './SplitFleet.svelte';

	const { split } = getGameContext();

	type Props = {
		show?: boolean;
		props: SplitFleetDialogEvent | undefined;
		onOk: OnOk<SplitFleetEvent>;
		onCancel: OnCancel;
	};

	let { show = false, props, onOk, onCancel }: Props = $props();
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
		{#if props && show}
			<SplitFleet src={props.src} dest={props.dest} {onOk} {onCancel} />
		{/if}
	</div>
</div>
