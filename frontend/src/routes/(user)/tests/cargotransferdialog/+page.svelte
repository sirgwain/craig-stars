<script lang="ts">
	import CargoTransferDialog from '../../games/(game)/[id]/dialogs/cargo/CargoTransferDialog.svelte';

	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	
	const planet: Planet = {
		type: MapObjectType.Planet,
		reportAge: 0,
		name: 'Obrin',
		position: { x: 0, y: 0 },
		num: 0,
		playerNum: undefined,
		mineralConcentration: { ironium: 10, boranium: 20, germanium: 30 },
		cargo: { ironium: 1000, boranium: 2000, germanium: 4000, colonists: 2500 }
	};

	const fleet: Fleet = {
		id: 5,
		createdAt: '2022-08-13T18:08:12.897309-07:00',
		updatedAt: '2022-08-14T12:33:39.894998-07:00',
		type: MapObjectType.Fleet,
		gameId: 1,
		position: {
			x: 266,
			y: 752
		},
		name: 'Santa Maria #3',
		num: 3,
		playerNum: 0,
		baseName: 'Santa Maria',
		cargo: {
			ironium: 11,
			colonists: 3
		},
		fuel: 200,
		damage: 0,
		battlePlan: 0,
		tokens: [
			{
				id: 5,
				createdAt: '2022-08-13T18:08:12.897362-07:00',
				updatedAt: '2022-08-14T12:33:39.894893-07:00',
				gameId: 5,
				designId: 4,
				quantity: 1
			}
		],
		waypoints: [
			{
				position: {
					x: 266,
					y: 752
				},
				warpFactor: 5,
				targetPlanetNum: 1,
				targetName: 'Megrez'
			}
		],
		spec: {
			engine: 'Quick Jump 5',
			cost: {
				ironium: 12,
				germanium: 14,
				resources: 21
			},
			mass: 24,
			armor: 20,
			fuelCapacity: 200,
			cargoCapacity: 25,
			scanRange: -1,
			scanRangePen: -1,
			immuneToOwnDetonation: false,
			totalShips: 1,
			massEmpty: 24,
			basePacketSpeed: 0,
			safePacketSpeed: 0,
			baseCloakedCargo: 24
		}
	};
</script>

<h1 class="text-xl">Cargo Transfer</h1>
<div class="md:w-[600px] md:h-[400px] w-full h-full">
	<CargoTransferDialog src={fleet} dest={planet} />
</div>
