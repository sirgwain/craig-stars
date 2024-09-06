<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { totalMinerals } from '$lib/types/Cost';
	import { absSum } from '$lib/types/Hab';
	import { CometSize, MessageType, type Message } from '$lib/types/Message';
	import type { Planet } from '$lib/types/Planet';
	import type { PlayerIntel } from '$lib/types/Player';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import { startCase } from 'lodash-es';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

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
		{$universe.getPlayerPluralName(message.spec.targetPlayerNum)}
		{message.spec.targetName} has bombed your planet {planet.name}
		{#if message.spec.bombing?.planetEmptied}
			killing off all its colonists.
		{:else}
			{#if bombing.colonistsKilled == 0 && bombing.minesDestroyed == 0 && bombing.factoriesDestroyed == 0 && bombing.defensesDestroyed == 0}
				doing no physical damage.
			{:else}
				killing {bombing.colonistsKilled ?? 0} colonists, and destroying {bombing.minesDestroyed ??
					0} mines,
				{bombing.factoriesDestroyed ?? 0} factories and {bombing.defensesDestroyed ?? 0} defenses.
			{/if}

			{#if absSum(bombing.unterraformAmount ?? {}) > 0}
				{#if bombing.numBombers ?? 0 > 1}
					The bombers have also retro-bombed the planet, undoing {absSum(
						bombing.unterraformAmount ?? {}
					)}% of its terraforming.
				{:else}
					The bomber has also retro-bombed the planet, undoing {absSum(
						bombing.unterraformAmount ?? {}
					)}% of its terraforming.
				{/if}
			{/if}
		{/if}
	{:else}
		<!-- Generic message, no bombing data (unexpected) -->
		Bombers have bombed planet ${planet.name}.
	{/if}
{:else if message.type === MessageType.PlanetBonusResearchArtifact}
	Your colonists settling {planet.name} have found a strange artifact boosting your research in {message
		.spec.field} by {message.spec.amount} resources.
{:else if message.type === MessageType.PlanetBuiltDefense}
	{#if message.spec.amount === 1}
		You have built a defense outpost on {planet.name}.
	{:else}
		You have built {message.spec.amount ?? 0} defense outposts on {planet.name}.
	{/if}
{:else if message.type === MessageType.PlanetBuiltFactory}
	{#if message.spec.amount === 1}
		You have built a factory on {planet.name}.
	{:else}
		You have built {message.spec.amount ?? 0} factories on {planet.name}.
	{/if}
{:else if message.type === MessageType.PlanetBuiltGensisDevice}
	Strong fundamental forces have rebirthed {planet.name}. All planetary installations have been wiped clean as its environment shifts drastically and newfound minerals spring forth from the ground.
{:else if message.type === MessageType.PlanetBuiltMineralAlchemy}
	Your scientists on {planet.name} have transmuted common materials into {message.spec.amount ??
		0}kT each of Ironium, Boranium and Germanium.
{:else if message.type === MessageType.PlanetBuiltMine}
	{#if message.spec.amount === 1}
		You have built a mine on {planet.name}.
	{:else}
		You have built {message.spec.amount ?? 0} mines on {planet.name}.
	{/if}
{:else if message.type === MessageType.PlanetBuiltInvalidItem}
	You have attempted to build a {message.spec.queueItemType?.toLowerCase()} on {planet.name}, but {planet.name}
	is unable to build any of these. The order has been canceled.
{:else if message.type === MessageType.PlanetBuiltInvalidMineralPacketNoMassDriver}
	You have attempted to build a mineral packet on {planet.name}, but you have no starbase equipped
	with a mass driver on this planet. The order has been canceled.
{:else if message.type === MessageType.PlanetBuiltInvalidMineralPacketNoTarget}
	You have attempted to build a mineral packet on {planet.name}, but you have failed to specify a
	planet to target. The order has been canceled.
{:else if message.type === MessageType.PlanetBuiltScanner}
	{planet.name} has built a new {message.spec.name} planetary scanner.
{:else if message.type === MessageType.PlanetBuiltStarbase}
	{planet.name} has built a new {message.spec.name}.
	{#if planet.spec.dockCapacity == UnlimitedSpaceDock}
		Ships of any size can now be built here.
	{:else if planet.spec.dockCapacity > 0}
		Ships up to {planet.spec.dockCapacity}kT in mass can now be built at this facility.
	{/if}
{:else if message.type === MessageType.PlanetCometStrike}
	{#if message.spec.comet?.size == CometSize.Small}
		A small comet has crashed into {planet.name} bringing new minerals and altering the planet's environment.
	{:else if message.spec.comet?.size == CometSize.Medium}
		A medium-sized comet has crashed into {planet.name} bringing a significant quantity of minerals and 
		significantly altering the planet's environment.
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
		).toLocaleString()} of your colonists. The comet has embedded vast stores of minerals and drastically
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
		You have found a planet occupied by someone else. {planet.name} is currently owned by the {owner.racePluralName}.
	{:else if $player.race.spec?.instaforming && ((planet.spec.terraformedHabitability && planet.spec.terraformedHabitability > 0) || 
	(planet.spec.habitability && planet.spec.habitability > 0))}
		You have found a new habitable planet. Your colonists will grow by up to {Math.max(
			1,
			((planet.spec.terraformedHabitability ?? (planet.spec.habitability ?? 0)) * growthRate) / 100
			).toFixed(2)}% per year if you colonize {planet.name}.
	{:else if planet.spec.habitability && planet.spec.habitability > 0}
		You have found a new habitable planet. Your colonists will grow by up to {Math.max(
			1,
			(planet.spec.habitability * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planet.name}.
	{:else if planet.spec.terraformedHabitability && planet.spec.terraformedHabitability > 0}
		You have found a new planet which you have the ability to make habitable. With terraforming,
		your colonists will grow by up to {Math.max(
			1,
			(planet.spec.terraformedHabitability * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planet.name}.
	{:else}
		You have found a new planet which unfortunately is not habitable by you. {Math.max(
			1,
			-(planet.spec.habitability ?? 0) / 10
		).toFixed(2)}% of your colonists will die per year if you colonize {planet.name}.
	{/if}
{:else if message.type === MessageType.PlanetPopulationDecreased}
	{#if message.spec.amount === message.spec.prevAmount}
		Your colonists on {planet.name} are suffering under the planet's hostile conditions but for now they
		are surviving.
	{:else}
		The population on {planet.name} has decreased from {(
			message.spec.prevAmount ?? 0
		).toLocaleString()} to {(message.spec.amount ?? 0).toLocaleString()}.
	{/if}
{:else if message.type === MessageType.PlanetPopulationDecreasedOvercrowding}
	The population on {planet.name} has decreased by {(-(message.spec.amount ?? 0)).toLocaleString()} colonists due
	to overcrowding.
{:else if message.type === MessageType.PlayerTechLevelGainedInvasion}
	Your colonists invading {planet.name} have picked through the defenders' remains looking for technology.
	In the process you have gained a level in {message.spec.field}.
{:else if message.type === MessageType.FleetScrapped}
	{#if planet.spec.hasStarbase}
		{message.spec.targetName} has been dismantled for {totalMinerals(message.spec.cost)}kT of
		minerals at the starbase orbiting {planet.name}.
	{:else}
		{message.spec.targetName} has been dismantled for {totalMinerals(message.spec.cost)}kT of
		minerals which have been deposited on {planet.name}.
	{/if}
	{#if message.spec.cost?.resources}
		&nbsp;Ultimate Recycling has also made {message.spec.cost?.resources} resources available for immediate
		use (less if other ships were scrapped here this year).
	{/if}
{:else if message.type === MessageType.PlayerTechLevelGainedScrapFleet}
	In the process of {message.targetName} being scrapped above {planet.name}, you have gained a level
	in {message.spec.field}.
{:else if message.type === MessageType.PlayerTechLevelGainedBattle}
	Wreckage from the battle that occurred in orbit of {planet.name} has boosted your research in {message
		.spec.field} by 1 level.
{:else if message.type === MessageType.FleetBuilt}
	<!-- TODO: remove this at some point. These are now fleet messages, not planet messages, but keeping this here so old savdes still target correctly -->
	Your starbase at {planet.name} has built {message.spec.amount ?? 'a'} new {message.spec.name}s.
{:else}
	<FallbackMessageDetail {message} />
{/if}
