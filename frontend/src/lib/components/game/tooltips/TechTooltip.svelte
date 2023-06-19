<script lang="ts" context="module">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { showTooltip } from '$lib/services/Stores';
	import type { Tech } from '$lib/types/Tech';
	import TechTooltip from './TechTooltip.svelte';

	export function onTechTooltip(e: PointerEvent, tech: Tech | undefined) {
		if (tech) {
			showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, { tech });
		}
	}


	export type TechTooltipProps = {
		tech: Tech;
	};
</script>

<script lang="ts">
	export let tech: Tech;
</script>

<div class="md:w-[380px] h-[380px]">
	<TechSummary {tech} />
</div>
