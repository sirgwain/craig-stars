<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { population, totalCargo } from '$lib/types/Cargo';
	import { MineralPacketDecayToNothing, ReportAgeUnexplored } from '$lib/types/Consts';
	import {
		PlayerMessageType,
		type MineralPacket,
		type PlayerIntel,
		type PlayerMessage
	} from '$lib/types/cs-proto';
	import { distance, emptyVector } from '$lib/types/Vector';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		message: PlayerMessage;
		mineralPacket: MineralPacket;
		owner: PlayerIntel;
	};

	let { message, mineralPacket, owner }: Props = $props();

	let target = $derived($universe.getPlanet(mineralPacket.targetPlanetNum));
	let eta = $derived(
		target
			? Math.ceil(
					distance(
						mineralPacket.mapObject?.position ?? emptyVector(),
						target.mapObject?.position ?? emptyVector()
					) /
						(mineralPacket.warpSpeed * mineralPacket.warpSpeed)
				)
			: ReportAgeUnexplored
	);
</script>

{#if message.text}
	{message.text}
{:else if message.type === PlayerMessageType.PLANET_BUILT_MINERAL_PACKET}
	Your starbase at {message.spec?.mapObjectTarget?.targetName} has built a new {message.spec
		?.amount}kT mineral packet targeting {target?.mapObject?.name ?? 'unknown'}.
{:else if message.type === PlayerMessageType.MINERAL_PACKET_DISCOVERED}
	A {owner.racePluralName} mineral packet containing {totalCargo(mineralPacket.cargo)}kT of minerals
	has been detected. It is travelling at warp {mineralPacket.warpSpeed} towards {$universe.getPlanet(
		mineralPacket.targetPlanetNum
	)?.mapObject?.name ?? 'unknown'}.
{:else if message.type === PlayerMessageType.MINERAL_PACKET_TARGETTING_PLAYER_DISCOVERED}
	{@const damage = message.spec?.mineralPacketDamage}
	A {owner.racePluralName} mineral packet containing {totalCargo(mineralPacket.cargo)}kT of minerals
	has been detected. It is travelling at warp {mineralPacket.warpSpeed} towards {$universe.getPlanet(
		mineralPacket.targetPlanetNum
	)?.mapObject?.name ?? 'unknown'}.

	<!-- for these messages, damage should never be null -->
	{#if damage}
		<!-- start with safe conditions, we have a catcher, we live on a starbase, etc -->
		{#if target?.spec?.planetStarbaseSpec?.hasStarbase && target.spec.planetStarbaseSpec.safePacketSpeed >= mineralPacket.warpSpeed}
			Fortunately, your starbase's mass driver is more than capable of safely catching this packet.
			Huzzah!
		{:else if damage.uncaught == MineralPacketDecayToNothing}
			Fortunately, this packet will decay into nothingness before it reaches you.
		{:else if $player.race.spec.livesOnStarbases}
			Though this packet will strike the planet, your race lives on starbases and will be unaffected
			by the ensuing collision.
		{:else if damage.killed > 0 || damage.defensesDestroyed > 0}
			<!-- uh oh, this packet will damage us. report how much and when -->
			{#if target?.spec?.planetStarbaseSpec?.hasStarbase}
				{#if damage.killed >= population(target.cargo)}
					Your starbase does not have a powerful enough mass driver to safely catch this packet. The
					entire planet will be annihilated when it strikes in {eta} years.
				{:else}
					Your starbase does not have a powerful enough mass driver to safely catch this packet.
					Approximately {damage.defensesDestroyed} defenses will be destroyed and {damage.killed}
					colonists will be killed when it strikes in {eta} years.
				{/if}
			{:else if damage.killed >= population(target?.cargo)}
				You have no starbase with a mass driver to catch this packet. The entire planet will be
				annihilated when it strikes in {eta} years.
			{:else}
				You have no starbase with a mass driver to catch this packet. Approximately {damage.defensesDestroyed}
				defenses will be destroyed and {damage.killed}
				colonists will be killed when it strikes in {eta} years.
			{/if}
		{:else}
			Fortunately, this packet will cause no damage. Hurrah!
		{/if}
	{/if}
{:else}
	<FallbackMessageDetail {message} />
{/if}
