<script lang="ts">
	import { RaceService } from '$lib/services/RaceService';
	import type { Race } from '$lib/types/Race';

	import { onMount } from 'svelte';

	let races: Race[];

	const raceService = new RaceService();

	onMount(async () => {
		races = await raceService.loadRaces();

		races.sort((a, b) => (b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0));
	});
</script>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<th>ID</th>
			<th>Name</th>
			<th>Created</th>
		</thead>
		<tbody>
			{#if races?.length}
				{#each races as race}
					<tr
						><td>{race.id}</td>
						<td><a href={`/races/${race.id}`}>{race.name}</a></td><td>{race.createdAt}</td></tr
					>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
