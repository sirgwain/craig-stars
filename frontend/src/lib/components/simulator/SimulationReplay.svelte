<script lang="ts">
	import { setContext } from 'svelte';
	import { gameKey } from '#lib/services/GameContext.js';
	import BattleView from '#lib/components/game/battle/BattleView.svelte';
	import { create } from '@bufbuild/protobuf';
	import { PlayerIntelSchema } from '#lib/types/cs-proto.js';
	import type { SimulationRun } from './scenario';
	let { run }: { run: SimulationRun } = $props();
	// Simulation player/design numbers are independent of the active game's numbers.
	setContext(gameKey, undefined);
	const finder = {
		getDesign: (p: number | undefined, n: number | undefined) =>
			run.response.designs.find((d) => d.playerNum === p && d.num === n),
		getMyDesign: (n: number | undefined) =>
			run.response.designs.find((d) => d.playerNum === 1 && d.num === n),
		getPlayerName: (n: number | undefined) => run.sides[(n ?? 1) - 1]?.race?.name ?? '',
		getPlayerPluralName: (n: number | undefined) => run.sides[(n ?? 1) - 1]?.race?.pluralName ?? '',
		getPlayerColor: (n: number | undefined) => (n === 1 ? '#0000FF' : '#FF0000'),
		getPlayerIntel: (n: number | undefined) =>
			create(PlayerIntelSchema, {
				num: n,
				racePluralName: run.sides[(n ?? 1) - 1]?.race?.pluralName
			})
	};
</script>

{#if run.response.record}<BattleView
		battleRecord={run.response.record}
		designFinder={finder}
		playerFinder={finder}
	/>{/if}
