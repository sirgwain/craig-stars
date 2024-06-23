<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MapObjectType, owned } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';
	import type { Readable, Writable } from 'svelte/store';
	import { clamp } from '$lib/services/Math';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import { type Fleet, idleFleetsFilter } from '$lib/types/Fleet';
	import { find } from 'lodash-es';

	const { settings } = getGameContext();
	const { game, player, universe } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const objectScale = getContext<Readable<number>>('objectScale');

	export let planet: Planet;
	export let commanded = false;

	let props = {};
	let ringProps: any | undefined = undefined;

	$: hasStarbase = planet.spec?.hasStarbase;
	$: hasMassDriver = planet.spec?.hasMassDriver;
	$: hasStargate = planet.spec?.hasStargate;

	$: radius = owned(planet) ? (commanded ? 6 : 3) : commanded ? 4 : 2;
	$: strokeWidth = commanded ? 1 : .5;
	$: ringRadius = radius * 2.5;
	$: ringWidth = commanded ? 2 : 1.5;

	$: starbaseWidth = commanded ? 6 : 4;
	$: starbaseXOffset = ringRadius * 0.75;
	$: starbaseYOffset = ringRadius + starbaseWidth;

	$: orbitingFleets = $universe
		.getMapObjectsByPosition(planet)
		.filter((mo) => mo.type === MapObjectType.Fleet)
		.filter((f) => idleFleetsFilter(f as Fleet, $settings.showIdleFleetsOnly));

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
				const enemies = find(playerNums, (n: number) => $player.isEnemy(n));
				const friends = find(playerNums, (n: number) => $player.isFriendOrNeutral(n));

				if (friends && !enemies) {
					ringColor = '#FFFF00';
				} else if (!friends && enemies) {
					ringColor = '#FF0000';
				} else {
					ringColor = '#FF00FF';
				}
				strokeDashArray = '10 6';
			}

			ringProps = {
				stroke: ringColor,
				'stroke-dasharray': strokeDashArray,
				'stroke-width': ringWidth,
				r: ringRadius,
				'fill-opacity': 0
			};
		} else {
			ringProps = undefined;
		}
	}
</script>

<MapObjectScaler mapObject={planet}>
	{#if ringProps}
		<circle {...ringProps} />
	{/if}
	<circle {...props} />
	{#if hasStarbase}
		<rect
			class="starbase"
			width={starbaseWidth}
			height={starbaseWidth}
			rx={0.5}
			x={starbaseXOffset}
			y={-starbaseYOffset}
		/>
	{/if}
	{#if hasStargate}
		<rect
			class="stargate"
			width={starbaseWidth}
			height={starbaseWidth}
			rx={0.5}
			x={-starbaseXOffset - starbaseWidth}
			y={-starbaseYOffset}
		/>
	{/if}
	{#if hasMassDriver}
		<rect
			class="massdriver"
			width={starbaseWidth}
			height={starbaseWidth}
			rx={0.5}
			x={-starbaseWidth / 2}
			y={-starbaseYOffset - starbaseWidth / 2}
		/>
	{/if}
</MapObjectScaler>
<ScannerFleetCount {planet} yOffset={ringRadius / 2} />
