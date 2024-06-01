<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const { game, player, universe, settings } = getGameContext();
	const scale = getContext<Writable<number>>('scale');

	export let planet: Planet;
	export let yOffset: number;

	$: orbitingFleets = $universe
		.getMapObjectsByPosition(planet)
		.filter((mo) => mo.type === MapObjectType.Fleet);

	$: orbitingTokens = orbitingFleets
		.map((of) => of as Fleet)
		.reduce(
			(count, f) =>
				count + (f.tokens ? f.tokens.reduce((tokenCount, t) => tokenCount + t.quantity, 0) : 0),
			0
		);
</script>

{#if $settings.showFleetTokenCounts && orbitingTokens}
	<!-- translate the group to the location of the fleet so when we scale the text it is around the center-->
	<g transform={`translate(${$xGet(planet)} ${$yGet(planet) + yOffset + 20 / $scale})`}>
		<text transform={`scale(${1 / $scale})`} text-anchor="middle" class="fill-base-content"
			>{orbitingTokens}</text
		>
	</g>
{/if}
