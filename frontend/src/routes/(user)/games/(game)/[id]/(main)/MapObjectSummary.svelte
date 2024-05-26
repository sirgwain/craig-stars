<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import Starbase from '$lib/components/icons/Starbase.svelte';
	import { getCarouselContext } from '$lib/services/CarouselContext';
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { MysteryTrader } from '$lib/types/MysteryTrader';
	import type { Planet } from '$lib/types/Planet';
	import type { Salvage } from '$lib/types/Salvage';
	import type { Wormhole } from '$lib/types/Wormhole';
	import { readable } from 'svelte/store';
	import FleetSummary from './FleetSummary.svelte';
	import MineFieldSummary from './MineFieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import MysteryTraderSummary from './MysteryTraderSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import SalvageSummary from './SalvageSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';
	import WormholeSummary from './WormholeSummary.svelte';
	import { ChevronUp, ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { universe, selectNextMapObject, selectedMapObject } = getGameContext();

	// if we are in a CommandPaneCarousel, show the disclosure chevrons and hide/show the command pane on click
	let carouselContext = getCarouselContext();
	let showDisclosure = carouselContext != undefined;
	let open = carouselContext ? carouselContext.open : readable<boolean>(true);

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
	$: selectedMysteryTrader =
		$selectedMapObject?.type == MapObjectType.MysteryTrader
			? ($selectedMapObject as MysteryTrader)
			: undefined;
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full">
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				{#if carouselContext}
					<button
						class:cursor-default={!showDisclosure}
						class="w-full"
						on:click={carouselContext?.onDisclosureClicked}
					>
						{$selectedMapObject?.name ?? ''}
					</button>
				{:else}
					{$selectedMapObject?.name ?? ''}
				{/if}
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
				{#if carouselContext}
					<button type="button" on:click={carouselContext.onDisclosureClicked}>
						{#if $open}
							<Icon src={ChevronUp} size="16" class="hover:stroke-accent" />
						{:else}
							<Icon src={ChevronDown} size="16" class="hover:stroke-accent" />
						{/if}
					</button>
				{/if}
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
		{:else if selectedMysteryTrader}
			<MysteryTraderSummary mysteryTrader={selectedMysteryTrader} />
		{:else}
			<UnknownSummary />
		{/if}
	</div>
</div>
