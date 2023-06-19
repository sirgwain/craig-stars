<script lang="ts">
	import Path from '$lib/components/icons/Path.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import Scanner from '$lib/components/icons/Scanner.svelte';
	import { nextMapObject, previousMapObject, settings } from '$lib/services/Context';
	import type { Game } from '$lib/types/Game';
	import type { Player } from '$lib/types/Player';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import { Envelope, ArrowLongLeft, ArrowLongRight } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MessagesPane from '../MessagesPane.svelte';
	import PlanetViewStates from './toolbar/PlanetViewStates.svelte';

	export let game: Game;
	export let player: Player;

	let showMessages = false;
</script>

<div class="flex-initial navbar bg-base-100">
	<div class="flex-none hidden sm:block">
		<PlanetViewStates class="menu menu-horizontal" />
	</div>
	<div class="flex-none block sm:hidden">
		<ul class="menu menu-horizontal">
			<!-- submenu -->
			<li>
				<a href="#planet-view-states" class="btn btn-xs h-full">
					{#if $settings.planetViewState == PlanetViewState.Normal}
						<PlanetWithStarbase class="w-6 h-6" />
					{:else if $settings.planetViewState == PlanetViewState.SurfaceMinerals}
						<span>S</span>
					{:else if $settings.planetViewState == PlanetViewState.MineralConcentration}
						<span>C</span>
					{:else if $settings.planetViewState == PlanetViewState.Percent}
						<span>%</span>
					{:else if $settings.planetViewState == PlanetViewState.Population}
						<Population class="w-6 h-6" />
					{:else}
						<PlanetWithStarbase class="w-6 h-6" />
					{/if}
				</a>
				<PlanetViewStates class="menu menu-vertical bg-base-100 z-20" />
			</li>
		</ul>
	</div>

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
		><Icon
			src={Envelope}
			class={`w-6 h-6 ${showMessages ? 'stroke-accent' : 'stroke-base-content'}`}
		/>
		<span class="indicator-item badge badge-secondary">{player.messages.length}</span>
	</a>

	<div class="ml-auto">
		<div class="tooltip" data-tip="previous">
			<button
				on:click={() => previousMapObject()}
				class="btn btn-outline btn-sm normal-case btn-secondary"
				title="previous"
				><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
			>
		</div>
		<div class="tooltip" data-tip="next">
			<button
				on:click={() => nextMapObject()}
				class="btn btn-outline btn-sm normal-case btn-secondary"
				title="next"
				><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
			>
		</div>
	</div>
</div>
<div class:hidden={!showMessages} class:block={showMessages}>
	<MessagesPane {game} {player} />
</div>
