<script lang="ts">
	import { MaxGameNameLength } from '#lib/types/Consts.js';
	import CheckboxInput from '#lib/components/CheckboxInput.svelte';
	import EnumSelect from '#lib/components/EnumSelect.svelte';
	import TextInput from '#lib/components/TextInput.svelte';
	import {
		Density,
		GameStartMode,
		PlayerPositions,
		Size,
		type GameSettings
	} from '#lib/types/cs-proto.js';
	import { GameStartModeFullNames } from '#lib/types/Enums.js';
	import SectionHeader from '#lib/components/SectionHeader.svelte';
	import PrivateGameLink from './PrivateGameLink.svelte';

	type Props = {
		settings: GameSettings;
		showInviteLink?: boolean;
	};

	let { settings = $bindable(), showInviteLink = false }: Props = $props();
</script>

<SectionHeader>Galaxy</SectionHeader>
<div class="grid grid-cols-2 gap-3">
	<div class="col-span-2">
		<TextInput
			name="name"
			stacked
			required
			bind:value={settings.name}
			maxlength={MaxGameNameLength}
		/>
	</div>
	<EnumSelect name="size" stacked enumType={Size} bind:value={settings.size} />
	<EnumSelect name="density" stacked enumType={Density} bind:value={settings.density} />
	<div class="col-span-2 sm:col-span-1">
		<EnumSelect
			name="playerPositions"
			stacked
			enumType={PlayerPositions}
			bind:value={settings.playerPositions}
		/>
	</div>
	<div class="col-span-2 sm:col-span-1">
		<EnumSelect
			name="startMode"
			stacked
			enumType={GameStartMode}
			bind:value={settings.startMode}
			typeTitle={(t) => GameStartModeFullNames[t]}
			showEmpty={true}
		/>
	</div>
</div>
<p class="mt-2 text-sm text-base-content/70">
	Accelerated BBS Play boosts starting population and minerals. Max Start begins with maximum
	technology and resources.
</p>
<SectionHeader>Rules</SectionHeader>
<div class="grid gap-1 sm:grid-cols-2">
	<CheckboxInput
		name="public"
		title="Public game"
		titleClass="cs-label-text grow text-left"
		bind:checked={settings.public}
	/>
	<CheckboxInput
		name="randomEvents"
		titleClass="cs-label-text grow text-left"
		bind:checked={settings.randomEvents}
	/>
	<CheckboxInput
		name="publicPlayerScores"
		titleClass="cs-label-text grow text-left"
		bind:checked={settings.publicPlayerScores}
	/>
	<CheckboxInput
		name="maxMinerals"
		title="Beginner: Max Minerals"
		titleClass="cs-label-text grow text-left"
		bind:checked={settings.maxMinerals}
	/>
	<CheckboxInput
		name="computerPlayersFormAlliances"
		title="AI players form alliances"
		titleClass="cs-label-text grow text-left"
		bind:checked={settings.computerPlayersFormAlliances}
	/>
	<CheckboxInput
		name="galaxyClumping"
		title="Galaxy Clumping"
		titleClass="cs-label-text grow text-left"
		bind:checked={settings.galaxyClumping}
	/>
</div>
<p class="mt-2 text-sm text-base-content/70">
	{settings.public
		? 'Public games appear in the open games list.'
		: 'Private games can be joined with an invite link.'}
</p>
{#if showInviteLink && !settings.public}<PrivateGameLink />{/if}
<div class="mt-4 border border-base-300 bg-base-200 p-3">
	<label for="tech-cost" class="flex items-center justify-between gap-3 font-semibold">
		Technology cost <span class="text-info">{settings.techCostMultiplier}×</span>
	</label>
	<input
		id="tech-cost"
		name="techCostMultiplier"
		aria-describedby="tech-cost-help"
		type="range"
		class="range range-info range-sm mt-3 w-full"
		min="1"
		max="10"
		step="0.25"
		bind:value={settings.techCostMultiplier}
	/>
	<div class="mt-1 flex justify-between text-xs text-base-content/70">
		<span>1× Normal</span><span>10×</span>
	</div>
	<p id="tech-cost-help" class="mt-2 text-sm text-base-content/70">
		Scales research costs in every field, after race modifiers. 2× matches Slower Tech Advances.
	</p>
	<div class="mt-2 flex gap-2">
		<button type="button" class="btn btn-sm" onclick={() => (settings.techCostMultiplier = 1)}
			>Normal (1×)</button
		><button type="button" class="btn btn-sm" onclick={() => (settings.techCostMultiplier = 2)}
			>Slower tech (2×)</button
		>
	</div>
</div>
