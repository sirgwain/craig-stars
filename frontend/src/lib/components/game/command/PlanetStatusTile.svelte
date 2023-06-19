<script lang="ts">
	import { getCommandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';

	const planet = getCommandedPlanet();
</script>

<CommandTile title="Planet Status">
	<div class="flex justify-between">
		<div>Population</div>
		<div>{planet.cargo.colonists ? planet.cargo.colonists * 100 : 0}</div>
	</div>
	<div class="flex justify-between">
		<div>Resources/Year</div>
		<div>30 of 35</div>
	</div>

	<div class="divider p-0 m-0" />

	<div class="flex justify-between">
		<div>Scanner Type</div>
		<div>{planet.spec?.scanner}</div>
	</div>
	<div class="flex justify-between">
		<div>Scanner Range</div>
		<div>{planet.spec?.scanRange} l.y.</div>
	</div>

	<div class="divider p-0 m-0" />

	<div class="flex justify-between">
		<div>Defenses</div>
		<div>{planet.defenses}</div>
	</div>
	<div class="flex justify-between">
		<div>Defense Type</div>
		<div>{planet.spec?.defense}</div>
	</div>
	<div class="flex justify-between">
		<div>Defense Coverage</div>
		<div>{planet.spec?.defenseCoverage}%</div>
	</div>
</CommandTile>
