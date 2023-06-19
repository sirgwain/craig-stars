<script lang="ts">
	import {
		player,
		commandedPlanet,
		commandMapObject,
		zoomToMapObject
	} from '$lib/services/Context';
	import { rollover } from '$lib/services/Math';
	import type { Planet } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	$: index =
		$player && $commandedPlanet ? $player.planets.findIndex((p) => p === $commandedPlanet) : 0;

	const previous = () => {
		if ($player) {
			index = rollover(index - 1, 0, $player.planets.length - 1);
			commandMapObject($player.planets[index]);
			zoomToMapObject($player.planets[index]);
		}
	};
	const next = () => {
		if ($player) {
			index = rollover(index + 1, 0, $player.planets.length - 1);
			commandMapObject($player.planets[index]);
			zoomToMapObject($player.planets[index]);
		}
	};
	const icon = (planet: Planet) => (planet ? `planet-${(planet.num - 1) % 26}` : '');
</script>

{#if $commandedPlanet}
	<CommandTile title={$commandedPlanet.name}>
		<div class="grid grid-cols-2">
			<div class="avatar ">
				<div class="border-2 border-neutral mr-2 p-2 bg-black">
					<div class="planet-avatar {icon($commandedPlanet)} bg-black" />
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
