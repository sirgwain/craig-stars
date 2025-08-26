<script lang="ts">
	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import GameSettingsEditor from '$lib/components/game/newgame/GameSettingsEditor.svelte';
	import { GameSettingsSchema, type GameSettings } from '$lib/types/cs-proto';
	import type { Player } from '$lib/types/cs-proto';
	import { PlayerType } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import { me } from '$lib/services/Stores';
	import { gameClient, playerClient } from '$lib/services/connect';
	import { create } from '@bufbuild/protobuf';
	import { CheckBadge, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onDestroy, onMount } from 'svelte';
	import RaceView from '../race/RaceView.svelte';
	import GuestLink from './GuestLink.svelte';
	import { addError } from '$lib/services/Errors';
	import type { ConnectError } from '@connectrpc/connect';

	const { game, cs, loadStatus, startPollingStatus, stopPollingStatus, updateGame } =
		getGameContext();

	let settings: GameSettings = $state(
		create(GameSettingsSchema, {
			name: $game.name,
			public: $game.public,
			quickStartTurns: $game.quickStartTurns,
			size: $game.size,
			density: $game.density,
			playerPositions: $game.playerPositions,
			randomEvents: $game.randomEvents,
			computerPlayersFormAlliances: $game.computerPlayersFormAlliances,
			publicPlayerScores: $game.publicPlayerScores,
			maxMinerals: $game.maxMinerals,
			startMode: $game.startMode,
			victoryConditions: $game.victoryConditions,
			players: []
		})
	);

	async function onLeave() {
		try {
			await gameClient.leaveGame({ gameId: $game.id });
			goto(`/games`);
		} catch (e) {
			addError(e as ConnectError);
		}
	}

	async function onUpdateGame() {
		try {
			const resp = await gameClient.updateGame({ gameId: $game.id, settings });
			if (!resp.game) {
				return;
			}

			updateGame(resp.game);
		} catch (e) {
			addError(e as ConnectError);
		}
	}

	async function onAddPlayer(playerType: PlayerType) {
		try {
			const resp = await gameClient.addPlayer({ gameId: $game.id, playerType });
			if (!resp.game) {
				return;
			}
			updateGame(resp.game);
		} catch (e) {
			addError(e as ConnectError);
		}
	}

	async function onDeletePlayer(playerNum: number) {
		try {
			const resp = await gameClient.deletePlayerSlot({ gameId: $game.id, playerNum: playerNum });
			if (!resp.game) {
				return;
			}
			updateGame(resp.game);
		} catch (e) {
			addError(e as ConnectError);
		}
	}

	async function onKickPlayer(playerNum: number) {
		try {
			const resp = await gameClient.kickPlayer({ gameId: $game.id, playerNum: playerNum });
			if (!resp.game) {
				return;
			}
			updateGame(resp.game);
		} catch (e) {
			addError(e as ConnectError);
		}
	}

	async function onStartGame() {
		try {
			stopPollingStatus();
			await gameClient.startGame({ gameId: $game.id });

			// force an update so the game reloads
			await loadStatus();
			goto(`/games/${$game.id}`);
		} catch (e) {
			addError(e as ConnectError);
		}
	}
	// let error = '';

	let player: Player | undefined = $state();

	onMount(async () => {
		try {
			const resp = await playerClient.getPlayer({ gameId: $game.id });
			if (!resp.player) {
				return;
			}
			player = resp.player;
			startPollingStatus();
		} catch (e) {
			addError(e as ConnectError);
		}
	});

	onDestroy(stopPollingStatus);

	let isHost = $derived($me.id === $game.hostId);
	let myPlayer = $derived($game.players.find((p) => p.userId == $me.id));
	let hasGuests = $derived(isHost && $game.players.find((p) => p.guest));
</script>

