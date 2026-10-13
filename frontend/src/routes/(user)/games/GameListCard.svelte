<script lang="ts">
	import { me } from '#lib/services/Stores.js';
	import type { GameWithPlayers } from '#lib/types/cs-proto.js';
	type Props = {
		game: GameWithPlayers;
		attention?: boolean;
		compact?: boolean;
		solo?: boolean;
		setup?: boolean;
		open?: boolean;
		onDelete?: () => void;
		onArchive?: () => void;
	};
	let {
		game,
		attention = false,
		compact = false,
		solo = false,
		setup = false,
		open = false,
		onDelete,
		onArchive
	}: Props = $props();
	let player = $derived(game.players.find((p) => p.userId === BigInt($me.id)));
	let submitted = $derived(game.players.filter((p) => p.submittedTurn).length);
	let host = $derived(game.game?.hostId === BigInt($me.id));
	let href = $derived(
		setup && !player?.ready ? `/join-game/${game.game?.id}` : `/games/${game.game?.id}`
	);
</script>

<article
	class="card bg-base-200 border border-base-300 min-w-0"
	class:border-l-4={attention}
	class:border-l-accent={attention}
>
	<div class="card-body p-4 gap-3" class:sm:flex-row={compact} class:sm:items-center={compact}>
		<div class="min-w-0 flex-1">
			<div class="flex flex-wrap items-center gap-2">
				<a
					class="text-2xl font-semibold hover:text-accent min-w-0"
					{href}
					data-type="game-link"
					data-id={game.game?.id}>{game.game?.name}</a
				>
				{#if attention}<span class="badge badge-accent">Turn ready</span
					>{:else if player?.submittedTurn && !solo && !setup}<span class="badge badge-success"
						>Submitted</span
					>{/if}
			</div>
			<p class="text-base-content/70 mt-1">
				{#if setup}{(game.game?.numPlayers ?? 0) - (game.game?.openPlayerSlots ?? 0)} / {game.game
						?.numPlayers} players{#if open}
						· {game.game?.openPlayerSlots} seats open{/if}{:else}Year {game.game?.year} · {game
						.players.length} players{/if}
			</p>
			{#if !solo && !setup}<div class="flex flex-wrap items-center gap-2 mt-3">
					<progress
						class="progress progress-secondary w-32 grow max-w-64"
						value={submitted}
						max={Math.max(1, game.players.length)}
						aria-label={`${submitted} of ${game.players.length} turns submitted`}
					></progress><span class="text-base-content/70"
						>{submitted} / {game.players.length} submitted</span
					>
				</div>{/if}
		</div>
		<div class="flex items-center gap-2 justify-end shrink-0" class:self-end={!compact}>
			<a class="btn" class:btn-accent={attention} class:btn-primary={!attention} {href}
				>{setup
					? open || !player?.ready
						? 'Join'
						: 'Open lobby'
					: solo
						? 'Continue'
						: attention
							? 'Play turn'
							: 'View'} →</a
			>
			{#if onArchive || (host && onDelete)}<details class="dropdown dropdown-end">
					<summary class="btn btn-ghost" aria-label={`Manage ${game.game?.name}`}>•••</summary>
					<ul
						class="dropdown-content menu bg-base-100 rounded-box border border-base-300 z-20 w-40 p-2 shadow-lg"
					>
						{#if onArchive}<li><button onclick={onArchive}>Archive game</button></li>{/if}
						{#if host && onDelete}<li>
								<button
									class="text-error"
									onclick={onDelete}
									data-type="delete-button"
									data-id={game.game?.id}>Delete game</button
								>
							</li>{/if}
					</ul>
				</details>{/if}
		</div>
	</div>
</article>
