<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import type { ChangeMassDriverSpeedProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import { UnlimitedSpaceDock } from '$lib/types/cs';
	import type { Fleet } from '$lib/types/cs';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import type { ShipDesign } from '$lib/types/cs';
	import CommandTile from './CommandTile.svelte';

	const { game, player, universe, settings } = getGameContext();

	type Props = {
		starbase: Fleet | undefined;
		planet: CommandedPlanet;
	} & ChangeMassDriverSpeedProps;

	let { starbase, planet, onChangeMassDriverSpeed }: Props = $props();

	let stargate = $derived(
		starbase?.spec?.stargate ? $techs.getHullComponent(starbase.spec.stargate) : undefined
	);

	let massDriver = $derived(
		starbase?.spec?.massDriver ? $techs.getHullComponent(starbase.spec.massDriver) : undefined
	);

	function showDesign(e: PointerEvent) {
		e.preventDefault();
		if (starbase?.tokens && starbase.tokens.length > 0) {
			onShipDesignTooltip(
				e,
				$universe.getDesign($player.num, starbase?.tokens[0].designNum) as ShipDesign | undefined
			);
		}
	}
</script>

{#if starbase?.spec}
	<CommandTile title={starbase.baseName}>
		<div class="cursor-help" onpointerdown={showDesign}>
			<div class="flex justify-between">
				<div class="text-tile-item-title">Dock Capacity</div>
				{#if starbase.spec.spaceDock === UnlimitedSpaceDock}
					<div>Unlimited</div>
				{:else if (starbase.spec.spaceDock ?? 0) > 0}
					<div>{starbase.spec.spaceDock}kT</div>
				{:else}
					<div>none</div>
				{/if}
			</div>
			<div class="flex justify-between">
				<div class="text-tile-item-title">Armor</div>
				<div>{starbase.spec.armor}dp</div>
			</div>
			<div class="flex justify-between">
				<div class="text-tile-item-title">Shields</div>
				<div>{starbase.spec.shields ? starbase.spec.shields + 'dp' : 'none'}</div>
			</div>
			{#if starbase.tokens && starbase.tokens.length > 0}
				<div class="flex justify-between">
					<div class="text-tile-item-title">Damage</div>
					{#if !starbase.tokens[0].damage}
						<div>none</div>
					{:else}
						<div>{starbase.tokens[0].damage}%</div>
					{/if}
				</div>
			{/if}
			<div class="divider p-0 m-0"></div>
		</div>
		<div>
			<div
				class="flex justify-between cursor-help"
				onpointerdown={(e) => stargate && onTechTooltip(e, stargate)}
			>
				<div class="text-tile-item-title">Stargate</div>
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
				onpointerdown={(e) => massDriver && onTechTooltip(e, massDriver)}
			>
				<div class="text-tile-item-title">Mass Driver</div>
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
					<div class="text-tile-item-title">Destination</div>
					<div>
						{$universe.getPlanet(planet.packetTargetNum)?.name ?? 'none'}
					</div>
				</div>
				<div class="flex justify-between mt-1 gap-1">
					<div class="w-32">
						<button
							onclick={() => ($settings.setPacketDest = !$settings.setPacketDest)}
							class:btn-accent={$settings.setPacketDest}
							type="button"
							class="btn btn-outline btn-sm normal-case btn-secondary p-2">Set Dest</button
						>
					</div>
					<div class="w-full my-auto">
						<WarpSpeedGauge
							value={planet.packetSpeed}
							isPacket={true}
							min={5}
							max={(planet.spec.basePacketSpeed ?? 0) + $game.rules.packetMaxOverwarpSpeed}
							warnSpeed={(planet.spec.safePacketSpeed ?? 0) + 1}
							dangerSpeed={(planet.spec.safePacketSpeed ?? 0) + 3}
							onValueDragged={(warpSpeed) => {
								planet.packetSpeed = warpSpeed;
							}}
							onValueChanged={(warpSpeed) => {
								planet.packetSpeed = warpSpeed;
								onChangeMassDriverSpeed?.({ planet, warpSpeed });
							}}
						/>
					</div>
				</div>
			{/if}
		</div>
	</CommandTile>
{/if}
