<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import type { Fleet } from '$lib/types/Fleet';
	import { ownedBy } from '$lib/types/MapObject';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { kebabCase, startCase } from 'lodash-es';

	const { player, universe } = getGameContext();

	export let fleet: Fleet;

	let design: ShipDesign | undefined;
	let icon = '';

	$: {
		icon = '';
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			design = $universe.getDesign(fleet.playerNum, designNum);
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
						on:pointerdown|preventDefault={(e) => onShipDesignTooltip(e, design)}
					/>
				</div>
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerName(fleet.playerNum)}</div>
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
		{#if ownedBy(fleet, $player.num)}
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
