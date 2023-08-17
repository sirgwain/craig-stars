<script lang="ts">
	import type { PopulationTooltipProps } from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import PopulationTooltip from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import ResourcesTooltip, {
		type ResourcesTooltipProps
	} from '$lib/components/game/tooltips/ResourcesTooltip.svelte';
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { showTooltip, techs } from '$lib/services/Stores';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	const { player, universe } = getGameContext();

	export let planet: CommandedPlanet;

	function onResourcesTooltip(e: PointerEvent) {
		showTooltip<ResourcesTooltipProps>(e.x, e.y, ResourcesTooltip, {
			planetName: planet.name,
			resourcesPerYear: planet.spec.resourcesPerYear,
			resourcesPerYearAvailable: planet.spec.resourcesPerYearAvailable,
			resourcesPerYearResearch: planet.spec.resourcesPerYearResearch,
			innateResources: $player.race.spec?.innateResources ?? false
		});
	}

	function onPopulationTooltip(e: PointerEvent) {
		showTooltip<PopulationTooltipProps>(e.x, e.y, PopulationTooltip, {
			playerFinder: $universe,
			player: $player,
			planet
		});
	}

	function onScannerPopup(e: PointerEvent) {
		onTechTooltip(e, $techs.getTech(planet.spec.scanner));
	}
	function onDefensePoopup(e: PointerEvent) {
		onTechTooltip(e, $techs.getTech(planet.spec.defense));
	}
</script>

{#if planet.spec && planet.cargo}
	<CommandTile title="Status">
		<div
			class="flex justify-between cursor-help"
			on:pointerdown|preventDefault={onPopulationTooltip}
		>
			<div>Population</div>
			<div>{(planet.cargo.colonists ?? 0) * 100}</div>
		</div>
		<div
			class="flex justify-between cursor-help"
			on:pointerdown|preventDefault={onResourcesTooltip}
		>
			<div>Resources/Year</div>
			<div>
				{planet.spec.resourcesPerYearAvailable} of {planet.spec.resourcesPerYear}
			</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onScannerPopup}>
			<div>Scanner Type</div>
			<div>{planet.spec.scanner ?? 'none'}</div>
		</div>
		<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onScannerPopup}>
			<div>Scanner Range</div>
			<div>{planet.spec.scanRange ?? '--'} l.y.</div>
		</div>

		{#if $player.race.spec?.canBuildDefenses}
			<div class="divider p-0 m-0" />

			<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onDefensePoopup}>
				<div>Defenses</div>
				<div>{planet.defenses} of {planet.spec.maxDefenses}</div>
			</div>
			<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onDefensePoopup}>
				<div>Defense Type</div>
				<div>{planet.spec.defense}</div>
			</div>
			<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onDefensePoopup}>
				<div>Defense Coverage</div>
				<div>
					{((planet.spec.defenseCoverage ?? 0) * 100).toFixed(1)}%
				</div>
			</div>
		{/if}
	</CommandTile>
{/if}
