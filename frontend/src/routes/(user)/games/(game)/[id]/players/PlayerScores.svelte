<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';

	const { universe } = getGameContext();

	let players = $derived($universe.playerIntels);
</script>

<table class="table table-zebra table-fixed mx-auto w-full sm:w-auto">
	<thead>
		<tr>
			<th></th>
			{#each players as player (player.num)}
				<th class="h-20 w-20"
					><div class="origin-bottom -rotate-45 text-ellipsis w-20 overflow-hidden">
						{$universe.getPlayerPluralName(player.num)}
					</div></th
				>
			{/each}
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>Planets</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.planets ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Starbases</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.starbases ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Unarmed Ships</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.unarmedShips ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Escort Ships</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.escortShips ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Capital Ships</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.capitalShips ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Tech Levels</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.techLevels ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Resources</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.resources ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Score</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.score ?? ''}</td>
			{/each}
		</tr>
		<tr>
			<td>Rank</td>
			{#each players as player (player.num)}
				<td>{$universe.getPlayerScore(player.num)?.rank ?? ''}</td>
			{/each}
		</tr>
	</tbody>
</table>
