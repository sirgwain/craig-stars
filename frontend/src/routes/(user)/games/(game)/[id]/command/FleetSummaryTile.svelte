<script lang="ts">
	import { commandedFleet, commandMapObject, player, zoomToMapObject } from '$lib/services/Context';
	import { rollover } from '$lib/services/Math';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { kebabCase } from 'lodash-es';
	import CommandTile from './CommandTile.svelte';

	let design: ShipDesign | undefined;
	let icon = '';

	$: {
		// console.log('loading icon of $commandedFleet: ', $commandedFleet);
		icon = '';
		if ($player && $commandedFleet && $commandedFleet?.tokens?.length > 0) {
			const designId = $commandedFleet.tokens[0].designId;
			design = $player.designs.find((d) => d.id == designId);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber}`;
			}
		}
	}

	$: index =
		$player && $commandedFleet ? $player.fleets.findIndex((f) => f === $commandedFleet) : 0;

	const previous = () => {
		if ($player) {
			index = rollover(index - 1, 0, $player.fleets.length - 1);
			commandMapObject($player.fleets[index]);
			zoomToMapObject($player.fleets[index]);
		}
	};
	const next = () => {
		if ($player) {
			index = rollover(index + 1, 0, $player.fleets.length - 1);
			commandMapObject($player.fleets[index]);
			zoomToMapObject($player.fleets[index]);
		}
	};

</script>

{#if $commandedFleet}
	<CommandTile title={$commandedFleet.name}>
		<div class="grid grid-cols-2">
			<div class="avatar ">
				<div class="border-2 border-neutral mr-2 p-2 bg-black">
					<div class="fleet-avatar {icon} bg-black" />
				</div>
			</div>

			<div class="flex flex-col gap-y-1">
				<button on:click={previous} class="btn btn-outline btn-sm normal-case btn-secondary"
					>Prev</button
				>
				<button on:click={next} class="btn btn-outline btn-sm normal-case btn-secondary"
					>Next</button
				>
			</div>
		</div>
	</CommandTile>
{/if}
