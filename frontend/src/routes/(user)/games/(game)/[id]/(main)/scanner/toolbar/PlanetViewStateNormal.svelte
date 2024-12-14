<script lang="ts">
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { PlanetViewState } from '$lib/types/PlayerSettings';

	const { settings } = getGameContext();
</script>

<a
	href="#normal-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.Normal}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.Normal}
	onclick={() => ($settings.planetViewState = PlanetViewState.Normal)}
	class="btn btn-xs h-full"><PlanetWithStarbase class="w-6 h-6" /></a
>
