<script lang="ts">
	import Path from '$lib/components/icons/Path.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import Scanner from '$lib/components/icons/Scanner.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import { ArrowLongLeft, ArrowLongRight, Envelope } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MessagesPane from '../MessagesPane.svelte';
	import PlanetViewStates from './toolbar/PlanetViewStates.svelte';
	import MineralConcentration from '$lib/components/icons/MineralConcentration.svelte';
	import SurfaceMinerals from '$lib/components/icons/SurfaceMinerals.svelte';
	import { clamp } from '$lib/services/Math';
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import { clickOutside } from '$lib/clickOutside';

	const { player, settings, nextMapObject, previousMapObject } = getGameContext();
	let menuDropdown: HTMLDetailsElement | undefined;

	function closeMenu() {
		menuDropdown?.removeAttribute('open');
	}
	settings.subscribe(closeMenu);
</script>

<div class="flex-initial navbar bg-base-200 py-0 my-0 min-h-0">
	<div class="flex-none hidden sm:block">
		<PlanetViewStates class="menu menu-horizontal" />
	</div>
	<div class="flex-none block sm:hidden">
		<ul class="menu menu-horizontal">
			<!-- submenu -->
			<li>
				<details bind:this={menuDropdown} use:clickOutside={closeMenu}>
					<summary>
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
					</summary>
					<PlanetViewStates class="menu menu-vertical bg-base-100 z-20 w-12" />
				</details>
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
					href="#show-fleet-token-count"
					class:fill-accent={$settings.showFleetTokenCounts}
					class:fill-current={!$settings.showFleetTokenCounts}
					class="btn btn-ghost btn-xs h-full border"
					on:click|preventDefault={() =>
						($settings.showFleetTokenCounts = !$settings.showFleetTokenCounts)}
					><FleetCount class="w-6 h-6" /></a
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
