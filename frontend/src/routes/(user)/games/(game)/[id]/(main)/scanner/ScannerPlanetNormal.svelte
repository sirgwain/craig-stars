<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, owned } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import ScannerFleetCount from './ScannerFleetCount.svelte';

	const { settings } = getGameContext();
	const { game, player, universe } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let planet: Planet;
	export let commanded = false;

	let props = {};
	let ringProps: any | undefined = undefined;

	$: hasStarbase = planet.spec?.hasStarbase;
	$: hasMassDriver = planet.spec?.hasMassDriver;
	$: hasStargate = planet.spec?.hasStargate;

	$: planetX = $xGet(planet);
	$: planetY = $yGet(planet);

	$: radius = owned(planet) ? (commanded ? 3 : 1.5) : commanded ? 2 : 1;
	$: ringRadius = radius * 1.5;
	$: strokeWidth = owned(planet) ? (commanded ? 0.5 : 0.3) : 0;

	$: starbaseWidth = commanded ? 3 : 2;
	$: starbaseXOffset = ringRadius * 0.75;
	$: starbaseYOffset = ringRadius + starbaseWidth;

	$: orbitingFleets = $universe
		.getMapObjectsByPosition(planet)
		.filter((mo) => mo.type === MapObjectType.Fleet);

	// setup props for planet circle
	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#999999';
		let strokeColor = '#999999';

		if (planet.playerNum === $player.num) {
			color = '#00FF00';
		} else if (planet.playerNum) {
			color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
		} else if (planet.reportAge !== Unexplored && !planet.playerNum) {
			color = '#FFF';
		}

		// setup the properties of our planet circle
		props = {
			r: radius,
			fill: color,
			stroke: strokeColor,
			'stroke-width': strokeWidth
		};
	}

	// setup props for the ring
	$: {
		// if anything is orbiting our planet, put a ring on it
		if (orbitingFleets?.length > 0) {
			let ringColor = '#FFFFFF';
			let strokeDashArray = '';
			const playerNums = new Set<number>(orbitingFleets.map((f) => f.playerNum));
			if (playerNums.size == 1) {
				if (orbitingFleets[0].playerNum === $player.num) {
					ringColor = '#FFFFFF';
				} else if ($player.isEnemy(orbitingFleets[0].playerNum)) {
					ringColor = '#FF0000';
				} else {
					ringColor = '#FFFF00';
				}
			} else {
				ringColor = '#6A0DAD'; // both
				strokeDashArray = '1 1';
			}

			ringProps = {
				cx: planetX,
				cy: planetY,
				stroke: ringColor,
				'stroke-dasharray': strokeDashArray,
				'stroke-width': 0.5,
				r: ringRadius,
				'fill-opacity': 0
			};
		} else {
			ringProps = undefined;
		}
	}
</script>

{#if ringProps}
	<circle {...ringProps} />
{/if}
<circle cx={planetX} cy={planetY} {...props} />
{#if hasStarbase}
	<rect
		class="starbase"
		width={starbaseWidth}
		height={starbaseWidth}
		rx={0.5}
		x={planetX + starbaseXOffset}
		y={planetY - starbaseYOffset}
	/>
{/if}
{#if hasStargate}
	<rect
		class="stargate"
		width={starbaseWidth}
		height={starbaseWidth}
		rx={0.5}
		x={planetX - starbaseXOffset - starbaseWidth}
		y={planetY - starbaseYOffset}
	/>
{/if}
{#if hasMassDriver}
	<rect
		class="massdriver"
		width={starbaseWidth}
		height={starbaseWidth}
		rx={0.5}
		x={planetX - starbaseWidth / 2}
		y={planetY - starbaseYOffset - starbaseWidth / 2}
	/>
{/if}
<ScannerFleetCount {planet} yOffset={ringRadius} />
