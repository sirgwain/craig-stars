<script lang="ts">
	import { me } from '$lib/services/Stores';
	import type { Game } from '$lib/types/Game';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let game: Game;
	export let showNumSubmitted = true;

	$: numSubmitted = game.players.filter((p) => p.submittedTurn).length;
</script>

<div class="col-span-6">
	<a class="text-primary text-2xl hover:text-accent w-full" href="/games/{game.id}">{game.name}</a>
</div>
<div class="col-span-2 text-md">
	{game.year}
</div>
<div class="col-span-2 text-md">
	{#if showNumSubmitted}
		{numSubmitted} / {game.numPlayers} Done
	{:else}
		{game.numPlayers}
	{/if}
</div>
{#if game.hostId == $me?.id}
	<div class="col-span-2">
		<button on:click={() => dispatch('delete')} class="float-right btn btn-error btn-danger btn-sm">
			<Icon src={XMark} size="16" class="hover:stroke-accent md:hidden" />
			<span class="hidden md:inline-block">Delete</span></button
		>
	</div>
{:else}
	<div class="col-span-2" />
{/if}
