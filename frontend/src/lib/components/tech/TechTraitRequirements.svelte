<script lang="ts">
	import type { Player } from '$lib/types/cs-proto';
	import { TechCategory, TechOrigin } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { getLabelForLRT, getLabelForPRT, lrts } from '$lib/types/Race';
	import type { TechLike } from '$lib/types/Tech';

	type Props = {
		tech: TechLike;
		player?: Player | undefined;
	};

	let { tech: techLike, player = undefined }: Props = $props();
</script>

<div class="flex flex-col text-base p-1">
	{#if techLike.tech}
		{@const tech = techLike.tech}
		{#if tech.requirements?.hullsAllowed?.length}
			<div class="text-warning">
				{#if tech.requirements?.hullsAllowed.length === 1}
					This {`${enumToString(TechCategory, tech.category).toLowerCase()}`} can only be mounted on
					the {tech.requirements.hullsAllowed[0]} Hull.
				{:else}
					This {`${enumToString(TechCategory, tech.category).toLowerCase()}`} can only be mounted on
					these hulls: {tech.requirements?.hullsAllowed.join(', ')}.
				{/if}
			</div>
		{/if}
		{#if tech.requirements?.hullsDenied?.length}
			<div class="text-warning">
				{#if tech.requirements?.hullsDenied.length === 1}
					This {`${enumToString(TechCategory, tech.category).toLowerCase()}`} cannot be mounted on the
					{tech.requirements.hullsDenied[0]} Hull.
				{:else}
					This {`${enumToString(TechCategory, tech.category).toLowerCase()}`} cannot be mounted on these
					hulls: {tech.requirements?.hullsDenied.join(', ')}.
				{/if}
			</div>
		{/if}
		{#if tech.requirements?.prtsRequired?.length}
			<div
				class:text-error={player?.race &&
					tech.requirements?.prtsRequired.indexOf(player.race.prt) == -1}
			>
				{#if tech.requirements?.prtsRequired?.length > 1}
					This part requires the Primary Racial traits {tech.requirements?.prtsRequired
						.map((prt) => getLabelForPRT(prt))
						.join(' or ')}.
				{:else}
					This part requires the Primary Racial trait {getLabelForPRT(
						tech.requirements?.prtsRequired[0]
					)}.
				{/if}
			</div>
		{/if}

		{#if tech.requirements?.prtsDenied?.length}
			<div
				class:text-error={player?.race &&
					tech.requirements?.prtsDenied.indexOf(player.race.prt) != -1}
			>
				{#if tech.requirements?.prtsDenied?.length > 1}
					This part will not be available to the Primary Racial traits {tech.requirements?.prtsDenied
						.map((prt) => getLabelForPRT(prt))
						.join(' or ')}.
				{:else}
					This part will not be available to the Primary Racial trait {getLabelForPRT(
						tech.requirements?.prtsDenied[0]
					)}.
				{/if}
			</div>
		{/if}

		{#each lrts as lrt (lrt)}
			{#if tech.requirements?.lrtsRequired && (tech.requirements?.lrtsRequired & lrt) > 0}
				<div
					class:text-error={player?.race && (!player.race.lrts || (player.race.lrts & lrt) == 0)}
				>
					This part requires the Lesser Racial trait {getLabelForLRT(lrt)}.
				</div>
			{/if}

			{#if tech.requirements?.lrtsDenied && (tech.requirements?.lrtsDenied & lrt) > 0}
				<div class:text-error={player?.race && player.race.lrts && (player.race.lrts & lrt) > 0}>
					This part will be unavailable if you have the Lesser Racial trait {getLabelForLRT(lrt)}.
				</div>
			{/if}
		{/each}

		{#if (tech.origin ?? TechOrigin.UNSPECIFIED) !== TechOrigin.UNSPECIFIED}
			<div>The origin of this part is unknown.</div>
		{/if}
	{/if}
</div>
