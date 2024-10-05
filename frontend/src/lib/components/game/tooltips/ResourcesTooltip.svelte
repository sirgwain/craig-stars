<script lang="ts" context="module">
	export type ResourcesTooltipProps = {
		planetName: string;
		resourcesPerYear: number;
		resourcesPerYearResearch: number;
		resourcesPerYearResearchEstimated: number;
		resourcesPerYearAvailable: number;
		innateResources: boolean;
	};
</script>

<script lang="ts">
	export let planetName: string;
	export let resourcesPerYear: number;
	export let resourcesPerYearResearch: number;
	export let resourcesPerYearResearchEstimated: number;
	export let resourcesPerYearAvailable: number;
	export let innateResources: boolean;
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div class="text-xl font-semibold text-center">Resources</div>
	<div>
		<span class="font-semibold">{planetName}</span> generates
		<span class="font-semibold">{resourcesPerYear}</span>
		resources each year.
		<span class="font-semibold"
			>{(resourcesPerYearResearch ?? 'None') == 0 ? 'None' : resourcesPerYearResearch}</span
		>
		of these resources have been allocated to research.
		{#if resourcesPerYearResearch != 0}
			That leaves <span class="font-semibold">{resourcesPerYearAvailable}</span> resources avaliable
			for use by the planet.
		{/if}

		{#if innateResources}
			<p>
				Your resources at this planet will increase with population, habitability and <span
					class="font-semibold">Energy</span
				> tech.
			</p>
		{/if}

		{#if resourcesPerYearResearchEstimated > 0}
			Planetary production will leave an estimated <span class="font-semibold"
				>{resourcesPerYearResearchEstimated}</span
			> resources leftover that will be spent on research.
		{/if}
	</div>
</div>
