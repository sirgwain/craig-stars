<script lang="ts" module>
	import type { HabType } from '$lib/types/Hab';
	import type { CommandedPlayer } from '$lib/types/Player';
	import type { Planet } from '$lib/types/cs-proto';

	export type HabTooltipProps = {
		player: CommandedPlayer;
		planet: Planet;
		habType: HabType;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import {
		add,
		emptyHab,
		getHabValue,
		getHabValueString,
		habTypeString,
		withHabValue
	} from '$lib/types/Hab';
	import { isImmune } from '$lib/types/Race';

	const { cs } = getGameContext();
	let { player, planet, habType }: HabTooltipProps = $props();

	const currentHab = getHabValue(planet.hab, habType);
	const terraformedHab = getHabValue(planet.spec?.terraformAmount ?? {}, habType);
	const habString = getHabValueString(habType, currentHab);
	const terraformedHabString = getHabValueString(
		habType,
		getHabValue(add(planet.hab ?? emptyHab(), withHabValue(habType, terraformedHab)), habType)
	);
	const habLowString = getHabValueString(habType, getHabValue(player.race.habLow, habType));
	const habHighString = getHabValueString(habType, getHabValue(player.race.habHigh, habType));
	const habCenter = getHabValue(player.race.spec.habCenter, habType);
	const habAfterTerraforming = add(planet.hab ?? emptyHab(), withHabValue(habType, terraformedHab));
	const habValueAfterTerraforming = getHabValue(habAfterTerraforming, habType);

	let habitabilityAfterTerraforming = $state(0);
	$effect(() => {
		cs.wasmService
			.getPlanetHabitability({ race: player.race, hab: habAfterTerraforming })
			.then((resp) => {
				habitabilityAfterTerraforming = resp.result ?? 0;
			});
	});
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div>
		{#if isImmune(player.race, habType)}
			<p>
				{habTypeString(habType)} is currently
				<span class="font-semibold"
					>{getHabValueString(habType, getHabValue(planet.hab, habType))}</span
				>. Your colonists are immune to the effects of {habTypeString(habType)}.
			</p>
		{:else if currentHab == habCenter}
			<p>
				{habTypeString(habType)} is currently
				<span class="font-semibold"
					>{getHabValueString(habType, getHabValue(planet.hab, habType))}</span
				>
				(perfect). Your colonists prefer planets where {habTypeString(habType)} is between
				<span class="font-semibold">{habLowString}</span> and
				<span class="font-semibold">{habHighString}</span>.
			</p>
		{:else}
			<p>
				{habTypeString(habType)} is currently
				<span class="font-semibold"
					>{getHabValueString(habType, getHabValue(planet.hab, habType))}
					({currentHab < habCenter ? habCenter - currentHab : currentHab - habCenter}%</span
				>
				away from ideal). Your colonists prefer planets where {habTypeString(habType)} is between
				<span class="font-semibold">{habLowString}</span> and
				<span class="font-semibold">{habHighString}</span>.
			</p>

			{#if terraformedHab != 0}
				<p>
					You currently possess the technology to modify the {habTypeString(habType)} on
					<span class="font-semibold">{planet.mapObject?.name}</span> within the range of
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
				</p>
			{/if}
		{/if}
	</div>
</div>
