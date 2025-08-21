<script lang="ts">
	import type { Planet } from "$lib/types/cs-proto";


	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();
</script>

<span
	class:text-habitable={(planet.spec?.habitability ?? 0) > 0}
	class:text-uninhabitable={(planet.spec?.habitability ?? 0) < 0}
	class:text-terraformable={(planet.spec?.habitability ?? 0) < 0 &&
		(planet.spec?.terraformedHabitability ?? 0) > 0}
	>{planet.spec?.habitability ?? 0}%{planet.spec?.terraformedHabitability &&
	planet.spec?.terraformedHabitability !== planet.spec?.habitability
		? ` (${planet.spec?.terraformedHabitability}%)`
		: ''}</span
>
