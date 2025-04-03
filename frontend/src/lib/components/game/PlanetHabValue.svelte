<script lang="ts">
	import type { AnyPlanet } from '$lib/services/Universe';

	type Props = {
		planet: AnyPlanet;
	};

	let { planet }: Props = $props();
</script>

<span
	class:text-habitable={(planet.spec.habitability ?? 0) > 0}
	class:text-uninhabitable={(planet.spec.habitability ?? 0) < 0}
	class:text-terraformable={(planet.spec.habitability ?? 0) < 0 &&
		(planet.spec.terraformedHabitability ?? 0) > 0}
	>{planet.spec.habitability ?? 0}%{planet.spec.terraformedHabitability &&
	planet.spec.terraformedHabitability !== planet.spec.habitability
		? ` (${planet.spec.terraformedHabitability}%)`
		: ''}</span
>
