<script lang="ts">
	import Habitability from '#lib/components/icons/Habitability.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlanetViewState } from '#lib/types/PlayerSettings.js';

	const { settings } = getGameContext();
</script>

<a
	href="#hab-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.Percent}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.Percent}
	onclick={() => ($settings.planetViewState = PlanetViewState.Percent)}
	class="btn btn-xs h-full text-xl font-semibold"><Habitability class="w-6 h-6" /></a
>
