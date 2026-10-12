<script lang="ts">
	import { clickOutside } from '#lib/clickOutside.js';
	import type { OnCancel, ShowCargoTransferDialogProps } from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MapObjectSummary from './MapObjectSummary.svelte';

	// The selected map object's summary as a full screen popup, for windows with no room to keep
	// it beside the scanner.
	type Props = {
		show?: boolean;
		onCancel: OnCancel;
	} & ShowCargoTransferDialogProps;

	let { show = false, onCancel, onShowCargoTransferDialog }: Props = $props();

	const { editingNote } = getGameContext();
</script>

<div class="modal" class:modal-open={show} data-type="map-object-summary-dialog">
	<div
		class="modal-box max-w-full max-h-max size-full p-2 flex flex-col gap-1"
		use:clickOutside={() => {
			// the note dialog opens over this one, taps in it aren't taps outside
			if (!$editingNote) {
				onCancel();
			}
		}}
	>
		{#if show}
			<div class="flex flex-row justify-end">
				<button
					type="button"
					class="btn btn-sm btn-outline btn-secondary"
					title="close"
					aria-label="Close summary"
					onclick={() => onCancel()}
				>
					<Icon src={XMark} size="16" class="hover:stroke-accent" />
				</button>
			</div>
			<div class="overflow-y-auto">
				<MapObjectSummary {onShowCargoTransferDialog} />
			</div>
		{/if}
	</div>
</div>
