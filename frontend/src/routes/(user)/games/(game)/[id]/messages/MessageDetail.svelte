<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getOurDead, getOurShips, getTheirDead } from '$lib/types/Battle';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Planet } from '$lib/types/Planet';
	import FleetEngineStrainMessageDetail from './FleetEngineStrainMessageDetail.svelte';
	import FleetMessageDetail from './FleetMessageDetail.svelte';
	import PlanetMessageDetail from './PlanetMessageDetail.svelte';
	import PlayerMessageDetail from './PlayerMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;

	$: target = $universe.getMapObject(message);
	$: owner = target && target.playerNum ? $universe.getPlayerIntel(target.playerNum) : undefined;
	$: planet = target?.type == MapObjectType.Planet ? (target as Planet) : undefined;
	$: fleet = target?.type == MapObjectType.Fleet ? (target as Fleet) : undefined;
	$: mineralPacket =
		target?.type == MapObjectType.MineralPacket ? (target as MineralPacket) : undefined;
	$: mineField = target?.type == MapObjectType.MineField ? (target as MineField) : undefined;

	function getBattleMessage(message: Message): string {
		const stats = message.spec.battle;
		const battle = $universe.getBattle(message.battleNum);
		if (battle) {
			const location = $universe.getBattleLocation(battle) ?? 'unknown';
			let text = `A battle took place at ${location}.`;

			const allies = new Set($player.getAllies());

			const ours = getOurShips(battle, allies);
			const theirs = getTheirDead(battle, allies);
			const ourDead = getOurDead(battle, allies);
			const theirDead = getTheirDead(battle, allies);
			const oursLeft = ours - ourDead;
			const theirsLeft = theirs - theirDead;

			if (ourDead === 0) {
				text += ' None of your forces were destroyed.';
			} else if (ourDead === ours) {
				text += ' All of your forces were destroyed by enemy forces.';
			} else if (oursLeft === 1) {
				text += ` Only one of your ships survived.`;
			} else if (oursLeft > 1) {
				text += ` ${oursLeft} of your ships surived.`;
			}

			if (theirDead === 0) {
				text += ' None of the enemy forces were destroyed.';
			} else if (theirDead === theirs) {
				text += ' All enemy forces were destroyed.';
			} else if (theirsLeft === 1) {
				text += ` Only one enemy ship survived.`;
			} else if (theirsLeft > 1) {
				text += ` ${theirsLeft} enemy ships surived.`;
			}

			return text;
		} else {
			return `A battle took place at an unknown location`;
		}
	}

	export function getMessageText(message: Message): string {
		switch (message.type) {
			case MessageType.Battle:
				return getBattleMessage(message);
			default:
				return message.text;
		}
	}
</script>

{#if message.type == MessageType.Battle}
	{getBattleMessage(message)}
{:else if planet}
	<PlanetMessageDetail {message} {planet} {owner} />
{:else if target?.type == MapObjectType.Fleet || fleet}
	<FleetMessageDetail {message} {fleet} {owner} />
{:else if message.type === MessageType.FleetShipExceededSafeSpeed}
	<!-- The fleet could have been destroyed, in which case we won't have a fleet for this message so capture it here -->
	<FleetEngineStrainMessageDetail {message} />
{:else if message.type === MessageType.FleetTransferGiven}
	{message.spec.name} has successfully been given to {$universe.getPlayerName(
		message.spec.destPlayerNum
	)}
{:else if message.type === MessageType.FleetTransferReceivedRefused}
	{$universe.getPlayerName(message.spec.sourcePlayerNum)} has attempted to gift you {message.spec
		.name}, but you have refused their offer. If you wish to receive gifts from this player in the
	future, make sure you are allies.
{:else if message.type === MessageType.TechLevelGainedBattle}
	{@const battle = $universe.getBattle(message.battleNum)}
	Wreckage from the battle that occurred in at {battle?.position.x ?? 0}, {battle?.position.y ?? 0} has
	boosted your research in {message.spec.field}.
{:else}
	<PlayerMessageDetail {message} />
{/if}
