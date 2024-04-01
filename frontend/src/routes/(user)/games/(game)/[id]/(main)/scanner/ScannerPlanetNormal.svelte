<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

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

	$: starbaseWidth = (commanded ? 6 : 4) / $scale;
	$: starbaseXOffset = (commanded ? 5 : 2) / $scale;
	$: starbaseYOffset = (commanded ? 11 : 6) / $scale;
	$: tokenCountOffset = (commanded ? 13 : 8) / $scale;
	$: radius = (commanded ? 7 : 3) / $scale;
	$: ringRadius = (commanded ? 10 : 5) / $scale;
	let orbitingTokens = 0;

	$: {
		const orbitingFleets = $universe
			.getMapObjectsByPosition(planet)
			.filter((mo) => mo.type === MapObjectType.Fleet);

		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#555';

		if (planet.playerNum === $player.num) {
			color = '#00FF00';
			strokeWidth = 1 / $scale;
		} else if (planet.playerNum) {
			color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
		} else if (planet.reportAge !== Unexplored && !planet.playerNum) {
			color = '#FFF';
		}

		// if anything is orbiting our planet, put a ring on it
		if (orbitingFleets?.length > 0) {
			let ringColor = '#fff';
			let strokeDashArray = '';
			const playerNums = new Set<number>(orbitingFleets.map((f) => f.playerNum));
			if (playerNums.size == 1) {
				if (orbitingFleets[0].playerNum === $player.num) {
					ringColor = '#fff';
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
				cx: $xGet(planet),
				cy: $yGet(planet),
				stroke: ringColor,
				'stroke-dasharray': strokeDashArray,
				'stroke-width': 1 / $scale,
				r: ringRadius,
				'fill-opacity': 0
			};

			orbitingTokens = orbitingFleets
				.map((of) => of as Fleet)
				.reduce(
					(count, f) =>
						count + (f.tokens ? f.tokens.reduce((tokenCount, t) => tokenCount + t.quantity, 0) : 0),
					0
				);
		} else {
			ringProps = undefined;
		}

		// setup the properties of our planet circle
		props = {
			r: radius,
			fill: color,
			stroke: strokeColor,
			'stroke-width': strokeWidth
		};
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
{#if $settings.showFleetTokenCounts && orbitingTokens}
	<!-- translate the group to the location of the fleet so when we scale the text it is around the center-->
	<g transform={`translate(${planetX - ringRadius} ${planetY + tokenCountOffset * 2.5})`}>
		<text transform={`scale(${1 / $scale})`} class="fill-base-content">{orbitingTokens}</text>
	</g>
{/if}
