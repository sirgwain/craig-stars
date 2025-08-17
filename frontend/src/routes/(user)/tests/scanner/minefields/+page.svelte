<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import {
		MapObjectType,
		MinefieldSchema,
		MinefieldType,
		type Minefield
	} from '$lib/types/cs-proto';
	import type { MapObjectLike } from '$lib/types/MapObject';
	import { create } from '@bufbuild/protobuf';
	import { LayerCake, Svg } from 'layercake';
	import ScannerMinefield from '../../../games/(game)/[id]/(main)/scanner/ScannerMinefield.svelte';
	import ScannerMinefieldPattern from '../../../games/(game)/[id]/(main)/scanner/ScannerMinefieldPattern.svelte';

	const { selectMapObject } = getGameContext();

	const minefields: Minefield[] = [
		create(MinefieldSchema, {
			mapObject: {
				type: MapObjectType.MINEFIELD,
				position: {
					x: 50,
					y: 50
				},
				name: `Humanoid Minefield #1`,
				num: 1,
				playerNum: 1
			},
			minefieldType: MinefieldType.STANDARD,
			numMines: 100
		}),
		create(MinefieldSchema, {
			mapObject: {
				type: MapObjectType.MINEFIELD,
				position: {
					x: 0,
					y: 50
				},
				name: `Humanoid Minefield #2`,
				num: 2,
				playerNum: 1
			},
			minefieldType: MinefieldType.STANDARD,
			numMines: 200
		})
	];

	const xGetter = (mo: MapObjectLike) => mo?.mapObject?.position?.x;
	const yGetter = (mo: MapObjectLike) => mo?.mapObject?.position?.y;

	selectMapObject(minefields[0]);
</script>

<div class="w-[300px] h-[300px] bg-black">
	<LayerCake
		data={minefields}
		x={xGetter}
		y={yGetter}
		xDomain={[0, 100]}
		yDomain={[0, 100]}
		xRange={[0, 300]}
		yRange={[300, 0]}
	>
		<Svg>
			<g>
				<ScannerMinefieldPattern />
				<ScannerMinefield minefield={minefields[0]} color="#FF0000" />
				<ScannerMinefield minefield={minefields[1]} color="#00FF00" />
			</g>
		</Svg>
	</LayerCake>
</div>
