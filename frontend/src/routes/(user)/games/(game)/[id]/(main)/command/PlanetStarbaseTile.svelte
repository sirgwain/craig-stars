<script lang="ts">
	import { game, showDesignPopup, showPopupTech, techs } from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import CommandTile from './CommandTile.svelte';

	export let starbase: Fleet | undefined;

	$: stargate = starbase?.spec?.stargate
		? $techs.getHullComponent(starbase.spec.stargate)
		: undefined;

	$: massDriver = starbase?.spec?.massDriver
		? $techs.getHullComponent(starbase.spec.massDriver)
		: undefined;

	function showDesign(e: MouseEvent) {
		if ($game && starbase?.tokens && starbase.tokens.length > 0) {
			showDesignPopup(
				$game?.player.getDesign($game.player.num, starbase?.tokens[0].designNum) as
					| ShipDesign
					| undefined,
				e.x,
				e.y
			);
		}
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
				<div>{starbase.spec.armor}</div>
			</div>
			<div class="flex justify-between">
				<div>Shields</div>
				<div>{starbase.spec.shield ?? 'none'}</div>
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
			<div class="flex justify-between">
				<div>Stargate</div>
				{#if stargate}
					<div>
						<button
							type="button"
							class="w-full h-full cursor-help"
							on:pointerdown|preventDefault={(e) =>
								stargate && showPopupTech($techs.getHullComponent(stargate.name), e.x, e.y)}
						>
							{stargate.name}
						</button>
					</div>
				{:else}
					<div>none</div>
				{/if}
			</div>
			<div class="flex justify-between">
				<div>Mass Driver</div>
				{#if starbase.spec.hasMassDriver}
					<div>
						<button
							type="button"
							class="w-full h-full cursor-help"
							on:pointerdown|preventDefault={(e) =>
								massDriver && showPopupTech($techs.getHullComponent(massDriver.name), e.x, e.y)}
						>
							Warp {starbase.spec.safePacketSpeed}
						</button>
					</div>
				{:else}
					<div>none</div>
				{/if}
			</div>
			<div class="flex justify-between">
				<div>Destination</div>
				<div>none</div>
			</div>
		</div>
	</CommandTile>
{/if}
