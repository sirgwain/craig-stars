<script lang="ts">
	import type { DesignFinder, PlayerFinder, Universe } from '$lib/services/Universe';
	import { Battle, type BattleRecord } from '$lib/types/Battle';
	import { setContext } from 'svelte/internal';
	import BattleBoard from './BattleBoard.svelte';
	import { designFinderKey, playerFinderKey } from '$lib/services/Contexts';

	export let designFinder: DesignFinder;
	export let playerFinder: PlayerFinder;
	export let battleRecord: BattleRecord;

	setContext<DesignFinder>(designFinderKey, designFinder);
	setContext<PlayerFinder>(playerFinderKey, playerFinder);

	$: battle = new Battle(battleRecord.num, battleRecord.position, battleRecord);
</script>

<BattleBoard {battle} />
