<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { onMount } from 'svelte';

	const { game, player, createDesign } = getGameContext();
	let hullName = $page.params.hull;

	let design: ShipDesign = {
		name: '',
		gameId: $game.id,
		playerNum: $player.num ?? 0,
		originalPlayerNum: 0,
		version: 0,
		hull: '',
		hullSetNumber: 0,
		slots: [],
		spec: {
			engine: {},
			techLevel: {}
		}
	};

	$: hull = $techs.getHull(hullName);
	$: design.hull = hull?.name ?? '';

	let error = '';

	onMount(() => {
		const copyParam = $page.url.searchParams.get('copy');
		if (copyParam) {
			const copyDesign = $game.universe.getMyDesign(parseInt(copyParam));
			if (copyDesign) {
				design.slots = copyDesign?.slots;
				design.spec = copyDesign.spec;
				design.hullSetNumber = copyDesign.hullSetNumber;
				design.version = copyDesign.version + 1;
			}
		}
	});

	const onSave = async () => {
		error = '';
		try {
			const { valid, reason } = $game.validateDesign(design);
			if (valid) {
				const created = await createDesign(design);
				goto(`/games/${$game.id}/designer/${created.num}`);
			} else {
				error = reason ?? 'invalid design';
			}
		} catch (e) {
			error = `${e}`;
		}
	};
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a class="cs-link" href={`/games/${$game.id}/designer`}>Ship Designs</a></li>
		<li><a class="cs-link" href={`/games/${$game.id}/designer/create`}>Choose Hull</a></li>
		<li>{design.name == '' ? 'new' : design.name}</li>
	</svelte:fragment>
	<div slot="end" class="flex justify-end mb-1">
		<button class="btn btn-success" type="submit" on:click={(e) => onSave()}>Save</button>
	</div>
</Breadcrumb>
{#if hull && $game}
	<ShipDesigner bind:design {hull} on:save={(e) => onSave()} bind:error />
{/if}
