<script lang="ts">
	import type { Tech, TechStore } from '$lib/types/Tech';

	import { onMount } from 'svelte';

	let techStore: TechStore;
	let techs: Tech[] = [];

	onMount(async () => {
		const response = await fetch(`/api/techs`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			techStore = (await response.json()) as TechStore;
			techs = [];
			techs = techs.concat(techStore.engines);
			techs = techs.concat(techStore.planetaryScanners);
		} else {
			console.error(response);
		}
	});
</script>

<h1 class="text-3xl mb-2">Techs</h1>

<table class="border-collapse table-auto w-full text-sm">
	<thead>
		<th
			class="border-b dark:border-slate-600 font-medium p-4 pl-8 pt-0 pb-3 text-slate-800 dark:text-slate-200 text-left"
			>Name</th
		>
		<th
			class="border-b dark:border-slate-600 font-medium p-4 pl-8 pt-0 pb-3 text-slate-800 dark:text-slate-200 text-left"
			>Category</th
		>
	</thead>
	<tbody class="bg-white dark:bg-slate-800">
		{#if techs?.length}
			{#each techs as tech}
				<tr>
					<td
						class="border-b border-slate-100 dark:border-slate-700 p-4 pl-8 text-slate-700 dark:text-slate-400"
						>{tech.name}</td
					><td
						class="border-b border-slate-100 dark:border-slate-700 p-4 pl-8 text-slate-700 dark:text-slate-400"
						>{tech.category}</td
					></tr
				>
			{/each}
		{/if}
	</tbody>
</table>
