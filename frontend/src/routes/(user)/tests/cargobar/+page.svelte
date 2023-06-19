<script lang="ts">
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import type { Cargo } from '$lib/types/Cargo';
import TestBreadcrumb from '../TestBreadcrumb.svelte';

	const cargo1: Cargo = {
		ironium: 1,
		boranium: 2,
		germanium: 3,
		colonists: 4
	};
</script>

<TestBreadcrumb title="Cargo Bar" />

<h1 class="text-xl">Cargo 1</h1>
<div class="flex justify-start">
	<div class="w-20">Cargo:</div>
	<div class="w-52"><CargoBar value={cargo1} capacity={10} /></div>
</div>

<h1 class="text-xl">Cargo 1 with more capacity</h1>
<CargoBar value={cargo1} capacity={20} />
