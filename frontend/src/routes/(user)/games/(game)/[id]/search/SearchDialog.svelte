<script lang="ts" context="module">
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

	export let show = false;

	function onOk(mo: MapObject | undefined) {
		show = false;
		dispatch('select-result', mo);
	}
</script>

<div class="modal" class:modal-open={show}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-2"
		use:clickOutside={() => (show = false)}
	>
		{#if show}
			<SearchResults on:ok={(e) => onOk(e.detail)} on:cancel={() => (show = false)} />
		{/if}
	</div>
</div>
