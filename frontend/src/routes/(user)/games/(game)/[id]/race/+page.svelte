<script lang="ts">
	import { game } from '$lib/services/Context';
</script>

{#if $game}
	<div class="w-full">
		You are playing as the {$game.player.race?.pluralName}
	</div>
{/if}
