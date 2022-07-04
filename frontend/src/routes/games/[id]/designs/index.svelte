<script lang="ts">
	import { game, player } from '$lib/services/Context';
</script>

<div class="w-full">
	<table class="table table-compact w-full">
		<thead>
			<th>Name</th>
			<th>Hull</th>
		</thead>
		<tbody>
			{#if $player.designs?.length}
				{#each $player.designs as design}
					<tr class="hover"> <td>{design.name}</td><td>{design.hull}</td></tr>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
