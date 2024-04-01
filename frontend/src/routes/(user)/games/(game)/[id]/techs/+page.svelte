<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { levelsAbove } from '$lib/types/TechLevel';

	const { game, player, universe } = getGameContext();

	const newTechs = $game.techs.techs.filter(
		(t) => $player.hasTech(t) && levelsAbove(t.requirements, $player.techLevels) == 0
	);
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Techs</li>
	</svelte:fragment>
</Breadcrumb>

<Techs techStore={$game.rules.techs} player={$player} />
