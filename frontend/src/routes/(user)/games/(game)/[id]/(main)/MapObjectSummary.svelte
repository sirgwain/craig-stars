<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import Starbase from '$lib/components/icons/Starbase.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Planet } from '$lib/types/Planet';
	import type { Salvage } from '$lib/types/Salvage';
	import type { Wormhole } from '$lib/types/Wormhole';
	import FleetSummary from './FleetSummary.svelte';
	import MineFieldSummary from './MineFieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import SalvageSummary from './SalvageSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';
	import WormholeSummary from './WormholeSummary.svelte';

	const { universe, selectNextMapObject, selectedMapObject } = getGameContext();

	function showStarbaseDesign(e: MouseEvent) {
		if (selectedPlanet?.spec.starbaseDesignNum) {
			onShipDesignTooltip(
				e,
				$universe.getDesign(selectedPlanet.playerNum, selectedPlanet.spec.starbaseDesignNum)
			);
		}
	}

	$: selectedPlanet =
		$selectedMapObject?.type == MapObjectType.Planet ? ($selectedMapObject as Planet) : undefined;
	$: selectedFleet =
		$selectedMapObject?.type == MapObjectType.Fleet ? ($selectedMapObject as Fleet) : undefined;
	$: selectedMineField =
		$selectedMapObject?.type == MapObjectType.MineField
			? ($selectedMapObject as MineField)
			: undefined;
	$: selectedMineralPacket =
		$selectedMapObject?.type == MapObjectType.MineralPacket
			? ($selectedMapObject as MineralPacket)
			: undefined;
	$: selectedSalvage =
		$selectedMapObject?.type == MapObjectType.Salvage ? ($selectedMapObject as Salvage) : undefined;
	$: selectedWormhole =
		$selectedMapObject?.type == MapObjectType.Wormhole
			? ($selectedMapObject as Wormhole)
			: undefined;
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full">
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				{$selectedMapObject?.name ?? ''}
			</div>
			<div>
				{#if selectedPlanet && selectedPlanet.spec.hasStarbase}
					<button type="button" on:pointerdown|preventDefault={showStarbaseDesign}>
						<Starbase class="w-4 h-4 starbase" /></button
					>
				{/if}
				<button type="button" on:pointerdown|preventDefault={selectNextMapObject}>
					<Cycle class="w-4 h-4 fill-base-content hover:stroke-accent" /></button
				>
			</div>
		</div>
		{#if selectedPlanet}
			<PlanetSummary planet={selectedPlanet} />
		{:else if selectedFleet}
			<FleetSummary fleet={selectedFleet} />
		{:else if selectedMineField}
			<MineFieldSummary mineField={selectedMineField} />
		{:else if selectedMineralPacket}
			<MineralPacketSummary mineralPacket={selectedMineralPacket} />
		{:else if selectedSalvage}
			<SalvageSummary salvage={selectedSalvage} />
		{:else if selectedWormhole}
			<WormholeSummary wormhole={selectedWormhole} />
		{:else}
			<UnknownSummary />
		{/if}
	</div>
</div>
