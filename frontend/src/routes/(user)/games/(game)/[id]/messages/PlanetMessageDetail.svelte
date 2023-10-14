<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { MessageType, type Message, CometSize } from '$lib/types/Message';
	import type { Planet, getQueueItemShortName } from '$lib/types/Planet';
	import type { PlayerIntel } from '$lib/types/Player';
	import { startCase } from 'lodash-es';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';
	import { totalMinerals } from '$lib/types/Cost';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let planet: Planet;
	export let owner: PlayerIntel | undefined;

	$: growthRate = $player.race.growthRate * ($player.race.spec?.growthFactor ?? 0);
</script>

{#if message.type === MessageType.HomePlanet}
	Your home planet is {planet.name}. Your people are ready to leave the nest and explore the
	universe. Good luck.
{:else if message.type === MessageType.BuildInvalidItem}
	You have attempted to build {startCase(message.spec.queueItemType)} on {planet.name}, but {planet.name}
	is unable to build any of these.
{:else if message.type === MessageType.BuildMineralPacketNoMassDriver}
	You have attempted to build a mineral packet on {planet.name}, but you have no Starbase equipped
	with a mass driver on this planet. Production for this planet has been cancelled.
{:else if message.type === MessageType.BuildMineralPacketNoTarget}
	You have attempted to build a mineral packet on {planet.name}, but you have not specified a
	target. The minerals have been returned to the planet and production has been cancelled.
{:else if message.type === MessageType.BuiltMineralAlchemy}
	Your scientists on {planet.name} have transmuted common materials into {message.spec.amount ??
		0}kT each of Ironium, Boranium and Germanium.
{:else if [MessageType.PlanetDiscovery, MessageType.PlanetDiscoveryHabitable, MessageType.PlanetDiscoveryTerraformable, MessageType.PlanetDiscoveryUninhabitable].indexOf(message.type) != -1}
	{#if owner}
		You have found a planet occupied by someone else. {planet.name} is currently owned by the {owner.racePluralName}
	{:else if planet.spec.habitability && planet.spec.habitability > 0}
		You have found a new habitable planet. Your colonists will grow by up {(
			(planet.spec.habitability * growthRate) /
			100
		).toFixed()}% per year if you colonize {planet.name}
	{:else if planet.spec.terraformedHabitability && planet.spec.terraformedHabitability > 0}
		You have found a new planet which you have the ability to make habitable. With terraforming,
		your colonists will grow by up to {(
			(planet.spec.terraformedHabitability * growthRate) /
			100
		).toFixed()}% per year if you colonize {planet.name}.
	{:else}
		You have found a new planet which unfortunately is not habitable by you. {-(
			((planet.spec.habitability ?? 0) * growthRate) /
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
{:else if message.type === MessageType.PlanetPopulationDecreased}
	The population on {planet.name} has decreased from {(
		message.spec.prevAmount ?? 0
	).toLocaleString()} to {(message.spec.amount ?? 0).toLocaleString()}.
{:else if message.type === MessageType.PlanetPopulationDecreasedOvercrowding}
	The population on {planet.name} has decreased by {(-(message.spec.amount ?? 0)).toLocaleString()} due
	to overcrowding.
{:else if message.type === MessageType.CometStrike}
	{#if message.spec.comet?.size == CometSize.Small}
		A small comet has crashed into {planet.name} bringing new minerals and altering the planet's environment.
	{:else if message.spec.comet?.size == CometSize.Medium}
		A medium-sized comet has crashed into {planet.name} bringing a significant quantity of minerals to
		the planet.
	{:else if message.spec.comet?.size == CometSize.Large}
		A large comet has crashed into {planet.name} bringing a wide variety of new minerals and drastically
		altering the planet's environment.
	{:else if message.spec.comet?.size == CometSize.Huge}
		A huge comet has crashed into {planet.name} embedding vast quantities of minerals in the planet and
		radically altering its environment.
	{:else}
		A comet has crashed into {planet.name} bringing new minerals and altering the planet's environment.
	{/if}
{:else if message.type === MessageType.CometStrikeMyPlanet}
	{#if message.spec.comet?.size == CometSize.Small}
		A small comet has crashed into your planet {planet.name}, killing {(
			message.spec.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet brought additional minerals and has slightly altered
		the planet's habitat.
	{:else if message.spec.comet?.size == CometSize.Medium}
		A medium-sized comet has crashed into your planet {planet.name}, killing {(
			message.spec.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet brought additional minerals and has altered the
		planet's environment.
	{:else if message.spec.comet?.size == CometSize.Large}
		A large comet has crashed into your planet {planet.name}, killing {(
			message.spec.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet brought significant quantities of minerals and has
		greatly altered the planet's environment.
	{:else if message.spec.comet?.size == CometSize.Huge}
		A huge comet has crashed into your planet {planet.name}, killing {(
			message.spec.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet has embedded vast stores of minerals and has drastically
		altered the planet's environment.
	{:else}
		A comet has crashed into {planet.name} bringing new minerals and altering the planet's environment.
	{/if}
{:else if message.type === MessageType.BonusResearchArtifact}
	Your colonists settling {planet.name} have found a strange artifact boosting your research in {message
		.spec.field} by {message.spec.amount} resources.
{:else if message.type === MessageType.TechLevelGainedInvasion}
	Your colonists invading {planet.name} have picked through the remains of the defenders looking for
	technology. In the process you have gained a level in {message.spec.field}.
{:else if message.type === MessageType.FleetScrapped}
	{#if planet.spec.hasStarbase}
		{message.spec.name} has been dismantled for {totalMinerals(message.spec.cost)}kT of minerals at
		the starbase orbiting {planet.name}.
	{:else}
		{message.spec.name} has been dismantled for {totalMinerals(message.spec.cost)}kT of minerals
		which have been deposited on {planet.name}.
	{/if}
	{#if message.spec.cost?.resources}
		&nbsp;Ultimate recycling has also made {message.spec.cost?.resources} resources available for immediate
		use (less if other ships were scrapped here this year).
	{/if}
{:else if message.type === MessageType.TechLevelGainedScrapFleet}
	In the process of {message.spec.name} being scrapped above {planet.name}, you have gained a level
	in {message.spec.field}
{:else if message.type === MessageType.TechLevelGainedBattle}
	Wreckage from the battle that occurred in orbit of {planet.name} has boosted your research in {message
		.spec.field}.
{:else}
	<FallbackMessageDetail {message} />
{/if}
