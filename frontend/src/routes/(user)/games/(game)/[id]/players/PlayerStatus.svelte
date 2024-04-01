<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { me } from '$lib/services/Stores';
	import { GameState } from '$lib/types/Game';
	import type { PlayerStatus } from '$lib/types/Player';
	import { CheckBadge, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GuestLink from '../(main)/GuestLink.svelte';

	const { game, player, universe } = getGameContext();

	export let playerStatus: PlayerStatus;

	$: isHost = $me.id === $game.hostId;
	$: hasGuests = isHost && $game.players.find((p) => p.guest);
</script>

<div class="flex flex-row h-10">
	{#if $game.state === GameState.Setup}
		{#if playerStatus.ready}
			<div class="w-20 my-auto">Ready</div>
			<div class="my-auto"><Icon src={CheckBadge} size="24" class="stroke-success" /></div>
		{:else}
			<div class="w-20 my-auto">Waiting</div>
			<div class="my-auto"><Icon src={XMark} size="24" class="stroke-error" /></div>
		{/if}
	{:else if playerStatus.submittedTurn}
		<div class="w-20 my-auto">Submitted</div>
		<div class="my-auto"><Icon src={CheckBadge} size="24" class="stroke-success" /></div>
	{:else}
		<div class="w-20 my-auto">Playing</div>
		<div class="my-auto"><Icon src={XMark} size="24" class="stroke-error" /></div>
	{/if}
	{#if hasGuests}
		<div class="ml-1">
			<GuestLink player={playerStatus} hideText={true} />
		</div>
	{/if}
</div>
