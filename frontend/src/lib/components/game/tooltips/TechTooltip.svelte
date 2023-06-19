<script lang="ts" context="module">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import type { Tech } from '$lib/types/Tech';

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
