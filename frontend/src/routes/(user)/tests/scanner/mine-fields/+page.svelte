<script lang="ts">
	import { selectMapObject } from '$lib/services/Stores';

	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import { MineFieldType, type MineField } from '$lib/types/MineField';
	import { LayerCake, Svg } from 'layercake';
	import ScannerMineField from '../../../games/(game)/[id]/(main)/scanner/ScannerMineField.svelte';
	import ScannerMineFieldPattern from '../../../games/(game)/[id]/(main)/scanner/ScannerMineFieldPattern.svelte';

	const mineFields: MineField[] = [
		{
			type: MapObjectType.MineField,
			position: {
				x: 50,
				y: 50
			},
			name: `Humanoid MineField #1`,
			num: 1,
			playerNum: 1,
			mineFieldType: MineFieldType.Standard,
			numMines: 100,
			spec: {
				radius: Math.sqrt(100)
			}
		},
		{
			type: MapObjectType.MineField,
			position: {
				x: 0,
				y: 50
			},
			name: `Humanoid MineField #2`,
			num: 2,
			playerNum: 1,
			mineFieldType: MineFieldType.Standard,
			numMines: 200,
			spec: {
				radius: Math.sqrt(200)
			}
		}
	];

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	selectMapObject(mineFields[0]);
</script>

<div class="w-[300px] h-[300px] bg-black">
	<LayerCake
		data={mineFields}
		x={xGetter}
		y={yGetter}
		xDomain={[0, 100]}
		yDomain={[0, 100]}
		xRange={[0, 300]}
		yRange={[300, 0]}
	>
		<Svg>
			<g>
				<ScannerMineFieldPattern />
				<ScannerMineField mineField={mineFields[0]} color="#FF0000" />
				<ScannerMineField mineField={mineFields[1]} color="#00FF00" />
			</g>
		</Svg>
	</LayerCake>
</div>
