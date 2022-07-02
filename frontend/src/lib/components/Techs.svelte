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
			techs = techs.concat(techStore.defenses);
			techs = techs.concat(techStore.hullComponents);
			techs = techs.concat(techStore.hulls);
		} else {
			console.error(response);
		}
	});
</script>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<th>Name</th>
			<th>Category</th>
		</thead>
		<tbody>
			{#if techs?.length}
				{#each techs as tech}
					<tr> <td>{tech.name}</td><td>{tech.category}</td></tr>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
