<script lang="ts">
	import type { ShipDesign } from '$lib/types/ShipDesign';

	export let design: ShipDesign;
</script>

<div class="flex flex-col min-w-[8rem]">
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Mass</div>
		<div>{design.spec.mass ?? 0}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Max Fuel</div>
		<div>{design.spec.fuelCapacity ?? 0}mg</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Armor</div>
		<div>{design.spec.armor ?? 0}dp</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Shields</div>
		<div>{design.spec.shield ?? 0}dp</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Rating</div>
		<div>{design.spec.powerRating ?? 0}</div>
	</div>

	<div class="flex justify-between">
		<div class="font-semibold mr-5">Cloak/Jam</div>
		<div>{design.spec.cloakPercent ?? 0}/0</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Initiative/Moves</div>
		<div>{design.spec.initiative ?? 0}/{design.spec.movement ?? 0}</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Scanner Range</div>
		<div>{design.spec.scanRange ?? 0}/{design.spec.scanRangePen ?? 0}</div>
	</div>
	{#if design.purpose}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Purpose</div>
			<div>{design.purpose}</div>
		</div>
	{/if}
</div>
