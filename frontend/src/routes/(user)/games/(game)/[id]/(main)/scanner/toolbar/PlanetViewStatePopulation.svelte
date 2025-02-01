<script lang="ts">
	import Population from '$lib/components/icons/Population.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { PlanetViewState } from '$lib/types/PlayerSettings';

	const { settings } = getGameContext();
</script>

<a
	href="#population-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.Population}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.Population}
	onclick={() => ($settings.planetViewState = PlanetViewState.Population)}
	class="btn btn-xs h-full fill-current"><Population class="w-6 h-6" /></a
>
