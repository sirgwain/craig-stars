<script lang="ts">
	import type { NextPrevMapObjectProps } from '$lib/services/Events';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	type Props = {
		planet: CommandedPlanet;
	} & NextPrevMapObjectProps;

	let { planet, onNextMapObject, onPreviousMapObject }: Props = $props();
	const icon = (planet: CommandedPlanet) => (planet ? `planet-${(planet.num - 1) % 26}` : '');
</script>

<CommandTile title={planet.name}>
	<div class="grid grid-cols-2">
		<div class="avatar">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="planet-avatar {icon(planet)} bg-black"></div>
			</div>
		</div>

		<div class="flex flex-col gap-y-1">
			<button onclick={onPreviousMapObject} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Prev</button
			>
			<button onclick={onNextMapObject} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Next</button
			>
		</div>
	</div>
</CommandTile>
