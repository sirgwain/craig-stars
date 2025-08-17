<script lang="ts">
	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { gameClient } from '$lib/services/connect';
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionOwnPlanets
	} from '$lib/types/Consts';
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
	} from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import { PlusCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GameSettingsEditor from './GameSettingsEditor.svelte';
	import NewGamePlayer from './NewGamePlayer.svelte';
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
			name,
			public: false,
			size: Size.SMALL,
			density: Density.NORMAL,
			playerPositions: PlayerPositions.MODERATE,
			randomEvents: true,
			computerPlayersFormAlliances: false,
			publicPlayerScores: false,
			maxMinerals: false,
			startMode: GameStartMode.UNSPECIFIED, // normal
			quickStartTurns: 0,
			players,
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
		const usedColors = new Set<string>(settings.players.map<string>((p) => p.color ?? ''));

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
	onsubmit={(e) => {
		e.preventDefault();
		onSubmit();
	}}
>
	<div class="w-full flex justify-end gap-2">
		<button type="submit" disabled={submitting} class="btn btn-success"
			>{submitting ? 'Creating...' : 'Create Game'}</button
		>
	</div>

	<ItemTitle>New Game</ItemTitle>

	<GameSettingsEditor bind:settings />

	<SectionHeader>
		<button type="button" class="btn-ghost w-full flex flex-row" onclick={addPlayer}>
			Players
			<div class="ml-auto">
				<Icon src={PlusCircle} size="24" class="hover:stroke-accent" />
			</div>
		</button></SectionHeader
	>

	{#each settings.players as player, i (i)}
		<NewGamePlayer
			bind:player={settings.players[i]}
			index={i + 1}
			onRemove={() => removePlayer(player)}
		/>
	{/each}

	<SectionHeader>Victory Conditions</SectionHeader>
	<VictoryConditions bind:settings />
</form>
