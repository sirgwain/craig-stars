<script lang="ts">
	import ItemTitle from '#lib/components/ItemTitle.svelte';
	import SectionHeader from '#lib/components/SectionHeader.svelte';
	import {
		NewGamePlayerType,
		type GameSettings,
		type Race,
		type NewGamePlayer as GamePlayer
	} from '#lib/types/cs-proto.js';
	import { PlusCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import GameSettingsEditor from './GameSettingsEditor.svelte';
	import NewGamePlayer from './NewGamePlayer.svelte';
	import { getFirstAvailableColor } from './playerColors';
	import VictoryConditions from './VictoryConditions.svelte';
	import { allPlayerTypes, newGamePlayer } from './defaults';

	type Props = {
		settings: GameSettings;
		races: Race[];
		allowedPlayerTypes?: NewGamePlayerType[];
		submitting?: boolean;
		submitLabel?: string;
		onSubmit: () => void | Promise<void>;
	};
	let {
		settings = $bindable(),
		races,
		allowedPlayerTypes = allPlayerTypes,
		submitting = false,
		submitLabel = 'Create Game',
		onSubmit
	}: Props = $props();

	function addPlayer() {
		const player = newGamePlayer(NewGamePlayerType.AI, settings.players.length);
		player.color = getFirstAvailableColor(new Set(settings.players.map((p) => p.color)));
		settings.players = [...settings.players, player];
	}

	function removePlayer(player: GamePlayer) {
		settings.players = settings.players.filter((p) => p !== player);
	}
</script>

<form
	onsubmit={(e) => {
		e.preventDefault();
		onSubmit();
	}}
>
	<div class="w-full flex justify-end gap-2">
		<button type="submit" disabled={submitting} class="btn btn-success"
			>{submitting ? 'Creating...' : submitLabel}</button
		>
	</div>

	<div class="flex items-center justify-between gap-2">
		<ItemTitle>New Game</ItemTitle>
		<a class="link text-sm" href="/docs/games/create#try-it">Game creation help</a>
	</div>

	{#key settings}
		<GameSettingsEditor bind:settings />

		<SectionHeader>
			<button
				type="button"
				aria-label="Add player"
				class="btn-ghost w-full flex flex-row"
				onclick={addPlayer}
			>
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
				{races}
				{allowedPlayerTypes}
				onRemove={() => removePlayer(player)}
			/>
		{/each}

		<SectionHeader>Victory Conditions</SectionHeader>
		<VictoryConditions bind:settings />
	{/key}
</form>
