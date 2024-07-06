<script lang="ts">
	import { designFinderKey, playerFinderKey } from '$lib/services/GameContext';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import { Battle, type BattleRecord } from '$lib/types/Battle';
	import { setContext } from 'svelte';
	import BattleBoard from './BattleBoard.svelte';

	export let designFinder: DesignFinder;
	export let playerFinder: PlayerFinder;
	export let battleRecord: BattleRecord;

	setContext<DesignFinder>(designFinderKey, designFinder);
	setContext<PlayerFinder>(playerFinderKey, playerFinder);

	$: battle = new Battle(battleRecord.num, battleRecord.position, designFinder, battleRecord);
</script>

<BattleBoard {battle} />
