<script lang="ts">
	import { onShipDesignTooltip } from '#lib/components/game/tooltips/ShipDesignTooltip.js';
	import Cycle from '#lib/components/icons/Cycle.svelte';
	import TagBadge from '#lib/components/game/TagBadge.svelte';
	import Starbase from '#lib/components/icons/Starbase.svelte';
	import type { ShowCargoTransferDialogProps } from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import type {
		Fleet,
		Minefield,
		MineralPacket,
		MysteryTrader,
		Planet,
		Salvage,
		Wormhole
	} from '#lib/types/cs-proto.js';
	import { MapObjectType } from '#lib/types/cs-proto.js';
	import { getTokenDamage } from '#lib/types/Fleet.js';
	import { getMapObjectName } from '#lib/types/MapObject.js';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import FleetSummary from './FleetSummary.svelte';
	import MinefieldSummary from './MinefieldSummary.svelte';
	import MineralPacketSummary from './MineralPacketSummary.svelte';
	import MysteryTraderSummary from './MysteryTraderSummary.svelte';
	import PlanetSummary from './PlanetSummary.svelte';
	import SalvageSummary from './SalvageSummary.svelte';
	import UnknownSummary from './UnknownSummary.svelte';
	import WormholeSummary from './WormholeSummary.svelte';

	const { universe, readOnly, selectNextMapObject, selectedMapObject, editNote } = getGameContext();

	type Props = {
		hideCycleButton?: boolean;
		hideTitle?: boolean;
	} & ShowCargoTransferDialogProps;

	let { onShowCargoTransferDialog, hideTitle, hideCycleButton }: Props = $props();

	function showStarbaseDesign(e: MouseEvent) {
		if (selectedPlanet?.spec?.planetStarbaseSpec?.starbaseDesignNum) {
			const design = $universe.getDesign(
				selectedPlanet.mapObject?.playerNum,
				selectedPlanet.spec.planetStarbaseSpec.starbaseDesignNum
			);
			// we only know the damage of our own starbases
			const starbase = $universe.getMyPlanetStarbase(selectedPlanet.mapObject?.num ?? 0);
			onShipDesignTooltip(e, design, getTokenDamage(starbase?.tokens[0], design));
		}
	}

	let selectedPlanet = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.PLANET
			? ($selectedMapObject as Planet)
			: undefined
	);
	// the player's private note and tags, for planets and their own fleets
	let noteTarget = $derived(
		$selectedMapObject && $universe.canHaveNote($selectedMapObject) ? $selectedMapObject : undefined
	);
	let note = $derived($universe.getNote(noteTarget));
	let selectedFleet = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.FLEET
			? ($selectedMapObject as Fleet)
			: undefined
	);
	let selectedMinefield = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.MINEFIELD
			? ($selectedMapObject as Minefield)
			: undefined
	);
	let selectedMineralPacket = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.MINERAL_PACKET
			? ($selectedMapObject as MineralPacket)
			: undefined
	);
	let selectedSalvage = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.SALVAGE
			? ($selectedMapObject as Salvage)
			: undefined
	);
	let selectedWormhole = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.WORMHOLE
			? ($selectedMapObject as Wormhole)
			: undefined
	);
	let selectedMysteryTrader = $derived(
		$selectedMapObject?.mapObject?.type === MapObjectType.MYSTERY_TRADER
			? ($selectedMapObject as MysteryTrader)
			: undefined
	);
</script>

<div
	class="card bg-base-200 shadow-sm rounded-xs border-2 border-base-300 w-full select-none"
	data-type="map-object-summary"
	data-id={getMapObjectName($selectedMapObject)}
>
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			{#if !hideTitle}
				<div
					class="flex-1 min-w-0 truncate text-center text-lg font-semibold text-secondary"
					title={getMapObjectName($selectedMapObject)}
				>
					{getMapObjectName($selectedMapObject)}
				</div>
			{/if}
			<div>
				{#if noteTarget && (note || !$readOnly)}
					<button
						type="button"
						title={note ? 'edit notes and tags' : 'add notes and tags'}
						aria-label={note ? 'Edit notes and tags' : 'Add notes and tags'}
						data-type="map-object-note-button"
						onclick={() => editNote(noteTarget)}
					>
						<Icon
							src={PencilSquare}
							class={`size-4 hover:stroke-accent ${note ? 'stroke-accent' : ''}`}
						/></button
					>
				{/if}
				{#if selectedPlanet && selectedPlanet.spec?.planetStarbaseSpec?.hasStarbase}
					<button
						type="button"
						onpointerdown={(e) => {
							e.preventDefault();
							showStarbaseDesign(e);
						}}
					>
						<Starbase class="size-4  starbase" /></button
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
						<Cycle class="size-4  fill-base-content hover:stroke-accent" /></button
					>
				{/if}
			</div>
		</div>
		{#if noteTarget && note}
			<button
				type="button"
				class="flex flex-col gap-0.5 text-left text-sm hover:text-accent"
				title={note.text}
				onclick={() => editNote(noteTarget)}
			>
				{#if note.tags.length}
					<span class="flex flex-row flex-wrap gap-1" data-type="map-object-tags">
						{#each note.tags as tag (tag)}
							<TagBadge {tag} />
						{/each}
					</span>
				{/if}
				{#if note.text}
					<span class="italic whitespace-pre-line line-clamp-2" data-type="map-object-note"
						>{note.text}</span
					>
				{/if}
			</button>
		{/if}
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
