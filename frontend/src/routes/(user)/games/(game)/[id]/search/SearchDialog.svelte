<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import type { OnCancel, OnOk } from '$lib/services/Events';
	import { type MapObject } from '$lib/types/MapObject';
	import SearchResults from './SearchResults.svelte';

	type Props = {
		show?: boolean;
		onOk: OnOk<MapObject | undefined>;
		onCancel: OnCancel;
	};

	let { show = false, onOk, onCancel }: Props = $props();
</script>

<div class="modal" class:modal-open={show}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full md:max-w-[40rem] md:max-h-[48rem] p-2"
		use:clickOutside={onCancel}
	>
		{#if show}
			<SearchResults {onOk} {onCancel} />
		{/if}
	</div>
</div>
