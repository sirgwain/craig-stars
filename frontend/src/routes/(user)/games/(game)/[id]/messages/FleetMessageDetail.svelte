<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { Fleet } from '$lib/types/Fleet';
	import type { PlayerIntel } from '$lib/types/Player';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';
	import FleetEngineStrainMessageDetail from './FleetEngineStrainMessageDetail.svelte';
	import { None } from '$lib/types/MapObject';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let fleet: Fleet | undefined;
	export let owner: PlayerIntel | undefined;
</script>

{#if message.type === MessageType.FleetRadiatingEngineDieoff}
	<!-- Colonist dieoff from engine radiation -->
	Engine radiation has killed {(message.spec.amount ?? 0) * -100} colonists traveling in {message.spec.name}.
{:else if message.type === MessageType.FleetDieoff}
	Due to the rigors of warp acceleration, {(message.spec.amount ?? 0) * -100} of your colonists on {message.spec.name}
	have died.
{:else if message.type === MessageType.FleetShipExceededSafeSpeed}
	<!-- Overwarp -->
	<FleetEngineStrainMessageDetail {message} />
{:else if message.type === MessageType.FleetTransferGivenFailed}
	<!-- Fleet Transfers -->
	{#if message.spec.destPlayerNum == undefined || message.spec.destPlayerNum == None || message.spec.destPlayerNum < 0 || message.spec.destPlayerNum >= $game.players.length}
		You cannot give {message.spec.name} away. No player to transfer to was specified.
	{:else}
		You cannot give {message.spec.name} to {$universe.getPlayerName(message.spec.destPlayerNum)}.
	{/if}
{:else if message.type === MessageType.FleetTransferGivenRefused}
	{$universe.getPlayerName(message.spec.destPlayerNum)} snub your attempted gift and refuse the fleet
	{message.spec.name}. Are you sure they your allies?
{:else if message.type === MessageType.FleetTransferGivenFailedColonists}
	You couldn't give {message.spec.name} away because there were some of your colonists on board.
{:else if message.type === MessageType.FleetTransferGiven}
	{message.spec.name} has successfully been given to {$universe.getPlayerName(
		message.spec.destPlayerNum
	)}
{:else if message.type === MessageType.FleetTransferReceived}
	{$universe.getPlayerName(message.spec.sourcePlayerNum)} has given you {message.spec.name}
{:else}
	<!-- Fallback for unknown message types -->
	<FallbackMessageDetail {message} />
{/if}
