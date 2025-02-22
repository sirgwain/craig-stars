<script lang="ts">
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import Popup from '$lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import type {
		AnyShipDesign,
		DesignFinder,
		PlayerFinder,
		PlayerUniverse
	} from '$lib/services/Universe';
	import type { BattleRecord, ShipDesign } from '$lib/types/cs';
	import { type Player, type PlayerIntel } from '$lib/types/cs';
	import { CommandedPlayer } from '$lib/types/Player';
	import { onMount } from 'svelte';

	type TestBattlePlayerResponse = Player & PlayerUniverse;
	let player: TestBattlePlayerResponse | undefined = $state();
	let battle: BattleRecord | undefined = $state();

	class TestPlayerFinder implements PlayerFinder {
		constructor(private player: TestBattlePlayerResponse) {}
		getPlayerIntel(num: number): PlayerIntel | undefined {
			return this.player.playerIntels?.find((p) => p.num == num);
		}
		getPlayerPluralName(playerNum: number | undefined): string {
			return (
				this.player.playerIntels?.find((p) => p.num == playerNum)?.name ?? 'Player ' + playerNum
			);
		}
		getPlayerColor(playerNum: number | undefined): string {
			return this.player.playerIntels?.find((p) => p.num == playerNum)?.color ?? '#FFFFFF';
		}
		getPlayerName(playerNum: number | undefined): string {
			return this.player.playerIntels?.find((p) => p.num == playerNum)?.name ?? '';
		}
	}

	class TestDesignFinder implements DesignFinder {
		constructor(private player: TestBattlePlayerResponse) {}

		getDesign(playerNum: number, num: number): AnyShipDesign | undefined {
			return (
				this.player.designs.find((d) => d && d.playerNum === playerNum && d.num === num) ??
				this.player.shipDesignIntels?.find((d) => d.playerNum === playerNum && d.num === num)
			);
		}
		getMyDesign(num: number | undefined): ShipDesign | undefined {
			return this.player.designs.find((d) => d && d.playerNum === this.player.num && d.num === num);
		}
	}

	let playerFinder: PlayerFinder | undefined = $state();
	let designFinder: DesignFinder | undefined = $state();

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
		player = new CommandedPlayer(
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
