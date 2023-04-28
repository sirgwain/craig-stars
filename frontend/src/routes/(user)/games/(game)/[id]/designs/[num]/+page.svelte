<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import Design from '$lib/components/game/design/Design.svelte';
	import { techs } from '$lib/services/Context';
	import type { ErrorResponse } from '$lib/types/ErrorResponse';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { onMount } from 'svelte';

	let gameId = $page.params.id;
	let num = $page.params.num;

	let design: ShipDesign;
	let error: ErrorResponse | undefined;

	onMount(async () => {
		const designResponse = await fetch(`/api/games/${gameId}/designs/${num}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (designResponse.ok) {
			design = (await designResponse.json()) as ShipDesign;
		} else {
			error = (await designResponse.json()) as ErrorResponse;
		}
	});

	$: hull = design && $techs.getHull(design.hull);
</script>

{#if error?.error}
	{error.error}
{:else if design}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${gameId}/designs`}>Designs</a></li>
			<li>{design?.name}</li>
			{#if !design.spec?.numInstances}
				<li><a class="cs-link" href={`/games/${gameId}/designs/${design.num}/edit`}>Edit</a></li>
			{/if}
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow h-full px-1 md:p-0">
		<Design {design} />
	</div>
{/if}
