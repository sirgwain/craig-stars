<script lang="ts">
	import type { Tech, TechStore } from '$lib/types/Tech';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';

	import { onMount } from 'svelte';
	import { kebabCase } from 'lodash-es';

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
			techs = techs.concat(techStore.terraforms);

// 			let css = '';
// 			for (let i = 0; i < techs.length; i++) {
// 				const tech = techs[i];
// 				css += `
// .${kebabCase(tech.name.replace("'", '').replace(' ', ''))} {
// 	background: url(/images/techs06.png) 0px 0px;
// }

// `;
// 			}

// 			navigator.clipboard.writeText(css);
// 			console.log(css);
		} else {
			console.error(response);
		}
	});
</script>

<!-- <div class="overflow-x-auto">
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
</div> -->
<div class="flex flex-wrap gap-2 justify-evenly">
	{#each techs as tech}
		<div>
			<TechSummary {tech} />
		</div>
	{/each}
</div>
