<script lang="ts">
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { None } from '$lib/types/cs';
	import { canTransferCargo, type CommandedFleet } from '$lib/types/Fleet';
	import { ownedBy } from '$lib/types/MapObject';
	import CommandTile from './CommandTile.svelte';

	const { player, universe, commandMapObject } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
	} & ShowCargoTransferDialogProps;

	let { fleet, onShowCargoTransferDialog }: Props = $props();

	let planet = $derived(
		fleet.orbitingPlanetNum != None && $universe.getPlanet(fleet.orbitingPlanetNum)
	);
	const transfer = () => {
		if (!onShowCargoTransferDialog) {
			return;
		}
		onShowCargoTransferDialog({
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
				onclick={gotoTarget}
				disabled={!planet || !ownedBy(planet, $player.num)}
				class="btn btn-outline btn-sm normal-case btn-secondary p-2"
				title="goto">Goto</button
			>
			<button
				onclick={transfer}
				class="btn btn-outline btn-sm normal-case btn-secondary p-2"
				title="transfer"
				disabled={!canTransferCargo(fleet)}
				>{planet ? 'Transfer' : 'Jettison'}
			</button>
		</div>
	</CommandTile>
{/if}
