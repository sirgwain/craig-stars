<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { getOurDead, getOurShips, getTheirDead } from '$lib/types/Battle';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Planet } from '$lib/types/Planet';
	import PlanetMessageDetail from './PlanetMessageDetail.svelte';

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
{:else}
	{message.text}
{/if}
