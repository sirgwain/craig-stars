<script lang="ts" module>
	import { type MapObject } from '$lib/types/MapObject';
	export type SearchDialogEvent = {
		'select-result': MapObject | undefined;
	};
</script>

<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import SearchResults from './SearchResults.svelte';
	import { clickOutside } from '$lib/clickOutside';

	const dispatch = createEventDispatcher<SearchDialogEvent>();

	interface Props {
		show?: boolean;
	}

	let { show = $bindable(false) }: Props = $props();

	function onOk(mo: MapObject | undefined) {
		show = false;
		dispatch('select-result', mo);
	}
</script>

<div class="modal" class:modal-open={show}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full md:max-w-[40rem] md:max-h-[48rem] p-2"
		use:clickOutside={() => (show = false)}
	>
		{#if show}
			<SearchResults on:ok={(e) => onOk(e.detail)} on:cancel={() => (show = false)} />
		{/if}
	</div>
</div>
