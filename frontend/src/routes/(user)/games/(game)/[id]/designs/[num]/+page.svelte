<script lang="ts">
	import { page } from '$app/stores';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import Cost from '$lib/components/game/Cost.svelte';
	import Design from '$lib/components/game/Design.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import TechHull from '$lib/components/tech/hull/TechHull.svelte';
	import { game, techs } from '$lib/services/Context';
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
		<ItemTitle>{design.name}</ItemTitle>
		<div class="px-1 md:p-0">
			<Design {design} />
		</div>
	</div>
{/if}
