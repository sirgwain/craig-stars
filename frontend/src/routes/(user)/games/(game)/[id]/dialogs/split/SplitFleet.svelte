<script lang="ts">
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import hotkeys from 'hotkeys-js';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let fleet: CommandedFleet;

	function ok() {}

	function cancel() {
		dispatch('cancel');
	}

	hotkeys('Esc', () => cancel());
</script>

<div class="flex flex-row justify-center px-1 w-full">
	<div class="flex flex-col grow">
		<div class="text-xl font-semibold w-full text-center">Split Fleet</div>
		<div
			class="border border-secondary bg-base-300 min-w-fit max-h-[26rem] h-full overflow-y-auto"
		/>
	</div>
	<div class="flex flex-col mt-7 ml-2 gap-2">
		<button on:click|preventDefault={ok} type="submit" class="btn btn-sm normal-case btn-primary"
			>OK</button
		>
		<button on:click={cancel} class="btn btn-outline btn-sm normal-case btn-secondary"
			>Cancel</button
		>
	</div>
</div>
