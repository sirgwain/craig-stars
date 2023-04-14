<script lang="ts">
	import { page } from '$app/stores';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { game, player, techs } from '$lib/services/Context';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { onMount } from 'svelte';

	let gameId = $page.params.id;

	let designs: ShipDesign[] = [];

	onMount(async () => {
		const response = await fetch(`/api/games/${gameId}/designs`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			designs = (await response.json()) as ShipDesign[];
		} else {
			console.error(response);
		}
	});
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<div class="w-full flex justify-end gap-2">
		<button class="btn btn-secondary" type="submit">Create Design</button>
	</div>

	<ItemTitle>Designs</ItemTitle>
	{#if designs?.length}
		<ul class="px-1">
			{#each designs as design}
				<li>
					<div class="flex flex-row place-items-center">
						<div class="mr-2 mb-2">
							<TechAvatar tech={$techs.getHull(design.hull)} hullSetNumber={design.hullSetNumber} />
						</div>
						<div>
							<a class="link" href={`/games/${gameId}/designs/${design.num}`}>{design.name}</a>
						</div>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</div>
