<script lang="ts">
	import Population from '#lib/components/icons/Population.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlanetViewState } from '#lib/types/PlayerSettings.js';

	const { settings } = getGameContext();
</script>

<a
	href="#population-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.Population}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.Population}
	onclick={() => ($settings.planetViewState = PlanetViewState.Population)}
	class="btn btn-xs size-full fill-current p-0"><Population class="size-6 " /></a
>
