<script lang="ts">
	import type { CostJson as Cost } from '$lib/types/cs-proto';

	type Props = {
		cost: Cost | undefined;
	};

	let { cost }: Props = $props();
</script>

{#if cost}
	<div class="flex flex-col min-w-[8rem]">
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
			<div>{cost.resources ?? 0}</div>
		</div>
	</div>
{/if}
