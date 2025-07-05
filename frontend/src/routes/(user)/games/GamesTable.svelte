<script lang="ts">
	import type { Game } from '$lib/types/cs';
	import { format } from 'date-fns';

	type Props = {
		games: Game[];
	};
	let { games }: Props = $props();

	let sortedGames = $derived(
		games.sort((a, b) => (b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0))
	);
</script>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<tr>
				<th class="sticky top-0 bg-base-200 z-10">ID</th>
				<th class="sticky top-0 bg-base-200 z-10">Name</th>
				<th class="sticky top-0 bg-base-200 z-10">Created</th>
			</tr>
		</thead>
		<tbody>
			{#if games?.length}
				{#each sortedGames as game (game.id)}
					<tr
						><td>{game.id}</td>
						<td><a href={`/games/${game.id}`}>{game.name}</a></td><td
							>{format(game.createdAt ?? '', 'E, MMM do yyyy hh:mm aaa')}</td
						></tr
					>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
