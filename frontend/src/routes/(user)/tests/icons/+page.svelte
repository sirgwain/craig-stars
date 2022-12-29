<script lang="ts">
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
</div>
