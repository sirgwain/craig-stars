<script lang="ts">
	import FactoriesTooltip, {
		type FactoriesTooltipProps
	} from '$lib/components/game/tooltips/FactoriesTooltip.svelte';
	import MineralTooltip, {
		type MineralTooltipProps
	} from '$lib/components/game/tooltips/MineralTooltip.svelte';
	import MinesTooltip, {
		type MinesTooltipProps
	} from '$lib/components/game/tooltips/MinesTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip } from '$lib/services/Stores';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import CommandTile from './CommandTile.svelte';

	const { player } = getGameContext();
	type Props = {
		planet: CommandedPlanet;
	};

	let { planet }: Props = $props();

	function onIroniumTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Ironium',
			surfaceAmount: planet.cargo.ironium,
			concentration: planet.mineralConcentration.ironium,
			miningRate: planet.spec.miningOutput?.ironium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	function onBoraniumTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Boranium',
			surfaceAmount: planet.cargo.boranium,
			concentration: planet.mineralConcentration.boranium,
			miningRate: planet.spec.miningOutput?.boranium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	function onGermaniumTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Germanium',
			surfaceAmount: planet.cargo.germanium,
			concentration: planet.mineralConcentration.germanium,
			miningRate: planet.spec.miningOutput?.germanium ?? 0,
			homeworld: !!planet.homeworld
		});
	}

	function onMinesTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<MinesTooltipProps>(e.x, e.y, MinesTooltip, {
			planetName: planet.mapObject.name,
			mines: planet.mines,
			maxMines: planet.spec.maxMines,
			maxPossibleMines: planet.spec.maxPossibleMines,
			canBuildMines: $player.race.spec.innateMining
		});
	}

	function onFactoriesTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<FactoriesTooltipProps>(e.x, e.y, FactoriesTooltip, {
			planetName: planet.mapObject.name,
			factories: planet.factories,
			maxFactories: planet.spec.maxFactories,
			maxPossibleFactories: planet.spec.maxPossibleFactories,
			canBuildFactories: $player.race.spec.innateResources
		});
	}
</script>

{#if planet}
	<CommandTile title="Minerals on Hand">
		<div class="flex justify-between cursor-help" onpointerdown={onIroniumTooltip}>
			<div class="text-tile-item-title text-ironium">Ironium</div>
			<div>{planet.cargo.ironium}kT</div>
		</div>
		<div class="flex justify-between cursor-help" onpointerdown={onBoraniumTooltip}>
			<div class="text-tile-item-title text-boranium">Boranium</div>
			<div>{planet.cargo.boranium}kT</div>
		</div>
		<div class="flex justify-between cursor-help" onpointerdown={onGermaniumTooltip}>
			<div class="text-tile-item-title text-germanium">Germanium</div>
			<div>{planet.cargo.germanium}kT</div>
		</div>

		<div class="divider p-0 m-0"></div>

		<div class="flex justify-between cursor-help" onpointerdown={onMinesTooltip}>
			<div class="text-tile-item-title">Mines</div>
			<div>
				{#if $player.race.spec.innateMining}
					{planet.mines}*
				{:else}
					{planet.mines} of {planet.spec.maxMines}
				{/if}
			</div>
		</div>
		<div class="flex justify-between cursor-help" onpointerdown={onFactoriesTooltip}>
			<div class="text-tile-item-title">Factories</div>
			<div>
				{#if $player.race.spec.innateResources}
					n/a
				{:else}
					{planet.factories} of {planet.spec.maxFactories}
				{/if}
			</div>
		</div>
	</CommandTile>
{/if}
