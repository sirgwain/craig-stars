<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { absSum } from '$lib/types/Hab';
	import { None } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';
	import FleetEngineStrainMessageDetail from './FleetEngineStrainMessageDetail.svelte';

	const { game, universe } = getGameContext();

	export let message: Message;
</script>

{#if message.text}
	{message.text}
{:else if message.type === MessageType.FleetBombedPlanet}
	{@const bombing = message.spec.bombing}
	{#if bombing}
		{#if bombing.numBombers == 1}
			Your fleet {message.targetName} has bombed {$universe.getPlayerPluralName(
				message.spec.targetPlayerNum
			)} planet
			{message.spec.targetName}
		{:else}
			Your fleets have bombed the {$universe.getPlayerPluralName(message.spec.targetPlayerNum)} planet
			{message.spec.targetName}
		{/if}
		{#if bombing.planetEmptied}
			killing off all colonists.
		{:else}
			killing {bombing.colonistsKilled ?? 0} colonists, and destroying {bombing.minesDestroyed ?? 0}
			mines,
			{bombing.factoriesDestroyed ?? 0} factories, and {bombing.defensesDestroyed ?? 0} defenses.

			{#if absSum(bombing.unterraformAmount ?? {}) > 0}
				Your bombers have retro-bombed the planet, undoing {absSum(
					bombing.unterraformAmount ?? {}
				)}% of its terraforming.
			{/if}
		{/if}
	{:else}
		<!-- Generic message, no bombing data (unexpected) -->
		Your fleet {message.targetName} has bombed {$universe.getPlayerPluralName(
			message.spec.targetPlayerNum
		)} planet
		{message.spec.targetName}.
	{/if}
{:else if message.type === MessageType.FleetBuilt}
	{#if message.spec.amount === 1}
	Your starbase at {message.spec.targetName} has built a new {message.spec.name}.
	{:else}
	Your starbase at {message.spec.targetName} has built {message.spec.amount ?? 'a'} new {message
		.spec.name} ships.
	{/if}
{:else if message.type === MessageType.FleetDieoff}
	Due to the rigors of warp acceleration, {(message.spec.amount ?? 0) * -100} of your colonists on {message.targetName}
	have died.
{:else if message.type === MessageType.FleetExceededSafeSpeed}
	<!-- Overwarp -->
	<FleetEngineStrainMessageDetail {message} />
{:else if message.type === MessageType.FleetPatrolTargeted}
	Your patrolling {message.targetName} has targeted {message.spec.targetName} for intercept.
{:else if message.type === MessageType.FleetRadiatingEngineDieoff}
	<!-- Colonist dieoff from engine radiation -->
	Engine radiation has killed {(message.spec.amount ?? 0) * -100} colonists traveling in {message.targetName}.
{:else if message.type === MessageType.FleetReproduce}
	{#if !message.spec.amount2 || !message.spec.targetNum}
		Your colonists in {message.targetName} have made good use of their time increasing their on-board
		number by {message.spec.amount} colonists.
	{:else}
		Breeding activities on {message.targetName} have overflowed living space. {message.spec.amount2}
		colonists have been beamed down to {message.spec.targetName}.
	{/if}
{:else if message.type === MessageType.FleetRemoteMined}
	{message.targetName} has remote mined {message.spec.targetName}, extracting {message.spec.mineral
		?.ironium ?? 0}kT of Ironium, {message.spec.mineral?.boranium ?? 0}kT of Boranium, and {message
		.spec.mineral?.germanium ?? 0}kT of Germanium.
{:else if message.type === MessageType.FleetTransferGiven}
	{message.targetName} has successfully been given to {$universe.getPlayerPluralName(
		message.spec.destPlayerNum
	)}.
{:else if message.type === MessageType.FleetScrapped}
	{message.targetName} has been dismantled. The scrap was left in deep space.
{:else if message.type === MessageType.FleetTransferInvalidPlayer}
	<!-- Fleet Transfers -->
	{#if message.spec.destPlayerNum == undefined || message.spec.destPlayerNum == None || message.spec.destPlayerNum < 0 || message.spec.destPlayerNum >= $game.players.length}
		You cannot give {message.targetName} away. No player to transfer to was specified.
	{:else}
		You cannot give {message.targetName} to {$universe.getPlayerPluralName(message.spec.destPlayerNum)}.
	{/if}
{:else if message.type === MessageType.FleetTransferInvalidColonists}
	You couldn't give {message.targetName} away because there were some of your colonists on board.
{:else if message.type === MessageType.FleetTransferInvalidGiveRefused}
	{$universe.getPlayerPluralName(message.spec.destPlayerNum)} snubbed your attempted gift and refused your offer of
	{message.targetName}. Are you sure they're still your allies?
{:else if message.type === MessageType.FleetTransferInvalidReceiveRefused}
	{$universe.getPlayerPluralName(message.spec.sourcePlayerNum)} has attempted to gift you {message.targetName},
	but you have refused their offer. If you wish to receive gifts from this player in the future,
	make sure to set them as allies.
{:else if message.type === MessageType.FleetTransferReceived}
	{$universe.getPlayerPluralName(message.spec.sourcePlayerNum)} has given you {message.targetName}.
{:else}
	<!-- Fallback for unknown message types -->
	<FallbackMessageDetail {message} />
{/if}
