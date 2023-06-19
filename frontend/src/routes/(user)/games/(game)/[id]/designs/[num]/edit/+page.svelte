<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { techs, designs } from '$lib/services/Context';
	import { DesignService } from '$lib/services/DesignService';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import type { TechHull } from '$lib/types/Tech';
	import { onMount } from 'svelte';

	let gameId = $page.params.id;
	let num = parseInt($page.params.num);

	$: design = $designs?.find((d) => d.num === num);
	$: hull = design && $techs.getHull(design.hull);

	let error = '';

	const onSave = async () => {
		error = '';

		try {
			if (design) {
				// update this design
				design = await DesignService.update(gameId, design);
				const filteredDesigns = $designs?.filter((d) => d != design) ?? []
				$designs = [...filteredDesigns, design];
				goto(`/games/${gameId}/designs/${design.num}`);
			}
		} catch (e) {
			error = (e as Error).message;
		}
	};
</script>

{#if design && hull}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${gameId}/designs`}>Designs</a></li>
			<li>{design.name == '' ? 'new' : design.name}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success" type="submit" on:click={(e) => onSave()}>Save</button>
		</div>
	</Breadcrumb>

	<ShipDesigner {gameId} bind:design {hull} on:save={(e) => onSave()} bind:error />
{/if}
