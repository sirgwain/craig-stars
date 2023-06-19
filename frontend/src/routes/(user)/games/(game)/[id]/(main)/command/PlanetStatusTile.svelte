<script lang="ts">
	import type { CommandedPlanet } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	export let planet: CommandedPlanet;
</script>

{#if planet.spec && planet.cargo}
	<CommandTile title="Planet Status">
		<div class="flex justify-between">
			<div>Population</div>
			<div>{(planet.cargo.colonists ?? 0) * 100}</div>
		</div>
		<div class="flex justify-between">
			<div>Resources/Year</div>
			<div>
				{planet.spec.resourcesPerYearAvailable} of {planet.spec.resourcesPerYear}
			</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Scanner Type</div>
			<div>{planet.spec.scanner}</div>
		</div>
		<div class="flex justify-between">
			<div>Scanner Range</div>
			<div>{planet.spec.scanRange} l.y.</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Defenses</div>
			<div>{planet.defenses} of {planet.spec.maxDefenses}</div>
		</div>
		<div class="flex justify-between">
			<div>Defense Type</div>
			<div>{planet.spec.defense}</div>
		</div>
		<div class="flex justify-between">
			<div>Defense Coverage</div>
			<div>
				{(planet.spec.defenseCoverage * 100).toFixed(1)}%
			</div>
		</div>
	</CommandTile>
{/if}
