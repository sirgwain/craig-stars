<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import type { Game } from '$lib/types/Game';
	import PlayersStatus from './players/PlayersStatus.svelte';

	export let game: Game;
	export let title: string;
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<ItemTitle>{title}</ItemTitle>

	<div class="flex flex-col justify-center gap-2 place-items-center">
		<div>
			<GameCard {game} href={`/games/${game.id}`} />
		</div>
		<div class="w-full bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1">
			<PlayersStatus />
		</div>
	</div>
	<slot />
</div>
