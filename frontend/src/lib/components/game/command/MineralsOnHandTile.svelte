<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';
</script>

<CommandTile title="Minerals on Hand">
	<div class="flex justify-between">
		<div class="text-ironium">Ironium</div>
		<div>{$commandedPlanet.cargo.ironium}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="text-boranium">Boranium</div>
		<div>{$commandedPlanet.cargo.boranium}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="text-germanium">Germanium</div>
		<div>{$commandedPlanet.cargo.germanium}kT</div>
	</div>

	<div class="divider p-0 m-0" />

	<div class="flex justify-between">
		<div>Mines</div>
		<div>{$commandedPlanet.mines} of {$commandedPlanet.spec?.maxMines}</div>
	</div>
	<div class="flex justify-between">
		<div>Factories</div>
		<div>{$commandedPlanet.factories} of {$commandedPlanet.spec?.maxFactories}</div>
	</div>
</CommandTile>
