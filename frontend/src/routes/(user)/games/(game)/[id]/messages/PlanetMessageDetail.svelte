<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { Planet } from '$lib/types/Planet';
	import type { PlayerIntel } from '$lib/types/Player';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let planet: Planet;
	export let owner: PlayerIntel | undefined;
</script>

{#if message.type === MessageType.HomePlanet}
	Your home planet is {planet.name}. Your people are ready to leave the nest and explore the
	universe. Good luck.
{:else if message.type === MessageType.BuiltMineralAlchemy}
	Your scientists on {planet.name} have transmuted common materials into {message.spec.amount ??
		0}kT each of Ironium, Boranium and Germanium.
{:else if message.type === MessageType.PlanetDiscovery}
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
{:else if message.type === MessageType.PlanetDiedOff}
	{#if $player.race.spec?.livesOnStarbases}
		All of your colonists orbiting {planet.name} have died off. Your starbase has been lost and you no
		longer control the planet.
	{:else}
		All of your colonists on {planet.name} have died off. You no longer control the planet.
	{/if}
{:else}
	{message.text}
{/if}
