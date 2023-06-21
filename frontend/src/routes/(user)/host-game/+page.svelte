<script lang="ts">
	import { goto } from '$app/navigation';
	import { PlusCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { $enum as eu } from 'ts-enum-util';

	import {
		Density,
		GameStartMode,
		NewGamePlayerType,
		PlayerPositions,
		Size,
		type GameSettings,
		VictoryCondition
	} from '$lib/types/Game';
	import NewGamePlayer from './NewGamePlayer.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import CheckboxInput from '$lib/components/CheckboxInput.svelte';

	let colors = ['#0000FF', '#C33232', '#1F8BA7', '#43A43E', '#8D29CB', '#B88628'];
	const getColor = (index: number) =>
		index < colors.length ? colors[index] : '#' + Math.floor(Math.random() * 16777215).toString(16);

	let settings: GameSettings = {
		name: 'A Barefoot Jaywalk',
		size: Size.Small,
		density: Density.Normal,
		playerPositions: PlayerPositions.Moderate,
		randomEvents: true,
		computerPlayersFormAlliances: false,
		publicPlayerScores: false,
		startMode: GameStartMode.Normal,
		players: [
			{ type: NewGamePlayerType.Host, color: getColor(0) },
			{ type: NewGamePlayerType.AI, color: getColor(1) },
			{ type: NewGamePlayerType.AI, color: getColor(2) }
		],
		victoryConditions: {
			conditions:
				VictoryCondition.OwnPlanets |
				VictoryCondition.AttainTechLevels |
				VictoryCondition.ExceedsScore,
			numCriteriaRequired: 1,
			yearsPassed: 50,
			ownPlanets: 60,
			attainTechLevel: 22,
			attainTechLevelNumFields: 4,
			exceedsScore: 11000,
			exceedsSecondPlaceScore: 100,
			productionCapacity: 100000,
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

		if (response.ok) {
			goto('/');
		} else {
			const resolvedResponse = await response?.json();
			error = resolvedResponse.error;
			console.error(error);
		}
	};

	const addPlayer = () => {
		settings.players = [
			...settings.players,
			{ type: NewGamePlayerType.AI, color: getColor(settings.players.length + 1) }
		];
	};

	let error = '';
</script>

<form on:submit|preventDefault={onSubmit}>
	<div class="w-full flex justify-end gap-2">
		<button class="btn btn-success" type="submit">Create Game</button>
	</div>

	<ItemTitle>Host New Game</ItemTitle>

	<div class="flex flex-row flex-wrap">
		<TextInput name="name" bind:value={settings.name} />
		<EnumSelect name="size" enumType={Size} bind:value={settings.size} />
		<EnumSelect name="density" enumType={Density} bind:value={settings.density} />
		<EnumSelect
			name="playerPositions"
			enumType={PlayerPositions}
			bind:value={settings.playerPositions}
		/>
		<CheckboxInput name="randomEvents" bind:checked={settings.randomEvents} />
		<CheckboxInput name="publicPlayerScores" bind:checked={settings.publicPlayerScores} />
		<CheckboxInput
			name="computerPlayersFormAlliances"
			bind:checked={settings.computerPlayersFormAlliances}
		/>
	</div>

	<SectionHeader>
		<button class="btn-ghost w-full flex flex-row" on:click|preventDefault={addPlayer}>
			Players
			<div class="ml-auto">
				<Icon src={PlusCircle} size="24" class="hover:stroke-accent" />
			</div>
		</button></SectionHeader
	>

	{#each settings.players as player, i}
		<NewGamePlayer bind:player index={i + 1} />
	{/each}
</form>
