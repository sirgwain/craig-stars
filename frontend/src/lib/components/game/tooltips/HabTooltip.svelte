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
	import {
		getHabValue,
		getHabValueString,
		withHabValue,
		type HabType,
		add,
		habTypeString
	} from '$lib/types/Hab';
	import { getPlanetHabitability, isImmune } from '$lib/types/Race';

	export let player: Player;
	export let planet: Planet;
	export let habType: HabType;

	const currentHab = getHabValue(planet.hab, habType);
	const terraformedHab = getHabValue(planet.spec.terraformAmount ?? {}, habType);
	const habString = getHabValueString(habType, currentHab);
	const terraformedHabString = getHabValueString(
		habType,
		getHabValue(add(planet.hab ?? {}, withHabValue(habType, terraformedHab)), habType)
	);
	const habLowString = getHabValueString(habType, getHabValue(player.race.habLow, habType));
	const habHighString = getHabValueString(habType, getHabValue(player.race.habHigh, habType));
	const habCenter = getHabValue(player.race.spec?.habCenter, habType);
	const habAfterTerraforming = add(planet.hab ?? {}, withHabValue(habType, terraformedHab));
	const habValueAfterTerraforming = getHabValue(habAfterTerraforming, habType);
	const habitabilityAfterTerraforming = getPlanetHabitability(player.race, habAfterTerraforming);
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div>
		{#if isImmune(player.race, habType)}
			{habTypeString(habType)} is currently
			<span class="font-semibold"
				>{getHabValueString(habType, getHabValue(planet.hab, habType))}</span
			>.<br />
			Your colonists are immune to the effects of {habTypeString(habType)}.
		{:else}
			{habTypeString(habType)} is currently
			<span class="font-semibold"
				>{getHabValueString(habType, getHabValue(planet.hab, habType))}</span
			>.<br />
			Your colonists prefer planets where {habTypeString(habType)} is between
			<span class="font-semibold">{habLowString}</span> and
			<span class="font-semibold">{habHighString}</span>.

			{#if currentHab != habCenter}
				<br />
				This value is currently
				<span class="font-semibold"
					>{currentHab < habCenter ? habCenter - currentHab : currentHab - habCenter}%</span
				>
				away from the ideal value for your race ({getHabValueString(habType, habCenter)}).

				{#if terraformedHab != 0}
					<br />
					You currently possess the technology to modify the {habTypeString(habType)} on
					<span class="font-semibold">{planet.name}</span> within the range of
					<span class="font-semibold"
						>{currentHab <= habValueAfterTerraforming ? habString : terraformedHabString}</span
					>
					to
					<span class="font-semibold"
						>{currentHab > habValueAfterTerraforming ? habString : terraformedHabString}</span
					>.
					<br />If you were to terraform <span class="font-semibold">{habTypeString(habType)}</span>
					to
					<span class="font-semibold">{terraformedHabString}</span>, the planet's value would
					improve to
					<span class="font-semibold">{habitabilityAfterTerraforming.toFixed()}%</span>.
				{/if}
			{:else}
				This value is perfect for your race.
			{/if}
		{/if}
	</div>
</div>
