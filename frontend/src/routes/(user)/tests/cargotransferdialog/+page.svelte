<script lang="ts">
	import CargoTransferDialog from '../../games/(game)/[id]/dialogs/cargo/CargoTransferDialog.svelte';

	import { CommandedFleet, type Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import { CommandedPlanet, type Planet } from '$lib/types/Planet';
	import techjson from '$lib/ssr/techs.json';

	const planet: Planet = new CommandedPlanet();
	Object.assign(planet, {
		type: MapObjectType.Planet,
		reportAge: 0,
		name: 'Obrin',
		mineralConcentration: { ironium: 10, boranium: 20, germanium: 30 },
		cargo: { ironium: 1000, boranium: 2000, germanium: 4000, colonists: 2500 }
	});

	const fleet = new CommandedFleet();
	Object.assign(fleet, {
		cargo: {
			ironium: 11,
			colonists: 3
		},
		fuel: 200,
		spec: {
			fuelCapacity: 200,
			cargoCapacity: 25
		}
	});
</script>

<h1 class="text-xl">Cargo Transfer</h1>
<div class="md:w-[600px] md:h-[400px] w-full h-full">
	<CargoTransferDialog src={fleet} dest={planet} />
</div>
