<script lang="ts">
	import Habitability from '$lib/components/icons/Habitability.svelte';
	import MineralConcentration from '$lib/components/icons/MineralConcentration.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import SurfaceMinerals from '$lib/components/icons/SurfaceMinerals.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { PlanetViewState } from '$lib/types/PlayerSettings';

	const { settings } = getGameContext();
</script>

<a
	href="#hab-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.Percent}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.Percent}
	on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Percent)}
	class="btn btn-xs h-full text-xl font-semibold"><Habitability class="w-6 h-6" /></a
>
