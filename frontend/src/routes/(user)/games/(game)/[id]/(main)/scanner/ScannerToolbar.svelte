<script lang="ts">
	import { clickOutside } from '#lib/clickOutside.js';
	import AddWaypoint from '#lib/components/icons/AddWaypoint.svelte';
	import AddWaypointFast from '#lib/components/icons/AddWaypointFast.svelte';
	import Habitability from '#lib/components/icons/Habitability.svelte';
	import MineralConcentration from '#lib/components/icons/MineralConcentration.svelte';
	import PlanetWithStarbase from '#lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '#lib/components/icons/Population.svelte';
	import SurfaceMinerals from '#lib/components/icons/SurfaceMinerals.svelte';
	import type { NextPrevMapObjectProps } from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { PlanetViewState } from '#lib/types/PlayerSettings.js';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		ArrowPath,
		Envelope,
		Funnel,
		MagnifyingGlass
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MobileViewSettings from './toolbar/MobileViewSettings.svelte';
	import PlanetViewStates from './toolbar/PlanetViewStates.svelte';
	import ScannerToolbarFilter from './toolbar/ScannerToolbarFilter.svelte';

	const { player, settings } = getGameContext();

	type Props = {
		onShowSearch: () => void;
		onCycleMapObject: () => void;
	} & NextPrevMapObjectProps;
	let { onShowSearch, onNextMapObject, onPreviousMapObject, onCycleMapObject }: Props = $props();

	let planetsViewMenuDropdown: HTMLDetailsElement | undefined = $state();

	function closePlanetsMenu() {
		planetsViewMenuDropdown?.removeAttribute('open');
	}

	// each planet view and filter button is this wide
	const itemWidth = 40;
	const viewCount = 5;
	// toolbar padding and the dividers between groups
	const paddingWidth = 34;
	// add waypoint and messages
	const actionsWidth = 80;
	// the scanner percent box
	const percentWidth = 72;
	// the button that opens the menu of collapsed items
	const menuWidth = 62;
	let toolbarWidth = $state(0);
	// search/previous/next, plus cycle on small screens
	let navWidth = $state(0);
	// 8 filters, plus 2 ally filters when anyone shares their map with us
	let shareMap = $derived($player.relations.some((r) => r.shareMap));
	let filterCount = $derived(8 + (shareMap ? 2 : 0));
	// the dropdowns beside the design filters
	let filterMenusWidth = $derived(12 * (2 + (shareMap ? 1 : 0)));

	// collapse into a menu in stages as the toolbar runs out of room: first the filters and
	// scanner percent, then the planet views
	let fixedWidth = $derived(paddingWidth + actionsWidth + navWidth);
	let showFilters = $derived(
		toolbarWidth >=
			fixedWidth + percentWidth + filterMenusWidth + itemWidth * (viewCount + filterCount)
	);
	let showViews = $derived(
		showFilters || toolbarWidth >= fixedWidth + menuWidth + itemWidth * viewCount
	);

	// close the planets menu when we switch views (but leave the filter menu open so we can filter more than one thing)
	settings.subscribe(closePlanetsMenu);
</script>

