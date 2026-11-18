<script module lang="ts">
	import type { TechTooltipProps } from './TechTooltip';
</script>

<script lang="ts">
	import TechSummary from '#lib/components/tech/TechSummary.svelte';
	import { gameKey, type GameContext } from '#lib/services/GameContext.js';
	import { getContext } from 'svelte';
	import { readable } from 'svelte/store';

	let { tech, showResearchCost = false, playerOverride, client }: TechTooltipProps = $props();

	const context = getContext<GameContext | undefined>(gameKey);
	const gamePlayer = context?.player ?? readable(undefined);
	let player = $derived(playerOverride ?? $gamePlayer);
	let cs = $derived(client ?? context?.cs);
</script>

<div class="max-w-[550px]">
	{#if tech}
		<TechSummary {tech} {showResearchCost} {player} hideNew={true} {cs} />
	{/if}
</div>
