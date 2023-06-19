<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { showDesignPopup } from '$lib/services/Context';
	import type { FullGame } from '$lib/services/FullGame';
	import type { Fleet } from '$lib/types/Fleet';
	import { ownedBy } from '$lib/types/MapObject';
	import type { Player } from '$lib/types/Player';
	import type { ShipDesign, ShipDesignIntel } from '$lib/types/ShipDesign';
	import { kebabCase, startCase } from 'lodash-es';

	export let game: FullGame;
	export let player: Player;
	export let fleet: Fleet;

	let design: ShipDesign | ShipDesignIntel | undefined;
	let icon = '';

	$: {
		// console.log('loading icon of fleet: ', fleet);
		icon = '';
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			design = player.getDesign(fleet.playerNum, designNum);
			if (design) {
				icon = `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
			}
		}
	}
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar ">
			<div class="border-2 border-neutral mr-2 p-2 bg-black">
				<div class="fleet-avatar {icon} bg-black">
					<button
						type="button"
						class="w-full h-full cursor-help"
						on:pointerdown={(e) => showDesignPopup(design, e.x, e.y)}
					/>
				</div>
			</div>
		</div>
		<div class="text-center">{game.getPlayerName(fleet.playerNum)}</div>
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
				{fleet.spec?.mass ?? fleet.mass ?? 0}kT
			</div>
		</div>
		{#if ownedBy(fleet, player.num)}
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
		{/if}
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
