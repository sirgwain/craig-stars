<script lang="ts">
	import type { Cost } from '$lib/types/Cost';

	export let cost: Cost | undefined;
</script>

{#if cost}
	<div>
		<span class="text-ironium">i: </span><span>{cost.ironium ?? 0}kT</span>
		<span class="text-boranium">b: </span><span>{cost.boranium ?? 0}kT</span>
		<span class="text-germanium">g: </span><span>{cost.germanium ?? 0}kT</span>
		<span>r: </span><span>{cost.resources ?? 0}kT</span>
	</div>
{/if}
