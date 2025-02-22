<script lang="ts">
	import CheckboxInput from '$lib/components/CheckboxInput.svelte';
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { Densities, GameStartModes, PlayerPositionses, Sizes } from '$lib/types/Game';
	import { GameStartModeNormal, type GameSettings } from '$lib/types/cs';
	import { startCase } from 'lodash-es';
	import PrivateGameLink from './PrivateGameLink.svelte';

	type Props = {
		settings: GameSettings;
		showInviteLink?: boolean;
	};

	let { settings = $bindable(), showInviteLink = false }: Props = $props();
</script>

<div class="flex flex-row flex-wrap">
	{#if showInviteLink && !settings.public}
		<PrivateGameLink />
	{/if}
	<TextInput name="name" bind:value={settings.name} />
	<EnumSelect name="size" options={Sizes} bind:value={settings.size} />
	<EnumSelect name="density" options={Densities} bind:value={settings.density} />
	<EnumSelect
		name="playerPositions"
		options={PlayerPositionses}
		bind:value={settings.playerPositions}
	/>
	<CheckboxInput name="public" bind:checked={settings.public} />
	<CheckboxInput name="randomEvents" bind:checked={settings.randomEvents} />
	<CheckboxInput name="publicPlayerScores" bind:checked={settings.publicPlayerScores} />
	<CheckboxInput
		title="Beginner: Max Minerals"
		name="maxMinerals"
		bind:checked={settings.maxMinerals}
	/>
	<CheckboxInput
		name="computerPlayersFormAlliances"
		bind:checked={settings.computerPlayersFormAlliances}
	/>
	<EnumSelect
		name="startMode"
		options={GameStartModes}
		bind:value={settings.startMode}
		typeTitle={(value) => (!value || value === GameStartModeNormal ? 'Normal' : startCase(value))}
		showEmpty={true}
		tooltip={`Setting mode to Max will create a game with all tech levels, max minerals, etc`}
	/>
</div>
