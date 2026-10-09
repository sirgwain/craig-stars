<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { UnlimitedSpaceDock } from '#lib/types/Consts.js';
	import { totalMinerals } from '#lib/types/Cost.js';
	import {
		CometSize,
		PlayerMessageType,
		ProductionQueueItemSchema,
		QueueItemType,
		TerraformHabType,
		type Planet,
		type PlayerIntel,
		type PlayerMessage
	} from '#lib/types/cs-proto.js';
	import { absSum, getHabValue, getTerraformHabValueString } from '#lib/types/Hab.js';
	import { getLongHabName } from '#lib/types/Tech.js';
	import { getFullName } from '#lib/types/QueueItemType.js';
	import { create } from '@bufbuild/protobuf';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		message: PlayerMessage;
		planet: Planet | undefined;
		owner: PlayerIntel | undefined;
	};

	let { message, planet, owner }: Props = $props();

	let planetName = $derived(planet?.mapObject?.name ?? message.target?.targetName ?? 'unknown');

	let growthRate = $derived($player.race.growthRate * $player.race.spec.growthFactor);
</script>

{#if message.text}
	{message.text}
{:else if message.type === PlayerMessageType.PLANET_COLONIZED}
	Your colonists are now in control of {planetName}.
{:else if message.type === PlayerMessageType.PLANET_PRODUCTION_QUEUE_EMPTY}
	The production queue on {planetName} is empty.
{:else if message.type === PlayerMessageType.PLANET_PRODUCTION_QUEUE_COMPLETE}
	{planetName} has completed its orders. The production queue is empty.
{:else if message.type === PlayerMessageType.PLANET_REMOTE_TERRAFORM}
	{@const spec = message.spec}
	{#if spec?.sourcePlayerNum === $player.num}
		Your fleet
	{:else}
		A {$universe.getPlayerPluralName(spec?.sourcePlayerNum)} fleet
	{/if}
	{spec?.mapObjectTarget?.targetName}
	{#if spec && absSum(spec.terraformAmount) > 0}
		{#if spec.amount !== spec.prevAmount}
			has {spec.amount > spec.prevAmount ? 'improved' : 'degraded'}
			{planetName} from a value of
			{spec.prevAmount}% to {spec.amount}%.
		{:else}
			is currently unable to {spec.amount2 < 0 ? 'degrade' : 'improve'} the value of {planetName}
			beyond {spec.amount}%.
		{/if}
		{#each [TerraformHabType.GRAV, TerraformHabType.TEMP, TerraformHabType.RAD] as habType (habType)}
			{@const change = getHabValue(spec.terraformAmount, habType - 1)}
			{#if change !== 0}
				{getLongHabName(habType)} has {change > 0 ? 'increased' : 'decreased'} by {Math.abs(
					change
				)}%.
			{/if}
		{/each}
	{:else}
		<!-- Preserve earlier remote terraforming reports containing an axis and raw habitat value. -->
		has remotely terraformed {planetName},
		{(spec?.amount ?? 0) > 0 ? 'increasing' : 'decreasing'} its {getLongHabName(spec?.habType ?? 0)}
		by {Math.abs(spec?.amount ?? 0)}% to {getTerraformHabValueString(
			spec?.habType ?? 0,
			spec?.amount2 ?? 0
		)}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_INVADE_INVALID_EMPTY}
	{message.spec?.mapObjectTarget?.targetName} has orders to beam colonists to {planetName}, but the
	planet is uninhabited. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_INVADE_INVALID_STARBASE}
	{message.spec?.mapObjectTarget?.targetName} has orders to invade {planetName}, but the planet is
	protected by a starbase. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_PACKET_LANDED}
	<!-- Amount is packet mineral mass in kT. -->
	Your mineral packet containing {message.spec?.amount ?? 0}kT of minerals has arrived at {planetName}.
{:else if message.type === PlayerMessageType.PLANET_PACKET_CAUGHT}
	Your mass accelerator at {planetName} has successfully captured a packet containing {message.spec
		?.amount ?? 0}kT of minerals.
{:else if message.type === PlayerMessageType.PLANET_PACKET_DAMAGE}
	{@const spec = message.spec}
	{#if spec?.mineralPacketDamage}
		{@const damage = spec.mineralPacketDamage}
		{#if spec.hasMassDriver}
			Your mass accelerator at {planetName} was partially successful at capturing a {spec.amount}kT
			mineral packet.
			{#if damage.defensesDestroyed === 0}
				Unable to completely slow the packet, {damage.killed.toLocaleString()} of your colonists were
				killed in the collision.
			{:else}
				Unfortunately, {damage.killed.toLocaleString()} of your colonists and {damage.defensesDestroyed}
				of your defenses were destroyed in the collision.
			{/if}
		{:else if spec.planetEmptied}
			{planetName} was annihilated by a mineral packet. All of your colonists were killed.
		{:else}
			{planetName} was bombarded with a {spec.amount}kT mineral packet.
			{#if damage.defensesDestroyed === 0}
				{damage.killed.toLocaleString()} of your colonists were killed in the collision.
			{:else}
				{damage.killed.toLocaleString()} of your colonists and {damage.defensesDestroyed} of your defenses
				were destroyed in the collision.
			{/if}
		{/if}
	{:else}
		<FallbackMessageDetail {message} />
	{/if}
{:else if [PlayerMessageType.PLANET_BUILT_TERRAFORM, PlayerMessageType.PLANET_PERMAFORM, PlayerMessageType.PLANET_PACKET_TERRAFORM, PlayerMessageType.PLANET_PACKET_PERMAFORM].includes(message.type)}
	<!-- Amount is the signed change; Amount2 is the resulting raw habitat value. -->
	{@const direction = (message.spec?.amount ?? 0) > 0 ? 'increased' : 'decreased'}
	{@const habName = getLongHabName(message.spec?.habType ?? 0)}
	{@const habValue = getTerraformHabValueString(
		message.spec?.habType ?? 0,
		message.spec?.amount2 ?? 0
	)}
	{#if message.type === PlayerMessageType.PLANET_BUILT_TERRAFORM}
		Your terraforming efforts on {planetName} have {direction} its {habName} to {habValue}.
	{:else if message.type === PlayerMessageType.PLANET_PERMAFORM}
		Your colonists have permanently {direction} the {habName} on {planetName} to {habValue}.
	{:else}
		Your mineral packet hitting {planetName} has {message.type ===
		PlayerMessageType.PLANET_PACKET_PERMAFORM
			? 'permanently '
			: ''}{direction} its {habName} to {habValue}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_HOMEWORLD}
	Your home planet is {planetName}. Your people are ready to leave the nest and explore the
	universe. Good luck.
{:else if message.type === PlayerMessageType.PLANET_BOMBED}
	{@const bombing = message.spec?.bombing}
	{#if bombing}
		{$universe.getPlayerPluralName(message.spec?.mapObjectTarget?.targetPlayerNum)}
		{message.spec?.mapObjectTarget?.targetName} has bombed your planet {planetName}
		{#if message.spec?.bombing?.planetEmptied}
			killing off all its colonists.
		{:else}
			{#if bombing.colonistsKilled == 0 && bombing.minesDestroyed == 0 && bombing.factoriesDestroyed == 0 && bombing.defensesDestroyed == 0}
				doing no physical damage.
			{:else}
				killing {bombing.colonistsKilled.toLocaleString()} colonists, and destroying {bombing.minesDestroyed}
				mines,
				{bombing.factoriesDestroyed} factories and {bombing.defensesDestroyed} defenses.
			{/if}

			{#if absSum(bombing.unterraformAmount) > 0}
				{#if bombing.numBombers > 1}
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
		Bombers have bombed planet {planetName}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BONUS_RESEARCH_ARTIFACT}
	Your colonists settling {planetName} have found a strange artifact boosting your research in {message
		.spec?.field} by {message.spec?.amount} resources.
{:else if message.type === PlayerMessageType.PLANET_BUILT_DEFENSE}
	{#if message.spec?.amount === 1}
		You have built a defense outpost on {planetName}.
	{:else}
		You have built {message.spec?.amount ?? 0} defense outposts on {planetName}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_FACTORY}
	{#if message.spec?.amount === 1}
		You have built a factory on {planetName}.
	{:else}
		You have built {message.spec?.amount ?? 0} factories on {planetName}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_GENESIS_DEVICE}
	Strong fundamental forces have rebirthed {planetName}. All planetary installations have been wiped
	clean as its environment shifts drastically and newfound minerals spring forth from the ground.
{:else if message.type === PlayerMessageType.PLANET_BUILT_MINERAL_ALCHEMY}
	Your scientists on {planetName} have transmuted common materials into {message.spec?.amount ??
		0}kT each of Ironium, Boranium and Germanium.
{:else if message.type === PlayerMessageType.PLANET_BUILT_MINE}
	{#if message.spec?.amount === 1}
		You have built a mine on {planetName}.
	{:else}
		You have built {message.spec?.amount ?? 0} mines on {planetName}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_ITEM}
	You have attempted to build a {getFullName(
		create(ProductionQueueItemSchema, { type: message.spec?.queueItemType }),
		$universe
	).toLowerCase()} on {planetName}, but {planetName}
	is unable to build any of these. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_SHIP}
	{#if planet?.spec?.planetStarbaseSpec?.hasStarbase}
		You have attempted to build a {message.spec?.name} on {planetName}, but the starbase can only
		build ships up to {message.spec?.amount2 ?? 0}kT. The order has been canceled.
	{:else}
		You have attempted to build a {message.spec?.name} on {planetName}, but
		{planetName} has no starbase to build it. The order has been canceled.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_BEYOND_MAXIMUM}
	{#if message.spec?.queueItemType === QueueItemType.TERRAFORM_ENVIRONMENT}
		{planetName} has orders to terraform beyond the maximum allowed. The orders have been reduced to the
		maximum allowable.
	{:else}
		{planetName} has orders to build planetary installations beyond the maximum allowed. The orders have
		been reduced to the maximum allowable.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_MINERAL_PACKET_NO_MASS_DRIVER}
	You have attempted to build a mineral packet on {planetName}, but you have no starbase equipped
	with a mass driver on this planet. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_BUILT_INVALID_MINERAL_PACKET_NO_TARGET}
	You have attempted to build a mineral packet on {planetName}, but you have failed to specify a
	planet to target. The order has been canceled.
{:else if message.type === PlayerMessageType.PLANET_BUILT_SCANNER}
	{planetName} has built a new {message.spec?.name} planetary scanner.
{:else if message.type === PlayerMessageType.PLANET_BUILT_STARBASE}
	{planetName} has built a new {message.spec?.name}.
	{#if planet?.spec?.planetStarbaseSpec?.dockCapacity == UnlimitedSpaceDock}
		Ships of any size can now be built here.
	{:else if (planet?.spec?.planetStarbaseSpec?.dockCapacity ?? 0) > 0}
		Ships up to {planet?.spec?.planetStarbaseSpec?.dockCapacity}kT in mass can now be built at this
		facility.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_COMET_STRIKE}
	{#if message.spec?.comet?.size === CometSize.SMALL}
		A small comet has crashed into {planetName} bringing new minerals and altering the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.MEDIUM}
		A medium-sized comet has crashed into {planetName} bringing a significant quantity of minerals and
		significantly altering the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.LARGE}
		A large comet has crashed into {planetName} bringing a wide variety of new minerals and drastically
		altering the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.HUGE}
		A huge comet has crashed into {planetName} embedding vast quantities of minerals in the planet and
		radically altering its environment.
	{:else}
		A comet has crashed into {planetName} bringing new minerals and altering the planet's environment.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_COMET_STRIKE_MY_PLANET}
	{#if message.spec?.comet?.size === CometSize.SMALL}
		A small comet has crashed into your planet {planetName}, killing {message.spec.comet.colonistsKilled.toLocaleString()}
		of your colonists. The comet brought additional minerals and has slightly altered the planet's habitat.
	{:else if message.spec?.comet?.size === CometSize.MEDIUM}
		A medium-sized comet has crashed into your planet {planetName}, killing {message.spec.comet.colonistsKilled.toLocaleString()}
		of your colonists. The comet brought additional minerals and has altered the planet's environment.
	{:else if message.spec?.comet?.size === CometSize.LARGE}
		A large comet has crashed into your planet {planetName}, killing {message.spec.comet.colonistsKilled.toLocaleString()}
		of your colonists. The comet brought significant quantities of minerals and has greatly altered the
		planet's environment.
	{:else if message.spec?.comet?.size === CometSize.HUGE}
		A huge comet has crashed into your planet {planetName}, killing {message.spec.comet.colonistsKilled.toLocaleString()}
		of your colonists. The comet has embedded vast stores of minerals and drastically altered the planet's
		environment.
	{:else}
		A comet has crashed into {planetName} bringing new minerals and altering the planet's environment.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_RANDOM_MINERAL_DEPOSIT}
	Your surveyors on {planetName} have discovered a previously unknown deposit of {message.spec
		?.name}, significantly increasing the planet's concentration.
{:else if message.type === PlayerMessageType.PLANET_DIED_OFF}
	{#if $player.race.spec.livesOnStarbases}
		All of your colonists orbiting {planetName} have died off. Your starbase has been lost and you no
		longer control the planet.
	{:else}
		All of your colonists on {planetName} have died off. You no longer control the planet.
	{/if}
{:else if [PlayerMessageType.PLANET_DISCOVERY, PlayerMessageType.PLANET_DISCOVERY_HABITABLE, PlayerMessageType.PLANET_DISCOVERY_TERRAFORMABLE, PlayerMessageType.PLANET_DISCOVERY_UNINHABITABLE].indexOf(message.type) != -1}
	{#if owner}
		You have found a planet occupied by someone else. {planetName} is currently owned by the
		{owner.racePluralName}.
	{:else if $player.race.spec.instaforming && ((planet?.spec?.terraformedHabitability && planet.spec.terraformedHabitability > 0) || (planet?.spec?.habitability && planet.spec.habitability > 0))}
		You have found a new habitable planet. Your colonists will grow by up to {Math.max(
			1,
			((planet.spec.terraformedHabitability || planet.spec.habitability) * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planetName}.
	{:else if planet?.spec?.habitability && planet.spec.habitability > 0}
		You have found a new habitable planet. Your colonists will grow by up to {Math.max(
			1,
			(planet.spec.habitability * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planetName}.
	{:else if planet?.spec?.terraformedHabitability && planet.spec.terraformedHabitability > 0}
		You have found a new planet which you have the ability to make habitable. With terraforming,
		your colonists will grow by up to {Math.max(
			1,
			(planet.spec.terraformedHabitability * growthRate) / 100
		).toFixed(2)}% per year if you colonize {planetName}.
	{:else}
		You have found a new planet which unfortunately is not habitable by you. {Math.max(
			1,
			-(planet?.spec?.habitability ?? 0) / 10
		).toFixed(2)}% of your colonists will die per year if you colonize {planetName}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_POPULATION_DECREASED}
	{#if message.spec?.amount === message.spec?.prevAmount}
		Your colonists on {planetName} are suffering under the planet's hostile conditions but for now they
		are surviving.
	{:else}
		The population on {planetName} has decreased from {(
			message.spec?.prevAmount ?? 0
		).toLocaleString()} to {(message.spec?.amount ?? 0).toLocaleString()}.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_INSTAFORM}
	Your race has instantly terraformed {planetName} up to optimal conditions. Its value is now {planet
		?.spec?.habitability ?? 0}%.
{:else if message.type === PlayerMessageType.FLEET_INVADED_PLANET}
	{@const invasion = message.spec?.invasion}
	{#if invasion}
		{#if invasion.successful}
			Your troops crush {$universe.getPlayerPluralName(invasion.defenderPlayerNum)}'s Colonists on
			{planetName}. You are now in control of the planet.
		{:else if invasion.attackersKilledByDefenses > 0}
			Of the {#if invasion.attackers > 0}{`${invasion.attackers.toLocaleString()} `}
			{/if}Colonists you dropped on {planetName},
			{#if invasion.attackers > 0}
				{((invasion.attackersKilledByDefenses / invasion.attackers) * 100).toLocaleString(
					undefined,
					{
						maximumFractionDigits: 2
					}
				)}%
			{:else}
				{invasion.attackersKilledByDefenses.toLocaleString()}
			{/if}
			were destroyed by planetary defenses, the rest were massacred by the ground troops of
			{$universe.getPlayerPluralName(invasion.defenderPlayerNum)}.
		{:else}
			The {#if invasion.attackers > 0}{`${invasion.attackers.toLocaleString()} `}
			{/if}Colonists you dropped on
			{planetName} were massacred by the ground troops of
			{$universe.getPlayerPluralName(invasion.defenderPlayerNum)}.
		{/if}
		{#if invasion.attackers > 0 && invasion.defenders > 0}
			The invasion began with {invasion.attackers.toLocaleString()} attacking colonists and
			{invasion.defenders.toLocaleString()} defending colonists.
		{/if}
		{#if invasion.successful}
			You lost {invasion.attackersKilled.toLocaleString()} colonists in the invasion.
		{:else}
			Your troops killed {invasion.defendersKilled.toLocaleString()} defending colonists.
		{/if}
	{:else}
		{planetName} was invaded, but your spies know nothing of the outcome.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_INVADED}
	{@const invasion = message.spec?.invasion}
	{#if invasion}
		{#if invasion.successful}
			{$universe.getPlayerPluralName(invasion.attackerPlayerNum)} have attacked you on {planetName}
			with {#if invasion.attackers > 0}{`${invasion.attackers.toLocaleString()} `}
			{/if}first-rate storm troopers. Though your colonists put up a spirited defense they are
			crushed.
		{:else if invasion.attackersKilledByDefenses > 0}
			Your planetary defenses and ground troops on {planetName} destroyed the
			{#if invasion.attackers > 0}{`${invasion.attackers.toLocaleString()} `}
			{/if}invading troops of
			{$universe.getPlayerPluralName(invasion.attackerPlayerNum)}.
		{:else}
			Your ground troops on {planetName} valiantly destroyed the
			{#if invasion.attackers > 0}{`${invasion.attackers.toLocaleString()} `}
			{/if}attacking barbarians of
			{$universe.getPlayerPluralName(invasion.attackerPlayerNum)}!
		{/if}
		{#if invasion.attackers > 0 && invasion.defenders > 0}
			The invasion began with {invasion.attackers.toLocaleString()} attacking colonists and
			{invasion.defenders.toLocaleString()} defending colonists.
		{/if}
		{#if invasion.successful}
			Your colonists killed {invasion.attackersKilled.toLocaleString()} invaders before being overrun.
		{:else}
			You lost {invasion.defendersKilled.toLocaleString()} colonists in the process.
		{/if}
	{:else}
		{planetName} was invaded, but your spies know nothing of the outcome.
	{/if}
{:else if message.type === PlayerMessageType.PLANET_POPULATION_DECREASED_OVERCROWDING}
	The population on {planetName} has decreased by {(-(message.spec?.amount ?? 0)).toLocaleString()}
	colonists due to overcrowding.
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_INVASION}
	Your colonists invading {planetName} have picked through the defenders' remains looking for technology.
	In the process you have gained a level in {message.spec?.field}.
{:else if message.type === PlayerMessageType.PLAYER_ACQUIRABLE_PART_GAINED_BATTLE}
	Your people have picked through the wreckage from the battle at {planetName} and have learned how to
	build {message.spec?.techGained}.
{:else if message.type === PlayerMessageType.FLEET_SCRAPPED}
	{#if planet?.spec?.planetStarbaseSpec?.hasStarbase}
		{message.spec?.mapObjectTarget?.targetName} has been dismantled for {totalMinerals(
			message.spec?.cost
		)}kT of minerals at the starbase orbiting {planetName}.
	{:else}
		{message.spec?.mapObjectTarget?.targetName} has been dismantled for {totalMinerals(
			message.spec?.cost
		)}kT of minerals which have been deposited on {planetName}.
	{/if}
	{#if message.spec?.cost?.resources}
		&nbsp;Ultimate Recycling has also made {message.spec.cost.resources} resources available for immediate
		use (less if other ships were scrapped here this year).
	{/if}
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_SCRAP_FLEET}
	In the process of {message.spec?.name} being scrapped above {planetName}, you have gained a level
	in {message.spec?.field}.
{:else if message.type === PlayerMessageType.PLAYER_ACQUIRABLE_PART_GAINED_SCRAP_FLEET}
	In the process of {message.spec?.name} being scrapped above {planetName}, you have learned to
	build {message.spec?.techGained}.
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_BATTLE}
	Wreckage from the battle that occurred in orbit of {planetName} has boosted your research in {message
		.spec?.field} by 1 level.
{:else if message.type === PlayerMessageType.FLEET_BUILT}
	<!-- TODO: remove this at some point. These are now fleet messages, not planet messages, but keeping this here so old savdes still target correctly -->
	Your starbase at {planetName} has built {message.spec?.amount ?? 'a'} new {message.spec?.name}s.
{:else}
	<FallbackMessageDetail {message} />
{/if}
