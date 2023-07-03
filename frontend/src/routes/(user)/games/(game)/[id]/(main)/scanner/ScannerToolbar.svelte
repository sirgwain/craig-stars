<script lang="ts">
	import Path from '$lib/components/icons/Path.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import Scanner from '$lib/components/icons/Scanner.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { nextMapObject, previousMapObject } from '$lib/services/Stores';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import { ArrowLongLeft, ArrowLongRight, Envelope } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MessagesPane from '../MessagesPane.svelte';
	import PlanetViewStates from './toolbar/PlanetViewStates.svelte';
	import MineralConcentration from '$lib/components/icons/MineralConcentration.svelte';
	import SurfaceMinerals from '$lib/components/icons/SurfaceMinerals.svelte';

	const { game, player, universe, settings } = getGameContext();
</script>

<div class="flex-initial navbar bg-base-100">
	<div class="flex-none hidden sm:block">
		<PlanetViewStates class="menu menu-horizontal" />
	</div>
	<div class="flex-none block sm:hidden">
		<ul class="menu menu-horizontal">
			<!-- submenu -->
			<li>
				<a href="#planet-view-states" class="btn btn-primary btn-xs w-12 h-12">
					{#if $settings.planetViewState == PlanetViewState.Normal}
						<PlanetWithStarbase class="w-6 h-6" />
					{:else if $settings.planetViewState == PlanetViewState.SurfaceMinerals}
						<SurfaceMinerals class="w-6 h-6" />
					{:else if $settings.planetViewState == PlanetViewState.MineralConcentration}
						<MineralConcentration class="w-6 h-6" />
					{:else if $settings.planetViewState == PlanetViewState.Percent}
						<span>%</span>
					{:else if $settings.planetViewState == PlanetViewState.Population}
						<Population class="w-6 h-6" />
					{:else}
						<PlanetWithStarbase class="w-6 h-6" />
					{/if}
				</a>
				<PlanetViewStates class="menu menu-vertical bg-base-100 z-20 w-12" />
			</li>
		</ul>
	</div>

	<div class="flex-none">
		<ul class="menu menu-horizontal">
			<li>
				<a
					href="#scanner-toggle"
					class:fill-accent={$settings.showScanners}
					class:fill-current={!$settings.showScanners}
					class="btn btn-ghost btn-xs h-full border"
					on:click|preventDefault={() => ($settings.showScanners = !$settings.showScanners)}
					><span><Scanner class="w-6 h-6" /></span></a
				>
			</li>
			<li>
				<a
					href="#add-waypoint"
					class:fill-accent={$settings.addWaypoint}
					class:fill-current={!$settings.addWaypoint}
					class="btn btn-ghost btn-xs h-full border"
					on:click|preventDefault={() => ($settings.addWaypoint = !$settings.addWaypoint)}
					><Path class="w-6 h-6" /></a
				>
			</li>
			<li>
				<a
					href="#messages"
					class="btn btn-ghost btn-xs h-full indicator"
					on:click|preventDefault={() => ($settings.showMessagePane = !$settings.showMessagePane)}
					><Icon
						src={Envelope}
						class={`w-6 h-6 ${$settings.showMessagePane ? 'stroke-accent' : 'stroke-current'}`}
					/>
					<span class="indicator-item badge badge-secondary">{$player.messages.length}</span>
				</a>
			</li>
		</ul>
	</div>

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
<MessagesPane bind:showMessages={$settings.showMessagePane} messages={$player.messages} />
