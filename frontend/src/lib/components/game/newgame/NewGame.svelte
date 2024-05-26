<script lang="ts">
	import { goto } from '$app/navigation';

	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { Service } from '$lib/services/Service';
	import {
		Density,
		GameStartMode,
		NewGamePlayerType,
		PlayerPositions,
		Size,
		VictoryCondition,
		type Game,
		type GameSettings,
		type NewGamePlayers,
		type NewGamePlayer as Player
	} from '$lib/types/Game';
	import { PlusCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GameSettingsEditor from './GameSettingsEditor.svelte';
	import NewGamePlayer from './NewGamePlayer.svelte';
	import { getColor, getFirstAvailableColor } from './playerColors';
	import VictoryConditions from './VictoryConditions.svelte';

	export let players = [
		{ type: NewGamePlayerType.Host, color: getColor(0) },
		{ type: NewGamePlayerType.AI, color: getColor(1) },
		{ type: NewGamePlayerType.AI, color: getColor(2) }
	];

	export let name = 'A Barefoot Jaywalk';

	let settings: GameSettings & NewGamePlayers = {
		name,
		public: false,
		size: Size.Small,
		density: Density.Normal,
		playerPositions: PlayerPositions.Moderate,
		randomEvents: true,
		computerPlayersFormAlliances: false,
		publicPlayerScores: false,
		maxMinerals: false,
		acceleratedPlay: false,
		startMode: GameStartMode.Normal,
		players,
		victoryConditions: {
			conditions:
				VictoryCondition.OwnPlanets |
				VictoryCondition.AttainTechLevels |
				VictoryCondition.ExceedsSecondPlaceScore,
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
	};

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
			{ type: NewGamePlayerType.AI, color: getFirstAvailableColor(usedColors) }
		];
	};

	const removePlayer = (player: Player) => {
		settings.players = settings.players.filter((p) => p !== player);
	};

	let error = '';
</script>

<form on:submit|preventDefault={onSubmit}>
	<div class="w-full flex justify-end gap-2">
		<button class="btn btn-success" type="submit">Create Game</button>
	</div>

	<ItemTitle>New Game</ItemTitle>

	<GameSettingsEditor bind:settings />

	<SectionHeader>
		<button class="btn-ghost w-full flex flex-row" on:click|preventDefault={addPlayer}>
			Players
			<div class="ml-auto">
				<Icon src={PlusCircle} size="24" class="hover:stroke-accent" />
			</div>
		</button></SectionHeader
	>

	{#each settings.players as player, i}
		<NewGamePlayer bind:player index={i + 1} on:remove={() => removePlayer(player)} />
	{/each}

	<SectionHeader>Victory Conditions</SectionHeader>
	<VictoryConditions bind:settings />
</form>
