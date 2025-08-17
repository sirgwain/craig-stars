<script lang="ts">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import {
		MinefieldType,
		Prt,
		TechCategory,
		TechHullComponentSchema,
		type TechHullComponent
	} from '$lib/types/cs-proto';
	import { CommandedPlayer } from '$lib/types/Player';
	import { HullSlotTypeShield } from '$lib/types/Consts';
	import { create } from '@bufbuild/protobuf';
	import TestBreadcrumb from '../TestBreadcrumb.svelte';

	const settlersDelight: TechHullComponent = create(TechHullComponentSchema, {
		tech: {
			name: "Settler's Delight",
			cost: {
				ironium: 1,
				germanium: 1,
				resources: 2
			},
			requirements: {
				prtsRequired: [Prt.HE],
				hullsAllowed: ['Mini-Colony Ship']
			},
			ranking: 69,
			category: TechCategory.ENGINE,
			tags: {
				Engine: true,
				Ramscoop: true
			}
		},
		hullSlotType: 2,
		mass: 2,
		idealSpeed: 6,
		freeSpeed: 6,
		maxSafeSpeed: 9,
		fuelUsage: [0, 0, 0, 0, 0, 0, 0, 150, 275, 480, 576]
	});

	const moleSkin: TechHullComponent = create(TechHullComponentSchema, {
		tech: {
			name: 'Mole-skin Shield',
			cost: { ironium: 1, germanium: 1, resources: 4 },
			requirements: {},
			ranking: 10,
			category: TechCategory.SHIELD
		},
		hullSlotType: HullSlotTypeShield,
		mass: 1,
		shield: 25
	});

	const techs: TechHullComponent[] = [
		create(TechHullComponentSchema, {
			tech: {
				name: 'Ferret Scanner',
				cost: {
					ironium: 2,
					germanium: 8,
					resources: 36
				},
				requirements: {
					techLevel: {
						energy: 3,
						electronics: 7,
						biotechnology: 2
					},
					lrtsDenied: 256
				},
				ranking: 80,
				category: TechCategory.SCANNER
			},
			hullSlotType: 4,
			mass: 6,
			scanRange: 185,
			scanRangePen: 50
		}),
		create(TechHullComponentSchema, {
			tech: {
				name: 'Mini Gun',
				cost: {
					boranium: 6,
					resources: 6
				},
				requirements: {
					techLevel: {
						weapons: 5
					},
					prtsRequired: [Prt.IS]
				},
				ranking: 20,
				category: TechCategory.BEAM_WEAPON
			},
			hullSlotType: 2048,
			mass: 3,
			power: 16,
			range: 2,
			initiative: 12,
			gatling: true,
			hitsAllTargets: true
		}),
		create(TechHullComponentSchema, {
			tech: {
				name: 'Mine Dispenser 40',
				cost: {
					ironium: 2,
					boranium: 9,
					germanium: 7,
					resources: 40
				},
				requirements: {
					prtsRequired: [Prt.SD]
				},
				ranking: 10,
				category: TechCategory.MINE_LAYER
			},
			hullSlotType: 8192,
			mass: 25,
			minefieldType: MinefieldType.STANDARD,
			mineLayingRate: 40
		}),
		create(TechHullComponentSchema, {
			tech: {
				name: 'Fuel Mizer',
				cost: {
					ironium: 8,
					resources: 11
				},
				requirements: {
					techLevel: {
						propulsion: 2
					},
					lrtsRequired: 1
				},
				ranking: 65,
				category: TechCategory.ENGINE,
				tags: {
					Engine: true,
					Ramscoop: true
				}
			},
			hullSlotType: 2,
			mass: 6,
			idealSpeed: 6,
			freeSpeed: 4,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 0, 0, 0, 35, 120, 175, 235, 360, 420]
		})
	];

	const testPlayer = new CommandedPlayer();
</script>

<TestBreadcrumb title="Techs" />
<div class="flex flex-wrap gap-2 justify-evenly">
	<div>
		<TechSummary tech={settlersDelight} />
	</div>

	<div>
		<TechSummary tech={moleSkin} />
	</div>

	{#each techs as tech (tech.tech?.name)}
		<div>
			<TechSummary {tech} player={testPlayer} />
		</div>
	{/each}
</div>
