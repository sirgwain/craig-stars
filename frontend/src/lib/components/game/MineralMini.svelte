<script lang="ts">
	import type { AnyPlanet } from '$lib/services/Universe';
	import type { Mineral } from '$lib/types/cs';
	import MineralTooltip, {
		type MineralTooltipProps
	} from '$lib/components/game/tooltips/MineralTooltip.svelte';
	import { showTooltip } from '$lib/services/Stores';

	type Props = {
		mineral: Mineral | undefined;
		planet?: AnyPlanet;
		showUnits?: boolean;
	};

	function onIroniumTooltip(e: PointerEvent, planet: AnyPlanet) {
		e.preventDefault();
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Ironium',
			surfaceAmount: planet.cargo?.ironium ?? 0,
			concentration: planet.mineralConcentration?.ironium ?? 0,
			miningRate: planet.spec.miningOutput?.ironium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	function onBoraniumTooltip(e: PointerEvent, planet: AnyPlanet) {
		e.preventDefault();
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Boranium',
			surfaceAmount: planet.cargo?.boranium ?? 0,
			concentration: planet.mineralConcentration?.boranium ?? 0,
			miningRate: planet.spec.miningOutput?.boranium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	function onGermaniumTooltip(e: PointerEvent, planet: AnyPlanet) {
		e.preventDefault();
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Germanium',
			surfaceAmount: planet.cargo?.germanium ?? 0,
			concentration: planet.mineralConcentration?.germanium ?? 0,
			miningRate: planet.spec.miningOutput?.germanium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	let { mineral, planet, showUnits = false }: Props = $props();
</script>

{#if mineral}
	<div class="tracking-wider text-center">
		<span
			class:cursor-help={!!planet}
			onpointerdown={(e) => planet && onIroniumTooltip(e, planet)}
			class="text-ironium">{Math.floor(mineral.ironium ?? 0).toFixed()}</span
		>{showUnits ? ' kT' : ''}
		<span
			class:cursor-help={!!planet}
			onpointerdown={(e) => planet && onBoraniumTooltip(e, planet)}
			class="text-boranium">{Math.floor(mineral.boranium ?? 0).toFixed()}</span
		>{showUnits ? ' kT' : ''}
		<span
			class:cursor-help={!!planet}
			onpointerdown={(e) => planet && onGermaniumTooltip(e, planet)}
			class="text-germanium">{Math.floor(mineral.germanium ?? 0).toFixed()}</span
		>{showUnits ? ' kT' : ''}
	</div>
{/if}
