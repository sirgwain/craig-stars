<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import PlayerScores from './PlayerScores.svelte';
	import PlayersStatus from './PlayersStatus.svelte';

	const { game, player, universe } = getGameContext();
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<ItemTitle>Players</ItemTitle>

	<PlayersStatus />
</div>
<div class="w-full">
	<ItemTitle>Score</ItemTitle>
	<PlayerScores />
</div>

<div class="mb-10"/>
