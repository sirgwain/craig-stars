<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { playerName } from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import type { Player } from '$lib/types/Player';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { kebabCase, startCase } from 'lodash-es';

	export let fleet: Fleet;
	export let player: Player;

	let design: ShipDesign | undefined;
	let icon = '';

	$: {
		// console.log('loading icon of fleet: ', fleet);
		icon = '';
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designUuid = fleet.tokens[0].designUuid;
			design = player.designs.find((d) => d.uuid == designUuid);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber}`;
			}
		}
	}
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="fleet-avatar {icon} bg-black" />
			</div>
		</div>
		<div class="text-center">{playerName(fleet.playerNum ?? 0)}</div>
	</div>
	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-24">Ship Count:</div>
			<div>
				{fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 'unknown'}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-24">Fleet Mass:</div>
			<div>
				{fleet.spec?.mass ?? 0}kT
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-24">Fuel:</div>
			<div class="grow">
				<FuelBar value={fleet.fuel ?? 0} capacity={fleet.spec?.fuelCapacity ?? 0} />
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-24">Cargo:</div>
			<div class="grow">
				<CargoBar value={fleet.cargo} capacity={fleet.spec?.cargoCapacity ?? 0} />
			</div>
		</div>
		{#if fleet.waypoints && fleet.waypoints.length > 1}
			<div class="flex flex-row">
				<div class="w-24">Next Waypoint:</div>
				<div>{fleet.waypoints[1].targetName}</div>
			</div>
			<div class="flex flex-row">
				<div class="w-24">Task:</div>
				<div>{startCase(fleet.waypoints[1].task)}</div>
			</div>
		{/if}
		<div class="flex flex-row">
			<div class="w-24">Warp Speed:</div>
			<div>{fleet.warpSpeed ?? 0}</div>
		</div>
	</div>
</div>
