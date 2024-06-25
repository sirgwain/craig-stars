<script lang="ts">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { MineFieldTypes } from '$lib/types/MineField';

	import { Player } from '$lib/types/Player';
	import { PRT } from '$lib/types/Race';
	import {
		HullSlotType,
		TechCategory,
		type TechEngine,
		type TechHullComponent
	} from '$lib/types/Tech';
	import TestBreadcrumb from '../TestBreadcrumb.svelte';

	const settlersDelight: TechEngine = {
		techStoreId: 1,
		name: "Settler's Delight",
		cost: { ironium: 1, germanium: 1, resources: 2 },
		requirements: { prtsRequired: [PRT.HE] },
		ranking: 10,
		category: TechCategory.Engine,
		hullSlotType: HullSlotType.Engine,
		mass: 2,
		idealSpeed: 6,
		freeSpeed: 6,
		fuelUsage: [0, 0, 0, 0, 0, 0, 0, 150, 275, 480, 576]
	};

	const moleSkin: TechHullComponent = {
		name: 'Mole-skin Shield',
		cost: { ironium: 1, germanium: 1, resources: 4 },
		requirements: {},
		ranking: 10,
		category: TechCategory.Shield,
		hullSlotType: HullSlotType.Shield,
		mass: 1,
		shield: 25
	};

	const techs: TechHullComponent[] | TechEngine[] = [
		{
			name: 'Ferret Scanner',
			cost: {
				ironium: 2,
				germanium: 8,
				resources: 36
			},
			requirements: {
				energy: 3,
				electronics: 7,
				biotechnology: 2,
				lrtsDenied: 256
			},
			ranking: 80,
			category: TechCategory.Scanner,
			hullSlotType: 4,
			mass: 6,
			scanRange: 185,
			scanRangePen: 50
		},
		{
			name: 'Mini Gun',
			cost: {
				boranium: 6,
				resources: 6
			},
			requirements: {
				weapons: 5,
				prtsRequired: [PRT.IS]
			},
			ranking: 20,
			category: TechCategory.BeamWeapon,
			hullSlotType: 2048,
			mass: 3,
			power: 16,
			range: 2,
			initiative: 12,
			gattling: true,
			hitsAllTargets: true
		},
		{
			name: 'Mine Dispenser 40',
			cost: {
				ironium: 2,
				boranium: 9,
				germanium: 7,
				resources: 40
			},
			requirements: {
				prtsRequired: [PRT.SD]
			},
			ranking: 10,
			category: TechCategory.MineLayer,
			hullSlotType: 8192,
			mass: 25,
			mineFieldType: MineFieldTypes.Standard,
			mineLayingRate: 40
		},
		{
			name: 'Fuel Mizer',
			cost: {
				ironium: 8,
				resources: 11
			},
			requirements: {
				propulsion: 2,
				lrtsRequired: 1
			},
			ranking: 40,
			category: TechCategory.Engine,
			hullSlotType: 0,
			mass: 6,
			idealSpeed: 6,
			freeSpeed: 4,
			fuelUsage: [0, 0, 0, 0, 0, 35, 120, 175, 235, 360, 420]
		}
	];

	const testPlayer = new Player();
</script>

<TestBreadcrumb title="Techs" />
<div class="flex flex-wrap gap-2 justify-evenly">
	<div>
		<TechSummary tech={settlersDelight} />
	</div>

	<!-- <div>
		<TechSummary tech={moleSkin} />
	</div>

	{#each techs as tech}
		<div>
			<TechSummary {tech} player={testPlayer} />
		</div>
	{/each} -->
</div>
