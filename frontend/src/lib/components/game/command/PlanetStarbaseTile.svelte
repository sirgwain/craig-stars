<script lang="ts">
	import { commandedPlanet } from '$lib/services/Context';
	import { PlanetService } from '$lib/services/PlanetService';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import CommandTile from './CommandTile.svelte';

	const planetService = new PlanetService();
</script>

{#if $commandedPlanet?.starbase}
	<CommandTile title={$commandedPlanet?.starbase.baseName}>
		<div class="flex justify-between">
			<div>Dock Capacity</div>
			{#if $commandedPlanet.starbase.spec.spaceDock === UnlimitedSpaceDock}
				<div>Unlimited</div>
			{:else if $commandedPlanet.starbase.spec.spaceDock > 0}
				<div>{$commandedPlanet.starbase.spec.spaceDock}kT</div>
			{:else}
				<div>none</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div>Armor</div>
			<div>{$commandedPlanet.starbase.spec.armor}</div>
		</div>
		<div class="flex justify-between">
			<div>Shields</div>
			<div>{$commandedPlanet.starbase.spec.shield ?? 'none'}</div>
		</div>
		<div class="flex justify-between">
			<div>Damage</div>
			{#if $commandedPlanet.starbase.damage === 0}
				<div>none</div>
			{:else}
				<div>{$commandedPlanet.starbase.damage}%</div>
			{/if}
		</div>
		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Stargate</div>
			{#if $commandedPlanet.starbase.spec.hasStargate}
				<div>{$commandedPlanet.starbase.spec.stargate}</div>
			{:else}
				<div>none</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div>Mass Driver</div>
			{#if $commandedPlanet.starbase.spec.hasMassDriver}
				<div>Warp {$commandedPlanet.starbase.spec.safePacketSpeed}</div>
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
