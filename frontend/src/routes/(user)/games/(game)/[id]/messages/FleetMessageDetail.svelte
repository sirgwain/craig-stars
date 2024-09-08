<script lang="ts">
	import { andCommaList } from '$lib/andCommandList';
	import { getGameContext } from '$lib/services/GameContext';
	import { absSum } from '$lib/types/Hab';
	import { None } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';
	import FleetEngineStrainMessageDetail from './FleetEngineStrainMessageDetail.svelte';

	const { game, universe, player } = getGameContext();

	export let message: Message;
</script>

{#if message.text}
	{message.text}
{:else if message.type === MessageType.FleetBombedPlanet}
	{@const bombing = message.spec.bombing}
	{#if bombing}
		{#if bombing.numBombers == 1}
			Your fleet {message.targetName} has bombed the {$universe.getPlayerPluralName(
				message.spec.targetPlayerNum
			)} settlement on
			{message.spec.targetName}
		{:else}
			Your fleets have bombed the {$universe.getPlayerPluralName(message.spec.targetPlayerNum)} settlement
			on
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
{:else if message.type === MessageType.FleetMineFieldHit}
	{@const damage = message.spec.mineFieldDamage}
	{@const mineFieldOwner = $universe.getPlayerPluralName(message.spec.targetPlayerNum)}
	{@const mineFieldPosition = `(${message.spec.targetPosition.x}, ${message.spec.targetPosition.y})`}
	{#if damage}
		{#if message.targetPlayerNum === $player.num}
			<!-- our fleet was hit -->
			{#if damage.fleetDestroyed}
				{message.spec.name} has been annihilated in a {mineFieldOwner} mine field at {mineFieldPosition}.
			{:else}
				{message.spec.name} has been stopped in a {mineFieldOwner} mine field at {mineFieldPosition}.
				{#if (damage.shipsDestroyed ?? 0) > 0}
					Your fleet has taken {damage.damage ?? 0} damage points and {damage.shipsDestroyed} ships were
					destroyed.
				{:else if (damage.damage ?? 0) > 0}
					Your fleet has taken {damage.damage ?? 0} damage points but none of your ships were destroyed.
				{/if}
			{/if}
		{:else}
			<!-- our minefield hit someone else's fleet -->
			{#if damage.fleetDestroyed}
				{message.spec.name} has been annihilated in your mine field at {mineFieldPosition}.
			{:else}
				{message.spec.name} has been stopped in your mine field at {mineFieldPosition}.
				{#if (damage.shipsDestroyed ?? 0) > 0}
					Your mines have inflicted {damage.damage ?? 0} damage points and destroyed {damage.shipsDestroyed}
					ships.
				{:else if (damage.damage ?? 0) > 0}
					Your mines have inflicted {damage.damage ?? 0} damage points, but you didn't manage to destroy
					any ships.
				{/if}
			{/if}
		{/if}
	{:else}
		Unknown damage was done
	{/if}
{:else if message.type === MessageType.FleetMineFieldSweptMines}
	{@const mineFieldPosition = `(${message.spec.targetPosition.x}, ${message.spec.targetPosition.y})`}
	{#if message.targetPlayerNum === $player.num}
		<!-- our fleet swept -->
		{message.spec.name} has has swept {message.spec.amount} mines from a mine field at {mineFieldPosition}
	{:else}
		<!-- our minefield was swept by fleet -->
		{message.spec.name} has has swept {message.spec.amount} mines from your mine field at {mineFieldPosition}
	{/if}
{:else if message.type === MessageType.FleetLaidMines}
	{@const mineField = $universe.getMineField(message.spec.targetPlayerNum, message.spec.targetNum)}
	{#if mineField?.numMines === message.spec.amount}
		{message.spec.name} has has dispensed {message.spec.amount} mines.
	{:else}
		{message.spec.name} has increased {mineField?.name} by {message.spec.amount} mines.
	{/if}
{:else if message.type === MessageType.FleetPatrolTargeted}
	Your patrolling {message.targetName} has targeted {message.spec.targetName} to intercept.
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
	<!-- Remote Mining messages -->
{:else if message.type === MessageType.FleetRemoteMined}
	{@const minerals = {
		ironium: message.spec.mineral?.ironium ?? 0,
		boranium: message.spec.mineral?.boranium ?? 0,
		germanium: message.spec.mineral?.germanium ?? 0
	}}
	{message.targetName} has remote mined {message.spec.targetName} extracting {andCommaList(
		[
			minerals.ironium > 0 ? `${minerals.ironium} kT of Ironium` : '',
			minerals.boranium > 0 ? `${minerals.boranium} kT of Boranium` : '',
			minerals.germanium > 0 ? `${minerals.germanium} kT of Germanium` : ''
		],
		'no minerals.'
	)}
{:else if message.type === MessageType.FleetScrapped}
	{message.targetName} has been dismantled. The scrap was left in deep space.
{:else if message.type === MessageType.FleetTransferGiven}
	{message.targetName} has successfully been given to {$universe.getPlayerPluralName(
		message.spec.destPlayerNum
	)}.
{:else if message.type === MessageType.FleetTransferInvalidPlayer}
	<!-- Fleet Transfers -->
	{#if message.spec.destPlayerNum == undefined || message.spec.destPlayerNum == None || message.spec.destPlayerNum < 0 || message.spec.destPlayerNum >= $game.players.length}
		You cannot give {message.targetName} away. No player to transfer to was specified.
	{:else}
		You cannot give {message.targetName} to {$universe.getPlayerPluralName(
			message.spec.destPlayerNum
		)}.
	{/if}
{:else if message.type === MessageType.FleetTransferInvalidColonists}
	You couldn't give {message.targetName} away because there were some of your colonists on board.
{:else if message.type === MessageType.FleetTransferInvalidGiveRefused}
	{$universe.getPlayerPluralName(message.spec.destPlayerNum)} snubbed your attempted gift and refused
	your offer of
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
