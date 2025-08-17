<script lang="ts">
	import { me } from '$lib/services/Stores';
	import type { GameWithPlayers } from '$lib/types/cs-proto';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		game: GameWithPlayers;
		onDelete?: () => void;
	};

	let { game, onDelete }: Props = $props();

	function ready(game: GameWithPlayers): boolean {
		return game.players.find((p) => p.userId == BigInt($me.id))?.ready ?? false;
	}
</script>

<div class="col-span-5">
	<a
		class="text-primary text-2xl hover:text-accent w-full"
		href={ready(game) ? `/games/${game.game?.id}` : `/join-game/${game.game?.id}`}
		>{game.game?.name}</a
	>
</div>
<div class="col-span-5 text-md">
	{(game.game?.numPlayers ?? 0) - (game.game?.openPlayerSlots ?? 0)} / {game.game?.numPlayers}
</div>

{#if game.game?.hostId == BigInt($me.id)}
	<div class="col-span-2 flex justify-center">
		<button onclick={onDelete} class="btn btn-error btn-sm rounded-md" title="Delete Game">
			<Icon src={XMark} size="16" class="hover:stroke-accent" />
		</button>
	</div>
{:else}
	<div class="col-span-2"></div>
{/if}
