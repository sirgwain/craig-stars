<script lang="ts">
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import PlanetWithStarbase from '$lib/components/icons/PlanetWithStarbase.svelte';
	import SelectedMapObject from '$lib/components/icons/SelectedMapObject.svelte';
</script>

<div class="flex flex-wrap justify-evenly">
	<div class="flex flex-col">
		<div class="text-center">PlanetWithStarbase</div>
		<div class="m-auto"><PlanetWithStarbase class="w-16 h-16 border" /></div>
	</div>
	<div class="flex flex-col">
		<div class="text-center">SelectedMapObject</div>
		<div class="m-auto"><SelectedMapObject class="w-16 h-16 border" /></div>
	</div>
	<div class="flex flex-col">
		<div class="text-center">Filter Fleets</div>
		<div class="m-auto"><FilterFleets class="w-16 h-16 border" /></div>
	</div>
	<div class="flex flex-col">
		<div class="text-center">Filter Enemy Fleets</div>
		<div class="m-auto"><FilterFleets primary="fill-fleet-enemy" class="w-16 h-16 border" /></div>
	</div>
	<div class="flex flex-col">
		<div class="text-center">Idle Fleets</div>
		<div class="m-auto"><IdleFleets class="w-16 h-16 border" /></div>
	</div>
</div>
