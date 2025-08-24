<script lang="ts">
	import type { Cost } from '$lib/types/cs-proto';

	type Props = {
		cost: Cost | undefined;
	};

	let { cost }: Props = $props();
</script>

{#if cost}
	<div>
		<span class="text-ironium">i: </span><span>{cost.ironium}kT</span>
		<span class="text-boranium">b: </span><span>{cost.boranium}kT</span>
		<span class="text-germanium">g: </span><span>{cost.germanium}kT</span>
		<span>r: </span><span>{cost.resources}</span>
	</div>
{/if}
