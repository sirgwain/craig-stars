<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { game, techs } from '$lib/services/Context';
	import { DesignService } from '$lib/services/DesignService';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	let gameId = $page.params.id;
	let num = parseInt($page.params.num);

	$: design = $game && ($game.player.getDesign($game.player.num, num) as ShipDesign);
	$: hull = design && $techs.getHull(design.hull);

	let error = '';

	const onSave = async () => {
		error = '';

		try {
			if (design) {
				// update this design
				design = await DesignService.update(gameId, design);
				$game?.player.updateDesign(design);
				goto(`/games/${gameId}/designs/${design.num}`);
			}
		} catch (e) {
			error = (e as Error).message;
		}
	};
</script>

{#if design && hull && $game}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${gameId}/designs`}>Ship Designs</a></li>
			<li>{design.name == '' ? 'new' : design.name}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success" type="submit" on:click={(e) => onSave()}>Save</button>
		</div>
	</Breadcrumb>

	<ShipDesigner
		player={$game.player}
		{gameId}
		bind:design
		{hull}
		on:save={(e) => onSave()}
		bind:error
	/>
{/if}
