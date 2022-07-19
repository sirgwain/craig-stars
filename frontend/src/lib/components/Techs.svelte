<script lang="ts">
	import type { Tech, TechStore } from '$lib/types/Tech';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import techjson from '$lib/ssr/techs.json';
	import { onMount } from 'svelte';
	import { kebabCase } from 'lodash-es';

	// for ssr, we start with techs from a json file
	export let techStore: TechStore = techjson as TechStore;
	export let techs: Tech[] = [
		...techStore.engines,
		...techStore.planetaryScanners,
		...techStore.defenses,
		...techStore.hullComponents,
		...techStore.hulls,
		...techStore.terraforms
	];

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
		} else {
			console.error(response);
		}
	});
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols4-4 gap-2">
	{#each techs as tech}
		<div>
			<TechSummary {tech} />
		</div>
	{/each}
</div>
