<script lang="ts">
	import type { Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import CommandTile from './CommandTile.svelte';

	export let player: Player;
	export let planet: Planet;
</script>

{#if planet.starbase?.spec}
	<CommandTile title={planet.starbase.baseName}>
		<div class="flex justify-between">
			<div>Dock Capacity</div>
			{#if planet.starbase.spec.spaceDock === UnlimitedSpaceDock}
				<div>Unlimited</div>
			{:else if (planet.starbase.spec.spaceDock ?? 0) > 0}
				<div>{planet.starbase.spec.spaceDock}kT</div>
			{:else}
				<div>none</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div>Armor</div>
			<div>{planet.starbase.spec.armor}</div>
		</div>
		<div class="flex justify-between">
			<div>Shields</div>
			<div>{planet.starbase.spec.shield ?? 'none'}</div>
		</div>
		<div class="flex justify-between">
			<div>Damage</div>
			{#if planet.starbase.damage === 0}
				<div>none</div>
			{:else}
				<div>{planet.starbase.damage}%</div>
			{/if}
		</div>
		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Stargate</div>
			{#if planet.starbase.spec.hasStargate}
				<div>{planet.starbase.spec.stargate}</div>
			{:else}
				<div>none</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div>Mass Driver</div>
			{#if planet.starbase.spec.hasMassDriver}
				<div>Warp {planet.starbase.spec.safePacketSpeed}</div>
			{:else}
				<div>none</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div>Destination</div>
			<div>none</div>
		</div>
	</CommandTile>
{/if}
