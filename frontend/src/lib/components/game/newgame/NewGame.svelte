<script lang="ts">
	import { untrack } from 'svelte';
	import { goto } from '$app/navigation';
	import ItemTitle from '#lib/components/ItemTitle.svelte';
	import SectionHeader from '#lib/components/SectionHeader.svelte';
	import { gameClient } from '#lib/services/connect.js';
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionOwnPlanets
	} from '#lib/types/Consts.js';
	import {
		AiDifficulty,
		Density,
		type NewGamePlayer as GamePlayer,
		type GameSettings,
		GameSettingsSchema,
		GameStartMode,
		NewGamePlayerSchema,
		NewGamePlayerType,
		PlayerPositions,
		Size
	} from '#lib/types/cs-proto.js';
	import { create } from '@bufbuild/protobuf';
	import { PlusCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GameSettingsEditor from './GameSettingsEditor.svelte';
	import NewGamePlayer from './NewGamePlayer.svelte';
	import { enumToString } from '#lib/types/Enums.js';
	import { getColor, getFirstAvailableColor } from './playerColors';
	import VictoryConditions from './VictoryConditions.svelte';

	type Props = {
		players?: GamePlayer[];
		name?: string;
	};

	let {
		players = [
			create(NewGamePlayerSchema, {
				type: NewGamePlayerType.HOST,
				color: getColor(0),
				aiDifficulty: AiDifficulty.UNSPECIFIED,
				defaultHullSet: 0
			}),
			create(NewGamePlayerSchema, {
				type: NewGamePlayerType.AI,
				color: getColor(1),
				aiDifficulty: AiDifficulty.NORMAL,
				defaultHullSet: 0
			}),
			create(NewGamePlayerSchema, {
				type: NewGamePlayerType.AI,
				color: getColor(2),
				aiDifficulty: AiDifficulty.NORMAL,
				defaultHullSet: 0
			})
		],
		name = 'A Barefoot Jaywalk'
	}: Props = $props();

	let settings: GameSettings = $state(
		create(GameSettingsSchema, {
			name: untrack(() => name),
			public: false,
			size: Size.SMALL,
			density: Density.NORMAL,
			playerPositions: PlayerPositions.MODERATE,
			randomEvents: true,
			computerPlayersFormAlliances: false,
			publicPlayerScores: false,
			maxMinerals: false,
			techCostMultiplier: 1,
			startMode: GameStartMode.UNSPECIFIED, // normal
			quickStartTurns: 0,
			players: untrack(() => players),
			victoryConditions: {
				conditions:
					VictoryConditionOwnPlanets |
					VictoryConditionAttainTechLevels |
					VictoryConditionExceedsSecondPlaceScore,
				numCriteriaRequired: 1,
				yearsPassed: 50,
				ownPlanets: 60,
				attainTechLevel: 22,
				attainTechLevelNumFields: 4,
				exceedsScore: 11000,
				exceedsSecondPlaceScore: 100,
				productionCapacity: 100,
				ownCapitalShips: 100,
				highestScoreAfterYears: 100
			}
		})
	);

	let submitting = $state(false);

	const onSubmit = async () => {
		submitting = true;
		try {
			const resp = await gameClient.createGame({ settings: settings });
			if (resp.game?.game) {
				goto(`/games/${resp.game.game.id}`);
			}
		} finally {
			submitting = false;
		}
	};

	const addPlayer = () => {
		const usedColors = new Set<string>(settings.players.map<string>((p) => p.color));

		settings.players = [
			...settings.players,
			create(NewGamePlayerSchema, {
				type: NewGamePlayerType.AI,
				color: getFirstAvailableColor(usedColors),
				aiDifficulty: AiDifficulty.NORMAL,
				defaultHullSet: 0
			})
		];
	};

	const removePlayer = (player: GamePlayer) => {
		settings.players = settings.players.filter((p) => p !== player);
	};
</script>

<form
	class="mx-auto w-full max-w-6xl pb-4"
	onsubmit={(e) => {
		e.preventDefault();
		onSubmit();
	}}
>
	<ItemTitle>New Game</ItemTitle>
	<div class="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_18rem]">
		<div class="min-w-0">
			<GameSettingsEditor bind:settings />
			<SectionHeader
				><span class="flex items-center justify-between"
					>Players <span class="text-sm font-normal">{settings.players.length} slots</span></span
				></SectionHeader
			>
			<div class="grid gap-2">
				{#each settings.players as player, i (player)}
					<NewGamePlayer
						bind:player={settings.players[i]}
						index={i + 1}
						onRemove={() => removePlayer(player)}
					/>
				{/each}
			</div>
			<button type="button" class="btn btn-outline btn-sm mt-3" onclick={addPlayer}
				><Icon src={PlusCircle} size="20" />Add player</button
			>
			<details class="mt-4 border border-secondary p-3">
				<summary class="cursor-pointer font-semibold"
					>Victory Conditions <span class="text-sm font-normal text-base-content/70"
						>· {settings.victoryConditions?.numCriteriaRequired} criteria after {settings
							.victoryConditions?.yearsPassed} years</span
					></summary
				>
				<VictoryConditions bind:settings />
			</details>
		</div>
		<aside
			class="sticky top-4 hidden border border-base-300 bg-base-200 p-4 lg:block"
			aria-label="Game summary"
		>
			<h3 class="mb-3 text-lg font-semibold">Game summary</h3>
			<p class="mb-3 wrap-break-word font-semibold">{settings.name}</p>
			<dl class="grid grid-cols-2 gap-2 text-sm">
				<dt>Galaxy</dt>
				<dd>{enumToString(Size, settings.size)}</dd>
				<dt>Density</dt>
				<dd>{enumToString(Density, settings.density)}</dd>
				<dt>Players</dt>
				<dd>{settings.players.length}</dd>
				<dt>Open slots</dt>
				<dd>{settings.players.filter((p) => p.type === NewGamePlayerType.OPEN).length}</dd>
				<dt>Visibility</dt>
				<dd>{settings.public ? 'Public' : 'Private'}</dd>
				<dt>Technology cost</dt>
				<dd>{settings.techCostMultiplier}×</dd>
			</dl>
		</aside>
	</div>
	<div class="sticky bottom-0 z-10 mt-4 flex justify-end border-t border-base-300 bg-base-100 py-3">
		<button type="submit" disabled={submitting} class="btn btn-success w-full sm:w-auto"
			>{submitting ? 'Creating...' : 'Create Game'}</button
		>
	</div>
</form>
