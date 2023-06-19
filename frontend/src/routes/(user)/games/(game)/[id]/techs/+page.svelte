<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import { getGameContext } from '$lib/services/Contexts';

	const { game, player, universe } = getGameContext();

</script>

{#if game}
	<div class="w-full">
		<Techs techStore={$game.rules.techs} player={$player} />
	</div>
{/if}