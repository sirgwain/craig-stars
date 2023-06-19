<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import ShipDesigner from '$lib/components/game/design/ShipDesigner.svelte';
	import { player, techs } from '$lib/services/Context';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	let gameId = $page.params.id;
	let hullName = $page.params.hull;

	let design: ShipDesign = {
		name: '',
		gameId: parseInt(gameId),
		playerNum: $player?.num ?? 0,
		version: 0,
		hull: '',
		hullSetNumber: 0,
		slots: [],
		spec: {}
	};

	$: hull = $techs.getHull(hullName);
	$: design.hull = hull?.name ?? '';
	let error = '';

	const onSave = async () => {
		error = '';
		const data = JSON.stringify(design);

		const response = await fetch(`/api/games/${gameId}/designs`, {
			method: 'post',
			headers: {
				accept: 'application/json'
			},
			body: data
		});

		if (response.ok) {
			const created = (await response.json()) as ShipDesign;
			goto(`/games/${gameId}/designs/${created.num}`);
		} else {
			const resolvedResponse = await response?.json();
			error = resolvedResponse.error;
			console.error(error);
		}
	};
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<ShipDesigner bind:design {hullName} on:save={(e) => onSave()} bind:error />
</div>
