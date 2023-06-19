<script lang="ts">
	import type { Game } from '$lib/types/Game';
	import { onMount } from 'svelte';

	export let games: Game[];

	$: games &&
		games.sort((a, b) => (b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0));
</script>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<th>ID</th>
			<th>Name</th>
			<th>Created</th>
		</thead>
		<tbody>
			{#if games?.length}
				{#each games as game}
					<tr
						><td>{game.id}</td>
						<td><a href={`/games/${game.id}`}>{game.name}</a></td><td>{game.createdAt}</td></tr
					>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
