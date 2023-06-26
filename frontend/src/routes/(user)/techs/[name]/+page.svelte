<script lang="ts">
	import { page } from '$app/stores';
	import TechHullSummary from '$lib/components/game/design/Hull.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { Service } from '$lib/services/Service';

	import { TechCategory, type Tech, type TechHull } from '$lib/types/Tech';
	import { startCase } from 'lodash-es';
	import { onMount } from 'svelte';

	let nameSlug = $page.params.name;
	let tech: Tech;

	$: hull = tech as TechHull;

	onMount(async () => {
		const name = startCase(nameSlug);
		const response = await fetch(`/api/techs/${name}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.raiseError(response);
		}
		tech = (await response.json()) as Tech;
	});
</script>

{#if tech}
	<TechSummary {tech} />
	{#if (hull && tech.category == TechCategory.ShipHull) || tech.category == TechCategory.StarbaseHull}
		<h1 class="my-3 text-lg text-center font-semibold">Hull</h1>
		<div
			class="card bg-base-200 shadow w-full max-h-fit min-h-fit rounded-sm border-2 border-base-300"
		>
			<div class="w-full flex flex-row justify-center">
				<TechHullSummary {hull} />
			</div>
		</div>
	{/if}
{/if}
