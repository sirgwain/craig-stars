<script lang="ts">
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { techs } from '$lib/services/Stores';
	import { canLearnTech } from '$lib/types/Player';
	import type { Tech } from '$lib/types/Tech';
	import { TechField, get, hasRequiredLevels, minus, sum } from '$lib/types/TechLevel';

	type FutureTech = {
		tech: Tech;
		distance: number;
	};

	const { game, player, universe } = getGameContext();

	export let field: TechField;

	$: currentLevel = get($player.techLevels, field);

	$: futureTechs = $techs.techs
		.filter(
			(tech) =>
				get(tech.requirements, field) > currentLevel &&
				canLearnTech($player, tech) &&
				!hasRequiredLevels($player.techLevels, tech.requirements)
		)
		.map((tech) => {
			const distanceToLearn = minus(tech.requirements, $player.techLevels);
			// zero out any level differences we have already achieved
			// i.e. if we are at level 5 for energy and this tech requires 3, distanceToLearn.Energy will equal -2
			// this makes it zero

			(distanceToLearn.energy = Math.max(0, distanceToLearn.energy ?? 0)),
				(distanceToLearn.weapons = Math.max(0, distanceToLearn.weapons ?? 0)),
				(distanceToLearn.propulsion = Math.max(0, distanceToLearn.propulsion ?? 0)),
				(distanceToLearn.construction = Math.max(0, distanceToLearn.construction ?? 0)),
				(distanceToLearn.electronics = Math.max(0, distanceToLearn.electronics ?? 0)),
				(distanceToLearn.biotechnology = Math.max(0, distanceToLearn.biotechnology ?? 0));

			if (sum(distanceToLearn) == get(distanceToLearn, field)) {
				// if the required tech difference is only in the field we care about
				// add it to our list of future techs
				return { tech, distance: get(distanceToLearn, field) };
			}
		})
		.filter((t) => t != undefined)
		.sort((t1, t2) => (t1?.distance ?? 0) - (t2?.distance ?? 0)) as FutureTech[];
</script>

<ul class="pl-1 pt-1">
	{#each futureTechs as futureTech}
		<li
			class:text-queue-item-this-year={futureTech.distance <= 1}
			class:text-queue-item-next-year={futureTech.distance == 2}
			class="cursor-help"
		>
			<button
				type="button"
				class="w-full h-full text-left"
				on:pointerdown|preventDefault={(e) => onTechTooltip(e, futureTech.tech, true)}
				>{futureTech.tech.name}</button
			>
		</li>
	{:else}
		None
	{/each}
</ul>
