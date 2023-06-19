<script lang="ts">
	import Cycle from '$lib/components/icons/Cycle.svelte';
	import { selectNextMapObject, selectedMapObject } from '$lib/services/Stores';
	import type { FullGame } from '$lib/services/FullGame';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import FleetSummary from './FleetSummary.svelte';
	import MineFieldSummary from './MineFieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';

	export let game: FullGame;
	export let player: Player;

	let selectedPlanet: Planet | undefined;
	let selectedFleet: Fleet | undefined;
	let selectedMineField: MineField | undefined;
	let selectedMineralPacket: MineralPacket | undefined;
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
	}
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300">
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
			<PlanetSummary {game} {player} planet={selectedPlanet} />
		{:else if selectedFleet}
			<FleetSummary {game} {player} fleet={selectedFleet} />
		{:else if selectedMineField}
			<MineFieldSummary {game} {player} mineField={selectedMineField} />
		{:else if selectedMineralPacket}
			<MineralPacketSummary {game} {player} mineralPacket={selectedMineralPacket} />
		{:else}
			<UnknownSummary />
		{/if}
	</div>
</div>
