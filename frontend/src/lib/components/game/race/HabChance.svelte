<script lang="ts">
	import { getHabChance, type Race } from '$lib/types/Race';

	export let race: Race;

	$: habChance = getHabChance(race);
	$: approximateHabitablePlanetRatio = Math.floor(1 / habChance);
</script>

{#if habChance == 1}
	All planets will be habitable to your race.
{:else if approximateHabitablePlanetRatio == 1}
	Virtually all planets will be habitable to your race.
{:else}
	{`You can expect that 1 in ${approximateHabitablePlanetRatio} planets will be habitable to your race.`}
{/if}
