<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { techs } from '$lib/services/Stores';
	import { ownedBy } from '$lib/types/MapObject';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { kebabCase } from 'lodash-es';

	const { game, player, universe, settings } = getGameContext();

	// filterable designs
	let filteredDesigns: ShipDesign[] = [];
	let search = '';

	$: filteredDesigns =
		$universe.designs
			.sort((a, b) =>
				a.playerNum != b.playerNum ? a.playerNum - b.playerNum : (a.num ?? 0) - (b.num ?? 0)
			)
			.filter(
				(i) =>
					i.name.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
					i.hull.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
					$universe.getPlayerName(i.playerNum).toLowerCase().indexOf(search.toLowerCase()) != -1
			) ?? [];

	const columns: SvelteTableColumn<ShipDesign>[] = [
		{
			key: 'playerNum',
			title: 'Player'
		},
		{
			key: 'num',
			title: 'ID'
		},
		{
			key: 'name',
			title: 'Name'
		},
		{
			key: 'hull',
			title: 'Hull'
		},
		{
			key: 'rating',
			title: 'Rating',
			sortBy: (a, b) => (a.spec.powerRating ?? 0) - (b.spec.powerRating ?? 0)
		},
		{
			key: 'armor',
			title: 'Armor',
			sortBy: (a, b) => (a.spec.armor ?? 0) - (b.spec.armor ?? 0)
		},
		{
			key: 'shields',
			title: 'Shields',
			sortBy: (a, b) => (a.spec.shields ?? 0) - (b.spec.shields ?? 0)
		},
		{
			key: 'initiative',
			title: 'Initiative',
			sortBy: (a, b) => (a.spec.initiative ?? 0) - (b.spec.initiative ?? 0)
		},
		{
			key: 'movement',
			title: 'Movement',
			sortBy: (a, b) => (a.spec.movement ?? 0) - (b.spec.movement ?? 0)
		},
		{
			key: 'mass',
			title: 'Mass',
			sortBy: (a, b) => (a.spec.mass ?? 0) - (b.spec.mass ?? 0)
		}
	];

	function icon(design: ShipDesign): string {
		if (design) {
			return `hull-${kebabCase(design.hull)}-${design.hullSetNumber ?? 0}`;
		}
		return '';
	}
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<SvelteTable
		{columns}
		rows={filteredDesigns}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full'
		}}
		let:column
		let:cell
		let:row
	>
		<span slot="head" let:isSorted let:sortDescending>
			<SortableTableHeader {column} {isSorted} {sortDescending} />
		</span>

		<span slot="cell">
			{#if column.key === 'name'}
				<button
					class="w-full h-full cursor-help text-left"
					on:pointerdown|preventDefault={(e) => onShipDesignTooltip(e, row)}
				>
					<div class="flex flex-row">
						<div class="avatar mr-2">
							<div
								class="border-2 border-neutral p-2 bg-black"
								style={`border-color: ${$universe.getPlayerColor(row.playerNum)};`}
							>
								<div class="fleet-avatar {icon(row)} bg-black" />
							</div>
						</div>
						<div class="my-auto">
							{cell}
							<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
						</div>
					</div>
				</button>
			{:else if column.key === 'num'}
				{#if row.playerNum === $player.num}
					<a href={`/games/${$game.id}/designer/${row.num}`} class="cs-link">{cell}</a>
				{:else}
					{cell}
				{/if}
			{:else if column.key === 'playerNum'}
				{$universe.getPlayerName(row.playerNum)}
			{:else if column.key === 'mass'}
				{row.spec.mass ?? ''}
			{:else if column.key === 'armor'}
				{row.spec.armor ?? ''}
			{:else if column.key === 'shields'}
				{row.spec.shields ?? ''}
			{:else if column.key === 'rating'}
				{row.spec.powerRating ?? ''}
			{:else if column.key === 'initiative'}
				{row.spec.initiative ?? ''}
			{:else if column.key === 'movement'}
				{row.spec.movement ?? ''}
			{:else if column.key === 'hull'}
				<button
					class="w-full h-full cursor-help text-left"
					on:pointerdown|preventDefault={(e) => onTechTooltip(e, $techs.getTech(row.hull))}
					>{cell}
					<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" /></button
				>
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>
