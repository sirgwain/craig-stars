<script lang="ts">
	import { commandMapObject, player, selectMapObject } from '$lib/services/Context';

	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import { humanoid } from '$lib/types/Race';
	import { normalized } from '$lib/types/Vector';
	import { LayerCake, Svg } from 'layercake';
	import ScannerFleets from '../../../games/(game)/[id]/(main)/scanner/ScannerFleets.svelte';
	import type { PlayerResponse } from '$lib/types/Player';

	type fleetPlacement = {
		x: number;
		y: number;
		headingX: number;
		headingY: number;
		playerNum?: number;
	};

	const fleetPlacements: fleetPlacement[] = [
		{
			x: 10,
			y: 10,
			headingX: 1,
			headingY: 1
		},
		{
			x: 20,
			y: 20,
			headingX: 0,
			headingY: 0
		},
		{
			x: 30,
			y: 30,
			headingX: 1,
			headingY: 0
		},
		{
			x: 40,
			y: 40,
			headingX: 0,
			headingY: 1
		},
		{
			x: 50,
			y: 50,
			headingX: -1,
			headingY: -1,
			playerNum: 1
		},
		{
			x: 60,
			y: 60,
			headingX: -1,
			headingY: 0,
			playerNum: 1
		},
		{
			x: 70,
			y: 70,
			headingX: 0,
			headingY: -1,
			playerNum: 1
		}
	];

	let num = 1;
	const fleets: Fleet[] = fleetPlacements.map((fp) => ({
		type: MapObjectType.Fleet,
		position: {
			x: fp.x,
			y: fp.y
		},
		name: 'Long Range Scout #1',
		num: num++,
		playerNum: fp.playerNum ?? 0,
		baseName: 'Long Range Scout',
		tokens: [
			{
				designNum: 1,
				quantity: 1
			}
		],
		heading: normalized({
			x: fp.headingX,
			y: fp.headingY
		}),
		waypoints: [],
		warpSpeed: 5,
		spec: {}
	}));

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	// create a test player
	player.update(
		() =>
			({
				gameId: 0,
				userId: 0,
				num: 0,
				race: humanoid,
				techLevels: {
					energy: 3,
					weapons: 3,
					propulsion: 3,
					construction: 3,
					electronics: 3,
					biotechnology: 3
				},
				techLevelsSpent: {},
				messages: [],
				designs: [],
				planets: [],
				fleets: fleets,
				planetIntels: [],
				fleetIntels: [],
				color: '#0000FF',
				spec: {}
			} as PlayerResponse)
	);
	selectMapObject(fleets[0]);
	commandMapObject(fleets[0]);
</script>

<div class="w-[300px] h-[300px] bg-black">
	<LayerCake
		data={fleets}
		x={xGetter}
		y={yGetter}
		xDomain={[0, 100]}
		yDomain={[0, 100]}
		xRange={[0, 300]}
		yRange={[300, 0]}
	>
		<Svg>
			<g>
				<ScannerFleets />
			</g>
		</Svg>
	</LayerCake>
</div>
