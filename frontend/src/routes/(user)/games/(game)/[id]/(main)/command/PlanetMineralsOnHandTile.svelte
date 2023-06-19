<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import type { CommandedPlanet, Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';
	import FactoriesTooltip, {
		type FactoriesTooltipProps
	} from '$lib/components/game/tooltips/FactoriesTooltip.svelte';
	import { showTooltip } from '$lib/services/Context';
	import MinesTooltip, {
		type MinesTooltipProps
	} from '$lib/components/game/tooltips/MinesTooltip.svelte';
	import type { Player } from '$lib/types/Player';

	export let player: Player;
	export let planet: CommandedPlanet;

	onMount(() => {
		const unsubscribe = EventManager.subscribeCargoTransferredEvent((mo) => {
			if (planet == mo) {
				// trigger a reaction
				planet.cargo = (mo as CommandedPlanet).cargo;
			}
		});

		return () => unsubscribe();
	});

	function onMinesTooltip(e: PointerEvent) {
		showTooltip<MinesTooltipProps>(e.x, e.y, MinesTooltip, {
			planetName: planet.name,
			mines: planet.mines,
			maxMines: planet.spec.maxMines,
			maxPossibleMines: planet.spec.maxPossibleMines,
			canBuildMines: player.race.spec?.innateMining ?? false
		});
	}

	function onFactoriesTooltip(e: PointerEvent) {
		showTooltip<FactoriesTooltipProps>(e.x, e.y, FactoriesTooltip, {
			planetName: planet.name,
			factories: planet.factories,
			maxFactories: planet.spec.maxFactories,
			maxPossibleFactories: planet.spec.maxPossibleFactories,
			canBuildFactories: player.race.spec?.innateResources ?? false
		});
	}
</script>

{#if planet}
	<CommandTile title="Minerals on Hand">
		<div class="flex justify-between">
			<div class="text-ironium">Ironium</div>
			<div>{planet.cargo.ironium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-boranium">Boranium</div>
			<div>{planet.cargo.boranium}kT</div>
		</div>
		<div class="flex justify-between">
			<div class="text-germanium">Germanium</div>
			<div>{planet.cargo.germanium}kT</div>
		</div>

		<div class="divider p-0 m-0" />

		<div class="flex justify-between cursor-help" on:pointerdown={onMinesTooltip}>
			<div>Mines</div>
			<div>
				{planet.mines} of {planet.spec.maxMines}
			</div>
		</div>
		<div class="flex justify-between cursor-help" on:pointerdown={onFactoriesTooltip}>
			<div>Factories</div>
			<div>
				{planet.factories} of {planet.spec.maxFactories}
			</div>
		</div>
	</CommandTile>
{/if}
