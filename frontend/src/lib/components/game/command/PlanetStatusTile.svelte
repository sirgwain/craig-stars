<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';
</script>

{#if $commandedPlanet && $commandedPlanet.spec && $commandedPlanet.cargo}
	<CommandTile title="Planet Status">
		<div class="flex justify-between">
			<div>Population</div>
			<div>{$commandedPlanet.cargo.colonists ? $commandedPlanet.cargo.colonists * 100 : 0}</div>
		</div>
		<div class="flex justify-between">
			<div>Resources/Year</div>
			<div>
				{$commandedPlanet.spec.resourcesPerYearAvailable} of {$commandedPlanet.spec
					?.resourcesPerYear}
			</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Scanner Type</div>
			<div>{$commandedPlanet.spec.scanner}</div>
		</div>
		<div class="flex justify-between">
			<div>Scanner Range</div>
			<div>{$commandedPlanet.spec.scanRange} l.y.</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Defenses</div>
			<div>{$commandedPlanet.defenses} of {$commandedPlanet.spec.maxDefenses}</div>
		</div>
		<div class="flex justify-between">
			<div>Defense Type</div>
			<div>{$commandedPlanet.spec.defense}</div>
		</div>
		<div class="flex justify-between">
			<div>Defense Coverage</div>
			<div>
				{($commandedPlanet.spec.defenseCoverage * 100).toFixed(1)}%
			</div>
		</div>
	</CommandTile>
{/if}
