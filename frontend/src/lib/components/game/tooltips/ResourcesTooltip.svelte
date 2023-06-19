<script lang="ts" context="module">
	export type ResourcesTooltipProps = {
		planetName: string;
		resourcesPerYear: number;
		resourcesPerYearResearch: number;
		resourcesPerYearAvailable: number;
		innateResources: boolean;
	};
</script>

<script lang="ts">
	export let planetName: string;
	export let resourcesPerYear: number;
	export let resourcesPerYearResearch: number;
	export let resourcesPerYearAvailable: number;
	export let innateResources: boolean;
</script>

<div class="flex flex-col sm:w-[26rem]">
	<div class="text-xl font-semibold text-center">Resources</div>
	<div>
		<span class="font-semibold">{planetName}</span> generates
		<span class="font-semibold">{resourcesPerYear}</span>
		resources each year. <span class="font-semibold">{resourcesPerYearResearch}</span> of these
		resources have been alloocated to research. That leaves
		<span class="font-semibold">{resourcesPerYearAvailable}</span>
		resources for use by the planet.
		{#if innateResources}
			Your resources at the planet increase with population and <span class="font-semibold"
				>Energy</span
			> tech.
		{/if}
	</div>
</div>
