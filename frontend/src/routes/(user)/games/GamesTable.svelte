<script lang="ts">
	import type { Game } from '$lib/types/Game';
	import { format, parseJSON } from 'date-fns';

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
				<th>ID</th>
				<th>Name</th>
				<th>Created</th>
			</tr>
		</thead>
		<tbody>
			{#if games?.length}
				{#each sortedGames as game}
					<tr
						><td>{game.id}</td>
						<td><a href={`/games/${game.id}`}>{game.name}</a></td><td
							>{format(parseJSON(game.createdAt), 'E, MMM do yyyy hh:mm aaa')}</td
						></tr
					>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
