<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { PlayerMessageType, type PlayerMessage } from '$lib/types/cs-proto';
	import { getMapObjectTarget } from '$lib/types/Message';
	import MessageDetail from './MessageDetail.svelte';

	const { game, player, universe, settings, gotoTarget } = getGameContext();

	function selectMessage(message: PlayerMessage) {
		gotoTarget(message, $game.id, $player.num, $universe);
	}

	function getTarget(message: PlayerMessage) {
		if (message.battleNum) {
			const battle = $universe.getBattle(message.battleNum);

			if (battle) {
				return `Battle at ${$universe.getBattleLocation(battle)}`;
			} else {
				return 'Battle';
			}
		}
		if (message.type === PlayerMessageType.PLAYER_GAIN_TECH_LEVEL) {
			return 'Research';
		}
		if (message.type === PlayerMessageType.BATTLE_REPORTS) {
			return 'Battle Reports';
		}
		if (message.type === PlayerMessageType.PLAYER_TECH_GAINED && message.spec) {
			return message.spec.techGained;
		}

		const target = $universe.getMapObject(getMapObjectTarget(message));
		if (target) {
			return target.mapObject?.name ?? '';
		}
		return '';
	}

	// filterable messages
	let search = $state('');
	let showAllMessages = $state(false);
	let filteredMessages: PlayerMessage[] = $derived(
		$player.messages.filter((m) => showAllMessages || $settings.isMessageVisible(m.type))
	);

	type TableMessage = PlayerMessage & { targetType?: never };
	const columns: TableColumn<TableMessage>[] = [
		{
			key: 'targetType',
			title: 'Target',
			sortBy: (a, b) => (a.target?.targetType ?? 0) - (b.target?.targetType ?? 0),
			filterBy: (value, row) => getTarget(row).toLowerCase().indexOf(value) != -1
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
				<input
					type="checkbox"
					class="toggle"
					class:toggle-accent={showAllMessages}
					bind:checked={showAllMessages}
				/>
			</label>
		</div>
	</div>
	<Table
		{columns}
		rows={filteredMessages}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full',
			th: 'sticky top-0 bg-base-200 z-10'
		}}
		filterBy={search.toLowerCase()}
	>
		{#snippet head({ isSorted, sortDescending, column })}
			<span>
				<SortableTableHeader {column} {isSorted} {sortDescending} />
			</span>
		{/snippet}

		{#snippet cell({ column, row })}
			<span>
				{#if column.key == 'targetType'}
					<button
						class="cs-link text-xl text-left"
						onclick={() => selectMessage(row)}
						data-type="goto-target-button"
						data-id={row.target?.targetName}>{getTarget(row)}</button
					>
				{:else}
					<MessageDetail message={row} />
				{/if}
			</span>
		{/snippet}
	</Table>
</div>
