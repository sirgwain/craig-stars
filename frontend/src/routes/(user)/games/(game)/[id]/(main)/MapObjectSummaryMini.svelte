<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import PlanetHabValue from '$lib/components/game/PlanetHabValue.svelte';
	import Starbase from '$lib/components/icons/Starbase.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { population } from '$lib/types/Cargo';
	import { type MapObject } from '$lib/types/cs';
	import { getUnderlyingMapObject, ownedBy } from '$lib/types/MapObject';
	import FleetSummary from './FleetSummary.svelte';
	import MapObjectIcon from './MapObjectIcon.svelte';
	import MinefieldSummary from './MinefieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import MysteryTraderSummary from './MysteryTraderSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import SalvageSummary from './SalvageSummary.svelte';
	import WormholeSummary from './WormholeSummary.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		mapObject: MapObject | undefined;
	};

	let { mapObject }: Props = $props();

	let { planet, fleet, wormhole, minefield, mysteryTrader, salvage, mineralPacket } = $derived(
		getUnderlyingMapObject(mapObject)
	);

	let playerFleet = $derived($universe.getMyFleet(fleet?.num));
</script>

<div class="flex flex-row justify-start gap-3 text-sm">
	{#if planet}
		{#if ownedBy(planet, $player.num)}
			<MapObjectIcon {mapObject} />
			<div class="flex flex-col w-full">
				<div class="flex flex-row justify-between">
					<div>
						Value: <PlanetHabValue {planet} />
					</div>
					<div>
						{#if population(planet.cargo)}
							Pop: {population(planet.cargo).toLocaleString()}
						{/if}
					</div>
				</div>
				<div class="mt-1">
					{#if planet.spec?.hasStarbase}
						<Starbase class="w-4 h-4 starbase" />
					{/if}
				</div>
				<div class="mt-1">
					<div class="flex justify-between">
						<div class="text-tile-item-title text-ironium">Ironium</div>
						<div>
							{(planet.cargo?.ironium ?? 0).toLocaleString()}kT ({planet.mineralConcentration
								?.ironium})
						</div>
					</div>
					<div class="flex justify-between">
						<div class="text-tile-item-title text-boranium">Boranium</div>
						<div>
							{(planet.cargo?.boranium ?? 0).toLocaleString()}kT ({planet.mineralConcentration
								?.boranium})
						</div>
					</div>
					<div class="flex justify-between">
						<div class="text-tile-item-title text-germanium">Germanium</div>
						<div>
							{(planet.cargo?.germanium ?? 0).toLocaleString()}kT ({planet.mineralConcentration
								?.germanium})
						</div>
					</div>
				</div>
			</div>
		{:else}
			<PlanetSummary {planet} />
		{/if}
	{:else if fleet}
		{#if playerFleet}
			<MapObjectIcon {mapObject} />
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
						{fleet.spec?.mass ?? 0}kT
					</div>
				</div>
				<div class="flex flex-row">
					<div class="w-32 text-tile-item-title">Fuel:</div>
					<div class="grow">
						<FuelBar value={playerFleet.fuel} capacity={playerFleet.spec?.fuelCapacity ?? 0} />
					</div>
				</div>
				<div class="flex flex-row">
					<div class="w-32 text-tile-item-title">Cargo:</div>
					<div class="grow">
						<CargoBar value={playerFleet.cargo} capacity={playerFleet.spec?.cargoCapacity} />
					</div>
				</div>
			</div>
		{:else}
			<FleetSummary {fleet} />
		{/if}
	{:else if mineralPacket}
		<MineralPacketSummary {mineralPacket} />
	{:else if minefield}
		<MinefieldSummary {minefield} />
	{:else if salvage}
		<SalvageSummary {salvage} />
	{:else if wormhole}
		<WormholeSummary {wormhole} />
	{:else if mysteryTrader}
		<MysteryTraderSummary {mysteryTrader} />
	{/if}
</div>
