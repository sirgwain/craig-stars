<script lang="ts">
	import type { PlayerResponse } from '$lib/types/Player';
	import { getLabelForLRT, getLabelForPRT, LRT, PRT } from '$lib/types/Race';
	import { $enum as eu } from 'ts-enum-util';

	import type { Tech } from '$lib/types/Tech';
	import { startCase } from 'lodash-es';

	export let tech: Tech;
	export let player: PlayerResponse | undefined = undefined;

	const lrts = eu(LRT).getValues();
</script>

<div class="flex flex-col text-base p-1">
	{#if tech.requirements.hullsAllowed}
		<div class="text-warning">
			{#if tech.requirements.hullsAllowed.length === 1}
				This {`${startCase(tech.category).toLowerCase()}`} can only be mounted on the {tech
					.requirements.hullsAllowed[0]} Hull.
			{:else}
				This {`${startCase(tech.category).toLowerCase()}`} can only be mounted on these hulls: {tech.requirements.hullsAllowed.join(
					', '
				)}.
			{/if}
		</div>
	{/if}
	{#if tech.requirements.hullsDenied}
		<div class="text-warning">
			{#if tech.requirements.hullsDenied.length === 1}
				This {`${startCase(tech.category).toLowerCase()}`} cannot be mounted on the {tech
					.requirements.hullsDenied[0]} Hull.
			{:else}
				This {`${startCase(tech.category).toLowerCase()}`} cannot be mounted on these hulls: {tech.requirements.hullsDenied.join(
					', '
				)}.
			{/if}
		</div>
	{/if}
	{#if tech.requirements.prtsRequired?.length}
		<div class:text-error={player && tech.requirements.prtsRequired.indexOf(player.race.prt) == -1}>
			{#if tech.requirements.prtsRequired?.length > 1}
				This part requires the Primary Racial traits {tech.requirements.prtsRequired
					.map((prt) => getLabelForPRT(prt))
					.join(' or ')}.
			{:else}
				This part requires the Primary Racial trait {getLabelForPRT(
					tech.requirements.prtsRequired[0]
				)}.
			{/if}
		</div>
	{/if}

	{#if tech.requirements.prtsDenied?.length}
		<div class:text-error={player && tech.requirements.prtsDenied.indexOf(player.race.prt) != -1}>
			{#if tech.requirements.prtsDenied?.length > 1}
				This part will not be available to the Primary Racial traits {tech.requirements.prtsDenied
					.map((prt) => getLabelForPRT(prt))
					.join(' or ')}.
			{:else}
				This part will not be available to the Primary Racial trait {getLabelForPRT(
					tech.requirements.prtsDenied[0]
				)}.
			{/if}
		</div>
	{/if}

	{#each lrts as lrt}
		{#if tech.requirements.lrtsRequired && (tech.requirements.lrtsRequired & lrt) > 0}
			<div class:text-error={player && (!player.race.lrts || (player.race.lrts & lrt) == 0)}>
				This part requires the Lesser Racial trait {getLabelForLRT(lrt)}.
			</div>
		{/if}

		{#if tech.requirements.lrtsDenied && (tech.requirements.lrtsDenied & lrt) > 0}
			<div class:text-error={player && player.race.lrts && (player.race.lrts & lrt) > 0}>
				This part will be unavailable if you have the Lesser Racial trait {getLabelForLRT(lrt)}.
			</div>
		{/if}
	{/each}

	{#if (tech.origin ?? '') != ''}
		<div>The origin of this part is unknown.</div>
	{/if}
</div>
