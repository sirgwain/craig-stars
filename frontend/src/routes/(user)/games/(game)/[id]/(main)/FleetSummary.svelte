<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { getHullIcon } from '$lib/techicon';
	import { StargateWarpSpeed, WaypointTaskNone, type Fleet, type FleetIntel } from '$lib/types/cs';
	import { canTransferCargo, CommandedFleet, getDamagePercentForToken } from '$lib/types/Fleet';
	import { ownedBy } from '$lib/types/MapObject';
	import type { ShipDesign } from '$lib/types/cs';
	import { startCase } from 'lodash-es';

	const { player, universe } = getGameContext();

	type Props = {
		fleet: Fleet;
	} & ShowCargoTransferDialogProps;

	let { fleet, onShowCargoTransferDialog }: Props = $props();

	const design: ShipDesign | undefined = $derived.by(() => {
		if (fleet.tokens && fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.playerNum, designNum);
		}
	});

	// get either warpSpeed as a number, or "stargate"
	function getWarpSpeed(fleet: Fleet): string {
		const warpSpeed: number =
			(fleet?.waypoints && fleet.waypoints.length > 1
				? fleet.waypoints[1].warpSpeed
				: fleet.warpSpeed) ?? 0;

		if (warpSpeed == StargateWarpSpeed) {
			return 'Use Stargate';
		}
		return `${warpSpeed}`;
	}

	function getMass(fleet: Fleet | FleetIntel): number {
		if ('spec' in fleet) {
			return fleet.spec?.mass ?? 0;
		}
		return fleet.mass ?? 0;
	}

	function transfer() {
		if (!onShowCargoTransferDialog) {
			return;
		}
		const f = new CommandedFleet(fleet);
		onShowCargoTransferDialog({ src: f, dest: f.getCargoTransferTarget($universe) });
	}
</script>

<div class="flex flex-row min-h-[11rem]">
	<div class="flex flex-col">
		<div class="avatar mr-2">
			<div
				class="border-2 border-neutral p-2 bg-black"
				style={`border-color: ${$universe.getPlayerColor(fleet.playerNum)};`}
			>
				{#if fleet.tokens && fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
					<div class="absolute -right-2 -top-1 text-xl w-6 h-6">+</div>
				{/if}

				<div class="fleet-avatar {getHullIcon(design)} bg-black">
					<button
						type="button"
						aria-label="Opens ship design tooltip"
						class="w-full h-full cursor-help"
						onpointerdown={(e) => onShipDesignTooltip(e, design)}
					></button>
				</div>
			</div>
		</div>
		<div class="text-center">{$universe.getPlayerPluralName(fleet.playerNum)}</div>
	</div>
	<div class="flex flex-col grow">
		<div class="flex flex-row">
			<div class="w-32 text-tile-item-title">Ship Count:</div>
			<div>
				{fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 'unknown'}
			</div>
		</div>
		<div class="flex flex-row">
			<div class="w-32 text-tile-item-title">Fleet Mass:</div>
			<div>
				{getMass(fleet)}kT
			</div>
		</div>
		{#if ownedBy(fleet, $player.num)}
			<div class="flex flex-row">
				<div class="w-32 text-tile-item-title">Fuel:</div>
				<div class="grow">
					<FuelBar value={fleet.fuel ?? 0} capacity={fleet.spec?.fuelCapacity ?? 0} />
				</div>
			</div>
			<div class="flex flex-row">
				<div class="w-32 text-tile-item-title">Cargo:</div>
				<div class="grow">
					<CargoBar
						onPointerDown={() => transfer()}
						canTransferCargo={canTransferCargo(fleet, $universe)}
						value={fleet.cargo}
						capacity={fleet.spec?.cargoCapacity}
					/>
				</div>
			</div>
		{/if}
		{#if fleet.waypoints && fleet.waypoints.length > 1}
			<div class="flex flex-row">
				<div class="w-32 text-tile-item-title">Next Waypoint:</div>
				<div>{$universe.getTargetName(fleet.waypoints[1])}</div>
			</div>
			{#if fleet.waypoints[1].task !== WaypointTaskNone}
				<div class="flex flex-row">
					<div class="w-32 text-tile-item-title">Task:</div>
					<div>{startCase(fleet.waypoints[1].task)}</div>
				</div>
			{/if}
			<div class="flex flex-row">
				<div class="w-32 text-tile-item-title">Warp Speed:</div>
				<div>{getWarpSpeed(fleet)}</div>
			</div>
		{:else if fleet.warpSpeed}
			<div class="flex flex-row">
				<div class="w-32 text-tile-item-title">Warp Speed:</div>
				<div>{getWarpSpeed(fleet)}</div>
			</div>
		{/if}

		{#if !ownedBy(fleet, $player.num) && fleet.tokens}
			<div class="text-tile-item-title">
				Fleet Composition:
				<div class="bg-base-100 h-16 overflow-y-auto mt-1 w-full md:w-60 font-normal">
					<ul class="w-full h-full">
						{#each fleet.tokens as token}
							<li class="pl-1">
								<button
									type="button"
									class="w-full cursor-help"
									onpointerdown={(e) =>
										onShipDesignTooltip(e, $universe.getDesign(fleet.playerNum, token.designNum))}
								>
									<span class="flex flex-row justify-between relative">
										{#if (token.damage ?? 0) > 0 && (token.quantityDamaged ?? 0) > 0}
											<div
												style={`width: ${getDamagePercentForToken(
													token,
													$universe.getDesign(fleet.playerNum, token.designNum)
												).toFixed()}%`}
												class="damage-bar h-full absolute opacity-50"
											></div>
										{/if}

										<span>
											{$universe.getDesign(fleet.playerNum, token.designNum)?.name}
										</span>
										<span>
											{token.quantity}
										</span>
									</span>
								</button>
							</li>
						{/each}
					</ul>
				</div>
			</div>
		{/if}
	</div>
</div>
