<script lang="ts">
	import { mapObjects, player, playerName } from '$lib/services/Context';
	import { Unexplored } from '$lib/types/Planet';
</script>

{#if $player && $mapObjects}
	<div class="w-full">
		<table class="table table-compact w-full">
			<caption>My Planets</caption>
			<thead>
				<th>Name</th>
				<th>Population</th>
			</thead>
			<tbody>
				{#if $mapObjects.planets?.length}
					{#each $mapObjects.planets as planet}
						<tr class="hover">
							<td>{planet.name}</td><td
								>{planet.cargo?.colonists ? planet.cargo.colonists * 100 : ''}</td
							></tr
						>
					{/each}
				{/if}
			</tbody>
		</table>

		<table class="table table-compact w-full">
			<caption>Other Planets</caption>
			<thead>
				<th>Name</th>
				<th>Hab</th>
				<th>Owner</th>
			</thead>
			<tbody>
				{#if $player.planetIntels?.length}
					{#each $player.planetIntels as planet}
						<tr class="hover">
							<td>{planet.name}</td><td
								>{planet.hab?.grav
									? `${planet.hab.grav}, ${planet.hab.temp}, ${planet.hab.rad}`
									: ''}</td
							>
							<td>
								{#if planet.playerNum}
									{playerName(planet.playerNum)}
								{:else if planet.reportAge != Unexplored}
									--
								{:else}
									Unknown
								{/if}
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
{/if}
