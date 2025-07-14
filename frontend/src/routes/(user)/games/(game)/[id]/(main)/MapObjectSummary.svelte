<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import Starbase from '$lib/components/icons/Starbase.svelte';
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet, AnyMineField, AnyMineralPacket } from '$lib/services/Universe';
	import type { MysteryTraderIntel, PlanetIntel, SalvageIntel, WormholeIntel } from '$lib/types/cs';
	import {
		MapObjectTypeFleet,
		MapObjectTypeMineField,
		MapObjectTypeMineralPacket,
		MapObjectTypeMysteryTrader,
		MapObjectTypePlanet,
		MapObjectTypeSalvage,
		MapObjectTypeWormhole
	} from '$lib/types/cs';
	import { getMapObjectName } from '$lib/types/MapObject';
	import FleetSummary from './FleetSummary.svelte';
	import MineFieldSummary from './MineFieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import MysteryTraderSummary from './MysteryTraderSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import SalvageSummary from './SalvageSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';
	import WormholeSummary from './WormholeSummary.svelte';

	const { universe, selectNextMapObject, selectedMapObject } = getGameContext();

	type Props = {
		hideCycleButton?: boolean;
		hideTitle?: boolean;
	} & ShowCargoTransferDialogProps;

	let { onShowCargoTransferDialog, hideTitle, hideCycleButton }: Props = $props();

	function showStarbaseDesign(e: MouseEvent) {
		if (selectedPlanet?.spec?.starbaseDesignNum) {
			onShipDesignTooltip(
				e,
				$universe.getDesign(selectedPlanet.playerNum, selectedPlanet.spec.starbaseDesignNum)
			);
		}
	}

	let selectedPlanet = $derived(
		$selectedMapObject?.type == MapObjectTypePlanet
			? ($selectedMapObject as PlanetIntel)
			: undefined
	);
	let selectedFleet = $derived(
		$selectedMapObject?.type == MapObjectTypeFleet ? ($selectedMapObject as AnyFleet) : undefined
	);
	let selectedMineField = $derived(
		$selectedMapObject?.type == MapObjectTypeMineField
			? ($selectedMapObject as AnyMineField)
			: undefined
	);
	let selectedMineralPacket = $derived(
		$selectedMapObject?.type == MapObjectTypeMineralPacket
			? ($selectedMapObject as AnyMineralPacket)
			: undefined
	);
	let selectedSalvage = $derived(
		$selectedMapObject?.type == MapObjectTypeSalvage
			? ($selectedMapObject as SalvageIntel)
			: undefined
	);
	let selectedWormhole = $derived(
		$selectedMapObject?.type == MapObjectTypeWormhole
			? ($selectedMapObject as WormholeIntel)
			: undefined
	);
	let selectedMysteryTrader = $derived(
		$selectedMapObject?.type == MapObjectTypeMysteryTrader
			? ($selectedMapObject as MysteryTraderIntel)
			: undefined
	);
</script>

<div
	class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full"
	data-type="map-object-summary"
	data-id={getMapObjectName($selectedMapObject)}
>
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			{#if !hideTitle}
				<div class="flex-1 text-center text-lg font-semibold text-secondary">
					{getMapObjectName($selectedMapObject)}
				</div>
			{/if}
			<div>
				{#if selectedPlanet && selectedPlanet.spec?.hasStarbase}
					<button
						type="button"
						onpointerdown={(e) => {
							e.preventDefault();
							showStarbaseDesign(e);
						}}
					>
						<Starbase class="w-4 h-4 starbase" /></button
					>
				{/if}
				{#if !hideCycleButton}
					<button
						type="button"
						data-type="cycle-selected-map-object-button"
						onpointerdown={(e) => {
							e.preventDefault();
							selectNextMapObject();
						}}
					>
						<Cycle class="w-4 h-4 fill-base-content hover:stroke-accent" /></button
					>
				{/if}
			</div>
		</div>
		{#if selectedPlanet}
			<PlanetSummary planet={selectedPlanet} />
		{:else if selectedFleet}
			<FleetSummary fleet={selectedFleet} {onShowCargoTransferDialog} />
		{:else if selectedMineField}
			<MineFieldSummary mineField={selectedMineField} />
		{:else if selectedMineralPacket}
			<MineralPacketSummary mineralPacket={selectedMineralPacket} />
		{:else if selectedSalvage}
			<SalvageSummary salvage={selectedSalvage} />
		{:else if selectedWormhole}
			<WormholeSummary wormhole={selectedWormhole} />
		{:else if selectedMysteryTrader}
			<MysteryTraderSummary mysteryTrader={selectedMysteryTrader} />
		{:else}
			<UnknownSummary />
		{/if}
	</div>
</div>
