<script lang="ts">
	import { game } from '$lib/services/Stores';
</script>

{#if $game}
	<div class="w-full">
		You are playing as the {$game.player.race?.pluralName}
	</div>
{/if}
