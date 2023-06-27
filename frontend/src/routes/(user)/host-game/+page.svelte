<script lang="ts">
	import NewGame from '$lib/components/game/newgame/NewGame.svelte';
	import { getColor } from '$lib/components/game/newgame/playerColors';
	import { me } from '$lib/services/Stores';
	import { NewGamePlayerType } from '$lib/types/Game';
	const players = [
		{ type: NewGamePlayerType.Host, color: getColor(0) },
		{ type: NewGamePlayerType.Open, color: getColor(1) },
		{ type: NewGamePlayerType.AI, color: getColor(2) }
	];

	const name = $me?.username ? `${$me.username}'s game` : undefined;
</script>

<NewGame {players} {name} />
