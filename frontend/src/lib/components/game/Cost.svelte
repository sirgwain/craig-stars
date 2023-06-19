<script lang="ts">
	import type { Cost } from '$lib/types/Cost';

	export let cost: Cost | undefined;
</script>

{#if cost}
	<div class="flex justify-between">
		<div class="text-ironium">Ironium</div>
		<div>{cost.ironium ?? 0}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="text-boranium">Boranium</div>
		<div>{cost.boranium ?? 0}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="text-germanium">Germanium</div>
		<div>{cost.germanium ?? 0}kT</div>
	</div>
	<div class="flex justify-between">
		<div>Resources</div>
		<div>{cost.resources ?? 0}kT</div>
	</div>
{/if}
