<script lang="ts">
	import { page } from '$app/stores';
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
	<div class="w-full mx-auto md:max-w-2xl">
		<div class="breadcrumbs border-primary border-b-2 mb-2">
			<ul>
				<li><a class="cs-link" href={`/games/${gameId}/designs`}>Designs</a></li>
				<li>{design?.name}</li>
				{#if !design.spec?.numInstances}
					<li><a class="cs-link" href={`/games/${gameId}/designs/${design.num}/edit`}>Edit</a></li>
				{/if}
			</ul>
		</div>

		<div class="px-1 md:p-0">
			<Design {design} />
		</div>
	</div>
{/if}
