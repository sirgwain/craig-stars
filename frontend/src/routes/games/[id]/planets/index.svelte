<script lang="ts">
	import { game, player } from '$lib/services/Context';
</script>

<div class="w-full">
	<table class="table table-compact w-full">
		<caption>My Planets</caption>
		<thead>
			<th>Name</th>
			<th>Population</th>
		</thead>
		<tbody>
			{#if $player.planets?.length}
				{#each $player.planets as planet}
					<tr class="hover">
						<td>{planet.name}</td><td
							>{planet.cargo.colonists ? planet.cargo.colonists * 100 : ''}</td
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
		</thead>
		<tbody>
			{#if $player.planetIntels?.length}
				{#each $player.planetIntels as planet}
					<tr class="hover">
						<td>{planet.name}</td><td
							>{planet.hab?.grav
								? `${planet.hab.grav}, ${planet.hab.temp}, ${planet.hab.rad}`
								: ''}</td
						></tr
					>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
