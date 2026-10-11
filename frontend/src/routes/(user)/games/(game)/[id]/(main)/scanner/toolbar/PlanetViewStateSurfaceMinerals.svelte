<script lang="ts">
	import SurfaceMinerals from '#lib/components/icons/SurfaceMinerals.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlanetViewState } from '#lib/types/PlayerSettings.js';

	const { settings } = getGameContext();
</script>

<a
	href="#surface-minerals-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.SurfaceMinerals}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.SurfaceMinerals}
	onclick={() => ($settings.planetViewState = PlanetViewState.SurfaceMinerals)}
	class="btn btn-xs size-full p-0"><SurfaceMinerals class="size-6 " /></a
>
