<script lang="ts">
	import MineralConcentration from '$lib/components/icons/MineralConcentration.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import SurfaceMinerals from '$lib/components/icons/SurfaceMinerals.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { PlanetViewState } from '$lib/types/PlayerSettings';

	const { game, player, universe, settings } = getGameContext();
</script>

<ul {...$$restProps}>
	<li>
		<a
			href="#normal-view"
			class:btn-primary={$settings.planetViewState == PlanetViewState.Normal}
			class:btn-ghost={$settings.planetViewState != PlanetViewState.Normal}
			on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Normal)}
			class="btn btn-xs h-full"><PlanetWithStarbase class="w-6 h-6" /></a
		>
	</li>
	<li>
		<a
			href="#surface-minerals-view"
			class:btn-primary={$settings.planetViewState == PlanetViewState.SurfaceMinerals}
			class:btn-ghost={$settings.planetViewState != PlanetViewState.SurfaceMinerals}
			on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.SurfaceMinerals)}
			class="btn btn-xs h-full"><SurfaceMinerals class="w-6 h-6"/></a
		>
	</li>
	<li>
		<a
			href="#mineral-concenctration-view"
			class:btn-primary={$settings.planetViewState == PlanetViewState.MineralConcentration}
			class:btn-ghost={$settings.planetViewState != PlanetViewState.MineralConcentration}
			on:click|preventDefault={() =>
				($settings.planetViewState = PlanetViewState.MineralConcentration)}
			class="btn btn-xs h-full"
		>
			<MineralConcentration class="w-6 h-6" />
		</a>
	</li>
	<li>
		<a
			href="#hab-view"
			class:btn-primary={$settings.planetViewState == PlanetViewState.Percent}
			class:btn-ghost={$settings.planetViewState != PlanetViewState.Percent}
			on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Percent)}
			class="btn btn-xs h-full text-xl font-semibold"><span>%</span></a
		>
	</li>
	<li>
		<a
			href="#population-view"
			class:btn-primary={$settings.planetViewState == PlanetViewState.Population}
			class:btn-ghost={$settings.planetViewState != PlanetViewState.Population}
			on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Population)}
			class="btn btn-xs h-full fill-current"><Population class="w-6 h-6" /></a
		>
	</li>
</ul>
