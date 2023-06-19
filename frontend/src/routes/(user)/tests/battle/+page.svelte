<script lang="ts">
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import {
		TokenActionType,
		type BattleRecord,
		BattleTactic,
		BattleTarget,
		BattleAttackWho
	} from '$lib/types/Battle';
	import { Player, type PlayerResponse } from '$lib/types/Player';

	const player = new Player(1, 1, {
		color: '#0000FF',
		techLevels: { energy: 5, weapons: 4, propulsion: 4, construction: 4, electronics: 4 },
		designs: [
			{
				id: 29,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 1,
				playerNum: 1,
				name: 'Long Range Scout',
				version: 0,
				hull: 'Scout',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'Rhino Scanner', hullSlotIndex: 2, quantity: 1 },
					{ hullComponent: 'Fuel Tank', hullSlotIndex: 3, quantity: 1 }
				],
				purpose: 'Scout',
				spec: {
					hullType: 'Scout',
					idealSpeed: 6,
					engine: 'Long Hump 6',
					fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080],
					numEngines: 1,
					cost: { ironium: 17, boranium: 2, germanium: 7, resources: 22 },
					mass: 25,
					armor: 20,
					fuelCapacity: 300,
					scanRange: 66,
					scanRangePen: 30,
					torpedoInaccuracyFactor: 1,
					initiative: 1,
					movement: 10,
					scanner: true,
					numInstances: 9,
					numBuilt: 9
				}
			},
			{
				id: 30,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 2,
				playerNum: 1,
				name: 'Santa Maria',
				version: 0,
				hull: 'Colony Ship',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'Colonization Module', hullSlotIndex: 2, quantity: 1 }
				],
				purpose: 'Colonizer',
				spec: {
					hullType: 'Colonizer',
					idealSpeed: 6,
					engine: 'Long Hump 6',
					fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080],
					numEngines: 1,
					cost: { ironium: 23, boranium: 8, germanium: 21, resources: 30 },
					mass: 61,
					armor: 20,
					fuelCapacity: 200,
					cargoCapacity: 25,
					scanRange: -1,
					scanRangePen: -1,
					torpedoInaccuracyFactor: 1,
					movement: 8,
					colonizer: true,
					numInstances: 1,
					numBuilt: 1
				}
			},
			{
				id: 31,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 3,
				playerNum: 1,
				name: 'Teamster',
				version: 0,
				hull: 'Medium Freighter',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'Rhino Scanner', hullSlotIndex: 2, quantity: 1 },
					{ hullComponent: 'Crobmnium', hullSlotIndex: 3, quantity: 1 }
				],
				purpose: 'Freighter',
				spec: {
					hullType: 'Freighter',
					idealSpeed: 6,
					engine: 'Long Hump 6',
					fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080],
					numEngines: 1,
					cost: { ironium: 34, germanium: 22, resources: 62 },
					mass: 130,
					armor: 125,
					fuelCapacity: 450,
					cargoCapacity: 210,
					scanRange: 50,
					torpedoInaccuracyFactor: 1,
					movement: 10,
					scanner: true,
					numInstances: 1,
					numBuilt: 1
				}
			},
			{
				id: 32,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 4,
				playerNum: 1,
				name: 'Cotton Picker',
				version: 0,
				hull: 'Mini-Miner',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'Rhino Scanner', hullSlotIndex: 2, quantity: 1 },
					{ hullComponent: 'Robo-Mini-Miner', hullSlotIndex: 3, quantity: 1 },
					{ hullComponent: 'Robo-Mini-Miner', hullSlotIndex: 4, quantity: 1 }
				],
				purpose: 'Miner',
				spec: {
					hullType: 'Miner',
					idealSpeed: 6,
					engine: 'Long Hump 6',
					fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080],
					numEngines: 1,
					cost: { ironium: 88, germanium: 23, resources: 243 },
					mass: 574,
					armor: 130,
					fuelCapacity: 210,
					scanRange: 50,
					torpedoInaccuracyFactor: 1,
					movement: 2,
					scanner: true,
					miningRate: 8,
					numInstances: 1,
					numBuilt: 1
				}
			},
			{
				id: 33,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 5,
				playerNum: 1,
				name: 'Armored Probe',
				version: 0,
				hull: 'Scout',
				hullSetNumber: 1,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'Rhino Scanner', hullSlotIndex: 2, quantity: 1 },
					{ hullComponent: 'X-Ray Laser', hullSlotIndex: 3, quantity: 1 }
				],
				purpose: 'FighterScout',
				spec: {
					hullType: 'Scout',
					idealSpeed: 6,
					engine: 'Long Hump 6',
					fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080],
					numEngines: 1,
					cost: { ironium: 12, boranium: 8, germanium: 7, resources: 24 },
					mass: 23,
					armor: 20,
					fuelCapacity: 50,
					scanRange: 66,
					scanRangePen: 30,
					torpedoInaccuracyFactor: 1,
					initiative: 1,
					movement: 10,
					powerRating: 16,
					scanner: true,
					mineSweep: 16,
					hasWeapons: true,
					weaponSlots: [{ hullComponent: 'X-Ray Laser', hullSlotIndex: 3, quantity: 1 }],
					numInstances: 1,
					numBuilt: 1
				}
			},
			{
				id: 34,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 6,
				playerNum: 1,
				name: 'Stalwart Defender',
				version: 0,
				hull: 'Destroyer',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'X-Ray Laser', hullSlotIndex: 2, quantity: 1 },
					{ hullComponent: 'X-Ray Laser', hullSlotIndex: 3, quantity: 1 },
					{ hullComponent: 'Rhino Scanner', hullSlotIndex: 4, quantity: 1 },
					{ hullComponent: 'Crobmnium', hullSlotIndex: 5, quantity: 2 },
					{ hullComponent: 'Fuel Tank', hullSlotIndex: 6, quantity: 1 },
					{ hullComponent: 'Battle Computer', hullSlotIndex: 7, quantity: 1 }
				],
				purpose: 'Fighter',
				spec: {
					hullType: 'Fighter',
					idealSpeed: 6,
					engine: 'Long Hump 6',
					fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080],
					numEngines: 1,
					cost: { ironium: 40, boranium: 15, germanium: 20, resources: 91 },
					mass: 162,
					armor: 350,
					fuelCapacity: 530,
					scanRange: 66,
					scanRangePen: 30,
					torpedoInaccuracyFactor: 0.8,
					initiative: 4,
					movement: 10,
					powerRating: 32,
					scanner: true,
					mineSweep: 32,
					hasWeapons: true,
					weaponSlots: [
						{ hullComponent: 'X-Ray Laser', hullSlotIndex: 2, quantity: 1 },
						{ hullComponent: 'X-Ray Laser', hullSlotIndex: 3, quantity: 1 }
					],
					numInstances: 1,
					numBuilt: 1
				}
			},
			{
				id: 35,
				gameId: 3,
				createdAt: '2023-05-06T23:12:53Z',
				updatedAt: '2023-05-06T23:12:53Z',
				num: 7,
				playerNum: 1,
				name: 'Starbase',
				version: 0,
				hull: 'Space Station',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Laser', hullSlotIndex: 2, quantity: 8 },
					{ hullComponent: 'Mole-skin Shield', hullSlotIndex: 3, quantity: 8 },
					{ hullComponent: 'Laser', hullSlotIndex: 4, quantity: 8 },
					{ hullComponent: 'Mole-skin Shield', hullSlotIndex: 6, quantity: 8 },
					{ hullComponent: 'Laser', hullSlotIndex: 8, quantity: 8 },
					{ hullComponent: 'Laser', hullSlotIndex: 10, quantity: 8 }
				],
				purpose: 'Starbase',
				spec: {
					hullType: 'Starbase',
					fuelUsage: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
					cost: { ironium: 122, boranium: 263, germanium: 236, resources: 752 },
					mass: 48,
					armor: 500,
					scanRange: -1,
					scanRangePen: -1,
					repairBonus: 0.15,
					torpedoInaccuracyFactor: 1,
					initiative: 14,
					powerRating: 320,
					shield: 400,
					starbase: true,
					spaceDock: -1,
					mineSweep: 640,
					hasWeapons: true,
					weaponSlots: [
						{ hullComponent: 'Laser', hullSlotIndex: 2, quantity: 8 },
						{ hullComponent: 'Laser', hullSlotIndex: 4, quantity: 8 },
						{ hullComponent: 'Laser', hullSlotIndex: 8, quantity: 8 },
						{ hullComponent: 'Laser', hullSlotIndex: 10, quantity: 8 }
					],
					numInstances: 1,
					numBuilt: 1
				}
			}
		],
		shipDesignIntels: [
			{
				id: 2,
				gameId: 1,
				num: 1,
				playerNum: 2,
				name: 'Teamster',
				hull: 'Medium Freighter',
				hullSetNumber: 0,
				slots: [
					{ hullComponent: 'Long Hump 6', hullSlotIndex: 1, quantity: 1 },
					{ hullComponent: 'Rhino Scanner', hullSlotIndex: 2, quantity: 1 },
					{ hullComponent: 'Crobmnium', hullSlotIndex: 3, quantity: 1 }
				]
			}
		]
	} as PlayerResponse);

	const battle: BattleRecord = {
		num: 3,
		position: {
			x: 0,
			y: 0
		},
		tokens: [
			{
				num: 1,
				playerNum: 1,
				token: {
					designNum: 6,
					quantity: 1
				},
				startingPosition: {
					x: 1,
					y: 4
				},
				initiative: 4,
				movement: 10,
				tactic: 'MaximizeDamageRatio',
				primaryTarget: 'ArmedShips',
				secondaryTarget: 'Any',
				attackWho: 'EnemiesAndNeutrals'
			},
			{
				num: 2,
				playerNum: 1,
				token: {
					designNum: 1,
					quantity: 1
				},
				startingPosition: {
					x: 1,
					y: 4
				},
				initiative: 1,
				movement: 10,
				tactic: 'MaximizeDamageRatio',
				primaryTarget: 'ArmedShips',
				secondaryTarget: 'Any',
				attackWho: 'EnemiesAndNeutrals'
			},
			{
				num: 3,
				playerNum: 2,
				token: {
					designNum: 1,
					quantity: 1
				},
				startingPosition: {
					x: 8,
					y: 5
				},
				movement: 10,
				tactic: 'MaximizeDamageRatio',
				primaryTarget: 'ArmedShips',
				secondaryTarget: 'Any',
				attackWho: 'EnemiesAndNeutrals'
			}
		],
		actionsPerRound: [
			[],
			[
				{
					type: 3,
					tokenNum: 1,
					round: 1,
					from: {
						x: 1,
						y: 4
					},
					to: {
						x: 2,
						y: 5
					}
				},
				{
					type: 3,
					tokenNum: 3,
					round: 1,
					from: {
						x: 8,
						y: 5
					},
					to: {
						x: 7,
						y: 4
					}
				},
				{
					type: 3,
					tokenNum: 2,
					round: 1,
					from: {
						x: 1,
						y: 4
					},
					to: {
						x: 0,
						y: 4
					}
				},
				{
					type: 3,
					tokenNum: 1,
					round: 1,
					from: {
						x: 2,
						y: 5
					},
					to: {
						x: 3,
						y: 4
					}
				},
				{
					type: 3,
					tokenNum: 3,
					round: 1,
					from: {
						x: 7,
						y: 4
					},
					to: {
						x: 8,
						y: 4
					}
				},
				{
					type: 3,
					tokenNum: 2,
					round: 1,
					from: {
						x: 0,
						y: 4
					},
					to: {
						x: 0,
						y: 3
					}
				}
			],
			[
				{
					type: 3,
					tokenNum: 1,
					round: 2,
					from: {
						x: 3,
						y: 4
					},
					to: {
						x: 4,
						y: 5
					}
				},
				{
					type: 3,
					tokenNum: 3,
					round: 2,
					from: {
						x: 8,
						y: 4
					},
					to: {
						x: 9,
						y: 4
					}
				},
				{
					type: 3,
					tokenNum: 2,
					round: 2,
					from: {
						x: 0,
						y: 3
					},
					to: {
						x: 0,
						y: 2
					}
				},
				{
					type: 3,
					tokenNum: 1,
					round: 2,
					from: {
						x: 4,
						y: 5
					},
					to: {
						x: 5,
						y: 4
					}
				},
				{
					type: 3,
					tokenNum: 3,
					round: 2,
					from: {
						x: 9,
						y: 4
					},
					to: {
						x: 9,
						y: 3
					}
				},
				{
					type: 3,
					tokenNum: 2,
					round: 2,
					from: {
						x: 0,
						y: 2
					},
					to: {
						x: 1,
						y: 2
					}
				},
				{
					type: 3,
					tokenNum: 1,
					round: 2,
					from: {
						x: 5,
						y: 4
					},
					to: {
						x: 6,
						y: 3
					}
				},
				{
					type: 3,
					tokenNum: 3,
					round: 2,
					from: {
						x: 9,
						y: 3
					},
					to: {
						x: 9,
						y: 2
					}
				},
				{
					type: 3,
					tokenNum: 2,
					round: 2,
					from: {
						x: 1,
						y: 2
					},
					to: {
						x: 1,
						y: 1
					}
				},
				{
					type: 2,
					tokenNum: 1,
					round: 2,
					from: {
						x: 6,
						y: 3
					},
					to: {
						x: 9,
						y: 2
					},
					slot: 2,
					targetNum: 3,
					torpedoMisses: 1
				}
			],
			[
				{
					type: 3,
					tokenNum: 1,
					round: 3,
					from: {
						x: 6,
						y: 3
					},
					to: {
						x: 7,
						y: 2
					}
				},
				{
					type: 3,
					tokenNum: 3,
					round: 3,
					from: {
						x: 9,
						y: 2
					},
					to: {
						x: 9,
						y: 1
					}
				},
				{
					type: 3,
					tokenNum: 2,
					round: 3,
					from: {
						x: 1,
						y: 1
					},
					to: {
						x: 2,
						y: 1
					}
				},
				{
					type: 3,
					tokenNum: 1,
					round: 3,
					from: {
						x: 7,
						y: 2
					},
					to: {
						x: 8,
						y: 1
					}
				},
				{
					type: 4,
					tokenNum: 3,
					round: 3,
					from: {
						x: 0,
						y: 0
					},
					to: {
						x: 0,
						y: 0
					}
				},
				{
					type: 4,
					tokenNum: 2,
					round: 3,
					from: {
						x: 0,
						y: 0
					},
					to: {
						x: 0,
						y: 0
					}
				}
			],
			[]
		]
	};
</script>

<h1 class="text-xl">Battle</h1>
<BattleView {player} battleRecord={battle} />
