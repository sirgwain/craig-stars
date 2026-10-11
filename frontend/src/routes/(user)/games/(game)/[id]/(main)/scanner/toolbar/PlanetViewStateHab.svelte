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
	class="btn btn-xs size-full text-xl font-semibold p-0"><Habitability class="size-6 " /></a
>
