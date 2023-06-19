<script lang="ts">
	import { selectedMapObject } from '$lib/services/Context';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import PlanetMineralsGraph from './PlanetMineralsGraph.svelte';
	import { Unexplored, type Planet } from '$lib/types/Planet';

	export let title = '';
	let planet = $selectedMapObject as Planet;

	$: $selectedMapObject && (planet = $selectedMapObject as Planet);
	$: $selectedMapObject && (title = $selectedMapObject.name);
</script>

<div class="flex flex-col min-h-[11rem]">
	{#if planet && planet.reportAge == Unexplored}
		<div class="m-auto">
			<Icon src={QuestionMarkCircle} size="64" class="hover:stroke-accent" />
		</div>
	{:else}
		<div class="flex justify-between">
			<div class="ml-[5.5rem]">Value: 100%</div>
			{#if planet?.population}
				<div>Population: {planet.population.toLocaleString()}</div>
			{/if}
		</div>
		<div class="ml-[5.5rem]">Report is current</div>

		<div class="flex flex-row">
			<div class="text-right w-[5.5rem]">Gravity</div>
			<div class="grow border-b border-base-300 bg-black mx-1" />
			<div class="w-[3rem]">1.00g</div>
		</div>
		<div class="flex flex-row">
			<div class="text-right w-[5.5rem]">Temperature</div>
			<div class="grow border-b border-base-300 bg-black mx-1" />
			<div class="w-[3rem]">0°C</div>
		</div>
		<div class="flex flex-row">
			<div class="text-right w-[5.5rem]">Radiation</div>
			<div class="grow bg-black mx-1" />
			<div class="w-[3rem]">50mR</div>
		</div>

		<div class="mb-1" />

		<PlanetMineralsGraph {planet} />
	{/if}
</div>
