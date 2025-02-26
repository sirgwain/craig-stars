<script lang="ts">
	import {
		GameStateGeneratingTurnError,
		GameStateSetup,
		type Game,
		type GameWithPlayers
	} from '$lib/types/cs';
	import { Check, Trash, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';

	type Props = {
		game: GameWithPlayers;
		href?: string | undefined;
		showDelete?: boolean;
		onDelete?: (game: Game) => void;
	};

	let { game, href = undefined, showDelete = false, onDelete }: Props = $props();

	const deleteGame = async (game: Game) => {
		if (game.name != undefined && confirm(`Are you sure you want to delete ${game.name}?`)) {
			onDelete?.(game);
		}
	};
</script>

<div
	class="card bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1 w-full sm:w-[350px]"
>
	<div class="card-body">
		<h2 class="card-title">
			{#if href}
				<a class="cs-link" {href}>{game.name}</a>
			{:else}
				{game.name}
			{/if}
		</h2>
		<div class="flex flex-col">
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">State</div>
				<div class:text-error={game.state === GameStateGeneratingTurnError}>
					{startCase(game.state)}
				</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Size</div>
				<div>{startCase(game.size)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Density</div>
				<div>{startCase(game.density)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Year</div>
				<div>{game.year ?? '2400'}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Players</div>
				<div>
					{#if (game.openPlayerSlots ?? 0) > 0}
						{(game.numPlayers ?? 0) - (game.openPlayerSlots ?? 0)}/ {game.numPlayers ?? 1}
					{:else}
						{game.players.length}
					{/if}
				</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Public Player Scores</div>
				<div class="my-auto">
					{#if game.publicPlayerScores}
						<Icon src={Check} size="20" class="stroke-success" />
					{:else}
						<Icon src={XMark} size="20" class="stroke-error" />
					{/if}
				</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Beginner: Maximum Minerals</div>
				<div class="my-auto">
					{#if game.maxMinerals}
						<Icon src={Check} size="20" class="stroke-success" />
					{:else}
						<Icon src={XMark} size="20" class="stroke-error" />
					{/if}
				</div>
			</div>
		</div>
		{#if game.state == GameStateSetup}
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-32">Public</div>
				<div>
					{#if game.public}
						<Icon src={Check} size="20" class="stroke-success" />
					{:else}
						<Icon src={XMark} size="20" class="stroke-error" />
					{/if}
				</div>
			</div>
		{/if}
		{#if showDelete}
			<div class="card-actions justify-start">
				<div>
					<button type="button" class="btn" onclick={() => deleteGame(game)}>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
