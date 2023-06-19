<script lang="ts">
	import { game } from '$lib/services/Context';
</script>

{#if $game}<div class="w-full">
		<table class="table table-compact w-full">
			<thead>
				<th>Text</th>
				<th>Target</th>
			</thead>
			<tbody>
				{#if $game.player.messages?.length}
					{#each $game.player.messages as message}
						<tr class="hover"> <td>{message.text}</td><td>{message.targetType}</td></tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
{/if}
