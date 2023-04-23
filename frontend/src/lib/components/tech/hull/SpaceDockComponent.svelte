<script lang="ts">
	import { UnlimitedSpaceDock } from '$lib/types/Tech';

	export let spaceDock: number;
	export let rounded = false;
</script>

<div
	class={`flex flex-col h-full w-full justify-between border border-slate-900 bg-white text-sm text-black ${
		rounded ? 'rounded-full p-3' : ''
	}`}
>
	{#if spaceDock === UnlimitedSpaceDock}
		<div class="text-center font-bold">Unlimited</div>
	{:else}
		<div class="text-center font-bold">{spaceDock}kT</div>
	{/if}
	<div class="text-center">Space Dock</div>
</div>
