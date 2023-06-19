<script lang="ts">
	import type { Universe } from '$lib/services/Universe';
	import { Battle, type BattleRecord } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import BattleBoard from './BattleBoard.svelte';

	export let universe: Universe;
	export let player: Player;
	export let battleRecord: BattleRecord;

	$: battle = new Battle(battleRecord.num, battleRecord.position, battleRecord);
</script>

<BattleBoard {universe} {player} {battle} />
