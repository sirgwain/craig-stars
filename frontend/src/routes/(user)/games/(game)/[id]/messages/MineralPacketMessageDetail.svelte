<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyMineralPacket } from '$lib/services/Universe';
	import { totalCargo } from '$lib/types/Cargo';
	import type { PlayerIntel } from '$lib/types/cs';
	import {
		MineralPacketDecayToNothing,
		PlayerMessageMineralPacketDiscovered,
		PlayerMessageMineralPacketTargettingPlayerDiscovered,
		PlayerMessagePlanetBuiltMineralPacket,
		ReportAgeUnexplored,
		type PlayerMessage
	} from '$lib/types/cs';
	import { distance } from '$lib/types/Vector';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		message: PlayerMessage;
		mineralPacket: AnyMineralPacket;
		owner: PlayerIntel;
	};

	let { message, mineralPacket, owner }: Props = $props();

	let target = $derived($universe.getPlanet(mineralPacket.targetPlanetNum));
	let eta = $derived(
		target
			? Math.ceil(
					distance(mineralPacket.position, target.position) /
						(mineralPacket.warpSpeed * mineralPacket.warpSpeed)
				)
			: ReportAgeUnexplored
	);
</script>

{#if message.text}
	{message.text}
{:else if message.type === PlayerMessagePlanetBuiltMineralPacket}
	Your starbase at {message.spec.targetName} has built a new {message.spec.amount}kT mineral packet
	targeting {target?.name ?? 'unknown'}.
{:else if message.type === PlayerMessageMineralPacketDiscovered}
	A {owner.racePluralName} mineral packet containing {totalCargo(mineralPacket.cargo)}kT of minerals
	has been detected. It is travelling at warp {mineralPacket.warpSpeed} towards {$universe.getPlanet(
		mineralPacket.targetPlanetNum
	)?.name ?? 'unknown'}.
{:else if message.type === PlayerMessageMineralPacketTargettingPlayerDiscovered}
	{@const damage = message.spec.mineralPacketDamage}
	A {owner.racePluralName} mineral packet containing {totalCargo(mineralPacket.cargo)}kT of minerals
	has been detected. It is travelling at warp {mineralPacket.warpSpeed} towards {$universe.getPlanet(
		mineralPacket.targetPlanetNum
	)?.name ?? 'unknown'}.

	<!-- for these messages, damage should never be null -->
	{#if damage}
		<!-- start with safe conditions, we have a catcher, we live on a starbase, etc -->
		{#if target?.spec.hasStarbase && (target.spec.safePacketSpeed ?? 0) >= mineralPacket.warpSpeed}
			Fortunately, your starbase's mass driver is more than capable of safely catching this packet.
			Huzzah!
		{:else if damage.uncaught == MineralPacketDecayToNothing}
			Fortunately, this packet will decay into nothingness before it reaches you.
		{:else if $player.race.spec?.livesOnStarbases}
			Though this packet will strike the planet, your race lives on starbases and will be unaffected
			by the ensuing collision.
		{:else if (damage.killed ?? 0) > 0 || (damage.defensesDestroyed ?? 0) > 0}
			<!-- uh oh, this packet will damage us. report how much and when -->
			{#if target?.spec.hasStarbase}
				{#if (damage.killed ?? 0) >= (target?.spec.population ?? 0)}
					Your starbase does not have a powerful enough mass driver to safely catch this packet. The
					entire planet will be annihilated when it strikes in {eta} years.
				{:else}
					Your starbase does not have a powerful enough mass driver to safely catch this packet.
					Approximately {damage.defensesDestroyed ?? 0} defenses will be destroyed and {damage.killed ??
						0}
					colonists will be killed when it strikes in {eta} years.
				{/if}
			{:else if (damage.killed ?? 0) >= (target?.spec.population ?? 0)}
				You have no starbase with a mass driver to catch this packet. The entire planet will be
				annihilated when it strikes in {eta} years.
			{:else}
				You have no starbase with a mass driver to catch this packet. Approximately {damage.defensesDestroyed ??
					0} defenses will be destroyed and {damage.killed ?? 0}
				colonists will be killed when it strikes in {eta} years.
			{/if}
		{:else}
			Fortunately, this packet will cause no damage. Hurrah!
		{/if}
	{/if}
{:else}
	<FallbackMessageDetail {message} />
{/if}
