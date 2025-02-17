<script lang="ts">
	import NewGame from '$lib/components/game/newgame/NewGame.svelte';
	import { getColor } from '$lib/components/game/newgame/playerColors';
	import { me } from '$lib/services/Stores';
	import {
		AIDifficultyNone,
		AIDifficultyNormal,
		NewGamePlayerTypeAI,
		NewGamePlayerTypeHost,
		NewGamePlayerTypeOpen,
		type NewGamePlayer
	} from '$lib/types/cs';

	const players: NewGamePlayer[] = [
		{
			type: NewGamePlayerTypeHost,
			color: getColor(0),
			aiDifficulty: AIDifficultyNone,
			hullSetNum: 0
		},
		{
			type: NewGamePlayerTypeOpen,
			color: getColor(1),
			aiDifficulty: AIDifficultyNone,
			hullSetNum: 0
		},
		{
			type: NewGamePlayerTypeAI,
			color: getColor(2),
			aiDifficulty: AIDifficultyNormal,
			hullSetNum: 0
		}
	];

	const name = $me.username ? `${$me.username}'s game` : undefined;
</script>

<NewGame {players} {name} />
