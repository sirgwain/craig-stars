<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';

	const { universe } = getGameContext();

	let players = $derived($universe.playerIntels);
</script>

<!-- names run up the page on phones so long race names can't overlap, and the table scrolls sideways if it has to -->
<div class="overflow-x-auto">
	<table class="table table-zebra table-pin-cols mx-auto w-auto">
		<thead>
			<tr>
				<th></th>
				{#each players as player (player.num)}
					<td class="px-2 align-bottom text-center sm:px-4">
						<div
							class="mx-auto max-h-28 rotate-180 truncate [writing-mode:vertical-rl] sm:max-h-none sm:max-w-28 sm:rotate-0 sm:[writing-mode:horizontal-tb]"
							title={$universe.getPlayerPluralName(player.num)}
						>
							{$universe.getPlayerPluralName(player.num)}
						</div>
					</td>
				{/each}
			</tr>
		</thead>
		<tbody>
			<tr>
				<th class="font-normal whitespace-nowrap">Planets</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.planets ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Starbases</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.starbases ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Unarmed Ships</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.unarmedShips ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Escort Ships</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.escortShips ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Capital Ships</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.capitalShips ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Tech Levels</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.techLevels ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Resources</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.resources ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Score</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.score ?? ''}</td
					>
				{/each}
			</tr>
			<tr>
				<th class="font-normal whitespace-nowrap">Rank</th>
				{#each players as player (player.num)}
					<td class="px-2 text-center sm:px-4"
						>{$universe.getPlayerScore(player.num)?.rank ?? ''}</td
					>
				{/each}
			</tr>
		</tbody>
	</table>
</div>
