<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { getOurDead, getOurShips, getTheirDead } from '$lib/types/Battle';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { MineField } from '$lib/types/MineField';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { Planet } from '$lib/types/Planet';

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
{:else if message.type === MessageType.PlanetDiscovery}
	{#if planet}
		{#if owner}
			You have found a planet occupied by someone else. {planet.name} is currently owned by the {owner.racePluralName}
		{:else if planet.spec.habitability && planet.spec.habitability > 0}
			You have found a new habitable planet. Your colonists will grow by up {(
				(planet.spec.habitability * $player.race.growthRate) /
				100
			).toFixed()}% per year if you colonize {planet.name}
		{:else if planet.spec.terraformedHabitability && planet.spec.terraformedHabitability > 0}
			You have found a new planet which you have the ability to make habitable. With terraforming,
			your colonists will grow by up to {(
				(planet.spec.terraformedHabitability * $player.race.growthRate) /
				100
			).toFixed()}% per year if you colonize {planet.name}.
		{:else}
			You have found a new planet which unfortunately is not habitable by you. {-(
				((planet.spec.habitability ?? 0) * $player.race.growthRate) /
				100
			).toFixed()}% of your colonists will die per year if you colonize {planet.name}
		{/if}
	{:else}
		You have discovered a new planet.
	{/if}
{:else if message.type === MessageType.PlanetDiedOff && target}
	{#if $player.race.spec?.livesOnStarbases}
		All of your colonists orbiting {target.name} have died off. Your starbase has been lost and you no
		longer control the planet.
	{:else}
		All of your colonists on {target.name} have died off. You no longer control the planet.
	{/if}
{:else}
	{message.text}
{/if}
