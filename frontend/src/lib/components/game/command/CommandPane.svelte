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
	<table class="border-collapse table-auto w-full text-sm">
		<thead>
			<th class="table-header">ID</th>
			<th class="table-header">Name</th>
			<th class="table-header">Population</th>
			<th class="table-header">Contributes Leftover<br />Resources to Research</th>
		</thead>
		<tbody class="bg-white dark:bg-slate-800">
			{#each player.planets.filter((p) => ownedBy(p, player.num)) as planet}
				<tr>
					<td class="table-cell">{planet.num}</td>
					<td class="table-cell">{planet.name}</td>
					<td class="table-cell">{planet.cargo.colonists ? planet.cargo.colonists * 100 : 0}</td>
					<td class="table-cell"
						><input
							type="checkbox"
							bind:checked={planet.contributesOnlyLeftoverToResearch}
							on:change={(e) => onContributesOnlyLeftoverToResearchChecked(e, planet)}
						/></td
					>
				</tr>
			{/each}
		</tbody>
	</table>
	<div class="grid grid-cols-2 gap-2 mt-5">
		<MineralsOnHandTile />
		<PlanetStatusTile />
	</div>
{/if}
