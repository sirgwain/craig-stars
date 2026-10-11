<script lang="ts">
	import type { Race, Rules } from '#lib/types/cs-proto.js';
	import { getHabitablePlanetOdds } from '#lib/types/Race.js';

	type Props = {
		race: Race;
		// the game's rules, or the default rules if the race isn't in a game
		rules?: Rules;
	};

	let { race, rules }: Props = $props();

	let odds = $derived(getHabitablePlanetOdds(race, rules));
</script>

{#if odds === 0}
	All planets will be habitable to your race.
{:else if odds === 1}
	Virtually all planets will be habitable to your race.
{:else}
	{`You can expect that 1 in ${odds} planets will be habitable to your race.`}
{/if}
