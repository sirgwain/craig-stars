<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import type { Fleet } from '$lib/types/Fleet';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import CommandTile from './CommandTile.svelte';

	const { game, player, universe, settings, updatePlanetOrders } = getGameContext();

	export let starbase: Fleet | undefined;
	export let planet: CommandedPlanet;

	$: stargate = starbase?.spec?.stargate
		? $techs.getHullComponent(starbase.spec.stargate)
		: undefined;

	$: massDriver = starbase?.spec?.massDriver
		? $techs.getHullComponent(starbase.spec.massDriver)
		: undefined;

	function showDesign(e: PointerEvent) {
		if (starbase?.tokens && starbase.tokens.length > 0) {
			onShipDesignTooltip(
				e,
				$universe.getDesign($player.num, starbase?.tokens[0].designNum) as ShipDesign | undefined
			);
		}
	}

	function updatePlanetOrdrers() {
		updatePlanetOrders(planet);
	}
</script>

{#if starbase?.spec}
	<CommandTile title={starbase.baseName}>
		<!-- svelte-ignore a11y-click-events-have-key-events -->
		<div class="cursor-help" on:pointerdown|preventDefault={showDesign}>
			<div class="flex justify-between">
				<div>Dock Capacity</div>
				{#if starbase.spec.spaceDock === UnlimitedSpaceDock}
					<div>Unlimited</div>
				{:else if (starbase.spec.spaceDock ?? 0) > 0}
					<div>{starbase.spec.spaceDock}kT</div>
				{:else}
					<div>none</div>
				{/if}
			</div>
			<div class="flex justify-between">
				<div>Armor</div>
				<div>{starbase.spec.armor}dp</div>
			</div>
			<div class="flex justify-between">
				<div>Shields</div>
				<div>{starbase.spec.shields ? starbase.spec.shields + 'dp' : 'none'}</div>
			</div>
			<div class="flex justify-between">
				<div>Damage</div>
				{#if !starbase.damage}
					<div>none</div>
				{:else}
					<div>{starbase.damage}%</div>
				{/if}
			</div>
			<div class="divider p-0 m-0" />
		</div>
		<div>
			<div
				class="flex justify-between cursor-help"
				on:pointerdown|preventDefault={(e) => stargate && onTechTooltip(e, stargate)}
			>
				<div>Stargate</div>
				{#if stargate}
					<div>
						<button type="button" class="w-full h-full">
							{stargate.name}
						</button>
					</div>
				{:else}
					<div>none</div>
				{/if}
			</div>
			<div
				class="flex justify-between cursor-help"
				on:pointerdown|preventDefault={(e) => massDriver && onTechTooltip(e, massDriver)}
			>
				<div>Mass Driver</div>
				{#if starbase.spec.hasMassDriver}
					<div>
						<button type="button" class="w-full h-full">
							Warp {starbase.spec.safePacketSpeed}
						</button>
					</div>
				{:else}
					<div>none</div>
				{/if}
			</div>
			{#if starbase.spec.hasMassDriver}
				<div class="flex justify-between">
					<div>Destination</div>
					<div>
						{$universe.getPlanet(planet.packetTargetNum)?.name ?? 'none'}
					</div>
				</div>
				<div class="flex justify-between mt-1 gap-1">
					<div class="w-32">
						<button
							on:click={() => ($settings.setPacketDest = !$settings.setPacketDest)}
							class:btn-accent={$settings.setPacketDest}
							type="button"
							class="btn btn-outline btn-sm normal-case btn-secondary p-2">Set Dest</button
						>
					</div>
					<div class="w-full my-auto">
						<WarpSpeedGauge
							bind:value={planet.packetSpeed}
							min={planet.spec.safePacketSpeed}
							max={(planet.spec.safePacketSpeed ?? 0) + 3}
							warnSpeed={(planet.spec.safePacketSpeed ?? 0) + 1}
							dangerSpeed={(planet.spec.safePacketSpeed ?? 0) + 3}
							on:valuechanged={() => updatePlanetOrdrers()}
						/>
					</div>
				</div>
			{/if}
		</div>
	</CommandTile>
{/if}
