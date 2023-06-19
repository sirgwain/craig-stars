<script lang="ts">
	import type { Game } from '$lib/types/Game';

	import { onMount } from 'svelte';

	let games: Game[];

	onMount(async () => {
		const response = await fetch(`/api/games`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			games = ((await response.json()) as Game[]).sort((a, b) =>
				b.createdAt.localeCompare(a.createdAt)
			);
		} else {
			console.error(response);
		}
	});
</script>

<h1 class="text-3xl mb-2">Games</h1>

<table class="border-collapse table-auto w-full text-sm">
	<thead>
		<th
			class="border-b dark:border-slate-600 font-medium p-4 pl-8 pt-0 pb-3 text-slate-800 dark:text-slate-200 text-left"
			>ID</th
		>
		<th
			class="border-b dark:border-slate-600 font-medium p-4 pl-8 pt-0 pb-3 text-slate-800 dark:text-slate-200 text-left"
			>Name</th
		>
		<th
			class="border-b dark:border-slate-600 font-medium p-4 pl-8 pt-0 pb-3 text-slate-800 dark:text-slate-200 text-left"
			>Created</th
		>
	</thead>
	<tbody class="bg-white dark:bg-slate-800">
		{#if games?.length}
			{#each games as game}
				<tr
					><td
						class="border-b border-slate-100 dark:border-slate-700 p-4 pl-8 text-slate-700 dark:text-slate-400"
						>{game.id}</td
					>
					<td
						class="border-b border-slate-100 dark:border-slate-700 p-4 pl-8 text-slate-700 dark:text-slate-400"
						><a href={`/games/${game.id}`}>{game.name}</a></td
					><td
						class="border-b border-slate-100 dark:border-slate-700 p-4 pl-8 text-slate-700 dark:text-slate-400"
						>{game.createdAt}</td
					></tr
				>
			{/each}
		{/if}
	</tbody>
</table>
