<script lang="ts">
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import { canLearnTech } from '$lib/types/Player';
	import type { TechLike } from '$lib/types/Tech';
	import { get, hasRequiredLevels, subtract, sum } from '$lib/types/TechLevel';
	import { type TechField } from '$lib/types/cs-proto';

	type FutureTech = {
		tech: TechLike;
		distance: number;
	};

	const { player } = getGameContext();

	type Props = {
		field: TechField;
	};

	let { field }: Props = $props();

	let currentLevel = $derived(get($player.techLevels, field));

	let futureTechs = $derived(
		$techs.techs
			.filter(
				(tech) =>
					get(tech.tech?.requirements?.techLevel, field) > currentLevel &&
					canLearnTech($player, tech) &&
					!hasRequiredLevels($player.techLevels, tech.tech?.requirements?.techLevel)
			)
			.map((tech) => {
				const distanceToLearn = subtract(tech.tech?.requirements?.techLevel, $player.techLevels);
				// zero out any level differences we have already achieved
				// i.e. if we are at level 5 for energy and this tech requires 3, distanceToLearn.Energy will equal -2
				// this makes it zero

				distanceToLearn.energy = Math.max(0, distanceToLearn.energy);
				distanceToLearn.weapons = Math.max(0, distanceToLearn.weapons);
				distanceToLearn.propulsion = Math.max(0, distanceToLearn.propulsion);
				distanceToLearn.construction = Math.max(0, distanceToLearn.construction);
				distanceToLearn.electronics = Math.max(0, distanceToLearn.electronics);
				distanceToLearn.biotechnology = Math.max(0, distanceToLearn.biotechnology);

				if (sum(distanceToLearn) == get(distanceToLearn, field)) {
					// if the required tech difference is only in the field we care about
					// add it to our list of future techs
					return { tech, distance: get(distanceToLearn, field) };
				}
			})
			.filter((t) => t != undefined)
			.sort((t1, t2) => t1.distance - t2.distance) as FutureTech[]
	);
</script>

<ul class="pl-1 pt-1">
	{#each futureTechs as futureTech (futureTech.tech.tech?.name)}
		<li
			class:text-queue-item-this-year={futureTech.distance <= 1}
			class:text-queue-item-next-year={futureTech.distance == 2}
			class="cursor-help"
		>
			<button
				type="button"
				class="w-full h-full text-left"
				onpointerdown={(e) => onTechTooltip(e, futureTech.tech, true)}
				>{futureTech.tech.tech?.name}</button
			>
		</li>
	{:else}
		None
	{/each}
</ul>
