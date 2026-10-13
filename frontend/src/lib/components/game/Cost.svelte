<script lang="ts">
	import type { CostJson as Cost } from '#lib/types/cs-proto.js';

	type Props = {
		cost: Cost | undefined;
		layout?: 'vertical' | 'horizontal';
	};

	let { cost, layout = 'vertical' }: Props = $props();
</script>

{#if cost}
	{#if layout === 'horizontal'}
		<div class="@container">
			<dl class="grid grid-cols-2 @[18rem]:grid-cols-4 gap-2">
				<div class="min-w-0">
					<dt class="whitespace-nowrap text-ironium">Ironium</dt>
					<dd class="mt-1 tabular-nums">{cost.ironium ?? 0}kT</dd>
				</div>
				<div class="min-w-0">
					<dt class="whitespace-nowrap text-boranium">Boranium</dt>
					<dd class="mt-1 tabular-nums">{cost.boranium ?? 0}kT</dd>
				</div>
				<div class="min-w-0">
					<dt class="whitespace-nowrap text-germanium">Germanium</dt>
					<dd class="mt-1 tabular-nums">{cost.germanium ?? 0}kT</dd>
				</div>
				<div class="min-w-0">
					<dt class="whitespace-nowrap">Resources</dt>
					<dd class="mt-1 tabular-nums">{cost.resources ?? 0}</dd>
				</div>
			</dl>
		</div>
	{:else}
		<div class="flex flex-col min-w-32">
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
{/if}
