<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Fleet } from '$lib/types/Fleet';
	import { getContext } from 'svelte';
	import ScannerFleet from './ScannerFleet.svelte';
	import type { LayerCake } from 'layercake';
	import { commandedFleet} from '$lib/services/Context';
	import type { FullGame } from '$lib/services/FullGame';

	const game = getContext<FullGame>('game');
	const { data } = getContext<LayerCake>('LayerCake');

	let fleets: Fleet[] = [];

	$: fleets = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.Fleet);
</script>

<!-- Fleets -->
{#each fleets as fleet}
	<ScannerFleet
		{fleet}
		color={game.getPlayerColor(fleet.playerNum)}
		commanded={$commandedFleet?.num === fleet.num && $commandedFleet?.playerNum === fleet.playerNum}
	/>
{/each}
