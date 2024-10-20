<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';

	const { game, universe, player, updateDesign } = getGameContext();
	let num = parseInt($page.params.num);

	$: design = $universe.designs.find((d) => d.playerNum == $player.num && d.num === num);
	$: hull = design && $techs.getHull(design.hull);

	let error = '';

	const onSave = async () => {
		error = '';

		try {
			if (design) {
				const { valid, reason } = $game.validateDesign(design);
				if (valid) {
					// update this design
					await updateDesign(design);
					goto(`/games/${$game.id}/designer/${design.num}`);
				}
			}
		} catch (e) {
			error = (e as Error).message;
		}
	};
</script>

{#if design && hull && $game}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${$game.id}/designer`}>Ship Designs</a></li>
			<li>{design.name == '' ? 'new' : design.name}</li>
		</svelte:fragment>
		<div slot="end" class="flex justify-end mb-1">
			<button class="btn btn-success mx-1" type="submit" on:click={(e) => onSave()}>Save</button>
		</div>
	</Breadcrumb>

	<ShipDesigner bind:design {hull} on:save={(e) => onSave()} bind:error />
{/if}
