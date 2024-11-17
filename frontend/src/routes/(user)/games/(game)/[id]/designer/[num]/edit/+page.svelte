<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	const { game, universe, player, updateDesign } = getGameContext();
	let num = parseInt($page.params.num);

	let design: ShipDesign | undefined = $state(
		$universe.designs.find((d) => d.playerNum == $player.num && d.num === num)
	);
	let hull = $derived(design && $techs.getHull(design.hull));
	let error = $state('');

	async function save() {
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
	}
</script>

{#if design && hull}
	<Breadcrumb>
		{#snippet crumbs()}
			<li><a class="cs-link" href={`/games/${$game.id}/designer`}>Ship Designs</a></li>
			<li>{design?.name === '' ? 'new' : design?.name}</li>
		{/snippet}
		{#snippet end()}
			<div class="flex justify-end mb-1">
				<button class="btn btn-success mx-1" type="submit" onclick={(e) => save()}>Save</button>
			</div>
		{/snippet}
	</Breadcrumb>

	<ShipDesigner bind:design {hull} onsave={save} {error} />
{/if}
