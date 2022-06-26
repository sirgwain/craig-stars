<script lang="ts">
	import { goto } from '$app/navigation';
	import { RaceService } from '$lib/services/RaceService';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { PlusCircle } from '@steeze-ui/heroicons';

	import {
		Density,
		GameStartMode,
		NewGamePlayerType,
		PlayerPositions,
		Size,
		type Game,
		type GameSettings
	} from '$lib/types/Game';
	import type { Race } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import NewGamePlayer from './_NewGamePlayer.svelte';

	let settings: GameSettings = {
		name: 'A Barefoot Jaywalk',
		size: Size.Tiny,
		density: Density.Normal,
		playerPositions: PlayerPositions.Moderate,
		randomEvents: true,
		computerPlayersFormAlliances: false,
		publicPlayerScores: false,
		startMode: GameStartMode.Normal,
		players: [{ type: NewGamePlayerType.Host }]
	};

	let sizes: string[] = Object.keys(Size).filter((v) => isNaN(Number(v)));

	const onSubmit = async () => {
		const data = JSON.stringify({ settings });

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
		settings.players = [...settings.players, { type: NewGamePlayerType.AI }];
	};

	const stringToSize = (size: string): Size => {
		return Size[size as keyof typeof Size];
	};

	let error = '';
</script>

<div class="prose"><h2>Host New Game</h2></div>
<form on:submit|preventDefault={onSubmit}>
	<fieldset name="settings" class="form-control">
		<label class="label" for="name">Name</label>
		<input name="name" class="input input-bordered" bind:value={settings.name} />

		<label class="label" for="size">Universe Size</label>
		<select name="size" class="select select-bordered" bind:value={settings.size}>
			{#each sizes as size}
				<option value={stringToSize(size)}>{size}</option>
			{/each}
		</select>
	</fieldset>

	<fieldset name="players" class="form-control mt-3">
		<h3 class="text-2xl items-start">
			Players <button on:click|preventDefault={addPlayer}
				><Icon src={PlusCircle} size="24" class="hover:stroke-accent" />
			</button>
		</h3>

		{#each settings.players as player, i}
			<NewGamePlayer bind:player index={i + 1} />
		{/each}
	</fieldset>
	<button class="btn btn-primary">Host</button>
</form>
