<script lang="ts">
	import MineralConcentration from '#lib/components/icons/MineralConcentration.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlanetViewState } from '#lib/types/PlayerSettings.js';

	const { settings } = getGameContext();
</script>

<a
	href="#mineral-concenctration-view"
	class:btn-primary={$settings.planetViewState == PlanetViewState.MineralConcentration}
	class:btn-ghost={$settings.planetViewState != PlanetViewState.MineralConcentration}
	onclick={() => ($settings.planetViewState = PlanetViewState.MineralConcentration)}
	class="btn btn-xs h-full"
>
	<MineralConcentration class="w-6 h-6" />
</a>
