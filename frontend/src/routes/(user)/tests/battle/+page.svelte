<script lang="ts">
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import Popup from '$lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { battleClient } from '$lib/services/connect';
	import { Universe } from '$lib/services/Universe';
	import { PlayerUniverseSchema, type BattleRecord } from '$lib/types/cs-proto';
	import { CommandedPlayer } from '$lib/types/Player';
	import { create } from '@bufbuild/protobuf';
	import { onMount } from 'svelte';

	let player: CommandedPlayer = $state(new CommandedPlayer());
	let universe: Universe = $state(new Universe());
	let battle: BattleRecord | undefined = $state();

	onMount(async () => {
		const resp = await battleClient.getTestBattle({});
		player = new CommandedPlayer(resp.player);
		battle = resp.battle;
		universe.setData(
			player.num,
			create(PlayerUniverseSchema, {
				designs: resp.designs ?? [],
				fleets: resp.fleets ?? []
			})
		);
	});
</script>

<h1 class="text-xl">Battle</h1>
{#if battle}
	<BattleView battleRecord={battle} playerFinder={universe} designFinder={universe} />
{/if}
<Tooltip />
<Popup />
