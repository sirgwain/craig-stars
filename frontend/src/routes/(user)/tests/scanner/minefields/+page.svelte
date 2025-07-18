<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import {
		MapObjectTypeMinefield,
		MinefieldTypeStandard,
		type MapObject,
		type Minefield
	} from '$lib/types/cs';
	import { LayerCake, Svg } from 'layercake';
	import ScannerMinefield from '../../../games/(game)/[id]/(main)/scanner/ScannerMinefield2.svelte';
	import ScannerMinefieldPattern from '../../../games/(game)/[id]/(main)/scanner/ScannerMinefieldPattern2.svelte';

	const { selectMapObject } = getGameContext();

	const minefields: Minefield[] = [
		{
			type: MapObjectTypeMinefield,
			position: {
				x: 50,
				y: 50
			},
			name: `Humanoid Minefield #1`,
			num: 1,
			playerNum: 1,
			minefieldType: MinefieldTypeStandard,
			numMines: 100,
			spec: {
				decayRate: 100,
				radius: Math.sqrt(100),
				canDetonate: false
			}
		},
		{
			type: MapObjectTypeMinefield,
			position: {
				x: 0,
				y: 50
			},
			name: `Humanoid Minefield #2`,
			num: 2,
			playerNum: 1,
			minefieldType: MinefieldTypeStandard,
			numMines: 200,
			spec: {
				decayRate: 100,
				radius: Math.sqrt(200),
				canDetonate: false
			}
		}
	];

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

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
