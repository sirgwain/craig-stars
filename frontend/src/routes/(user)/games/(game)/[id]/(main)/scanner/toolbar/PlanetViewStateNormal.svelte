<script lang="ts">
	import PlanetWithStarbase from '#lib/components/icons/PlanetWithStarbase.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlanetViewState } from '#lib/types/PlayerSettings.js';

	const { settings } = getGameContext();
</script>

<a
	href="#normal-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.Normal}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.Normal}
	onclick={() => ($settings.planetViewState = PlanetViewState.Normal)}
	class="btn btn-xs size-full p-0"><PlanetWithStarbase class="size-6 " /></a
>
