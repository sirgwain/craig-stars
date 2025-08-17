<script lang="ts">
	import Archive from '$lib/components/icons/Archive.svelte';
	import { me } from '$lib/services/Stores';
	import type { GameWithPlayers } from '$lib/types/cs-proto';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		game: GameWithPlayers;
		showNumSubmitted?: boolean;
		onArchive?: () => void;
		onDelete?: () => void;
	};

	let { game, showNumSubmitted = true, onArchive, onDelete }: Props = $props();

	let numSubmitted = $derived(game.players.filter((p) => p.submittedTurn).length);
</script>

<div class="col-span-5">
	<a
		class="text-primary text-2xl hover:text-accent w-full"
		href="/games/{game.game?.id}"
		data-type="game-link"
		data-id={game.game?.id}>{game.game?.name}</a
	>
</div>
<div class="col-span-2 text-md">
	{game.game?.year}
</div>
<div class="col-span-3 text-md">
	{#if showNumSubmitted}
		{numSubmitted} / {game.players.length} Submitted
	{:else}
		{game.players.length}
	{/if}
</div>
{#if game.game?.hostId == BigInt($me.id)}
	<div class="col-span-2 flex justify-center join">
		<button onclick={onArchive} class="btn btn-error btn-sm rounded-l-md" title="Archive Game">
			<Archive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
		</button>
		<button
			onclick={onDelete}
			class="btn btn-error btn-sm border-l-secondary rounded-r-md"
			title="Delete Game"
			data-type="delete-button"
			data-id={game.game?.id}
		>
			<Icon src={XMark} size="16" class="hover:stroke-accent" />
		</button>
	</div>
{:else}
	<div class="col-span-2 flex justify-center">
		<button onclick={onArchive} class="btn btn-error btn-sm rounded-md" title="Archive Game">
			<Archive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
		</button>
	</div>
{/if}
