<script lang="ts" context="module">
	export type ToolbarEvent = {
		'show-search': void;
	};
</script>

<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import Habitability from '$lib/components/icons/Habitability.svelte';
	import MineralConcentration from '$lib/components/icons/MineralConcentration.svelte';
	import Path from '$lib/components/icons/Path.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import SurfaceMinerals from '$lib/components/icons/SurfaceMinerals.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		Envelope,
		MagnifyingGlass
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import MessagesPane from '../MessagesPane.svelte';
	import MobileViewSettings from './toolbar/MobileViewSettings.svelte';
	import PlanetViewStates from './toolbar/PlanetViewStates.svelte';
	import ScannerToolbarFilter from './toolbar/ScannerToolbarFilter.svelte';

	const { player, settings, nextMapObject, previousMapObject } = getGameContext();
	const dispatch = createEventDispatcher<ToolbarEvent>();

	let planetsViewMenuDropdown: HTMLDetailsElement | undefined;

	function closePlanetsMenu() {
		planetsViewMenuDropdown?.removeAttribute('open');
	}

	// close the planets menu when we switch views (but leave the filter menu open so we can filter more than one thing)
	settings.subscribe(closePlanetsMenu);
</script>

<div class="flex-initial navbar bg-base-200 py-0 my-0 min-h-0">
	<div class="flex-none hidden lg:block">
		<PlanetViewStates class="menu menu-horizontal" />
	</div>
	<div class="flex-none hidden lg:block">
		<ScannerToolbarFilter class="menu menu-horizontal" />
	</div>
	<div class="flex-none block lg:hidden">
		<ul class="menu menu-horizontal">
			<!-- submenu -->
			<li>
				<details bind:this={planetsViewMenuDropdown} use:clickOutside={closePlanetsMenu}>
					<summary class="p-0">
						<a
							href="#planet-view-states"
							class="btn btn-xs w-12 h-12"
							on:click|preventDefault={() => planetsViewMenuDropdown?.toggleAttribute('open')}
						>
							{#if $settings.planetViewState == PlanetViewState.Normal}
								<PlanetWithStarbase class="w-6 h-6" />
							{:else if $settings.planetViewState == PlanetViewState.SurfaceMinerals}
								<SurfaceMinerals class="w-6 h-6" />
							{:else if $settings.planetViewState == PlanetViewState.MineralConcentration}
								<MineralConcentration class="w-6 h-6" />
							{:else if $settings.planetViewState == PlanetViewState.Percent}
								<Habitability class="w-6 h-6" />
							{:else if $settings.planetViewState == PlanetViewState.Population}
								<Population class="w-6 h-6" />
							{:else}
								<PlanetWithStarbase class="w-6 h-6" />
							{/if}
						</a>
					</summary>
					<ul class="bg-base-200 z-20 p-1 w-64">
						<li>
							<MobileViewSettings />
						</li>
					</ul>
				</details>
			</li>
		</ul>
	</div>

	<div class="flex-none">
		<ul class="menu menu-horizontal">
			<li class="hidden sm:block">
				<div class="w-full px-1">
					<input
						class="input input-sm input-bordered w-16 pr-0 pl-1"
						type="number"
						min={0}
						max={100}
						step={10}
						value={$settings.scannerPercent}
						on:change={(e) => {
							const val = parseInt(e.currentTarget.value);
							if (val) {
								$settings.scannerPercent = clamp(val, 0, 100);
								e.currentTarget.value = `${$settings.scannerPercent}`;
							} else {
								$settings.scannerPercent = 0;
								e.currentTarget.value = '0';
							}
						}}
					/>
				</div>
			</li>
			<li>
				<a
					href="#add-waypoint"
					id="add-waypoint"
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
		<button
			on:click={() => dispatch('show-search')}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="previous"
			><Icon src={MagnifyingGlass} size="16" class="hover:stroke-accent inline" /></button
		>

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
