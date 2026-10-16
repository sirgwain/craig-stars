<script lang="ts">
	import { clone } from '@bufbuild/protobuf';
	import GameCreationForm from '#lib/components/game/newgame/GameCreationForm.svelte';
	import {
		hostedGameSettings,
		singlePlayerSettings
	} from '#lib/components/game/newgame/defaults.js';
	import { GameSettingsSchema, NewGamePlayerType } from '#lib/types/cs-proto.js';
	import { humanoid } from '#lib/types/Race.js';
	import { enumToString } from '#lib/types/Enums.js';
	import { Size, Density } from '#lib/types/cs-proto.js';
	import DocExample from '../DocExample.svelte';

	// Example choices only; these are never fetched from or saved to an account.
	const races = [
		Object.assign(humanoid(), { id: 1n }),
		Object.assign(humanoid(), { id: 2n, name: 'Example Settler', pluralName: 'Example Settlers' })
	];
	let mode = $state('single');
	let settings = $state(singlePlayerSettings());
	let preview = $state<ReturnType<typeof singlePlayerSettings>>();

	function reset() {
		settings = mode === 'host' ? hostedGameSettings() : singlePlayerSettings();
		preview = undefined;
	}

	function showPreview() {
		preview = clone(GameSettingsSchema, settings);
	}
</script>

<DocExample title="Try game creation" onReset={reset}>
	<label class="mb-4 flex flex-wrap items-center gap-3">
		<span class="font-semibold">Example setup</span>
		<select class="select select-sm" bind:value={mode} onchange={reset}>
			<option value="single">Single player</option>
			<option value="host">Host a game</option>
		</select>
	</label>
	<GameCreationForm bind:settings {races} submitLabel="Preview setup" onSubmit={showPreview} />
	{#if preview}
		<div
			role="status"
			class="mt-5 rounded-box border border-primary bg-base-200 p-4"
			data-type="setup-preview"
		>
			<h4 class="text-xl font-semibold">{preview.name}</h4>
			<p class="mt-2">
				{enumToString(Size, preview.size)} galaxy · {enumToString(Density, preview.density)} density ·
				{preview.players.length} players
			</p>
			<ul class="mt-2 space-y-1">
				{#each preview.players as player, index (index)}
					<li>
						Player {index + 1}: {enumToString(
							NewGamePlayerType,
							player.type
						)}{#if player.type === NewGamePlayerType.HOST}
							— {player.race?.pluralName}{/if}
					</li>
				{/each}
			</ul>
			<p class="mt-3">
				Example only. No game was created. Change the controls and preview again to update this
				summary.
			</p>
		</div>
	{/if}
</DocExample>
