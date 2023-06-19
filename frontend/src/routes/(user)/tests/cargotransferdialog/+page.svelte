<script lang="ts">
	import CargoTransferDialog from '../../games/(game)/[id]/dialogs/cargo/CargoTransferDialog.svelte';

	import { CommandedFleet, type Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { CommandedPlanet, type Planet } from '$lib/types/Planet';

	const planet: Planet = new CommandedPlanet();
	Object.assign(planet, {
		type: MapObjectType.Planet,
		reportAge: 0,
		name: 'Obrin',
		mineralConcentration: { ironium: 10, boranium: 20, germanium: 30 },
		cargo: { ironium: 1000, boranium: 2000, germanium: 4000, colonists: 2500 }
	});

	const fleet = new CommandedFleet({
		name: 'Santa Maria #3',
		baseName: 'Santa Maria',
		cargo: {
			ironium: 11,
			colonists: 3
		},
		fuel: 200,
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
	});
</script>

<h1 class="text-xl">Cargo Transfer</h1>
<div class="md:w-[600px] md:h-[400px] w-full h-full">
	<CargoTransferDialog src={fleet} dest={planet} />
</div>
