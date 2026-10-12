<script lang="ts">
	import { goto } from '$app/navigation';
	import TagBadge from '#lib/components/game/TagBadge.svelte';
	import SortableTableHeader from '#lib/components/table/SortableTableHeader.svelte';
	import type { TableColumn } from '#lib/components/table/Table.js';
	import Table from '#lib/components/table/Table.svelte';
	import TableSearchInput from '#lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import {
		MapObjectSchema,
		MapObjectType,
		type Planet,
		type PlayerNote,
		type ProductionQueueItem
	} from '#lib/types/cs-proto.js';
	import { getMapObjectName, type MapObjectLike } from '#lib/types/MapObject.js';
	import { getFullName, getQuantityText, hasQuantity } from '#lib/types/QueueItemType.js';
	import { sameTag } from '#lib/types/Tags.js';
	import { create } from '@bufbuild/protobuf';
	import NoteDialog from '../dialogs/note/NoteDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';

	const {
		game,
		player,
		universe,
		commandedPlanet,
		commandMapObject,
		nextMapObject,
		previousMapObject,
		selectMapObject,
		zoomToMapObject,
		editNote,
		editLabel,
		deletePlayerNote,
		updatePlanetOrders
	} = getGameContext();

	type Kind = 'Planet' | 'Fleet' | 'Label' | 'Production';
	const kinds: Kind[] = ['Planet', 'Fleet', 'Label', 'Production'];

	// Everything the player has written down: notes and tags on map objects, labels on the map and
	// notes on the things their planets are building.
	type Row = {
		id: string;
		kind: Kind;
		subject: string;
		text: string;
		tags: string[];
		// what the note is about, a label isn't about anything
		mo?: MapObjectLike;
		note?: PlayerNote;
		// production notes belong to an item in a planet's queue
		planet?: Planet;
		item?: ProductionQueueItem;
	};

	function noteRow(note: PlayerNote): Row | undefined {
		const id = `note-${note.gameDbObject?.id}`;
		if (note.targetType === MapObjectType.UNSPECIFIED) {
			const { x = 0, y = 0 } = note.position ?? {};
			return { id, kind: 'Label', subject: `(${x}, ${y})`, text: note.text, tags: [], note };
		}
		const mo = $universe.getMapObject(note);
		if (!mo) {
			return undefined;
		}
		return {
			id,
			kind: note.targetType === MapObjectType.FLEET ? 'Fleet' : 'Planet',
			subject: getMapObjectName(mo),
			text: note.text,
			tags: note.tags,
			mo,
			note
		};
	}

	function queueRows(planet: Planet): Row[] {
		return (planet.planetOrders?.productionQueue ?? []).flatMap((item, index) =>
			item.note
				? [
						{
							id: `queue-${planet.mapObject?.num}-${index}`,
							kind: 'Production' as Kind,
							subject: `${planet.mapObject?.name}: ${getFullName(item, $universe)}${
								hasQuantity(item.type) ? ` x ${getQuantityText(item.type, item.quantity)}` : ''
							}`,
							text: item.note,
							tags: [],
							mo: planet,
							planet,
							item
						}
					]
				: []
		);
	}

	let rows = $derived([
		...$universe.notes.map(noteRow).filter((row) => row !== undefined),
		...$universe.getMyPlanets('name', false).flatMap(queueRows)
	]);

	let search = $state('');
	let kind = $state<Kind | ''>('');
	// rows must have all of these tags
	let tagFilter = $state<string[]>([]);
	let sortKey = $state<keyof Row>('subject');
	let sortDescending = $state(false);

	function toggleTag(tag: string) {
		tagFilter = tagFilter.some((t) => sameTag(t, tag))
			? tagFilter.filter((t) => !sameTag(t, tag))
			: [...tagFilter, tag];
	}

	// tags that are no longer on anything can't be filtered by
	let activeTagFilter = $derived(
		tagFilter.filter((t) => $universe.tags.some((tc) => sameTag(tc.tag, t)))
	);

	let filteredRows = $derived(
		rows
			.filter((row) => !kind || row.kind === kind)
			.filter((row) => activeTagFilter.every((t) => row.tags.some((tag) => sameTag(tag, t))))
			.filter((row) => {
				const s = search.trim().toLowerCase();
				return (
					!s ||
					row.subject.toLowerCase().includes(s) ||
					row.text.toLowerCase().includes(s) ||
					row.tags.some((tag) => tag.toLowerCase().includes(s))
				);
			})
			.sort((a, b) => {
				const [x, y] = sortDescending ? [b, a] : [a, b];
				const compared =
					sortKey === 'tags'
						? x.tags.join(' ').localeCompare(y.tags.join(' '))
						: `${x[sortKey] ?? ''}`.localeCompare(`${y[sortKey] ?? ''}`, undefined, {
								numeric: true
							});
				return compared || a.subject.localeCompare(b.subject, undefined, { numeric: true });
			})
	);

	type TableRow = Row & { actions?: never };
	const columns: TableColumn<TableRow>[] = [
		{ key: 'subject', title: 'About' },
		{ key: 'kind', title: 'Type' },
		{ key: 'text', title: 'Note' },
		{ key: 'tags', title: 'Tags' },
		{ key: 'actions', title: '', sortable: false }
	];

	function onSorted(column: TableColumn<TableRow>, descending: boolean) {
		sortKey = column.key as keyof Row;
		sortDescending = descending;
	}

	// show what the note is about on the map
	function gotoRow(row: Row) {
		if (row.mo) {
			if (row.mo.mapObject?.playerNum === $player.num) {
				commandMapObject(row.mo);
				if (row.kind === 'Fleet') {
					// a fleet in orbit doesn't select itself, and the summary is where its note is shown
					selectMapObject(row.mo);
				}
			} else {
				selectMapObject(row.mo);
			}
			zoomToMapObject(row.mo);
		} else if (row.note) {
			// there is nothing to select at a label, just look at it
			zoomToMapObject({
				mapObject: create(MapObjectSchema, { position: row.note.position, name: row.note.text })
			});
		}
		goto(`/games/${$game.id}`);
	}

	let showProductionQueueDialog = $state(false);

	// save the queue being edited and move the dialog to another planet
	async function changePlanet(change: () => void) {
		if ($commandedPlanet) {
			await updatePlanetOrders($commandedPlanet);
		}
		change();
	}

	function edit(row: Row) {
		if (row.planet) {
			// the note is edited along with the queue it's in
			commandMapObject(row.planet);
			showProductionQueueDialog = true;
		} else if (row.kind === 'Label' && row.note) {
			editLabel(row.note);
		} else if (row.mo) {
			editNote(row.mo);
		}
	}
