<script lang="ts">
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import type { BattleRecord } from '$lib/types/Battle';
	import { Player, type PlayerResponse } from '$lib/types/Player';
	import { onMount } from 'svelte';

	let player: Player | undefined;
	let battle: BattleRecord | undefined;

	onMount(async () => {
		const response = await fetch(`/api/battles/test`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		const json = (await response.json()) as { player: PlayerResponse; battle: BattleRecord };
		player = new Player(0, json.player.num, json.player);
		battle = json.battle;
	});
</script>

<h1 class="text-xl">Battle</h1>
{#if player && battle}
	<BattleView {player} battleRecord={battle} />
{/if}
