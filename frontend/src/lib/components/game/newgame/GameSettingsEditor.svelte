<script lang="ts">
	import CheckboxInput from '$lib/components/CheckboxInput.svelte';
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { Density, PlayerPositions, Size, type GameSettings } from '$lib/types/Game';
	import PrivateGameLink from '../../../../routes/(user)/games/(game)/[id]/(main)/PrivateGameLink.svelte';

	export let settings: GameSettings;
	export let showInviteLink = false;
</script>

<div class="flex flex-row flex-wrap">
	{#if showInviteLink && !settings.public}
		<PrivateGameLink />
	{/if}
	<TextInput name="name" bind:value={settings.name} />
	<EnumSelect name="size" enumType={Size} bind:value={settings.size} />
	<EnumSelect name="density" enumType={Density} bind:value={settings.density} />
	<EnumSelect
		name="playerPositions"
		enumType={PlayerPositions}
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
</div>
