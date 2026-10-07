<script lang="ts">
	import { andCommaList } from '$lib/andCommandList';
	import { getGameContext } from '$lib/services/GameContext';
	import { absSum } from '$lib/types/Hab';
	import { None } from '$lib/types/Consts';
	import {
		CargoTransferStatus,
		MapObjectType,
		PlayerMessageType,
		ResourceType,
		type PlayerMessage
	} from '$lib/types/cs-proto';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';
	import FleetEngineStrainMessageDetail from './FleetEngineStrainMessageDetail.svelte';
	import { enumToString } from '$lib/types/Enums';
	import { cargoDescription, population, totalCargo } from '$lib/types/Cargo';

	const { game, universe, player } = getGameContext();

	type Props = {
		message: PlayerMessage;
	};

	let { message }: Props = $props();

	let fleetName = $derived(message.target?.targetName ?? 'your fleet');
	let sourceName = $derived(message.spec?.mapObjectTarget?.targetName ?? 'unknown');
	let destinationName = $derived(message.spec?.routeTarget?.targetName ?? 'unknown');

	let fleet = $derived(
		message.target?.targetPlayerNum && message.target.targetNum
			? $universe.getFleet(message.target.targetPlayerNum, message.target.targetNum)
			: undefined
	);
</script>

{#if message.text}
	{message.text}
{:else if message.type === PlayerMessageType.FLEET_ORDERS_COMPLETE}
	{fleetName} has completed its assigned orders.
{:else if message.type === PlayerMessageType.FLEET_ENGINE_FAILURE}
	{fleetName} was unable to engage its engines due to balky equipment. Engineers think they have the problem
	fixed for the time being.
{:else if message.type === PlayerMessageType.FLEET_OUT_OF_FUEL}
	<!-- Amount is the resulting warp speed. -->
	{fleetName} has run out of fuel. The fleet's speed has been decreased to Warp {message.spec
		?.amount ?? 0}.
{:else if message.type === PlayerMessageType.FLEET_MERGED}
	{message.spec?.name} has been merged into {fleetName}.
{:else if message.type === PlayerMessageType.FLEET_MERGE_INVALID_NOT_FLEET}
	{fleetName} was unable to complete its merge orders as the waypoint destination wasn't a fleet.
{:else if message.type === PlayerMessageType.FLEET_MERGE_INVALID_UNOWNED}
	{fleetName} was unable to complete its merge orders as the destination fleet wasn't one of yours.
{:else if message.type === PlayerMessageType.FLEET_ROUTE_INVALID_NOT_PLANET}
	{fleetName} could not be routed because it is not orbiting a planet.
{:else if message.type === PlayerMessageType.FLEET_ROUTE_INVALID_NOT_FRIENDLY_PLANET}
	{fleetName} could not be routed because you are not friends with the inhabitants of {sourceName}.
{:else if message.type === PlayerMessageType.FLEET_ROUTE_INVALID_NO_ROUTE_TARGET}
	{fleetName} could not be routed at {sourceName} as the planet has no route set.
{:else if message.type === PlayerMessageType.FLEET_ROUTE}
	{fleetName} has been routed by the citizens of {sourceName} to {destinationName}.
{:else if message.type === PlayerMessageType.FLEET_COLONIZE_INVALID_NOT_PLANET}
	{fleetName} has orders to colonize, but is not currently orbiting a planet. The order has been canceled.
{:else if message.type === PlayerMessageType.FLEET_COLONIZE_INVALID_OWNED_PLANET}
	{fleetName} has orders to colonize {sourceName}, but {sourceName} is already populated. The order has
	been canceled.
{:else if message.type === PlayerMessageType.FLEET_COLONIZE_INVALID_NO_MODULE}
	{fleetName} has orders to colonize a planet without a colonization module. The order has been canceled.
{:else if message.type === PlayerMessageType.FLEET_COLONIZE_INVALID_NO_COLONISTS}
	{fleetName} has orders to colonize a planet, but has failed to bring along any colonists. The order
	has been canceled.
{:else if message.type === PlayerMessageType.FLEET_LAY_MINES_INVALID_NO_MINE_LAYERS}
	{fleetName} has orders to lay mines, but has no mine layers. The order has been canceled.
{:else if message.type === PlayerMessageType.FLEET_REMOTE_MINE_INVALID_NO_MINERS}
	{fleetName} has orders to remote mine {sourceName}, but the fleet doesn't have any remote mining
	modules. The order has been canceled.
{:else if message.type === PlayerMessageType.FLEET_REMOTE_MINE_INVALID_INHABITED}
	{fleetName} has orders to remote mine {sourceName}, but the planet is already inhabited. The order
	has been canceled.
{:else if message.type === PlayerMessageType.FLEET_REMOTE_MINE_INVALID_DEEP_SPACE}
	{fleetName} has orders to remote mine in deep space. The order has been canceled.
{:else if message.type === PlayerMessageType.FLEET_TRANSFERRED_CARGO}
	{@const transfer = message.spec?.cargoTransfer}
	{#if transfer}
		{@const loading = transfer.transfered < 0}
		{#if transfer.cargoType === ResourceType.COLONISTS}
			{fleetName} has beamed {cargoDescription(transfer.cargoType, Math.abs(transfer.transfered))} colonists
			{loading ? 'from' : 'to'}
			{sourceName}.
		{:else}
			{fleetName} has {loading ? 'loaded' : 'unloaded'}
			{cargoDescription(transfer.cargoType, Math.abs(transfer.transfered))} of {enumToString(
				ResourceType,
				transfer.cargoType
			)}
			{loading ? 'from' : 'to'}
			{sourceName}.
		{/if}
	{:else}
		<FallbackMessageDetail {message} />
	{/if}
{:else if message.type === PlayerMessageType.FLEET_TARGET_LOST}
	{#if message.spec?.lostTargetType === MapObjectType.FLEET}
		The fleet you were tracking with {fleetName}, {message.spec.name}, appears to have outrun the
		range of your scanners. Orders for your fleet have been changed to go to the last known location
		of that fleet.
	{:else}
		The {message.spec?.name} that you were tracking with {fleetName}, appears to have disappeared.
		Orders for your fleet have been changed to go to the last known location of the target.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_SOURCE}
	{fleetName} attempted to use a stargate at {sourceName}, but no stargate exists there.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_SOURCE_OWNER}
	{fleetName} attempted to use a stargate at {sourceName}, but could not because the starbase is not
	owned by you or your allies.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_DEST}
	{fleetName} attempted to use a stargate at {sourceName} to reach {destinationName}, but no
	stargate could be detected at the destination.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_DEST_OWNER}
	{fleetName} attempted to use a stargate at {sourceName} to reach {destinationName}, but could not
	because the destination starbase is not owned by you or your allies.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_RANGE}
	{fleetName} attempted to use a stargate at {sourceName} to reach {destinationName}, but the
	distance of {(message.spec?.distance ?? 0).toFixed(1)} Ly. was far beyond the max range of the stargates.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_MASS}
	{fleetName} attempted to use a stargate at {sourceName} to reach {destinationName}, but your ships
	are far too massive for the gate's limits.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_INVALID_COLONISTS}
	{fleetName} attempted to use a stargate at {sourceName} to reach {destinationName}, but you are
	carrying colonists and can't drop them off as you don't own the planet.
{:else if message.type === PlayerMessageType.FLEET_DUMPED_CARGO}
	{@const cargo = message.spec?.cargo}
	{@const minerals = totalCargo(cargo) - (cargo?.colonists ?? 0)}
	{fleetName} has unloaded
	{#if population(cargo) > 0 && minerals > 0}
		{population(cargo).toLocaleString()} colonists and {minerals.toLocaleString()}kt of minerals
	{:else if population(cargo) > 0}
		{population(cargo).toLocaleString()} colonists
	{:else}
		{minerals.toLocaleString()}kt of minerals
	{/if}
	in preparation for jumping through the stargate at {sourceName} to reach {destinationName}.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_DESTROYED}
	Heedless to the danger, {fleetName} attempted to use the stargate at {sourceName} to reach {destinationName}.
	The fleet never arrived. The distance or mass must have been too great.
{:else if message.type === PlayerMessageType.FLEET_STARGATE_DAMAGED}
	<!-- Amount is damage points; Amount2 is total ships lost. -->
	{@const shipsLost = message.spec?.amount2 ?? 0}
	{fleetName} used the stargate at {sourceName} to reach {destinationName}
	{#if shipsLost === 0}
		losing no ships but suffering {message.spec?.amount ?? 0} dp of damage. They exceeded the capability
		of the gates.
	{:else if shipsLost < 5}
		losing only {shipsLost}
		{shipsLost === 1 ? 'ship' : 'ships'} to the treacherous void. They were fortunate. They exceeded the
		capability of the gates.
	{:else if shipsLost <= 10}
		losing {shipsLost} ships to the unforgiving void. Exceeding the capability of your stargates can be
		dangerous...
	{:else if shipsLost <= 50}
		losing {shipsLost} ships to the great unknown. Such disregard for stargates' capabilities is not recommended...
	{:else}
		losing an unbelievable {shipsLost} ships to the cosmic ocean. The jump was far in excess of the capabilities
		of stargates involved...
	{/if}
{:else if message.type === PlayerMessageType.FLEET_BOMBED_PLANET}
	{@const bombing = message.spec?.bombing}
	{#if bombing}
		{#if bombing.numBombers == 1}
			Your fleet {message.target?.targetName} has bombed the {$universe.getPlayerPluralName(
				message.spec?.mapObjectTarget?.targetPlayerNum
			)} settlement on
			{message.spec?.mapObjectTarget?.targetName}
		{:else}
			Your fleets have bombed the {$universe.getPlayerPluralName(
				message.spec?.mapObjectTarget?.targetPlayerNum
			)} settlement on
			{message.spec?.mapObjectTarget?.targetName}
		{/if}
		{#if bombing.planetEmptied}
			killing off all colonists.
		{:else}
			killing {bombing.colonistsKilled} colonists, and destroying {bombing.minesDestroyed}
			mines,
			{bombing.factoriesDestroyed} factories, and {bombing.defensesDestroyed} defenses.

			{#if absSum(bombing.unterraformAmount) > 0}
				Your bombers have retro-bombed the planet, undoing {absSum(bombing.unterraformAmount)}% of
				its terraforming.
			{/if}
		{/if}
	{:else}
		<!-- Generic message, no bombing data (unexpected) -->
		Your fleet {message.target?.targetName} has bombed {$universe.getPlayerPluralName(
			message.spec?.mapObjectTarget?.targetPlayerNum
		)} planet
		{message.spec?.mapObjectTarget?.targetName}.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_BUILT}
	{@const routeTarget =
		message.spec?.routeTarget && $universe.getMapObject(message.spec.routeTarget)}
	{#if message.spec?.amount === 1}
		Your starbase at {message.spec.mapObjectTarget?.targetName} has built a new {message.spec.name}.
	{:else}
		Your starbase at {message.spec?.mapObjectTarget?.targetName} has built {message.spec?.amount ??
			'a'} new {message.spec?.name} ships.
	{/if}
	{#if routeTarget}
		It will be routed to {routeTarget.mapObject?.name}.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_DIEOFF}
	Due to the rigors of warp acceleration, {-(message.spec?.amount ?? 0)} of your colonists on {message
		.target?.targetName}
	have died.
{:else if message.type === PlayerMessageType.FLEET_EXCEEDED_SAFE_SPEED}
	<!-- Overwarp -->
	<FleetEngineStrainMessageDetail {message} />
{:else if message.type === PlayerMessageType.FLEET_ENGINE_STRAIN_DESTROYED}
	{fleetName} was destroyed in a massive reactor accident due to unsafe operating procedures.
{:else if message.type === PlayerMessageType.FLEET_GENERATED_FUEL}
	{@const hasRamscoops = fleet?.tokens.some(
		(t) =>
			($universe.getDesign(fleet.mapObject?.playerNum, t.designNum)?.spec?.engine?.freeSpeed ?? 0) >
			1
	)}
	{#if hasRamscoops}
		{message.target?.targetName}'s ramscoops have produced {message.spec?.amount}mg of fuel from
		interstellar hydrogen.
	{:else}
		{message.target?.targetName} has generated {message.spec?.amount}mg of fuel.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_MINEFIELD_HIT}
	{@const damage = message.spec?.minefieldDamage}
	{@const minefieldOwner = $universe.getPlayerPluralName(
		message.spec?.mapObjectTarget?.targetPlayerNum
	)}
	{@const minefieldPosition = `(${message.spec?.mapObjectTarget?.targetPosition?.x ?? 0}, ${message.spec?.mapObjectTarget?.targetPosition?.y ?? 0})`}
	{#if damage}
		{#if message.target?.targetPlayerNum === $player.num}
			<!-- our fleet was hit -->
			{#if damage.fleetDestroyed}
				{message.target.targetName} has been annihilated in a {minefieldOwner} minefield at {minefieldPosition}.
			{:else}
				{message.target.targetName} has been stopped in a {minefieldOwner} minefield at {minefieldPosition}.
				{#if damage.shipsDestroyed > 0}
					Your fleet has taken {damage.damage} damage points and {damage.shipsDestroyed} ships were destroyed.
				{:else if damage.damage > 0}
					Your fleet has taken {damage.damage} damage points but none of your ships were destroyed.
				{/if}
			{/if}
		{:else}
			<!-- our minefield hit someone else's fleet -->
			{#if damage.fleetDestroyed}
				{$universe.getPlayerName(message.target?.targetPlayerNum)}
				{message.target?.targetName} has been annihilated in your minefield at {minefieldPosition}.
			{:else}
				{$universe.getPlayerName(message.target?.targetPlayerNum)}
				{message.target?.targetName} has been stopped in your minefield at {minefieldPosition}.
				{#if damage.shipsDestroyed > 0}
					Your mines have inflicted {damage.damage} damage points and destroyed {damage.shipsDestroyed}
					ships.
				{:else if damage.damage > 0}
					Your mines have inflicted {damage.damage} damage points, but you didn't manage to destroy any
					ships.
				{/if}
			{/if}
		{/if}
	{:else}
		Unknown damage was done.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_MINEFIELD_SWEPT_MINES}
	{@const minefieldPosition = `(${message.spec?.mapObjectTarget?.targetPosition?.x ?? 0}, ${message.spec?.mapObjectTarget?.targetPosition?.y || 0})`}
	{#if message.target?.targetPlayerNum === $player.num}
		<!-- our fleet swept -->
		{message.target.targetName} has has swept {message.spec?.amount ?? 0} mines from a minefield at
		{minefieldPosition}
	{:else}
		<!-- our minefield was swept by fleet -->
		{$universe.getPlayerName(message.target?.targetPlayerNum)}
		{message.target?.targetName} has has swept {message.spec?.amount ?? 0} mines from your minefield at
		{minefieldPosition}
	{/if}
{:else if message.type === PlayerMessageType.FLEET_LAID_MINES}
	{@const minefield = $universe.getMinefield(
		message.spec?.mapObjectTarget?.targetPlayerNum,
		message.spec?.mapObjectTarget?.targetNum
	)}
	{#if minefield?.numMines === message.spec?.amount}
		{message.target?.targetName} has has dispensed {message.spec?.amount} mines.
	{:else}
		{message.target?.targetName} has increased {message.spec?.mapObjectTarget?.targetName} by {message
			.spec?.amount} mines.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_PATROL_TARGETED}
	Your patrolling {message.target?.targetName} has targeted {message.spec?.mapObjectTarget
		?.targetName} to intercept.
{:else if message.type === PlayerMessageType.FLEET_RADIATING_ENGINE_DIEOFF}
	<!-- Colonist dieoff from engine radiation -->
	Engine radiation has killed {(message.spec?.amount ?? 0).toLocaleString()} colonists traveling in {message
		.target?.targetName}.
{:else if message.type === PlayerMessageType.FLEET_REPRODUCE}
	{#if !message.spec?.amount2 || !message.spec.mapObjectTarget?.targetNum}
		Your colonists in {message.target?.targetName} have made good use of their time increasing their on-board
		number by {message.spec?.amount} colonists.
	{:else}
		<!-- TODO: actually fix bug non jankily by multiplying message.amount2 by 100 during assignment-->
		Breeding activities on {message.target?.targetName} have overflowed living space. {message.spec
			.amount2 * 100}
		colonists have been beamed down to {message.spec?.mapObjectTarget?.targetName}.
	{/if}
	<!-- Remote Mining messages -->
{:else if message.type === PlayerMessageType.FLEET_REMOTE_MINED}
	{@const minerals = {
		ironium: message.spec?.mineral?.ironium ?? 0,
		boranium: message.spec?.mineral?.boranium ?? 0,
		germanium: message.spec?.mineral?.germanium ?? 0
	}}
	{message.target?.targetName} has remote mined {message.spec?.mapObjectTarget?.targetName} extracting
	{andCommaList(
		[
			minerals.ironium > 0 ? `${minerals.ironium}kT of Ironium` : '',
			minerals.boranium > 0 ? `${minerals.boranium}kT of Boranium` : '',
			minerals.germanium > 0 ? `${minerals.germanium}kT of Germanium` : ''
		],
		'no minerals'
	)}.
{:else if message.type === PlayerMessageType.FLEET_SCRAPPED}
	{message.target?.targetName} has been dismantled. The scrap was left in deep space.
{:else if message.type === PlayerMessageType.FLEET_BY_HAND_TRANSFER_INCOMPLETE}
	{@const transfer = message.spec?.cargoTransfer}
	{#if transfer}
		{@const cargoType = enumToString(ResourceType, transfer.cargoType)}
		{@const fromTo = transfer.wanted < 0 ? 'from' : 'to'}
		{message.target?.targetName} has attempted to transfer {cargoDescription(
			transfer.cargoType,
			Math.abs(transfer.wanted)
		)} of {cargoType}
		{fromTo}
		{message.spec?.mapObjectTarget?.targetName}, but was
		{#if transfer.transfered === 0}
			unable to transfer any cargo.
		{:else}
			only able to transfer {cargoDescription(transfer.cargoType, Math.abs(transfer.transfered))}.
		{/if}
		{#if transfer.status === CargoTransferStatus.CARGO}
			{message.target?.targetName} did not have enough {cargoType}.
		{:else if transfer.status === CargoTransferStatus.CARGO_CAPACITY}
			{message.target?.targetName} did not have enough space in their hold.
		{:else if transfer.status === CargoTransferStatus.DEST_CARGO}
			{message.spec?.mapObjectTarget?.targetName} did not have enough {cargoType}.
		{:else if transfer.status === CargoTransferStatus.DEST_CARGO_CAPACITY}
			{message.spec?.mapObjectTarget?.targetName} did not have enough space in their hold.
		{:else if transfer.status === CargoTransferStatus.DEST_STARBASE}
			A starbase in orbit prevented the transfer.
		{:else if transfer.status === CargoTransferStatus.DEST_UNOWNED}
			The planet is unoccupied. Your colonists refuse to be beamed down without a colonization
			module.
		{:else if transfer.status === CargoTransferStatus.OWNED}
			{message.spec?.mapObjectTarget?.targetName} is owned by another player and {message.target
				?.targetName}
			does not have the required technology to bypass their sensors.
		{/if}
	{:else}
		<!-- Generic failure message -->
		{message.target?.targetName} has attempted to transfer cargo from {message.spec?.mapObjectTarget
			?.targetName}, but the cargo transfer was unsuccessful.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_TRANSPORT_INVALID}
	{@const transfer = message.spec?.cargoTransfer}
	{#if transfer}
		{@const cargoType = enumToString(ResourceType, transfer.cargoType)}
		{@const fromTo = transfer.wanted < 0 ? 'from' : 'to'}
		{message.target?.targetName} has attempted to transfer {cargoDescription(
			transfer.cargoType,
			Math.abs(transfer.wanted)
		)} of {cargoType}
		{fromTo}
		{message.spec?.mapObjectTarget?.targetName}, but was
		{#if transfer.transfered === 0}
			unable to transfer any cargo.
		{:else}
			only able to transfer {cargoDescription(transfer.cargoType, Math.abs(transfer.transfered))}.
		{/if}
		{#if transfer.status === CargoTransferStatus.CARGO}
			{message.target?.targetName} did not have enough {cargoType}.
		{:else if transfer.status === CargoTransferStatus.CARGO_CAPACITY}
			{message.target?.targetName} did not have enough space in their hold.
		{:else if transfer.status === CargoTransferStatus.DEST_CARGO}
			{message.spec?.mapObjectTarget?.targetName} did not have enough {cargoType}.
		{:else if transfer.status === CargoTransferStatus.DEST_CARGO_CAPACITY}
			{message.spec?.mapObjectTarget?.targetName} did not have enough space in their hold.
		{:else if transfer.status === CargoTransferStatus.DEST_STARBASE}
			A starbase in orbit prevented the transfer.
		{:else if transfer.status === CargoTransferStatus.DEST_UNOWNED}
			The planet is unoccupied. Your colonists refuse to be beamed down without a colonization
			module.
		{:else if transfer.status === CargoTransferStatus.OWNED}
			{#if $player.race.spec.livesOnStarbases}
				{message.spec?.mapObjectTarget?.targetName} is owned by another player and your people
			{:else}
				{message.spec?.mapObjectTarget?.targetName} is owned by another player and {message.target
					?.targetName}
				does not have the required technology to bypass their sensors.
			{/if}
		{/if}
	{:else}
		<!-- Generic failure message -->
		{message.target?.targetName} has attempted to transfer cargo from {message.spec?.mapObjectTarget
			?.targetName}, but the cargo transfer was unsuccessful.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_TRANSFER_GIVEN}
	{message.target?.targetName} has successfully been given to {$universe.getPlayerPluralName(
		message.spec?.destPlayerNum
	)}.
{:else if message.type === PlayerMessageType.FLEET_TRANSFER_INVALID_PLAYER}
	<!-- Fleet Transfers -->
	{#if message.spec?.destPlayerNum == undefined || message.spec.destPlayerNum == None || message.spec.destPlayerNum < 0 || message.spec.destPlayerNum >= $game.players.length}
		You cannot give {message.target?.targetName} away. No player to transfer to was specified.
	{:else}
		You cannot give {message.target?.targetName} to {$universe.getPlayerPluralName(
			message.spec?.destPlayerNum
		)}.
	{/if}
{:else if message.type === PlayerMessageType.FLEET_TRANSFER_INVALID_COLONISTS}
	You couldn't give {message.target?.targetName} away because there were some of your colonists on board.
{:else if message.type === PlayerMessageType.FLEET_TRANSFER_INVALID_GIVE_REFUSED}
	{$universe.getPlayerPluralName(message.spec?.destPlayerNum)} snubbed your attempted gift and refused
	your offer of
	{message.target?.targetName}. Are you sure they're still your allies?
{:else if message.type === PlayerMessageType.FLEET_TRANSFER_INVALID_RECEIVE_REFUSED}
	{$universe.getPlayerPluralName(message.spec?.sourcePlayerNum)} has attempted to gift you {message
		.target?.targetName}, but you have refused their offer. If you wish to receive gifts from this
	player in the future, make sure to set them as allies.
{:else if message.type === PlayerMessageType.FLEET_TRANSFER_RECEIVED}
	{$universe.getPlayerPluralName(message.spec?.sourcePlayerNum)} has given you {message.target
		?.targetName}.
{:else}
	<!-- Fallback for unknown message types -->
	<FallbackMessageDetail {message} />
{/if}
