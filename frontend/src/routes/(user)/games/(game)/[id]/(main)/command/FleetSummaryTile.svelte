<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { nextMapObject, previousMapObject } from '$lib/services/Stores';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { kebabCase } from 'lodash-es';
	import CommandTile from './CommandTile.svelte';

	const { player, universe } = getGameContext();

	export let fleet: CommandedFleet;

	let icon = '';

	$: {
		// console.log('loading icon of fleet: ', fleet);
		icon = '';
		if (fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			const design = $universe.getDesign($player.num, designNum);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
			}
		}
	}
</script>

<CommandTile title={fleet.name}>
	<div class="grid grid-cols-2">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="fleet-avatar {icon} bg-black" />
			</div>
		</div>

		<div class="flex flex-col gap-y-1">
			<button
				on:click={() => previousMapObject()}
				class="btn btn-outline btn-sm normal-case btn-secondary">Prev</button
			>
			<button
				on:click={() => nextMapObject()}
				class="btn btn-outline btn-sm normal-case btn-secondary">Next</button
			>
		</div>
	</div>
</CommandTile>
