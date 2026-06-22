<script module lang="ts">
	import type { TechTooltipProps } from './TechTooltip';
</script>

<script lang="ts">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	let { tech, showResearchCost = false }: TechTooltipProps = $props();

	const { player, cs } = getGameContext();
</script>

<div class="max-w-[550px]">
	{#if tech}
		<TechSummary {tech} {showResearchCost} player={$player} hideNew={true} {cs} />
	{/if}
</div>
