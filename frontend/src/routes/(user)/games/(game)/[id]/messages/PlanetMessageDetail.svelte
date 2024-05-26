<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MessageType, type Message, CometSize } from '$lib/types/Message';
	import type { Planet, getQueueItemShortName } from '$lib/types/Planet';
	import type { PlayerIntel } from '$lib/types/Player';
	import { startCase } from 'lodash-es';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';
	import { totalMinerals } from '$lib/types/Cost';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import { absSum } from '$lib/types/Hab';
	import { text } from '@sveltejs/kit';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let planet: Planet;
	export let owner: PlayerIntel | undefined;

	$: growthRate = $player.race.growthRate * ($player.race.spec?.growthFactor ?? 0);
</script>

{#if message.text}
	{message.text}
{:else if message.type === MessageType.PlanetHomeworld}
	Your home planet is {planet.name}. Your people are ready to leave the nest and explore the
	universe. Good luck.
{:else if message.type === MessageType.PlanetBombed}
	{@const bombing = message.spec.bombing}
	{#if bombing}
		{#if bombing.numBombers == 1}
			{$universe.getPlayerName(message.spec.targetPlayerNum)}
			{message.spec.targetName} has bombed your planet {planet.name}
			{#if message.spec.bombing?.planetEmptied}
				killing off all colonists.
			{:else}
				killing {bombing.colonistsKilled ?? 0} colonists, and destroying {bombing.minesDestroyed ??
					0} mines,
				{bombing.factoriesDestroyed ?? 0} factories and {bombing.defensesDestroyed ?? 0} defenses.

				{#if absSum(bombing.unterraformAmount ?? {}) > 0}
					The bombers have retro-bombed the planet, undoing {absSum(
						bombing.unterraformAmount ?? {}
					)}% of its terraforming.
				{/if}
			{/if}
		{:else}
			{$universe.getPlayerName(message.spec.targetPlayerNum)}
			{message.spec.targetName} has bombed your planet {planet.name} killing off all colonists
		{/if}
	{:else}
		<!-- Generic message, no bombing data (unexpected) -->
		Bombers have bombed planet ${planet.name}.
	{/if}
{:else if message.type === MessageType.PlanetBonusResearchArtifact}
	Your colonists settling {planet.name} have found a strange artifact boosting your research in {message
		.spec.field} by {message.spec.amount} resources.
{:else if message.type === MessageType.PlanetBuiltDefense}
	You have built {message.spec.amount ?? 0} defense(s) on {planet.name}.
{:else if message.type === MessageType.PlanetBuiltFactory}
	You have built {message.spec.amount ?? 0} factory(s) on {planet.name}.
{:else if message.type === MessageType.PlanetBuiltMineralAlchemy}
	Your scientists on {planet.name} have transmuted common materials into {message.spec.amount ??
		0}kT each of Ironium, Boranium and Germanium.
{:else if message.type === MessageType.PlanetBuiltMine}
	You have built {message.spec.amount ?? 0} mine(s) on {planet.name}.
{:else if message.type === MessageType.PlanetBuiltInvalidItem}
	You have attempted to build {startCase(message.spec.queueItemType)} on {planet.name}, but {planet.name}
	is unable to build any of these.
{:else if message.type === MessageType.PlanetBuiltInvalidMineralPacketNoMassDriver}
	You have attempted to build a mineral packet on {planet.name}, but you have no Starbase equipped
	with a mass driver on this planet. Production for this planet has been cancelled.
{:else if message.type === MessageType.PlanetBuiltInvalidMineralPacketNoTarget}
	You have attempted to build a mineral packet on {planet.name}, but you have not specified a
	target. The minerals have been returned to the planet and production has been cancelled.
{:else if message.type === MessageType.PlanetBuiltScanner}
	{planet.name} has built a new {message.spec.name} planetary scanner.
{:else if message.type === MessageType.PlanetBuiltStarbase}
	{planet.name} has built a new {message.spec.name}.
	{#if planet.spec.dockCapacity == UnlimitedSpaceDock}
		Ships of any size can now be built here.
	{:else if planet.spec.dockCapacity > 0}
		Ships up to {planet.spec.dockCapacity}kT in total hull weight can now be built at this facility.
	{/if}
{:else if message.type === MessageType.PlanetCometStrike}
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
{:else if message.type === MessageType.PlanetCometStrikeMyPlanet}
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
{:else if message.type === MessageType.PlanetDiedOff}
	{#if $player.race.spec?.livesOnStarbases}
		All of your colonists orbiting {planet.name} have died off. Your starbase has been lost and you no
		longer control the planet.
	{:else}
		All of your colonists on {planet.name} have died off. You no longer control the planet.
	{/if}
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
{:else if message.type === MessageType.PlanetPopulationDecreased}
	The population on {planet.name} has decreased from {(
		message.spec.prevAmount ?? 0
	).toLocaleString()} to {(message.spec.amount ?? 0).toLocaleString()}.
{:else if message.type === MessageType.PlanetPopulationDecreasedOvercrowding}
	The population on {planet.name} has decreased by {(-(message.spec.amount ?? 0)).toLocaleString()} due
	to overcrowding.
{:else if message.type === MessageType.PlayerTechLevelGainedInvasion}
	Your colonists invading {planet.name} have picked through the remains of the defenders looking for
	technology. In the process you have gained a level in {message.spec.field}.
{:else if message.type === MessageType.FleetScrapped}
	{#if planet.spec.hasStarbase}
		{message.spec.targetName} has been dismantled for {totalMinerals(message.spec.cost)}kT of minerals at
		the starbase orbiting {planet.name}.
	{:else}
		{message.spec.targetName} has been dismantled for {totalMinerals(message.spec.cost)}kT of minerals
		which have been deposited on {planet.name}.
	{/if}
	{#if message.spec.cost?.resources}
		&nbsp;Ultimate recycling has also made {message.spec.cost?.resources} resources available for immediate
		use (less if other ships were scrapped here this year).
	{/if}
{:else if message.type === MessageType.PlayerTechLevelGainedScrapFleet}
	In the process of {message.targetName} being scrapped above {planet.name}, you have gained a level
	in {message.spec.field}
{:else if message.type === MessageType.PlayerTechLevelGainedBattle}
	Wreckage from the battle that occurred in orbit of {planet.name} has boosted your research in {message
		.spec.field}.
{:else if message.type === MessageType.FleetBuilt}
	<!-- TODO: remove this at some point. These are now fleet messages, not planet messages, but keeping this here so old savdes still target correctly -->
	Your starbase at {planet.name} has built {message.spec.amount ?? 'a'} new {message.spec.name}s.
{:else}
	<FallbackMessageDetail {message} />
{/if}
