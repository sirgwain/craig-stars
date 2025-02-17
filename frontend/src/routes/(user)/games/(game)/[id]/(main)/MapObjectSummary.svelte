<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import Starbase from '$lib/components/icons/Starbase.svelte';
	import { getCarouselContext } from '$lib/services/CarouselContext';
	import type { ShowCargoTransferDialogProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { type Fleet } from '$lib/types/cs';
	import { getMapObjectName, MapObjectType } from '$lib/types/MapObject';
	import type { MineField } from '$lib/types/cs';
	import type { MineralPacket } from '$lib/types/cs';
	import type { MysteryTrader } from '$lib/types/cs';
	import type { Planet } from '$lib/types/cs';
	import type { Salvage } from '$lib/types/cs';
	import type { Wormhole } from '$lib/types/cs';
	import { ChevronDown, ChevronUp } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { readable } from 'svelte/store';
	import FleetSummary from './FleetSummary.svelte';
	import MineFieldSummary from './MineFieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import MysteryTraderSummary from './MysteryTraderSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import SalvageSummary from './SalvageSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';
	import WormholeSummary from './WormholeSummary.svelte';

	const { universe, selectNextMapObject, selectedMapObject } = getGameContext();

	let { onShowCargoTransferDialog }: ShowCargoTransferDialogProps = $props();

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

	let selectedPlanet = $derived(
		$selectedMapObject?.type == MapObjectType.Planet ? ($selectedMapObject as Planet) : undefined
	);
	let selectedFleet = $derived(
		$selectedMapObject?.type == MapObjectType.Fleet ? ($selectedMapObject as Fleet) : undefined
	);
	let selectedMineField = $derived(
		$selectedMapObject?.type == MapObjectType.MineField
			? ($selectedMapObject as MineField)
			: undefined
	);
	let selectedMineralPacket = $derived(
		$selectedMapObject?.type == MapObjectType.MineralPacket
			? ($selectedMapObject as MineralPacket)
			: undefined
	);
	let selectedSalvage = $derived(
		$selectedMapObject?.type == MapObjectType.Salvage ? ($selectedMapObject as Salvage) : undefined
	);
	let selectedWormhole = $derived(
		$selectedMapObject?.type == MapObjectType.Wormhole
			? ($selectedMapObject as Wormhole)
			: undefined
	);
	let selectedMysteryTrader = $derived(
		$selectedMapObject?.type == MapObjectType.MysteryTrader
			? ($selectedMapObject as MysteryTrader)
			: undefined
	);
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full">
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				{#if carouselContext}
					<button
						class:cursor-default={!showDisclosure}
						class="w-full"
						onclick={() => carouselContext?.onDisclosureClicked}
					>
						{getMapObjectName($selectedMapObject)}
					</button>
				{:else}
					{getMapObjectName($selectedMapObject)}
				{/if}
			</div>
			<div>
				{#if selectedPlanet && selectedPlanet.spec.hasStarbase}
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
				<button
					type="button"
					onpointerdown={(e) => {
						e.preventDefault();
						selectNextMapObject();
					}}
				>
					<Cycle class="w-4 h-4 fill-base-content hover:stroke-accent" /></button
				>
				{#if carouselContext}
					<button type="button" onclick={(_) => carouselContext.onDisclosureClicked()}>
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
