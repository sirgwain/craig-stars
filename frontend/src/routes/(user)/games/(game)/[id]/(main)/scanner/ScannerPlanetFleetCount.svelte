<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { filterFleet } from '$lib/types/Filter';
	import { type Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { getEnemiesAndFriends, getScannerContext } from './Scanner';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const { game, player, universe, settings } = getGameContext();
	const { scale } = getScannerContext();

	type Props = {
		planet: Planet;
		yOffset: number;
	};

	let { planet, yOffset }: Props = $props();

	let orbitingFleets = $derived(
		$universe.getMapObjectsByPosition(planet).filter((mo) => mo.type === MapObjectType.Fleet)
	);

	let orbitingTokens = $derived(
		orbitingFleets
			.map((of) => of as Fleet)
			.filter((f: Fleet) => filterFleet($player, f, $settings))
			.reduce(
				(count, f) =>
					count + (f.tokens ? f.tokens.reduce((tokenCount, t) => tokenCount + t.quantity, 0) : 0),
				0
			)
	);
	let { enemies, friends } = $derived(getEnemiesAndFriends(orbitingFleets, $player));

	let textColor = $derived.by(() => {
		if (friends && !enemies) {
			return 'fill-orbit-friends';
		} else if (!friends && enemies) {
			return 'fill-orbit-enemies';
		} else if (friends && enemies) {
			return 'fill-orbit-friends-and-enemies';
		}
		return 'fill-orbit';
	});
</script>

{#if $settings.showFleetTokenCounts && orbitingTokens}
	<!-- translate the group to the location of the fleet so when we scale the text it is around the center-->
	<g transform={`translate(${$xGet(planet)} ${$yGet(planet) + yOffset + 20 / $scale})`}>
		<text transform={`scale(${1 / $scale})`} text-anchor="middle" class={textColor}
			>{orbitingTokens}</text
		>
	</g>s
{/if}
