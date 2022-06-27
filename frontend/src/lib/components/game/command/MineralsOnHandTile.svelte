<script lang="ts">
	import { getCommandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';

	const planet = getCommandedPlanet();
</script>

<CommandTile title="Minerals on Hand">
	<div class="flex justify-between">
		<div class="text-ironium">Ironium</div>
		<div>{planet.cargo.ironium}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="text-boranium">Boranium</div>
		<div>{planet.cargo.boranium}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="text-germanium">Germanium</div>
		<div>{planet.cargo.germanium}kT</div>
	</div>

	<div class="divider p-0 m-0" />

	<div class="flex justify-between">
		<div>Mines</div>
		<div>{planet.mines}</div>
	</div>
	<div class="flex justify-between">
		<div>Factories</div>
		<div>{planet.factories}</div>
	</div>
</CommandTile>
