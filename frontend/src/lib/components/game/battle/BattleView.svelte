<script lang="ts">
	import { designFinderKey, playerFinderKey } from '#lib/services/GameContext.js';
	import type { DesignFinder, PlayerFinder } from '#lib/services/Universe.js';
	import { Battle } from '#lib/types/Battle.js';
	import type { BattleRecord } from '#lib/types/cs-proto.js';
	import { setContext } from 'svelte';
	import BattleBoard from './BattleBoard.svelte';

	type Props = {
		designFinder: DesignFinder;
		playerFinder: PlayerFinder;
		battleRecord: BattleRecord;
	};

	let { designFinder, playerFinder, battleRecord }: Props = $props();

	setContext<DesignFinder>(designFinderKey, {
		getDesign: (playerNum, num) => designFinder.getDesign(playerNum, num),
		getMyDesign: (num) => designFinder.getMyDesign(num)
	});
	setContext<PlayerFinder>(playerFinderKey, {
		getPlayerIntel: (num) => playerFinder.getPlayerIntel(num),
		getPlayerName: (num) => playerFinder.getPlayerName(num),
		getPlayerPluralName: (num) => playerFinder.getPlayerPluralName(num),
		getPlayerColor: (num) => playerFinder.getPlayerColor(num)
	});

	let battle = $derived(
		new Battle(battleRecord.num, battleRecord.position, designFinder, battleRecord)
	);
</script>

<BattleBoard {battle} />
