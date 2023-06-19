<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	const { game, player, universe } = getGameContext();
</script>

{#if $game}
	<div class="w-full">
		You are playing as the {$player.race.pluralName} with PRT {$player.race.prt}
	</div>
{/if}
