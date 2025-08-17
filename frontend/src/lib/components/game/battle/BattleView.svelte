<script lang="ts">
	import type { BattleRecord } from '$lib/types/cs-proto';
	import { designFinderKey, playerFinderKey } from '$lib/services/GameContext';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import { Battle } from '$lib/types/Battle';
	import { setContext } from 'svelte';
	import BattleBoard from './BattleBoard.svelte';

	type Props = {
		designFinder: DesignFinder;
		playerFinder: PlayerFinder;
		battleRecord: BattleRecord;
	};

	let { designFinder, playerFinder, battleRecord }: Props = $props();

	setContext<DesignFinder>(designFinderKey, designFinder);
	setContext<PlayerFinder>(playerFinderKey, playerFinder);

	let battle = $derived(
		new Battle(battleRecord.num, battleRecord.position, designFinder, battleRecord)
	);
</script>

<BattleBoard {battle} />
