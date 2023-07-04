<script lang="ts" context="module">
	import type { Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	export type HabTooltipProps = {
		player: Player;
		planet: Planet;
		habType: HabType;
	};
</script>

<script lang="ts">
	import { getHabValue, getHabValueString, withHabValue, type HabType, add } from '$lib/types/Hab';
	import { getPlanetHabitability } from '$lib/types/Race';

	export let player: Player;
	export let planet: Planet;
	export let habType: HabType;

	const currentHab = getHabValue(planet.hab, habType);
	const terraformedHab = getHabValue(planet.spec.terraformAmount, habType);
	const habString = getHabValueString(habType, currentHab);
	const terraformedHabString = getHabValueString(habType, terraformedHab);
	const habLowString = getHabValueString(habType, getHabValue(player.race.habLow, habType));
	const habHighString = getHabValueString(habType, getHabValue(player.race.habHigh, habType));
	const habAfterTerraforming = add(planet.hab ?? {}, withHabValue(habType, terraformedHab));
	const habitabilityAfterTerraforming = getPlanetHabitability(player.race, habAfterTerraforming);
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div>
		{habType} is currently
		<span class="font-semibold">{getHabValueString(habType, getHabValue(planet.hab, habType))}</span
		>. Your colonists prefer planets where {habType}
		is between <span class="font-semibold">{habLowString}</span> and
		<span class="font-semibold">{habHighString}</span>

		{#if terraformedHab != 0}
			<div>
				You currently possess the technology to modify the {habType} on
				<span class="font-semibold">{planet.name}</span>
				within the range of within the range of
				<span class="font-semibold"
					>{currentHab <= terraformedHab ? habString : terraformedHabString}</span
				>
				to
				<span class="font-semibold"
					>{currentHab > terraformedHab ? habString : terraformedHabString}</span
				>. If you were to terraform <span class="font-semibold">{habType}</span> to
				<span class="font-semibold">{terraformedHabString}</span>, the planet's value would improve
				to
				<span class="font-semibold">{habitabilityAfterTerraforming.toFixed()}%</span>
			</div>
		{/if}
	</div>
</div>
