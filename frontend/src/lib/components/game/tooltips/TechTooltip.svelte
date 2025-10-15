<script lang="ts" module>
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip } from '$lib/services/Stores';
	import type { TechLike } from '$lib/types/Tech';
	import TechTooltip from './TechTooltip.svelte';

	export function onTechTooltip(
		e: PointerEvent | MouseEvent,
		tech: TechLike | undefined,
		showResearchCost = false
	) {
		e.preventDefault();
		if (tech) {
			showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, { tech, showResearchCost });
		}
	}

	export type TechTooltipProps = {
		tech: TechLike | undefined;
		showResearchCost?: boolean;
	};
</script>

<script lang="ts">
	let { tech, showResearchCost = false }: TechTooltipProps = $props();

	const { player, cs } = getGameContext();
</script>

<div class="max-w-[550px]">
	{#if tech}
		<TechSummary {tech} {showResearchCost} player={$player} hideNew={true} {cs} />
	{/if}
</div>
