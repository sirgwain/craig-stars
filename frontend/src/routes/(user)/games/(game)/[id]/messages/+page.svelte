<script lang="ts">
	import { game, player } from '$lib/services/Context';
</script>

{#if $player}<div class="w-full">
		<table class="table table-compact w-full">
			<thead>
				<th>Text</th>
				<th>Target</th>
			</thead>
			<tbody>
				{#if $player.messages?.length}
					{#each $player.messages as message}
						<tr class="hover"> <td>{message.text}</td><td>{message.targetType}</td></tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
{/if}
