<script lang="ts">
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
		<div class="font-semibold mr-5">Armor</div>
		<div>{spec.armor ?? 0}dp</div>
	</div>
	<div class="flex justify-between">
		<div class="font-semibold mr-5">Shields</div>
		<div>{spec.shields ?? 0}dp</div>
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
		<div>{scanRange(spec.scanRange)}/{scanRange(spec.scanRangePen)}</div>
	</div>
</div>
