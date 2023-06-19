<script lang="ts">
	import { NewGamePlayerType, type NewGamePlayer } from '$lib/types/Game';
	import Host from './_Host.svelte';

	export let player: NewGamePlayer;
	export let index: number;

	let types: string[] = Object.keys(NewGamePlayerType).filter((v) => isNaN(Number(v)));

	const stringToNewGamePlayerType = (value: string): NewGamePlayerType => {
		return NewGamePlayerType[value as keyof typeof NewGamePlayerType];
	};
</script>

{#if player}
	<div class="block">
		Player {index}

		<select name="type" class="select select-bordered" bind:value={player.type}>
			{#each types as type}
				<option value={stringToNewGamePlayerType(type)}>{type}</option>
			{/each}
		</select>

		{#if player.type === NewGamePlayerType.Host}
			<Host {player} />
		{/if}
	</div>
{/if}
