<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet, MineralPacket, Planet } from '$lib/types/cs';
	import {
		MapObjectTypeFleet,
		MapObjectTypeMineralPacket,
		MapObjectTypeMysteryTrader,
		MapObjectTypePlanet,
		PlayerMessageBattle,
		PlayerMessageBattleAlly,
		type PlayerMessage
	} from '$lib/types/cs';
	import BattleMessageDetail from './BattleMessageDetail.svelte';
	import FleetMessageDetail from './FleetMessageDetail.svelte';
	import MineralPacketMessageDetail from './MineralPacketMessageDetail.svelte';
	import MysteryTraderMessageDetail from './MysteryTraderMessageDetail.svelte';
	import PlanetMessageDetail from './PlanetMessageDetail.svelte';
	import PlayerMessageDetail from './PlayerMessageDetail.svelte';

	const { universe } = getGameContext();

	let { message }: { message: PlayerMessage } = $props();

	let target = $derived($universe.getMapObject(message));
	let owner = $derived(
		target && target.playerNum ? $universe.getPlayerIntel(target.playerNum) : undefined
	);
	let planet = $derived(target?.type == MapObjectTypePlanet ? (target as Planet) : undefined);
	let fleet = $derived(target?.type == MapObjectTypeFleet ? (target as Fleet) : undefined);
	let mineralPacket = $derived(
		target?.type == MapObjectTypeMineralPacket ? (target as MineralPacket) : undefined
	);
</script>

{#if message.type === PlayerMessageBattle || message.type === PlayerMessageBattleAlly}
	<BattleMessageDetail {message} />
{:else if planet}
	<PlanetMessageDetail {message} {planet} {owner} />
{:else if message.targetType === MapObjectTypeMysteryTrader}
	<MysteryTraderMessageDetail {message} />
{:else if mineralPacket && owner}
	<MineralPacketMessageDetail {message} {mineralPacket} {owner} />
{:else if message.targetType === MapObjectTypeFleet || fleet}
	<FleetMessageDetail {message} />
{:else}
	<PlayerMessageDetail {message} />
{/if}
