<script lang="ts">
	import { game, player } from '$lib/services/Context';
</script>

<div class="w-full">
	<table class="table table-compact w-full">
		<thead>
			<th>Name</th>
			<th>Tokens</th>
		</thead>
		<tbody>
			{#if $player.fleets?.length}
				{#each $player.fleets as design}
					<tr class="hover"> <td>{design.name}</td><td>{design.tokens.length}</td></tr>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
