<script lang="ts">
	import Path from '$lib/components/icons/Path.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import Scanner from '$lib/components/icons/Scanner.svelte';
	import { settings } from '$lib/services/Context';
	import type { Game } from '$lib/types/Game';
	import type { Player } from '$lib/types/Player';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import { Envelope } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MessagesPane from '../MessagesPane.svelte';

	export let game: Game;
	export let player: Player;

	let showMessages = false;
</script>

<div class="flex-initial navbar bg-base-100">
	<a
		href="#normal-view"
		class:btn-primary={$settings.planetViewState == PlanetViewState.Normal}
		class:btn-ghost={$settings.planetViewState != PlanetViewState.Normal}
		on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Normal)}
		class="btn btn-xs h-full"><PlanetWithStarbase class="w-6 h-6" /></a
	>
	<a
		href="#surface-minerals-view"
		class:btn-primary={$settings.planetViewState == PlanetViewState.SurfaceMinerals}
		class:btn-ghost={$settings.planetViewState != PlanetViewState.SurfaceMinerals}
		on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.SurfaceMinerals)}
		class="btn btn-xs h-full"><span>S</span></a
	>
	<a
		href="#mineral-concenctration-view"
		class:btn-primary={$settings.planetViewState == PlanetViewState.MineralConcentration}
		class:btn-ghost={$settings.planetViewState != PlanetViewState.MineralConcentration}
		on:click|preventDefault={() =>
			($settings.planetViewState = PlanetViewState.MineralConcentration)}
		class="btn btn-xs h-full"><span>C</span></a
	>
	<a
		href="#hab-view"
		class:btn-primary={$settings.planetViewState == PlanetViewState.Percent}
		class:btn-ghost={$settings.planetViewState != PlanetViewState.Percent}
		on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Percent)}
		class="btn btn-xs h-full"><span>%</span></a
	>
	<a
		href="#population-view"
		class:btn-primary={$settings.planetViewState == PlanetViewState.Population}
		class:btn-ghost={$settings.planetViewState != PlanetViewState.Population}
		on:click|preventDefault={() => ($settings.planetViewState = PlanetViewState.Population)}
		class="btn btn-xs h-full fill-base-content"><span><Population class="w-6 h-6" /></span></a
	>
	<a
		href="#scanner-toggle"
		class:fill-accent={$settings.showScanners}
		class:fill-base-content={!$settings.showScanners}
		class="btn btn-ghost btn-xs h-full border"
		on:click|preventDefault={() => ($settings.showScanners = !$settings.showScanners)}
		><span><Scanner class="w-6 h-6" /></span></a
	>
	<a
		href="#add-waypoint"
		class:fill-accent={$settings.addWaypoint}
		class:fill-base-content={!$settings.addWaypoint}
		class="btn btn-ghost btn-xs h-full border"
		on:click|preventDefault={() => ($settings.addWaypoint = !$settings.addWaypoint)}
		><Path class="w-6 h-6" /></a
	>
	<a
		href="#add-waypoint"
		class="btn btn-ghost btn-xs h-full indicator"
		on:click|preventDefault={() => (showMessages = !showMessages)}
		><Icon src={Envelope} class="w-6 h-6" />
		<span class="indicator-item badge badge-secondary">{player.messages.length}</span>
	</a>
</div>
<div class:hidden={!showMessages} class:block={showMessages}>
	<MessagesPane {game} {player} />
</div>
