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

	let hull = $derived($techs.getHull(hullName));

	let design: ShipDesign = $state({
		name: '',
		gameId: $game.id,
		playerNum: $player.num ?? 0,
		originalPlayerNum: 0,
		version: 0,
		hull: hullName ?? '',
		hullSetNumber: 0,
		slots: [],
		spec: {
			engine: {},
			techLevel: {}
		}
	});

	let error = $state('');

	onMount(() => {
		const copyParam = $page.url.searchParams.get('copy');
		if (copyParam) {
			const copyDesign = $game.universe.getMyDesign(parseInt(copyParam));
			if (copyDesign) {
				design.slots = copyDesign?.slots.map((s) => Object.assign({}, s));
				design.spec = Object.assign({}, copyDesign.spec);
				design.hullSetNumber = copyDesign.hullSetNumber;
				design.version = copyDesign.version + 1;
				design.name = copyDesign.name;
			}
		}
	});

	async function save() {
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
	}
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li><a class="cs-link" href={`/games/${$game.id}/designer`}>Ship Designs</a></li>
		<li><a class="cs-link" href={`/games/${$game.id}/designer/create`}>Choose Hull</a></li>
		<li>{design.name == '' ? 'new' : design.name}</li>
	{/snippet}
	{#snippet end()}
		<div class="flex justify-end mb-1">
			<button class="btn btn-success mx-1" type="submit" onclick={(e) => save()}>Save</button>
		</div>
	{/snippet}
</Breadcrumb>
{#if hull && $game}
	<ShipDesigner bind:design {hull} onsave={save} {error} />
{/if}
