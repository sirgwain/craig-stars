<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getOurDead, getOurShips, getTheirDead } from '$lib/types/Battle';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Planet } from '$lib/types/Planet';
	import BattleMessageDetail from './BattleMessageDetail.svelte';
	import FleetMessageDetail from './FleetMessageDetail.svelte';
	import PlanetMessageDetail from './PlanetMessageDetail.svelte';
	import PlayerMessageDetail from './PlayerMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;

	$: target = $universe.getMapObject(message);
	$: specTarget = $universe.getMapObject(message.spec);
	$: owner = target && target.playerNum ? $universe.getPlayerIntel(target.playerNum) : undefined;
	$: planet = target?.type == MapObjectType.Planet ? (target as Planet) : undefined;
	$: fleet = target?.type == MapObjectType.Fleet ? (target as Fleet) : undefined;
	$: mineralPacket =
		target?.type == MapObjectType.MineralPacket ? (target as MineralPacket) : undefined;
	$: mineField = target?.type == MapObjectType.MineField ? (target as MineField) : undefined;
</script>

{#if message.type == MessageType.Battle || message.type == MessageType.BattleAlly}
	<BattleMessageDetail {message} />
{:else if planet}
	<PlanetMessageDetail {message} {planet} {owner} />
{:else if message.targetType == MapObjectType.Fleet || fleet}
	<FleetMessageDetail {message} />
{:else}
	<PlayerMessageDetail {message} />
{/if}
