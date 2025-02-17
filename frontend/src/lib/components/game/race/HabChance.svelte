<script lang="ts">
	import { getHabChance } from '$lib/types/Race';
	import { type Race } from '$lib/types/cs';

	type Props = {
		race: Race;
	};

	let { race }: Props = $props();

	let habChance = $derived(getHabChance(race));
	let approximateHabitablePlanetRatio = $derived(Math.floor(1 / habChance));
</script>

{#if habChance == 1}
	All planets will be habitable to your race.
{:else if approximateHabitablePlanetRatio == 1}
	Virtually all planets will be habitable to your race.
{:else}
	{`You can expect that 1 in ${approximateHabitablePlanetRatio} planets will be habitable to your race.`}
{/if}
