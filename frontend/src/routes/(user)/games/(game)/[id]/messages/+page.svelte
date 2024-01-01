<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { MessageType, gotoTarget, type Message } from '$lib/types/Message';
	import MessageDetail from './MessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	function selectMessage(message: Message) {
		gotoTarget(message, $game.id, $player.num, $universe);
	}

	function getTarget(message: Message) {
		if (message.battleNum) {
			const battle = $universe.getBattle(message.battleNum);

			if (battle) {
				return `Battle at ${$universe.getBattleLocation(battle)}`;
			} else {
				return 'Battle';
			}
		}
		if (message.type === MessageType.GainTechLevel) {
			return 'Research';
		}
		if (message.type === MessageType.TechGained) {
			return message.spec.techGained;
		}

		const target = $universe.getMapObject(message);
		if (target) {
			return target.name;
		}
		return '';
	}

	// filterable messages
	let filteredMessages: Message[] = [];
	let search = '';
	let showAllMessages = false;

	$: filteredMessages =
		$player.messages.filter((m) => showAllMessages || $settings.isMessageVisible(m.type)) ?? [];
	// .filter(
	// 	(m) =>
	// 		m.text?.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
	// 		m.type.toString().toLowerCase().indexOf(search.toLowerCase()) != -1
	// ) ?? [];

	const columns: TableColumn<Message>[] = [
		{
			key: 'target',
			title: 'Target',
			sortBy: (a, b) => (a.targetType ?? '').localeCompare(b.targetType ?? ''),
			filterBy: (value, row) => (getTarget(row) ?? '').toLowerCase().indexOf(value) != -1
		},
		{
			key: 'text',
			title: 'Text'
		}
	];
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
		<div class="form-control">
			<label class="label cursor-pointer">
				<span class="label-text mr-1">Show All Messages</span>
				<input type="checkbox" class="toggle" bind:checked={showAllMessages} />
			</label>
		</div>
	</div>
	<Table
		{columns}
		rows={filteredMessages}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full'
		}}
		filterBy={search.toLowerCase()}
	>
		<span slot="head" let:isSorted let:sortDescending let:column>
			<SortableTableHeader {column} {isSorted} {sortDescending} />
		</span>

		<span slot="cell" let:column let:row>
			{#if column.key == 'target'}
				<button class="cs-link text-2xl" on:click={() => selectMessage(row)}
					>{getTarget(row)}</button
				>
			{:else}
				<MessageDetail message={row} />
			{/if}
		</span>
	</Table>
</div>
