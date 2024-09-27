<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { None, ownedBy } from '$lib/types/MapObject';
	import { createEventDispatcher } from 'svelte';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTranfserDialog.svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher<CargoTransferDialogEvent>();
	const { player, universe, commandMapObject } = getGameContext();

	export let fleet: CommandedFleet;

	$: planet = fleet.orbitingPlanetNum != None && $universe.getPlanet(fleet.orbitingPlanetNum);
	const transfer = () => {
		dispatch('cargo-transfer-dialog', {
			src: fleet,
			dest: planet ? planet : fleet.getCargoTransferTarget($universe)
		});
	};
	const gotoTarget = () => {
		if (planet && ownedBy(planet, $player.num)) {
			commandMapObject(planet);
		}
	};
</script>

{#if fleet}
	<CommandTile title={planet ? `Orbiting ${planet.name}` : 'In Deep Space'}>
		<div class="flex justify-between my-1 btn-group">
			<button
				on:click={gotoTarget}
				disabled={!planet || !ownedBy(planet, $player.num)}
				class="btn btn-outline btn-sm normal-case btn-secondary p-2"
				title="goto">Goto</button
			>
			<button
				on:click={transfer}
				class="btn btn-outline btn-sm normal-case btn-secondary p-2"
				title="goto"
				>{planet ? 'Transfer' : 'Jettison'}
			</button>
		</div>
	</CommandTile>
{/if}
