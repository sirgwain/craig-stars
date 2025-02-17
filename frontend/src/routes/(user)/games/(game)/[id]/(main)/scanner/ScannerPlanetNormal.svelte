<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MapObjectTypeFleet, ReportAgeUnexplored, type Fleet, type Planet } from '$lib/types/cs';
	import { filterFleet } from '$lib/types/Filter';
	import { owned } from '$lib/types/MapObject';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import { getEnemiesAndFriends } from './Scanner';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';

	const { settings } = getGameContext();
	const { player, universe } = getGameContext();

	type Props = {
		planet: Planet;
		commanded?: boolean;
	};

	let { planet, commanded = false }: Props = $props();

	let hasStarbase = planet.spec?.hasStarbase;
	let hasMassDriver = planet.spec?.hasMassDriver;
	let hasStargate = planet.spec?.hasStargate;

	let radius = $derived(owned(planet) ? (commanded ? 6 : 3) : commanded ? 4 : 2);
	let strokeWidth = $derived(commanded ? 1 : 0.5);
	let ringRadius = $derived(radius * 2.5);
	let ringWidth = $derived(commanded ? 2 : 1.5);

	let starbaseWidth = $derived(commanded ? 6 : 4);
	let starbaseXOffset = $derived(ringRadius * 0.75);
	let starbaseYOffset = $derived(ringRadius + starbaseWidth);

	let orbitingFleets = $derived(
		$universe
			.getMapObjectsByPosition(planet)
			.filter((mo) => mo.type === MapObjectTypeFleet)
			.filter((f) => filterFleet($player, f as Fleet, $settings))
	);

	// setup props for planet circle
	let circleProps = $derived.by(() => {
		// green for us, gray for unexplored, white for explored
		let color = '#999999';
		let strokeColor = '#999999';

		if (planet.playerNum === $player.num) {
			color = '#00FF00';
		} else if (planet.playerNum) {
			color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
		} else if (planet.reportAge !== ReportAgeUnexplored && !planet.playerNum) {
			color = '#FFF';
		}

		// setup the properties of our planet circle
		return {
			r: radius,
			fill: color,
			stroke: strokeColor,
			'stroke-width': strokeWidth
		};
	});

	// setup props for the ring
	let ringProps = $derived.by(() => {
		// if anything is orbiting our planet, put a ring on it
		if (orbitingFleets?.length > 0) {
			const { enemies, friends } = getEnemiesAndFriends(orbitingFleets, $player);

			let ringColor = 'stroke-orbit';
			let strokeDashArray = '';

			if (friends && !enemies) {
				ringColor = 'stroke-orbit-friends';
			} else if (!friends && enemies) {
				ringColor = 'stroke-orbit-enemies';
			} else if (friends && enemies) {
				ringColor = 'stroke-orbit-friends-and-enemies';
				strokeDashArray = '10 6';
			}

			return {
				class: ringColor,
				'stroke-dasharray': strokeDashArray,
				'stroke-width': ringWidth,
				r: ringRadius,
				'fill-opacity': 0
			};
		}
	});
</script>

<MapObjectScaler mapObject={planet}>
	{#if ringProps}
		<circle {...ringProps} />
	{/if}
	<circle {...circleProps} />
	{#if hasStarbase}
		<rect
			class:starbase={planet.spec?.dockCapacity}
			class:starbase-fort={!planet.spec?.dockCapacity}
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
