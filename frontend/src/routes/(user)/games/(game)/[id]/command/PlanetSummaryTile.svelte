<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
import type { Planet } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	const previous = () => {};
	const next = () => {};
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
