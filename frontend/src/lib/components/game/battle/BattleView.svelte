<script lang="ts">
	import { Battle, type BattleRecord } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import BattleBoard from './BattleBoard.svelte';

	export let player: Player;
	export let battleRecord: BattleRecord;

	$: battle = new Battle(battleRecord.num, battleRecord.position, battleRecord);
</script>

<BattleBoard {player} {battle} />
