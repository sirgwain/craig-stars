<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import type { Snippet } from 'svelte';
	import PlayersStatus from './players/PlayersStatus.svelte';
	import type { GameWithPlayers } from '$lib/types/cs-proto';
	import { getGameWithPlayersFlat } from '$lib/types/Game';

	type Props = {
		game: GameWithPlayers;
		title: string;
		children?: Snippet;
	};

	let { game, title, children }: Props = $props();
	let flatGame = $derived(getGameWithPlayersFlat(game));
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<ItemTitle>{title}</ItemTitle>

	<div class="flex flex-col justify-center gap-2 place-items-center">
		<div>
			{#if flatGame}
				<GameCard game={flatGame} href={`/games/${game.game?.id}`} />
			{/if}
		</div>
		<div class="w-full bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1">
			<PlayersStatus />
		</div>
	</div>
	{@render children?.()}
</div>
