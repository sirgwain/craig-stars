<script lang="ts">
	import { goto } from '$app/navigation';

	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { Service } from '$lib/services/Service';
	import {
		AIDifficultyNone,
		AIDifficultyNormal,
		DensityNormal,
		GameStartModeNormal,
		NewGamePlayerTypeAI,
		NewGamePlayerTypeHost,
		PlayerPositionsModerate,
		SizeSmall,
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionOwnPlanets,
		type Game,
		type GameSettings,
		type NewGamePlayer as Player
	} from '$lib/types/cs';
	import { PlusCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GameSettingsEditor from './GameSettingsEditor.svelte';
	import { getColor, getFirstAvailableColor } from './playerColors';
	import VictoryConditions from './VictoryConditions.svelte';
	import NewGamePlayer from './NewGamePlayer.svelte';

	type Props = {
		players?: Player[];
		name?: string;
	};

	let {
		players = [
			{
				type: NewGamePlayerTypeHost,
				color: getColor(0),
				aiDifficulty: AIDifficultyNone,
				hullSetNum: 0
			},
			{
				type: NewGamePlayerTypeAI,
				color: getColor(1),
				aiDifficulty: AIDifficultyNormal,
				hullSetNum: 0
			},
			{
				type: NewGamePlayerTypeAI,
				color: getColor(2),
				aiDifficulty: AIDifficultyNormal,
				hullSetNum: 0
			}
		],
		name = 'A Barefoot Jaywalk'
	}: Props = $props();

	let settings: GameSettings = $state({
		name,
		public: false,
		size: SizeSmall,
		density: DensityNormal,
		playerPositions: PlayerPositionsModerate,
		randomEvents: true,
		computerPlayersFormAlliances: false,
		publicPlayerScores: false,
		maxMinerals: false,
		acceleratedPlay: false,
		quickStartTurns: 0,
		startMode: GameStartModeNormal,
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
	});

	const onSubmit = async () => {
		const data = JSON.stringify(settings);

		const response = await fetch(`/api/games`, {
			method: 'post',
			headers: {
				accept: 'application/json'
			},
			body: data
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		const game = (await response.json()) as Game;
		goto(`/games/${game.id}`);
	};

	const addPlayer = () => {
		const usedColors = new Set<string>(settings.players.map<string>((p) => p.color ?? ''));

		settings.players = [
			...settings.players,
			{
				type: NewGamePlayerTypeAI,
				color: getFirstAvailableColor(usedColors),
				aiDifficulty: AIDifficultyNormal,
				hullSetNum: 0
			}
		];
	};

	const removePlayer = (player: Player) => {
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
		<button class="btn btn-success" type="submit">Create Game</button>
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
