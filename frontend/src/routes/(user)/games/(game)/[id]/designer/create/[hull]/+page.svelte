<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { GameDBObjectSchema } from '$lib/types/cs-proto';
	import { ShipDesignSchema, type ShipDesign } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import { clone, create } from '@bufbuild/protobuf';
	import { onMount } from 'svelte';

	const { game, universe, player, createDesign } = getGameContext();
	let hullName = page.params.hull;

	let hull = $derived($techs.getHull(hullName));

	let design: ShipDesign | undefined = $state();

	let error = $state('');

	onMount(() => {
		design = create(ShipDesignSchema, {
			gameDbObject: {
				gameId: $game.id
			},
			playerNum: $player.num,
			hull: hullName ?? ''
		});

		const copyParam = page.url.searchParams.get('copy');
		if (copyParam) {
			const copyDesign = $universe.getMyDesign(parseInt(copyParam));
			if (copyDesign) {
				design = clone(ShipDesignSchema, copyDesign);
				design.gameDbObject = create(GameDBObjectSchema);
				design.version++;
			}
		}
	});

	async function save() {
		if (!design) {
			return;
		}
		error = '';
		try {
			const { valid, reason } = $universe.validateDesign(design);
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
		<li>{design?.name == '' ? 'new' : design?.name}</li>
	{/snippet}
	{#snippet end()}
		<div class="flex justify-end mb-1">
			<button class="btn btn-success mx-1" type="submit" onclick={save}>Save</button>
		</div>
	{/snippet}
</Breadcrumb>
{#if hull && design}
	<ShipDesigner bind:design {hull} onSave={save} {error} />
{/if}
