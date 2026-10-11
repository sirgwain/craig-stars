<script lang="ts">
	import PlanetNames from '#lib/components/icons/PlanetNames.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';

	const { settings } = getGameContext();
</script>

<a
	href="#planet-names-toggle"
	title="Show Planet Names"
	class:btn-primary={$settings.showNames}
	class:btn-ghost={!$settings.showNames}
	class="btn btn-ghost btn-xs size-full p-0"
	onclick={() => ($settings.showNames = !$settings.showNames)}
	><PlanetNames class="size-6  hover:text-accent" /></a
>
