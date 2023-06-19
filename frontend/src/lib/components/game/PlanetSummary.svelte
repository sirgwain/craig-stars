<script lang="ts">
	import { selectedPlanet } from '$lib/services/Context';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { PaperAirplane } from '@steeze-ui/heroicons';
	import PlanetMineralsGraph from './PlanetMineralsGraph.svelte';

	export let title = '';
	$: $selectedPlanet && (title = $selectedPlanet.name);
</script>

<div class="flex flex-col">
	<div class="flex justify-between">
		<div class="ml-[5.5rem]">Value: 100%</div>
		{#if $selectedPlanet.cargo?.colonists}
			<div>Population: {($selectedPlanet.cargo.colonists * 100).toLocaleString()}</div>
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

	<div class="mb-1"></div>

	<PlanetMineralsGraph planet={$selectedPlanet} />
</div>
