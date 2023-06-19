<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
	import CommandTile from './CommandTile.svelte';

	const previous = () => {};
	const next = () => {};
</script>

<CommandTile title={$commandedPlanet.name}>
	<div class="grid grid-cols-2">
		<div class="avatar">
			<div class="border border-neutral">
				<img alt="Planet icon" src="../../planet/planet01.png" />
			</div>
		</div>
		<div class="flex flex-col gap-y-1 ml-1">
			<button on:click={previous} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Prev</button
			>
			<button on:click={next} class="btn btn-outline btn-sm normal-case btn-secondary">Next</button>
		</div>
	</div>
</CommandTile>
