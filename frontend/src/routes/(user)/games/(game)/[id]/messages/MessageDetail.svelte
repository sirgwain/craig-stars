<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { MysteryTrader } from '$lib/types/MysteryTrader';
	import type { Planet } from '$lib/types/Planet';
	import BattleMessageDetail from './BattleMessageDetail.svelte';
	import FleetMessageDetail from './FleetMessageDetail.svelte';
	import MineralPacketMessageDetail from './MineralPacketMessageDetail.svelte';
	import PlanetMessageDetail from './PlanetMessageDetail.svelte';
	import PlayerMessageDetail from './PlayerMessageDetail.svelte';
	import MysteryTraderMessageDetail from './MysteryTraderMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;

	$: target = $universe.getMapObject(message);
	$: specTarget = $universe.getMapObject(message.spec);
	$: owner = target && target.playerNum ? $universe.getPlayerIntel(target.playerNum) : undefined;
	$: planet = target?.type == MapObjectType.Planet ? (target as Planet) : undefined;
	$: fleet = target?.type == MapObjectType.Fleet ? (target as Fleet) : undefined;
	$: mysteryTrader =
		target?.type == MapObjectType.MysteryTrader ? (target as MysteryTrader) : undefined;
	$: mineralPacket =
		target?.type == MapObjectType.MineralPacket ? (target as MineralPacket) : undefined;
	$: mineField = target?.type == MapObjectType.MineField ? (target as MineField) : undefined;
</script>

{#if message.type == MessageType.Battle || message.type == MessageType.BattleAlly}
	<BattleMessageDetail {message} />
{:else if planet}
	<PlanetMessageDetail {message} {planet} {owner} />
{:else if mysteryTrader}
	<MysteryTraderMessageDetail {message} {mysteryTrader} />
{:else if mineralPacket && owner}
	<MineralPacketMessageDetail {message} {mineralPacket} {owner} />
{:else if message.targetType == MapObjectType.Fleet || fleet}
	<FleetMessageDetail {message} />
{:else}
	<PlayerMessageDetail {message} />
{/if}
