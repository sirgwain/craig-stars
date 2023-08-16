<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { gotoMessageTarget } from '$lib/services/Stores';
	import { MessageType, type Message } from '$lib/types/Message';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	const { game, player, universe } = getGameContext();

	const selectMessage = (message: Message) => {
		gotoMessageTarget($game, $player, message);
	};

	const getTarget = (message: Message) => {
		if (message.battleNum) {
			return 'Battle';
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
		return ''
	};

	// filterable messages
	let filteredMessages: Message[] = [];
	let search = '';

	$: filteredMessages =
		$player.messages.filter(
			(i) =>
				i.text.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
				i.type.toString().toLowerCase().indexOf(search.toLowerCase()) != -1
		) ?? [];

	const columns: SvelteTableColumn<Message>[] = [
		{
			key: 'target',
			title: 'Target',
			sortBy: (a, b) => (a.targetType ?? '').localeCompare(b.targetType ?? '')
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
	</div>
	<SvelteTable
		{columns}
		rows={filteredMessages}
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
			{#if column.key == 'target'}
				<button class="cs-link text-2xl" on:click={() => selectMessage(row)}
					>{getTarget(row)}</button
				>
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>