<div class="w-full mx-auto md:max-w-3xl">
	<ItemTitle>{$game.name}</ItemTitle>

	<div>
		Welcome, {myPlayer?.name}. You are playing as the {player?.race?.pluralName}.
	</div>
	<form class="mt-2">
		<div class="flex flex-col justify-center gap-2 place-items-center">
			<div>
				{#if isHost}
					<GameSettingsEditor
						bind:settings
						showInviteLink={$game.isMultiplayer() && $game.openPlayerSlots > 0}
					/>
				{:else}
					<GameCard game={$game} />
				{/if}
			</div>
			<div class="w-full bg-base-200 shadow rounded-sm border-2 border-base-300 py-2 m-2">
				<div class="grid grid-cols-2 gap-x-5 px-2" class:grid-cols-3={hasGuests}>
					<div class="text-center border-b border-b-secondary mb-1">Player</div>
					<div class="text-center border-b border-b-secondary mb-1 font-semibold text-xl">
						Status
					</div>
					{#if hasGuests}
						<div class="text-center border-b border-b-secondary mb-1">Invite Link</div>
					{/if}

					{#each $game.players as playerStatus, index (index)}
						<div class="flex flex-row h-8 my-auto">
							<div class="w-4">
								{playerStatus.num}
							</div>
							<div
								class="h-4 w-4 my-auto border border-secondary mx-2"
								style={`background-color: ${playerStatus.color}`}
							></div>
							<div class="h-8">
								{playerStatus.name}
							</div>
						</div>
						<div class="flex flex-row" class:justify-center={!hasGuests}>
							<div class="flex flex-row my-auto">
								{#if playerStatus.ready}
									<div class="w-20">Ready</div>
									<Icon src={CheckBadge} size="24" class="stroke-success" />
								{:else}
									<div class="w-20">Waiting</div>
									<Icon src={XMark} size="24" class="stroke-error" />
								{/if}
							</div>

							{#if isHost}
								<div class="flex grow justify-center mx-2">
									{#if !playerStatus.aiControlled && playerStatus.userId && index != 0}
										<button
											onclick={() => onKickPlayer(playerStatus.num)}
											type="button"
											class="w-full btn btn-outline btn-sm my-1 normal-case">Kick</button
										>
									{:else if playerStatus.aiControlled || !playerStatus.userId}
										<button
											onclick={() => onDeletePlayer(playerStatus.num)}
											type="button"
											class="w-full btn btn-outline btn-sm my-1 normal-case">Delete</button
										>
									{/if}
								</div>
							{/if}
						</div>
						{#if hasGuests}
							<div>
								<GuestLink player={playerStatus} />
							</div>
						{/if}
					{/each}
				</div>
			</div>
		</div>
		<div class="flex flex-row gap-1 mt-1">
			{#if isHost}
				<button
					type="submit"
					class="btn btn-primary"
					onclick={(e) => {
						e.preventDefault();
						onUpdateGame();
					}}>Update Game</button
				>
				<button type="button" class="btn btn-secondary" onclick={() => onAddPlayer(PlayerType.AI)}
					>Add AI</button
				>
				<button type="button" class="btn btn-secondary" onclick={() => onAddPlayer(PlayerType.OPEN)}
					>Add Open Slot</button
				>
				<button
					type="button"
					class="btn btn-secondary"
					onclick={() => onAddPlayer(PlayerType.GUEST)}>Add Guest</button
				>
				<button
					disabled={$game.players.findIndex((p) => !p.ready) != -1}
					type="button"
					class="btn btn-secondary ml-auto"
					onclick={onStartGame}>Start Game</button
				>
			{:else if $game.players.findIndex((p) => p.userId === $me.id) != -1}
				<button type="button" class="btn btn-secondary" onclick={onLeave}>Leave Game</button>
			{/if}
		</div>
	</form>

	{#if player?.race}
		<ItemTitle>Your Race - {player.race.pluralName}</ItemTitle>
		<RaceView wasmClient={cs.wasmService} race={player.race} />
	{/if}
</div>
