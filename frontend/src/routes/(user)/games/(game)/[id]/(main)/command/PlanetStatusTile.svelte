<script lang="ts">
	import type { MinesTooltipProps } from '$lib/components/game/tooltips/MinesTooltip.svelte';
	import ResourcesTooltip, {
		type ResourcesTooltipProps
	} from '$lib/components/game/tooltips/ResourcesTooltip.svelte';
	import TechTooltip, {
		type TechTooltipProps
	} from '$lib/components/game/tooltips/TechTooltip.svelte';
	import { showTooltip, techs } from '$lib/services/Context';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import CommandTile from './CommandTile.svelte';

	export let player: Player;
	export let planet: CommandedPlanet;

	function onResourcesTooltip(e: PointerEvent) {
		showTooltip<ResourcesTooltipProps>(e.x, e.y, ResourcesTooltip, {
			planetName: planet.name,
			resourcesPerYear: planet.spec.resourcesPerYear,
			resourcesPerYearAvailable: planet.spec.resourcesPerYearAvailable,
			resourcesPerYearResearch: planet.spec.resourcesPerYearAvailable,
			innateResources: player.race.spec?.innateResources ?? false
		});
	}

	function onScannerPopup(e: PointerEvent) {
		const tech = $techs.getTech(planet.spec.scanner);
		if (tech) {
			showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, { tech });
		}
	}
	function onDefensePoopup(e: PointerEvent) {
		const tech = $techs.getTech(planet.spec.defense);
		if (tech) {
			showTooltip<TechTooltipProps>(e.x, e.y, TechTooltip, { tech });
		}
	}
</script>

{#if planet.spec && planet.cargo}
	<CommandTile title="Planet Status">
		<div class="flex justify-between">
			<div>Population</div>
			<div>{(planet.cargo.colonists ?? 0) * 100}</div>
		</div>
		<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onResourcesTooltip}>
			<div>Resources/Year</div>
			<div>
				{planet.spec.resourcesPerYearAvailable} of {planet.spec.resourcesPerYear}
			</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onScannerPopup}>
			<div>Scanner Type</div>
			<div>{planet.spec.scanner}</div>
		</div>
		<div class="flex justify-between cursor-help" on:pointerdown|preventDefault={onScannerPopup}>
			<div>Scanner Range</div>
			<div>{planet.spec.scanRange} l.y.</div>
		</div>

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
				{(planet.spec.defenseCoverage * 100).toFixed(1)}%
			</div>
		</div>
	</CommandTile>
{/if}
