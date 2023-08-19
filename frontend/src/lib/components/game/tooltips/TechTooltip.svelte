<script lang="ts" context="module">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { showTooltip } from '$lib/services/Stores';
	import type { Tech } from '$lib/types/Tech';
	import TechTooltip from './TechTooltip.svelte';

	export function onTechTooltip(
		e: PointerEvent | MouseEvent,
		tech: Tech | undefined,
		showResearchCost = false
	) {
		if (tech) {
			showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, { tech, showResearchCost });
		}
	}

	export type TechTooltipProps = {
		tech: Tech;
		showResearchCost?: boolean;
	};
</script>

<script lang="ts">
	export let tech: Tech;
	export let showResearchCost = false;

	const { game, player } = getGameContext();
</script>

<div class="md:w-[380px] h-[420px]">
	<TechSummary {tech} {showResearchCost} player={$player} game={$game} />
</div>
