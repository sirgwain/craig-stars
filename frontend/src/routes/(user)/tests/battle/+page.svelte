<script lang="ts">
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import Popup from '$lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import type { BattleRecord } from '$lib/types/Battle';
	import {
		Player,
		type PlayerIntel,
		type PlayerResponse,
		type PlayerUniverse
	} from '$lib/types/Player';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { onMount } from 'svelte';

	type TestBattlePlayerResponse = PlayerResponse &
		PlayerUniverse & {
			playerIntels: PlayerIntel[];
			shipDesignIntels: ShipDesign[];
		};
	let player: TestBattlePlayerResponse | undefined;
	let battle: BattleRecord | undefined;

	class TestPlayerFinder implements PlayerFinder {
		constructor(private player: TestBattlePlayerResponse) {}
		getPlayerIntel(num: number): PlayerIntel | undefined {
			return this.player.playerIntels?.find((p) => p.num == num);
		}
		getPlayerName(playerNum: number | undefined): string {
			return (
				this.player.playerIntels?.find((p) => p.num == playerNum)?.name ?? 'Player ' + playerNum
			);
		}
		getPlayerColor(playerNum: number | undefined): string {
			return this.player.playerIntels?.find((p) => p.num == playerNum)?.color ?? '#FFFFFF';
		}
	}

	class TestDesignFinder implements DesignFinder {
		constructor(private player: TestBattlePlayerResponse) {}

		getDesign(playerNum: number, num: number): ShipDesign | undefined {
			return (
				this.player.designs.find((d) => d.playerNum === playerNum && d.num === num) ??
				this.player.shipDesignIntels.find((d) => d.playerNum === playerNum && d.num === num)
			);
		}
		getMyDesign(num: number | undefined): ShipDesign | undefined {
			return this.player.designs.find((d) => d.playerNum === this.player.num && d.num === num);
		}
	}

	let playerFinder: PlayerFinder | undefined;
	let designFinder: DesignFinder | undefined;

	onMount(async () => {
		const response = await fetch(`/api/battles/test`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		const json = (await response.json()) as {
			player: TestBattlePlayerResponse;
			battle: BattleRecord;
		};
		player = new Player(
			json.player as TestBattlePlayerResponse
		) as unknown as TestBattlePlayerResponse;
		battle = json.battle;
		playerFinder = new TestPlayerFinder(player);
		designFinder = new TestDesignFinder(player);
	});
</script>

<h1 class="text-xl">Battle</h1>
{#if player && battle && designFinder && playerFinder}
	<BattleView battleRecord={battle} {playerFinder} {designFinder} />
{/if}
<Tooltip />
<Popup />
