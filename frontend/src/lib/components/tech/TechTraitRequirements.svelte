<script lang="ts">
	import type { PlayerResponse } from '$lib/types/Player';
	import { getLabelForLRT, getLabelForPRT, LRT, PRT } from '$lib/types/Race';
	import { $enum as eu } from 'ts-enum-util';

	import type { Tech } from '$lib/types/Tech';

	export let tech: Tech;
	export let player: PlayerResponse | undefined = undefined;

	const lrts = eu(LRT).getValues();
</script>

<div class="flex flex-col text-base">
	{#if tech.requirements.prtRequired && tech.requirements.prtRequired != PRT.None}
		<div class:text-error={player && player.race.prt != tech.requirements.prtRequired}>
			This part requires the Primary Racial trait {getLabelForPRT(tech.requirements.prtRequired)}
		</div>
	{/if}

	{#if tech.requirements.prtDenied && tech.requirements.prtDenied != PRT.None}
		<div class:text-error={player && player.race.prt == tech.requirements.prtDenied}>
			This part will not be available to the Primary Racial trait {getLabelForPRT(
				tech.requirements.prtDenied
			)}
		</div>
	{/if}

	{#each lrts as lrt}
		{#if tech.requirements.lrtsRequired && (tech.requirements.lrtsRequired & lrt) > 0}
			<div class:text-error={player && (!player.race.lrts || (player.race.lrts & lrt) == 0)}>
				This part requires the Lesser Racial trait {getLabelForLRT(lrt)}
			</div>
		{/if}

		{#if tech.requirements.lrtsDenied && (tech.requirements.lrtsDenied & lrt) > 0}
			<div class:text-error={player && player.race.lrts && (player.race.lrts & lrt) > 0}>
				This part will be unavailable if you have the Lesser Racial trait {getLabelForLRT(lrt)}
			</div>
		{/if}
	{/each}
</div>
