<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import Starbase from '$lib/components/icons/Starbase.svelte';
	import { MapObjectType } from '$lib/types/cs-proto';
	import type {
		MysteryTraderIntel,
		PlanetIntel,
		SalvageIntel,
		WormholeIntel
	} from '$lib/types/cs-proto';
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet, AnyMinefield, AnyMineralPacket } from '$lib/services/Universe';
	import { getMapObjectName } from '$lib/types/MapObject';
	import FleetSummary from './FleetSummary.svelte';
	import MinefieldSummary from './MinefieldSummary.svelte';
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
		if (selectedPlanet?.spec?.planetStarbaseSpec?.starbaseDesignNum) {
			onShipDesignTooltip(
				e,
				$universe.getDesign(
					selectedPlanet.mapObject?.playerNum,
					selectedPlanet.spec.planetStarbaseSpec.starbaseDesignNum
				)
			);
		}
	}

	let selectedPlanet = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.PLANET
			? ($selectedMapObject as PlanetIntel)
			: undefined
	);
	let selectedFleet = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.FLEET
			? ($selectedMapObject as AnyFleet)
			: undefined
	);
	let selectedMinefield = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.MINEFIELD
			? ($selectedMapObject as AnyMinefield)
			: undefined
	);
	let selectedMineralPacket = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.MINERAL_PACKET
			? ($selectedMapObject as AnyMineralPacket)
			: undefined
	);
	let selectedSalvage = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.SALVAGE
			? ($selectedMapObject as SalvageIntel)
			: undefined
	);
	let selectedWormhole = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.WORMHOLE
			? ($selectedMapObject as WormholeIntel)
			: undefined
	);
	let selectedMysteryTrader = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.MYSTERY_TRADER
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
				{#if selectedPlanet && selectedPlanet.spec?.planetStarbaseSpec?.hasStarbase}
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
		{:else if selectedMinefield}
			<MinefieldSummary minefield={selectedMinefield} />
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
