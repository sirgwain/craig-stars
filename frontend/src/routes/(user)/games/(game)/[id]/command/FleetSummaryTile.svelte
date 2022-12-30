<script lang="ts">
	import { nextMapObject, previousMapObject } from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import type { Player } from '$lib/types/Player';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { kebabCase } from 'lodash-es';
	import CommandTile from './CommandTile.svelte';

	export let player: Player;
	export let fleet: Fleet;

	let design: ShipDesign | undefined;
	let icon = '';

	$: {
		// console.log('loading icon of fleet: ', fleet);
		icon = '';
		if (player && fleet && fleet?.tokens?.length > 0) {
			const designUuid = fleet.tokens[0].designUuid;
			design = player.designs.find((d) => d.uuid == designUuid);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber}`;
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
