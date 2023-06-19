<script lang="ts">
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import PlanetMineralsGraph from '../../games/(game)/[id]/(main)/PlanetMineralsGraph.svelte';

	const planet1: Planet = {
		type: MapObjectType.Planet,
		reportAge: 0,
		name: 'New Planet',
		position: { x: 0, y: 0 },
		num: 0,
		playerNum: 0,
		mineralConcentration: { ironium: 10, boranium: 20, germanium: 30 },
		cargo: { ironium: 1000, boranium: 2000, germanium: 4000, colonists: 2500 }
	};

	const planet2: Planet = {
		type: MapObjectType.Planet,
		reportAge: 0,
		name: 'New Planet2',
		position: { x: 0, y: 0 },
		num: 0,
		playerNum: 0,
		mineralConcentration: { ironium: 130, boranium: 20, germanium: 30 },
		cargo: { ironium: 10000, boranium: 2000, germanium: 4000, colonists: 2500 }
	};
</script>

<h1 class="text-xl">Planet 1</h1>
<PlanetMineralsGraph planet={planet1} />

<h1 class="text-xl">Planet 2</h1>
<PlanetMineralsGraph planet={planet2} />

<h1 class="text-xl">Planet 2 - Small scale</h1>
<PlanetMineralsGraph planet={planet2} max={1000} />
<h1 class="text-xl">Planet 2 - Large scale</h1>
<PlanetMineralsGraph planet={planet2} max={25000} />
