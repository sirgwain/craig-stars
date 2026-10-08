<script lang="ts">
	import { playerClient } from '#lib/services/connect.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { me } from '#lib/services/Stores.js';
	import type { PlayerStatus } from '#lib/types/cs-proto.js';
	import { GameState } from '#lib/types/cs-proto.js';
	import { isHotSeat } from '#lib/types/HotSeat.js';
	import { CheckBadge, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GuestLink from '../(main)/GuestLink.svelte';

	const { game, player, choosingPlayer, loadStatus, updateGame } = getGameContext();

	type Props = {
		playerStatus: PlayerStatus;
	};

	let { playerStatus }: Props = $props();

	let isHost = $derived($me.id === $game.hostId);
	let hasGuests = $derived(isHost && $game.players.find((p) => p.guest));
	// in hot seat games, let the user switch to their other players that haven't played yet
	let isHotSeatPlayer = $derived(
		$game.state === GameState.WAITING_FOR_PLAYERS &&
			playerStatus.userId === $me.id &&
			isHotSeat($game.players, $me.id)
	);
	let canPlay = $derived(
		isHotSeatPlayer &&
			!playerStatus.submittedTurn &&
			(playerStatus.num !== $player.num || $choosingPlayer)
	);
	let canUnsubmit = $derived(isHotSeatPlayer && playerStatus.submittedTurn);

	async function onUnsubmitTurn() {
		await playerClient.unsubmitTurn(
			{ gameId: $game.id },
			{ headers: { 'X-As-Player': `${playerStatus.num}` } }
		);
		if (playerStatus.num === $player.num) {
			$player.submittedTurn = false;
			// trigger reactivity on the layout so our hotkeys are wired up again
			updateGame($game);
		}
		await loadStatus();
	}
</script>

<div class="flex flex-row h-10">
	{#if $game.state === GameState.SETUP}
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
	{#if canPlay}
		<a
			data-type="play-button"
			data-id={playerStatus.name}
			href={`/games/${$game.id}?asPlayer=${playerStatus.num}`}
			class="btn btn-primary btn-sm my-auto ml-1">Play</a
		>
	{/if}
	{#if canUnsubmit}
		<button
			type="button"
			data-type="unsubmit-button"
			data-id={playerStatus.name}
			class="btn btn-secondary btn-sm my-auto ml-1"
			onclick={onUnsubmitTurn}>Unsubmit</button
		>
	{/if}
	{#if hasGuests}
		<div class="ml-1">
			<GuestLink player={playerStatus} hideText={true} />
		</div>
	{/if}
</div>
