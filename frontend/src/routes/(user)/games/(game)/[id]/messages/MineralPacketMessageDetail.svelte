<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { totalCargo } from '$lib/types/Cargo';
	import { Unknown } from '$lib/types/MapObject';
	import { MessageType, type Message } from '$lib/types/Message';
	import { MineralPacketDecayToNothing, type MineralPacket } from '$lib/types/MineralPacket';
	import type { PlayerIntel } from '$lib/types/Player';
	import { distance } from '$lib/types/Vector';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let mineralPacket: MineralPacket;
	export let owner: PlayerIntel;

	$: target = $universe.getPlanet(mineralPacket.targetPlanetNum);
	$: eta = target
		? Math.ceil(
				distance(mineralPacket.position, target.position) /
					(mineralPacket.warpSpeed * mineralPacket.warpSpeed)
			)
		: Unknown;
</script>

{#if message.text}
	{message.text}
{:else if message.type === MessageType.PlanetBuiltMineralPacket}
	Your starbase at {message.spec.targetName} has built a new {message.spec.amount}kT mineral packet
	targeting {target?.name ?? 'unknown'}.
{:else if message.type === MessageType.MineralPacketDiscovered}
	A {owner.racePluralName} mineral packet containing {totalCargo(mineralPacket.cargo)}kT of minerals
	has been detected. It is travelling at warp {mineralPacket.warpSpeed} towards {$universe.getPlanet(
		mineralPacket.targetPlanetNum
	)?.name ?? 'unknown'}.
{:else if message.type === MessageType.MineralPacketTargettingPlayerDiscovered}
	{@const damage = message.spec.mineralPacketDamage}
	A {owner.racePluralName} mineral packet containing {totalCargo(mineralPacket.cargo)}kT of minerals
	has been detected. It is travelling at warp {mineralPacket.warpSpeed} towards {$universe.getPlanet(
		mineralPacket.targetPlanetNum
	)?.name ?? 'unknown'}.

	<!-- for these messages, damage should never be null -->
	{#if damage}
		<!-- start with safe conditions, we have a catcher, we live on a starbase, etc -->
		{#if target?.spec.hasStarbase && (target.spec.safePacketSpeed ?? 0) >= mineralPacket.warpSpeed}
			Thankfully, your starbase's mass driver is more than capable of safely catching this packet.
			Huzzah!
		{:else if $player.race.spec?.livesOnStarbases}
			Thankfully, your race lives on starbases and will be unaffected by the ensuing collision.
		{:else if damage.uncaught == MineralPacketDecayToNothing}
			Thankfully, this packet will decay into nothingness before it reaches you.
		{:else if (damage.killed ?? 0) > 0 || (damage.defensesDestroyed ?? 0) > 0}
			<!-- uh oh, this packet will damage us. report how much and when -->
			{#if target?.spec.hasStarbase}
				{#if (damage.killed ?? 0) >= (target?.spec.population ?? 0)}
					Your starbase does not have a powerful enough mass driver to safelt catch this packet. The
					entire planet will be annihilated when it strikes in {eta} years.
				{:else}
					Your starbase does not have a powerful enough mass driver to safely catch this packet.
					Approximately {damage.defensesDestroyed ?? 0} defenses will be destroyed and {damage.killed}
					colonists will be killed when it strikes in {eta} years.
				{/if}
			{:else if (damage.killed ?? 0) >= (target?.spec.population ?? 0)}
				You have no starbase with a mass driver to catch this packet. The entire planet will be
				annihilated when it strikes in {eta} years.
			{:else}
				You have no starbase with a mass driver to catch this packet. Approximately {damage.defensesDestroyed ??
					0} defenses will be destroyed and {damage.killed}
				colonists will be killed when it strikes in {eta} years.
			{/if}
		{:else}
			Thankfully, this packet will cause no damage. Hurrah!
		{/if}
	{/if}
{:else}
	<FallbackMessageDetail {message} />
{/if}
