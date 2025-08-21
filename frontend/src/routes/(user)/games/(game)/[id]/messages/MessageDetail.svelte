<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet, MineralPacket, Planet, PlayerMessage } from '$lib/types/cs-proto';
	import { MapObjectType, PlayerMessageTargetType, PlayerMessageType } from '$lib/types/cs-proto';
	import { getMapObjectTarget } from '$lib/types/Message';
	import BattleMessageDetail from './BattleMessageDetail.svelte';
	import FleetMessageDetail from './FleetMessageDetail.svelte';
	import MineralPacketMessageDetail from './MineralPacketMessageDetail.svelte';
	import MysteryTraderMessageDetail from './MysteryTraderMessageDetail.svelte';
	import PlanetMessageDetail from './PlanetMessageDetail.svelte';
	import PlayerMessageDetail from './PlayerMessageDetail.svelte';

	const { universe } = getGameContext();

	let { message }: { message: PlayerMessage } = $props();

	let target = $derived($universe.getMapObject(getMapObjectTarget(message)));
	let owner = $derived(
		target && target.mapObject?.playerNum
			? $universe.getPlayerIntel(target.mapObject?.playerNum)
			: undefined
	);
	let planet = $derived(
		target?.mapObject?.type === MapObjectType.PLANET ? (target as Planet) : undefined
	);
	let fleet = $derived(
		target?.mapObject?.type === MapObjectType.FLEET ? (target as Fleet) : undefined
	);
	let mineralPacket = $derived(
		target?.mapObject?.type === MapObjectType.MINERAL_PACKET ? (target as MineralPacket) : undefined
	);
</script>

{#if message.type === PlayerMessageType.BATTLE || message.type === PlayerMessageType.BATTLE_ALLY}
	<BattleMessageDetail {message} />
{:else if planet}
	<PlanetMessageDetail {message} {planet} {owner} />
{:else if message.target?.targetType === PlayerMessageTargetType.MYSTERY_TRADER}
	<MysteryTraderMessageDetail {message} />
{:else if mineralPacket && owner}
	<MineralPacketMessageDetail {message} {mineralPacket} {owner} />
{:else if message.target?.targetType === PlayerMessageTargetType.FLEET || fleet}
	<FleetMessageDetail {message} />
{:else}
	<PlayerMessageDetail {message} />
{/if}
