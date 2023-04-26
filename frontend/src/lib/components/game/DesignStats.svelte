<script lang="ts">
	import type { ShipDesign, Spec } from '$lib/types/ShipDesign';

	export let spec: Spec;
</script>

<div class="flex flex-col min-w-[8rem]">
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Mass</div>
		<div>{spec.mass ?? 0}kT</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Max Fuel</div>
		<div>{spec.fuelCapacity ?? 0}mg</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Armor</div>
		<div>{spec.armor ?? 0}dp</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Shields</div>
		<div>{spec.shield ?? 0}dp</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Rating</div>
		<div>{spec.powerRating ?? 0}</div>
	</div>

	<div class="flex justify-between">
		<div class="font-semibold mr-5">Cloak/Jam</div>
		<div>{spec.cloakPercent ?? 0}/0</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Initiative/Moves</div>
		<div>{spec.initiative ?? 0}/{spec.movement ?? 0}</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Scanner Range</div>
		<div>{spec.scanRange ?? 0}/{spec.scanRangePen ?? 0}</div>
	</div>
	<!-- {#if design.purpose}
		<div class="flex justify-between">
			<div class="font-semibold mr-5">Purpose</div>
			<div>{design.purpose}</div>
		</div>
	{/if} -->
</div>
