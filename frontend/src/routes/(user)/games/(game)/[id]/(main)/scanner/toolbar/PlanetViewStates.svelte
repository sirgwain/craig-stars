<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import PlanetViewStateHab from './PlanetViewStateHab.svelte';
	import PlanetViewStateMineralConc from './PlanetViewStateMineralConc.svelte';
	import PlanetViewStatesNormal from './PlanetViewStateNormal.svelte';
	import PlanetViewStatePopulation from './PlanetViewStatePopulation.svelte';
	import PlanetViewStateSurfaceMinerals from './PlanetViewStateSurfaceMinerals.svelte';

	const { settings } = getGameContext();
</script>

<ul {...$$restProps}>
	<li class="h-10 w-10">
		<PlanetViewStatesNormal />
	</li>
	<li class="h-10 w-10">
		<PlanetViewStateSurfaceMinerals />
	</li>
	<li class="h-10 w-10">
		<PlanetViewStateMineralConc />
	</li>
	<li class="h-10 w-10">
		<PlanetViewStateHab />
	</li>
	<li class="h-10 w-10">
		<PlanetViewStatePopulation />
	</li>
</ul>
