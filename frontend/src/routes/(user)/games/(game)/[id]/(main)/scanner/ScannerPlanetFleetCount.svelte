<script lang="ts">
	import { MapObjectType } from '$lib/types/cs-proto';
	import type { PlanetIntel } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet } from '$lib/services/Universe';
	import { filterFleet } from '$lib/types/Filter';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { getEnemiesAndFriends, getScannerContext } from './Scanner';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const { player, universe, settings } = getGameContext();
	const { scale } = getScannerContext();

	type Props = {
		planet: PlanetIntel;
		yOffset: number;
	};

	let { planet, yOffset }: Props = $props();

	let orbitingFleets = $derived(
		$universe
			.getMapObjectsByPosition(planet.mapObject?.position)
			.filter((mo) => mo.mapObject?.type === MapObjectType.FLEET)
	);

	let orbitingTokens = $derived(
		orbitingFleets
			.map((of) => of as AnyFleet)
			.filter((f: AnyFleet) => filterFleet($player, f, $settings))
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
	<g
		transform={`translate(${$xGet(planet.mapObject)} ${$yGet(planet.mapObject) + yOffset + 20 / $scale})`}
	>
		<text transform={`scale(${1 / $scale})`} text-anchor="middle" class={textColor}
			>{orbitingTokens}</text
		>
	</g>s
{/if}