<!-- the menu of everything that doesn't fit in the toolbar -->
{#snippet collapsedMenu()}
	<div class="flex-none">
		<ul class="menu menu-horizontal px-0">
			<li>
				<details bind:this={planetsViewMenuDropdown} use:clickOutside={closePlanetsMenu}>
					<summary class="p-0">
						<a
							href="#planet-view-states"
							class="btn btn-ghost btn-xs size-10 p-0"
							onclick={(e) => {
								e.preventDefault();
								planetsViewMenuDropdown?.toggleAttribute('open');
							}}
						>
							{#if showViews}
								<Icon src={Funnel} class="size-6 " />
							{:else if $settings.planetViewState == PlanetViewState.SurfaceMinerals}
								<SurfaceMinerals class="size-6 " />
							{:else if $settings.planetViewState == PlanetViewState.MineralConcentration}
								<MineralConcentration class="size-6 " />
							{:else if $settings.planetViewState == PlanetViewState.Percent}
								<Habitability class="size-6 " />
							{:else if $settings.planetViewState == PlanetViewState.Population}
								<Population class="size-6 " />
							{:else}
								<PlanetWithStarbase class="size-6 " />
							{/if}
						</a>
					</summary>
					<ul class="bg-base-200 z-20 p-1 w-64">
						<li>
							<MobileViewSettings hideViews={showViews} />
						</li>
					</ul>
				</details>
			</li>
		</ul>
	</div>
{/snippet}

<!-- the toolbar is in the same order as the original's. Add waypoint and messages are always shown -->
<div class="flex-initial navbar bg-base-200 py-0 my-0 min-h-0" bind:clientWidth={toolbarWidth}>
	{#if showViews}
		<div class="flex-none">
			<PlanetViewStates class="menu menu-horizontal px-0" />
		</div>
		<div class="flex-none w-px h-6 mx-1 bg-base-content/20"></div>
	{:else}
		{@render collapsedMenu()}
	{/if}
	<div class="flex-none">
		<ul class="menu menu-horizontal items-center px-0">
			<li class="size-10">
				<a
					href="#add-waypoint"
					id="add-waypoint"
					class:fill-accent={$settings.addWaypoint}
					class:fill-current={!$settings.addWaypoint}
					class="btn btn-ghost btn-xs size-full p-0 border"
					onclick={(e) => {
						e.preventDefault();
						// 3 state toggle
						if ($settings.addWaypoint && $settings.fastestWaypoint) {
							$settings.addWaypoint = $settings.fastestWaypoint = false;
						} else if ($settings.addWaypoint && !$settings.fastestWaypoint) {
							$settings.addWaypoint = true;
							$settings.fastestWaypoint = true;
						} else if (!$settings.addWaypoint && !$settings.fastestWaypoint) {
							$settings.addWaypoint = true;
							$settings.fastestWaypoint = false;
						} else {
							// weird state, reset
							$settings.addWaypoint = $settings.fastestWaypoint = false;
						}
					}}
				>
					{#if $settings.fastestWaypoint}
						<AddWaypointFast class="size-6 " />
					{:else}
						<AddWaypoint class="size-6 " />
					{/if}
				</a>
			</li>
		</ul>
	</div>
	<div class="flex-none w-px h-6 mx-1 bg-base-content/20"></div>
	{#if showFilters}
		<div class="flex-none">
			<ScannerToolbarFilter class="flex flex-row flex-none items-center" />
		</div>
	{:else if showViews}
		{@render collapsedMenu()}
	{/if}
	<div class="flex-none">
		<ul class="menu menu-horizontal items-center px-0">
			<li class="size-10">
				<a
					href="#messages"
					class="btn btn-ghost btn-xs size-full p-0"
					onclick={(e) => {
						e.preventDefault();
						$settings.showMessagePane = !$settings.showMessagePane;
					}}
				>
					<!-- anchor the count to the corner of the icon so it stays inside the toolbar -->
					<span class="indicator">
						<Icon
							src={Envelope}
							class={`size-6  ${$settings.showMessagePane ? 'stroke-accent' : 'stroke-current'}`}
						/>
						<span class="indicator-item badge badge-secondary badge-xs px-1"
							>{$player.messages.length}</span
						>
					</span>
				</a>
			</li>
		</ul>
	</div>

	<div class="ml-auto flex-none" bind:clientWidth={navWidth}>
		<button
			onclick={onShowSearch}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="search"
			><Icon src={MagnifyingGlass} size="16" class="hover:stroke-accent inline" /></button
		>
		<button
			onclick={onCycleMapObject}
			class="md:hidden btn btn-outline btn-sm normal-case btn-secondary"
			title="cycle"><Icon src={ArrowPath} size="16" class="hover:stroke-accent inline" /></button
		>

		<div class="tooltip" data-tip="previous">
			<button
				onclick={onPreviousMapObject}
				class="btn btn-outline btn-sm normal-case btn-secondary"
				title="previous"
				><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
			>
		</div>
		<div class="tooltip" data-tip="next">
			<button
				onclick={onNextMapObject}
				class="btn btn-outline btn-sm normal-case btn-secondary"
				title="next"
				><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
			>
		</div>
	</div>
</div>
