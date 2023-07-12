<script lang="ts">
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { selectNextMapObject, selectedMapObject } from '$lib/services/Stores';
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

	const { game, player, universe } = getGameContext();

	let selectedPlanet: Planet | undefined;
	let selectedFleet: Fleet | undefined;
	let selectedMineField: MineField | undefined;
	let selectedMineralPacket: MineralPacket | undefined;
	let selectedSalvage: Salvage | undefined;
	let selectedWormhole: Wormhole | undefined;
	$: {
		selectedPlanet =
			$selectedMapObject?.type == MapObjectType.Planet ? ($selectedMapObject as Planet) : undefined;
		selectedFleet =
			$selectedMapObject?.type == MapObjectType.Fleet ? ($selectedMapObject as Fleet) : undefined;
		selectedMineField =
			$selectedMapObject?.type == MapObjectType.MineField
				? ($selectedMapObject as MineField)
				: undefined;
		selectedMineralPacket =
			$selectedMapObject?.type == MapObjectType.MineralPacket
				? ($selectedMapObject as MineralPacket)
				: undefined;
		selectedSalvage =
			$selectedMapObject?.type == MapObjectType.Salvage
				? ($selectedMapObject as Salvage)
				: undefined;
		selectedWormhole =
			$selectedMapObject?.type == MapObjectType.Wormhole
				? ($selectedMapObject as Wormhole)
				: undefined;
	}
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full">
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				{$selectedMapObject?.name ?? ''}
			</div>
			<button type="button" on:click|preventDefault={() => selectNextMapObject()}>
				<Cycle class="w-4 h-4 fill-base-content hover:stroke-accent" /></button
			>
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
