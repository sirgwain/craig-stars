<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { GameService } from '$lib/services/GameService';
	import type { PlayerResponse } from '$lib/types/Player';
	import { CheckBadge, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte/internal';

	const { game, player, universe } = getGameContext();

	$: players = $universe.players;
	$: playerStatuses = [] as PlayerResponse[];

	async function loadPlayerStatuses() {
		playerStatuses = await GameService.loadPlayerStatuses($game.id);
	}

	function getPlayerStatus(num: number): PlayerResponse | undefined {
		if (num > 0 && num <= playerStatuses.length) {
			return playerStatuses[num - 1];
		}
	}

	onMount(async () => {
		await loadPlayerStatuses();
	});
</script>

<div class="grid grid-cols-4 px-2">
	<div>Player</div>
	<div class="font-semibold text-xl">Name</div>
	<div class="font-semibold text-xl">Race</div>
	<div class="font-semibold text-xl">Status</div>
	{#each players as player}
		<div class="flex flex-row">
			<div class="w-4">
				{player.num}
			</div>
			<div
				class="h-4 w-4 my-auto border border-secondary ml-2"
				style={`background-color: ${$universe.getPlayerColor(player.num)}`}
			/>
		</div>
		<div>{player.racePluralName}</div>
		<div>{player.name}</div>
		<div>
			{#if playerStatuses.length > 0 && getPlayerStatus(player.num)?.submittedTurn}
				<div class="flex flex-row">
					<div class="w-14">Done</div>
					<Icon src={CheckBadge} size="24" class="stroke-success" />
				</div>
			{:else}
				<div class="flex flex-row">
					<div class="w-14">Playing</div>
					<Icon src={XMark} size="24" class="stroke-error" />
				</div>
			{/if}
		</div>
	{/each}
</div>
