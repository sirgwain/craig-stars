<script lang="ts">
	import type { Game } from '$lib/types/Game';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let game: Game;
	export let href: string | undefined = undefined;
	export let showDelete = false;

	const deleteGame = async (game: Game) => {
		if (game.name != undefined && confirm(`Are you sure you want to delete ${game.name}?`)) {
			dispatch('delete', { game });
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
				<div class="text-right font-semibold mr-2 w-28">State</div>
				<div>{startCase(game.state)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Size</div>
				<div>{startCase(game.size)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Density</div>
				<div>{startCase(game.density)}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Year</div>
				<div>{game.year ?? '2400'}</div>
			</div>
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Players</div>
				<div>
					{#if game.openPlayerSlots > 0}
						{game.numPlayers - game.openPlayerSlots}/ {game.numPlayers}
					{:else}
						{game.numPlayers}
					{/if}
				</div>
			</div>
		</div>
		{#if showDelete}
			<div class="card-actions justify-start">
				<div>
					<button class="btn" on:click={(e) => deleteGame(game)}>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
