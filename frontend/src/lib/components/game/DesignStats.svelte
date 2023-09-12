<script lang="ts">
	import { Infinite } from '$lib/types/MapObject';
	import { MineFieldType } from '$lib/types/MineField';
	import type { ShipDesign, Spec } from '$lib/types/ShipDesign';
	import { NoScanner } from '$lib/types/Tech';

	export let spec: Spec;

	function scanRange(range: number | undefined) {
		return !range || range === NoScanner ? '-' : range;
	}
</script>

<div class="flex flex-col min-w-[8rem] mr-2">
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Mass</div>
		<div>{spec.mass ?? 0}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Max Fuel</div>
		<div>{spec.fuelCapacity ?? 0}mg</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Est Range</div>

		<div>
			{#if spec.estimatedRange == Infinite}
				Infinite
			{:else if spec.cargoCapacity}
				{spec.estimatedRange ?? 0}ly/{spec.estimatedRangeFull ?? 0}ly
			{:else}
				{spec.estimatedRange ?? 0}ly
			{/if}
		</div>
	</div>

	{#if spec.cargoCapacity}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Cargo Capacity</div>
			<div>{spec.cargoCapacity}kT</div>
		</div>
	{/if}
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Armor</div>
		<div>{spec.armor ?? 0}dp</div>
	</div>
	{#if spec.shields}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Shields</div>
			<div>{spec.shields}dp</div>
		</div>
	{/if}

	{#if spec.powerRating}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Rating</div>
			<div>{spec.powerRating}</div>
		</div>
	{/if}
	{#if spec.cloakPercent || spec.torpedoJamming}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Cloak/Jam</div>
			<div>{spec.cloakPercent ?? 0}%/{((spec.torpedoJamming ?? 0) * 100).toFixed()}%</div>
		</div>
	{/if}
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Initiative/Moves</div>
		<div>{spec.initiative ?? 0}/{spec.movement ?? 0}</div>
	</div>
	{#if spec.scanRange != NoScanner || spec.scanRangePen != NoScanner}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Scanner Range</div>
			<div>{scanRange(spec.scanRange)}/{scanRange(spec.scanRangePen)}</div>
		</div>
	{/if}

	{#if spec.mineLayingRateByMineType && spec.mineLayingRateByMineType[MineFieldType.Standard]}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Mine Laying</div>
			<div>{spec.mineLayingRateByMineType[MineFieldType.Standard]} std/yr</div>
		</div>
	{/if}
	{#if spec.mineLayingRateByMineType && spec.mineLayingRateByMineType[MineFieldType.Heavy]}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Mine Laying</div>
			<div>{spec.mineLayingRateByMineType[MineFieldType.Heavy]} hvy/yr</div>
		</div>
	{/if}
	{#if spec.mineLayingRateByMineType && spec.mineLayingRateByMineType[MineFieldType.SpeedBump]}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Mine Laying</div>
			<div>{spec.mineLayingRateByMineType[MineFieldType.SpeedBump]} spd/yr</div>
		</div>
	{/if}
	{#if spec.miningRate}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Remote Mining</div>
			<div>{spec.miningRate}kT/yr</div>
		</div>
	{/if}
</div>
