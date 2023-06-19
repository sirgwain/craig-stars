<script lang="ts">
	import { ownedBy } from '$lib/types/MapObject';
	import type { GameContext } from '$lib/types/GameContext';
	import type { Planet } from '$lib/types/Planet';
	import { createEventDispatcher, getContext } from 'svelte';
	import { PlanetService } from '$lib/services/PlanetService';
	import MineralsOnHandTile from './MineralsOnHandTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';

	const { game, player } = getContext<GameContext>('game');

	const planetService = new PlanetService();

	async function onContributesOnlyLeftoverToResearchChecked(e: Event, planet: Planet) {
		console.log(
			'Checked ContributesOnlyLeftoverToResearchChecked',
			planet.contributesOnlyLeftoverToResearch
		);
		planet = await planetService.updatePlanet(planet);
	}
</script>

{#if game && player.planets?.length}
	<!-- <table class="table w-full">
		<thead>
			<th>ID</th>
			<th>Name</th>
			<th>Population</th>
			<th>Contributes Leftover<br />Resources to Research</th>
		</thead>
		<tbody>
			{#each player.planets.filter((p) => ownedBy(p, player.num)) as planet}
				<tr>
					<td>{planet.num}</td>
					<td>{planet.name}</td>
					<td>{planet.cargo.colonists ? planet.cargo.colonists * 100 : 0}</td>
					<td
						><input
							type="checkbox"
							bind:checked={planet.contributesOnlyLeftoverToResearch}
							on:change={(e) => onContributesOnlyLeftoverToResearchChecked(e, planet)}
						/></td
					>
				</tr>
			{/each}
		</tbody>
	</table> -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
		<MineralsOnHandTile />
		<PlanetStatusTile />
	</div>
{/if}
