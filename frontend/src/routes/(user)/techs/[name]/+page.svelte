<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import TechHullSummary from '$lib/components/game/design/Hull.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { Service } from '$lib/services/Service';
	import techjson from '$lib/ssr/techs.json';

	import { TechCategory, type Tech, type TechHull, type TechStore } from '$lib/types/Tech';
	import { startCase } from 'lodash-es';
	import { onMount } from 'svelte';

	// for ssr, we start with techs from a json file
	let techStore: TechStore = techjson as TechStore;
	let techs: Tech[] = [
		...techStore.engines,
		...techStore.planetaryScanners,
		...techStore.defenses,
		...techStore.planetaries,
		...techStore.hullComponents,
		...techStore.hulls,
		...techStore.terraforms
	];

	let nameSlug = $page.params.name;
	let tech = techs.find((t) => t.name === startCase(nameSlug));

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
			await Service.throwError(response);
		}
		tech = (await response.json()) as Tech;
	});
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a href={`/techs`}>Techs</a></li>
		<li>{tech?.name ?? '<unknown>'}</li>
	</svelte:fragment>
</Breadcrumb>

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
