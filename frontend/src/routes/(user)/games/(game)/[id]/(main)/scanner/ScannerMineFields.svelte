<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerMineField from './ScannerMineField.svelte';

	const game = getContext<FullGame>('game');
	const { data } = getContext<LayerCake>('LayerCake');

	$: minefields = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.MineField);
</script>

<!-- MineFields -->
{#each minefields as mineField}
	<ScannerMineField {mineField} color={game.getPlayerColor(mineField.playerNum)} />
{/each}
