<script lang="ts">
	import Archive from '$lib/components/icons/Archive.svelte';
	import { me } from '$lib/services/Stores';
	import type { Game } from '$lib/types/Game';
	import { XMark, ArchiveBox } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	interface Props {
		game: Game;
		showNumSubmitted?: boolean;
	}

	let { game, showNumSubmitted = true }: Props = $props();

	let numSubmitted = $derived(game.players.filter((p) => p.submittedTurn).length);
</script>

<div class="col-span-5">
	<a class="text-primary text-2xl hover:text-accent w-full" href="/games/{game.id}">{game.name}</a>
</div>
<div class="col-span-2 text-md">
	{game.year}
</div>
<div class="col-span-3 text-md">
	{#if showNumSubmitted}
		{numSubmitted} / {game.players.length} Submitted
	{:else}
		{game.players.length}
	{/if}
</div>
{#if game.hostId == $me.id}
	<div class="col-span-2 flex justify-center join">
		<button
			onclick={() => dispatch('archive')}
			class="btn btn-error btn-sm rounded-l-md"
			title="Archive Game"
		>
			<Archive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
		</button>
		<button
			onclick={() => dispatch('delete')}
			class="btn btn-error btn-sm border-l-secondary rounded-r-md"
			title="Delete Game"
		>
			<Icon src={XMark} size="16" class="hover:stroke-accent" />
		</button>
	</div>
{:else}
	<div class="col-span-2 flex justify-center">
		<button
			onclick={() => dispatch('archive')}
			class="btn btn-error btn-sm rounded-md"
			title="Archive Game"
		>
			<Archive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
		</button>
	</div>
{/if}
