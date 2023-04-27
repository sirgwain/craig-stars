<script lang="ts">
	import type { CommandedPlanet, Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import { UnlimitedSpaceDock } from '$lib/types/Tech';
	import CommandTile from './CommandTile.svelte';

	export let player: Player;
	export let planet: CommandedPlanet;

	$: starbase = player.starbases.find((sb) => sb.planetNum == planet.num);
</script>

{#if starbase?.spec}
	<CommandTile title={starbase.baseName}>
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
			{#if starbase.damage === 0}
				<div>none</div>
			{:else}
				<div>{starbase.damage}%</div>
			{/if}
		</div>
		<div class="divider p-0 m-0" />

		<div class="flex justify-between">
			<div>Stargate</div>
			{#if starbase.spec.hasStargate}
				<div>{starbase.spec.stargate}</div>
			{:else}
				<div>none</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div>Mass Driver</div>
			{#if starbase.spec.hasMassDriver}
				<div>Warp {starbase.spec.safePacketSpeed}</div>
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
