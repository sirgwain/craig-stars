<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { totalMinerals } from '$lib/types/Cost';
	import { absSum } from '$lib/types/Hab';
	import { UnlimitedSpaceDock } from '$lib/types/Consts';
	import {
		CometSize,
		PlayerMessageType,
		type Planet,
		type PlayerIntel,
		type PlayerMessage
	} from '$lib/types/cs-proto';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		message: PlayerMessage;
		planet: Planet;
		owner: PlayerIntel | undefined;
	};

	let { message, planet, owner }: Props = $props();

	let growthRate = $derived($player.race.growthRate * ($player.race.spec.growthFactor ?? 0));
</script>

{#if message.text}
	{message.text}
{:else if message.type === PlayerMessageType.PLANET_HOMEWORLD}
	Your home planet is {planet.mapObject?.name}. Your people are ready to leave the nest and explore
	the universe. Good luck.
{:else if message.type === PlayerMessageType.PLANET_BOMBED}
	{@const bombing = message.spec?.bombing}
	{#if bombing}
		{$universe.getPlayerPluralName(message.spec?.target?.targetPlayerNum)}
		{message.spec?.target?.targetName} has bombed your planet {planet.mapObject?.name}
		{#if message.spec?.bombing?.planetEmptied}
			killing off all its colonists.
		{:else}
			{#if bombing.colonistsKilled == 0 && bombing.minesDestroyed == 0 && bombing.factoriesDestroyed == 0 && bombing.defensesDestroyed == 0}
				doing no physical damage.
			{:else}
				killing {(bombing.colonistsKilled ?? 0).toLocaleString()} colonists, and destroying {bombing.minesDestroyed ??
					0} mines,
				{bombing.factoriesDestroyed ?? 0} factories and {bombing.defensesDestroyed ?? 0} defenses.
			{/if}

			{#if absSum(bombing.unterraformAmount) > 0}
				{#if bombing.numBombers ?? 0 > 1}
					The bombers have also retro-bombed the planet, undoing {absSum(
						bombing.unterraformAmount
					)}% of its terraforming.
				{:else}
					The bomber has also retro-bombed the planet, undoing {absSum(bombing.unterraformAmount)}%
					of its terraforming.
				{/if}
			{/if}
		{/if}
	{:else}
		<!-- Generic message, no bombing data (unexpected) -->
		Bombers have bombed planet ${planet.mapObject?.name}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BONUS_RESEARCH_ARTIFACT}
	Your colonists settling {planet.mapObject?.name} have found a strange artifact boosting your research
	in {message.spec?.field} by {message.spec?.amount} resources.
{:else if message.type === PlayerMessageType.PLANET_BUILT_DEFENSE}
	{#if message.spec?.amount === 1}
		You have built a defense outpost on {planet.mapObject?.name}.
	{:else}
		You have built {message.spec?.amount ?? 0} defense outposts on {planet.mapObject?.name}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_FACTORY}
	{#if message.spec?.amount === 1}
		You have built a factory on {planet.mapObject?.name}.
	{:else}
		You have built {message.spec?.amount ?? 0} factories on {planet.mapObject?.name}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_GENESIS_DEVICE}
	Strong fundamental forces have rebirthed {planet.mapObject?.name}. All planetary installations
	have been wiped clean as its environment shifts drastically and newfound minerals spring forth
	from the ground.
{:else if message.type === PlayerMessageType.PLANET_BUILT_MINERAL_ALCHEMY}
	Your scientists on {planet.mapObject?.name} have transmuted common materials into {message.spec
		?.amount ?? 0}kT each of Ironium, Boranium and Germanium.
{:else if message.type === PlayerMessageType.PLANET_BUILT_MINE}
	{#if message.spec?.amount === 1}
		You have built a mine on {planet.mapObject?.name}.
	{:else}
		You have built {message.spec?.amount ?? 0} mines on {planet.mapObject?.name}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_ITEM}
	You have attempted to build a {message.spec?.queueItemType?.toLowerCase()} on {planet.mapObject
		?.name}, but {planet.mapObject?.name}
	is unable to build any of these. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_MINERAL_PACKET_NO_MASS_DRIVER}
	You have attempted to build a mineral packet on {planet.mapObject?.name}, but you have no starbase
	equipped with a mass driver on this planet. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_MINERAL_PACKET_NO_TARGET}
	You have attempted to build a mineral packet on {planet.mapObject?.name}, but you have failed to
	specify a planet to target. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_BUILT_SCANNER}
	{planet.mapObject?.name} has built a new {message.spec?.name} planetary scanner.
{:else if message.type === PlayerMessageType.PLANET_BUILT_STARBASE}
	{planet.mapObject?.name} has built a new {message.spec?.name}.
	{#if planet.spec?.planetStarbaseSpec?.dockCapacity == UnlimitedSpaceDock}
		Ships of any size can now be built here.
	{:else if (planet.spec?.planetStarbaseSpec?.dockCapacity ?? 0) > 0}
		Ships up to {planet.spec?.planetStarbaseSpec?.dockCapacity}kT in mass can now be built at this
		facility.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_COMET_STRIKE}
	{#if message.spec?.comet?.size === CometSize.SMALL}
		A small comet has crashed into {planet.mapObject?.name} bringing new minerals and altering the planet's
		environment.
	{:else if message.spec?.comet?.size === CometSize.MEDIUM}
		A medium-sized comet has crashed into {planet.mapObject?.name} bringing a significant quantity of
		minerals and significantly altering the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.LARGE}
		A large comet has crashed into {planet.mapObject?.name} bringing a wide variety of new minerals and
		drastically altering the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.HUGE}
		A huge comet has crashed into {planet.mapObject?.name} embedding vast quantities of minerals in the
		planet and radically altering its environment.
	{:else}
		A comet has crashed into {planet.mapObject?.name} bringing new minerals and altering the planet's
		environment.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_COMET_STRIKE_MY_PLANET}
	{#if message.spec?.comet?.size === CometSize.SMALL}
		A small comet has crashed into your planet {planet.mapObject?.name}, killing {(
			message.spec?.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet brought additional minerals and has slightly altered
		the planet's habitat.
	{:else if message.spec?.comet?.size === CometSize.MEDIUM}
		A medium-sized comet has crashed into your planet {planet.mapObject?.name}, killing {(
			message.spec?.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet brought additional minerals and has altered the
		planet's environment.
	{:else if message.spec?.comet?.size === CometSize.LARGE}
		A large comet has crashed into your planet {planet.mapObject?.name}, killing {(
			message.spec?.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet brought significant quantities of minerals and has
		greatly altered the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.HUGE}
		A huge comet has crashed into your planet {planet.mapObject?.name}, killing {(
			message.spec?.comet?.colonistsKilled ?? 0
		).toLocaleString()} of your colonists. The comet has embedded vast stores of minerals and drastically
		altered the planet's environment.
	{:else}
		A comet has crashed into {planet.mapObject?.name} bringing new minerals and altering the planet's
		environment.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_DIED_OFF}
	{#if $player.race.spec.livesOnStarbases}
		All of your colonists orbiting {planet.mapObject?.name} have died off. Your starbase has been lost
		and you no longer control the planet.
	{:else}
		All of your colonists on {planet.mapObject?.name} have died off. You no longer control the planet.
	{/if}
{:else if [PlayerMessageType.PLANET_DISCOVERY, PlayerMessageType.PLANET_DISCOVERY_HABITABLE, PlayerMessageType.PLANET_DISCOVERY_TERRAFORMABLE, PlayerMessageType.PLANET_DISCOVERY_UNINHABITABLE].indexOf(message.type) != -1}
	{#if owner}
		You have found a planet occupied by someone else. {planet.mapObject?.name} is currently owned by
		the {owner.racePluralName}.
	{:else if $player.race.spec.instaforming && ((planet.spec?.terraformedHabitability && planet.spec?.terraformedHabitability > 0) || (planet.spec?.habitability && planet.spec?.habitability > 0))}
		You have found a new habitable planet. Your colonists will grow by up to {Math.max(
			1,
			((planet.spec?.terraformedHabitability ?? planet.spec?.habitability ?? 0) * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planet.mapObject?.name}.
	{:else if planet.spec?.habitability && planet.spec?.habitability > 0}
		You have found a new habitable planet. Your colonists will grow by up to {Math.max(
			1,
			(planet.spec?.habitability * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planet.mapObject?.name}.
	{:else if planet.spec?.terraformedHabitability && planet.spec?.terraformedHabitability > 0}
		You have found a new planet which you have the ability to make habitable. With terraforming,
		your colonists will grow by up to {Math.max(
			1,
			(planet.spec?.terraformedHabitability * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planet.mapObject?.name}.
	{:else}
		You have found a new planet which unfortunately is not habitable by you. {Math.max(
			1,
			-(planet.spec?.habitability ?? 0) / 10
		).toFixed(2)}% of your colonists will die per year if you colonize {planet.mapObject?.name}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_POPULATION_DECREASED}
	{#if message.spec?.amount === message.spec?.prevAmount}
		Your colonists on {planet.mapObject?.name} are suffering under the planet's hostile conditions but
		for now they are surviving.
	{:else}
		The population on {planet.mapObject?.name} has decreased from {(
			message.spec?.prevAmount ?? 0
		).toLocaleString()} to {(message.spec?.amount ?? 0).toLocaleString()}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_INSTAFORM}
	Your race has instantly terraformed {planet.mapObject?.name} up to optimal conditions. Its value is
	now {planet.spec?.habitability ?? 0}%.
{:else if message.type === PlayerMessageType.FLEET_INVADED_PLANET}
	{@const invasion = message.spec?.invasion}
	{#if invasion}
		{#if invasion.successful}
			Your troops beaming down from {invasion.fleetName ?? 'multiple fleets'} have successfully wrested
			{planet.mapObject?.name}
			from {$universe.getPlayerName(invasion.defenderPlayerNum)} control, killing off all their colonists
			with only {invasion.attackersKilled} causalties.
		{:else}
			Your troops beaming down from {invasion.fleetName ?? 'multiple fleets'} tried to invade {planet
				.mapObject?.name}, but all of them were massacred by the {$universe.getPlayerName(
				invasion.defenderPlayerNum
			)}. Your valiant fighters managed to kill {invasion.defendersKilled} of their colonists in return.
		{/if}
	{:else}
		{planet.mapObject?.name} was invaded, but your spies no nothing of the outcome.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_INVADED}
	{@const invasion = message.spec?.invasion}
	{#if invasion}
		{#if invasion.successful}
			{$universe.getPlayerName(invasion.attackerPlayerNum)}'s {invasion.fleetName ??
				'multiple fleets'} have successfully invaded {planet.mapObject?.name} and wrested it from your
			control. Your colonists managed to defeat {invasion.attackersKilled} of their invaders before being
			overrun. Your troops beaming down from {invasion.fleetName ?? 'multiple fleets'} have successfully
			wrested
		{:else}
			{$universe.getPlayerName(invasion.attackerPlayerNum)}'s {invasion.fleetName ??
				'multiple fleets'} tried to invade {planet.mapObject?.name}, but your troops were able to
			fend them off. You lost {invasion.defendersKilled} colonists in the process.
		{/if}
	{:else}
		{planet.mapObject?.name} was invaded, but your spies no nothing of the outcome.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_POPULATION_DECREASED_OVERCROWDING}
	The population on {planet.mapObject?.name} has decreased by {(-(
		message.spec?.amount ?? 0
	)).toLocaleString()}
	colonists due to overcrowding.
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_INVASION}
	Your colonists invading {planet.mapObject?.name} have picked through the defenders' remains looking
	for technology. In the process you have gained a level in {message.spec?.field}.
{:else if message.type === PlayerMessageType.PLAYER_ACQUIRABLE_PART_GAINED_BATTLE}
	Your people have picked through the wreckage from the battle at {planet.mapObject?.name} and have learned
	how to build {message.spec?.techGained}.
{:else if message.type === PlayerMessageType.FLEET_SCRAPPED}
	{#if planet.spec?.planetStarbaseSpec?.hasStarbase}
		{message.spec?.target?.targetName} has been dismantled for {totalMinerals(message.spec?.cost)}kT
		of minerals at the starbase orbiting {planet.mapObject?.name}.
	{:else}
		{message.spec?.target?.targetName} has been dismantled for {totalMinerals(message.spec?.cost)}kT
		of minerals which have been deposited on {planet.mapObject?.name}.
	{/if}
	{#if message.spec?.cost?.resources}
		&nbsp;Ultimate Recycling has also made {message.spec?.cost?.resources} resources available for immediate
		use (less if other ships were scrapped here this year).
	{/if}
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_SCRAP_FLEET}
	In the process of {message.spec?.name} being scrapped above {planet.mapObject?.name}, you have
	gained a level in {message.spec?.field}.
{:else if message.type === PlayerMessageType.PLAYER_ACQUIRABLE_PART_GAINED_SCRAP_FLEET}
	In the process of {message.spec?.name} being scrapped above {planet.mapObject?.name}, you have
	learned to build {message.spec?.techGained}.
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_BATTLE}
	Wreckage from the battle that occurred in orbit of {planet.mapObject?.name} has boosted your research
	in {message.spec?.field} by 1 level.
{:else if message.type === PlayerMessageType.FLEET_BUILT}
	<!-- TODO: remove this at some point. These are now fleet messages, not planet messages, but keeping this here so old savdes still target correctly -->
	Your starbase at {planet.mapObject?.name} has built {message.spec?.amount ?? 'a'} new {message
		.spec?.name}s.
{:else}
	<FallbackMessageDetail {message} />
{/if}