</script>

<div class="w-full" data-type="notes-report">
	<div class="flex flex-row flex-wrap items-center gap-2 m-2">
		<div><TableSearchInput bind:value={search} /></div>
		<select class="select select-sm w-auto" bind:value={kind} aria-label="Type of note">
			<option value="">All types</option>
			{#each kinds as k (k)}
				<option value={k}>{k}</option>
			{/each}
		</select>
		{#if $universe.tags.length}
			<div class="flex flex-row flex-wrap items-center gap-1" data-type="tag-filter">
				{#each $universe.tags as { tag, count } (tag)}
					{@const on = activeTagFilter.some((t) => sameTag(t, tag))}
					<button
						type="button"
						class="cursor-pointer"
						aria-pressed={on}
						aria-label={`Filter by ${tag}`}
						title={`${count} tagged`}
						onclick={() => toggleTag(tag)}
					>
						<TagBadge {tag} outline={!on} />
					</button>
				{/each}
			</div>
		{/if}
	</div>
	<Table
		{columns}
		rows={filteredRows}
		externalSortAndFilter={true}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full',
			th: 'sticky top-0 bg-base-200 z-10'
		}}
	>
		{#snippet head({ column })}
			<span>
				<SortableTableHeader
					{column}
					isSorted={sortKey === column.key}
					{sortDescending}
					{onSorted}
				/>
			</span>
		{/snippet}

		{#snippet cell({ column, row })}
			{#if column.key == 'subject'}
				<button
					class="cs-link text-xl text-left max-w-[20rem] truncate align-bottom"
					title="Go to this on the map"
					onclick={() => gotoRow(row)}
					>{row.kind === 'Label' ? 'Map label' : ''} {row.subject}</button
				>
			{:else if column.key == 'kind'}
				{row.kind}
			{:else if column.key == 'text'}
				<div class="min-w-48 max-w-xl whitespace-pre-wrap wrap-break-word" data-type="note-text">
					{row.text}
				</div>
			{:else if column.key == 'tags'}
				<div class="flex flex-row flex-wrap gap-1 max-w-[16rem]">
					{#each row.tags as tag (tag)}
						<button
							type="button"
							class="cursor-pointer"
							title={`Filter by ${tag}`}
							onclick={() => toggleTag(tag)}
						>
							<TagBadge {tag} />
						</button>
					{/each}
				</div>
			{:else if column.key == 'actions'}
				<div class="flex flex-row gap-1 justify-end">
					<button class="btn btn-xs btn-outline btn-secondary" onclick={() => edit(row)}>
						{row.planet ? 'Queue' : 'Edit'}
					</button>
					{#if row.note}
						{@const note = row.note}
						<button
							class="btn btn-xs btn-outline btn-error"
							onclick={() => {
								if (confirm(`Delete the note about ${row.subject}?`)) {
									deletePlayerNote(note);
								}
							}}
						>
							Delete
						</button>
					{/if}
				</div>
			{/if}
		{/snippet}

		{#snippet empty()}
			<tr>
				<td colspan={columns.length} class="p-4 text-center opacity-70">
					{#if rows.length}
						No notes match.
					{:else}
						You haven't written any notes. Add notes and tags to planets and fleets from their
						summary, label the map from the scanner toolbar, or write a note on an item in a
						planet's production queue.
					{/if}
				</td>
			</tr>
		{/snippet}
	</Table>
</div>

<NoteDialog />

<ProductionQueueDialog
	show={showProductionQueueDialog}
	onNext={() => changePlanet(nextMapObject)}
	onPrev={() => changePlanet(previousMapObject)}
	onOk={(planet) => {
		showProductionQueueDialog = false;
		updatePlanetOrders(planet);
	}}
	onCancel={() => (showProductionQueueDialog = false)}
/>
