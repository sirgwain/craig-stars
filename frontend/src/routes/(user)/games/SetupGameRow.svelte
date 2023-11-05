<script lang="ts">
	import { me } from '$lib/services/Stores';
	import type { Game } from '$lib/types/Game';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let game: Game;

	function ready(game: Game): boolean {
		return game.players.find((p) => p.userId == $me.id)?.ready ?? false;
	}
</script>

<div class="col-span-5">
	<a
		class="text-primary text-2xl hover:text-accent w-full"
		href={ready(game) ? `/games/${game.id}` : `/join-game/${game.id}`}>{game.name}</a
	>
</div>
<div class="col-span-5 text-md">
	{(game.numPlayers ?? 0) - (game.openPlayerSlots ?? 0)} / {game.numPlayers}
</div>

{#if game.hostId == $me.id}
	<div class="col-span-2">
		<button on:click={() => dispatch('delete')} class="float-right btn btn-error btn-danger btn-sm">
			<Icon src={XMark} size="16" class="hover:stroke-accent md:hidden" />
			<span class="hidden md:inline-block">Delete</span></button
		>
	</div>
{:else}
	<div class="col-span-2" />
{/if}
