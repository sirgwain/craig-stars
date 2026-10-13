<script lang="ts">
	// A made-up battle with every kind of shot (beams and torpedoes, shield and armor damage,
	// misses and kills) to check the battle viewer and its playback without a game or a server.
	import BattleView from '#lib/components/game/battle/BattleView.svelte';
	import Popup from '#lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '#lib/components/game/tooltips/Tooltip.svelte';
	import { Universe } from '#lib/services/Universe.js';
	import {
		BattleRecordSchema,
		PlayerUniverseSchema,
		ShipDesignSchema,
		type BattleRecordTokenActionJson,
		type ShipDesign
	} from '#lib/types/cs-proto.js';
	import { create, fromJson } from '@bufbuild/protobuf';

	type Shooter = {
		weapon?: string;
		at: [number, number];
		targets: number[];
		shields?: number;
		armor?: number;
		hits?: number;
		misses?: number;
		destroyed?: number;
	};

	// targets: player 2 tokens 1..10 in column 5 (row = token - 1)
	const shooters: Shooter[] = [
		{ weapon: 'Laser', at: [1, 0], targets: [1], armor: 20 },
		{ weapon: 'Yakimora Light Phaser', at: [1, 1], targets: [2], shields: 20, armor: 30 },
		{ weapon: 'Heavy Blaster', at: [1, 2], targets: [3], armor: 120 },
		{ weapon: 'Disruptor', at: [1, 3], targets: [4], shields: 60, armor: 100 },
		{ weapon: 'Anti-Matter Pulverizer', at: [1, 4], targets: [5], armor: 400, destroyed: 1 },
		{ weapon: 'Pulsed Sapper', at: [1, 5], targets: [6], shields: 80 },
		{ weapon: 'Gatling Gun', at: [3, 6], targets: [6, 7, 8], shields: 10, armor: 30 },
		{ weapon: 'Bludgeon', at: [5, 7], targets: [8], armor: 230 },
		{ weapon: 'Alpha Torpedo', at: [0, 8], targets: [9], armor: 10, hits: 2, misses: 2 },
		{ weapon: 'Omega Torpedo', at: [1, 8], targets: [9], armor: 600, hits: 6, misses: 2 },
		{
			weapon: 'Jihad Missile',
			at: [0, 9],
			targets: [10],
			shields: 40,
			armor: 80,
			hits: 2,
			misses: 1
		},
		{
			weapon: 'Armageddon Missile',
			at: [1, 9],
			targets: [10],
			armor: 1500,
			hits: 3,
			misses: 1,
			destroyed: 2
		},
		// an enemy design we know nothing about: the generic look
		{ at: [8, 2], targets: [3], armor: 40 }
	];

	const design = (playerNum: number, num: number, weapon: string | undefined): ShipDesign =>
		create(ShipDesignSchema, {
			playerNum,
			num,
			name: weapon ?? 'Target',
			hull: weapon ? 'Destroyer' : 'Battleship',
			slots: weapon ? [{ hullSlotIndex: 2, hullComponent: weapon, quantity: 2 }] : []
		});

	const designs = [
		...shooters.filter((s) => s.weapon).map((s, i) => design(1, i + 1, s.weapon)),
		design(2, 1, undefined),
		design(2, 2, undefined)
	];

	const targets = Array.from({ length: 10 }, (_, i) => ({
		num: i + 1,
		playerNum: 2,
		designNum: 1,
		position: { x: 5, y: i },
		armor: 500,
		startingQuantity: 4
	}));
	let designNum = 0;
	const shooterTokens = shooters.map((s, i) => ({
		num: targets.length + i + 1,
		playerNum: s.weapon ? 1 : 2,
		designNum: s.weapon ? ++designNum : 2,
		position: { x: s.at[0], y: s.at[1] },
		armor: 200,
		startingQuantity: 2
	}));

	const actions: BattleRecordTokenActionJson[] = shooters.flatMap((s, i) =>
		s.targets.map((targetNum) => {
			const torpedo = s.hits !== undefined;
			const unknownShooter = !s.weapon;
			return {
				type: torpedo
					? 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE'
					: 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE',
				tokenNum: shooterTokens[i].num,
				round: 1,
				slot: 2,
				from: shooterTokens[i].position,
				to: unknownShooter ? { x: 1, y: 2 } : targets[targetNum - 1].position,
				targetNum: unknownShooter ? 3 + targets.length : targetNum,
				target: { designNum: 1, quantity: 4 - (s.destroyed ?? 0) },
				damageDoneShields: s.shields ?? 0,
				damageDoneArmor: s.armor ?? 0,
				tokensDestroyed: s.destroyed ?? 0,
				torpedoHits: s.hits ?? 0,
				torpedoMisses: s.misses ?? 0
			} satisfies BattleRecordTokenActionJson;
		})
	);

	const battle = fromJson(BattleRecordSchema, {
		num: 1,
		tokens: [...targets, ...shooterTokens],
		actionsPerRound: [{ actions: [] }, { actions }]
	});

	const universe = new Universe();
	universe.setData(1, create(PlayerUniverseSchema, { designs }));
</script>

<h1 class="text-xl">Battle</h1>
<BattleView battleRecord={battle} playerFinder={universe} designFinder={universe} />
<Tooltip />
<Popup />
