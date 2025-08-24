<script lang="ts">
	import { page } from '$app/state';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import TechHullSummary from '$lib/components/game/design/Hull.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { techClient } from '$lib/services/connect';
	import { TechService } from '$lib/services/TechService';
	import techjson from '$lib/ssr/techs.json';
	import {
		GetTechsResponseSchema,
		TechCategory,
		type GetTechsResponseJson,
		type TechHull
	} from '$lib/types/cs-proto';

	import type { TechLike } from '$lib/types/Tech';
	import { fromJson } from '@bufbuild/protobuf';
	import { startCase } from 'lodash-es';
	import { onMount } from 'svelte';

	// for ssr, we start with techs from a json file
	let techStore = new TechService(
		fromJson(GetTechsResponseSchema, techjson as unknown as GetTechsResponseJson)
	);
	let techs: TechLike[] = [
		...techStore.planetaryScanners,
		...techStore.defenses,
		...techStore.planetaries,
		...techStore.hullComponents,
		...techStore.hulls,
		...techStore.terraforms
	];

	let nameSlug = page.params.name;
	let tech = $state(techs.find((t) => t.tech?.name === startCase(nameSlug)));

	let hull = $derived(tech as TechHull);

	onMount(async () => {
		const name = startCase(nameSlug);
		const resp = await techClient.getTech({ name });
		if (resp.tech.value) {
			tech = resp.tech.value;
		}
	});
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li><a href="/techs">Techs</a></li>
		<li>{tech?.tech?.name ?? '<unknown>'}</li>
	{/snippet}
</Breadcrumb>

{#if tech}
	<TechSummary {tech} />
	{#if tech.tech?.category === TechCategory.SHIP_HULL || tech.tech?.category === TechCategory.STARBASE_HULL}
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
