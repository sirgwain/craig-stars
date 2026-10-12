<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { MapObjectType, type PlayerNote } from '#lib/types/cs-proto.js';
	import { ConnectError } from '@connectrpc/connect';
	import hotkeys from 'hotkeys-js';
	import { onMount, untrack } from 'svelte';
	import TagPicker from './TagPicker.svelte';

	const { universe, readOnly, savePlayerNote, deletePlayerNote } = getGameContext();

	// the server rejects longer notes
	const maxLength = 2000;

	type Props = {
		note: PlayerNote;
		onClose: () => void;
	};

	let { note, onClose }: Props = $props();

	// the dialog edits a copy, the note itself only changes when it's saved
	let text = $state(untrack(() => note.text));
	let showOnMap = $state(untrack(() => note.showOnMap));
	let tags = $state(untrack(() => [...note.tags]));
	let saving = $state(false);
	let error = $state('');
	let textarea: HTMLTextAreaElement | undefined = $state();

	// labels aren't attached to anything, they only exist on the map
	let isLabel = $derived(note.targetType === MapObjectType.UNSPECIFIED);
	// fleets move, so only notes on planets can be drawn on the map
	let isPlanet = $derived(note.targetType === MapObjectType.PLANET);
	let exists = $derived(!!note.gameDbObject?.id);
	let targetName = $derived(
		note.targetType === MapObjectType.FLEET
			? ($universe.getMyFleet(note.targetNum)?.mapObject?.name ?? 'Fleet')
			: ($universe.getPlanet(note.targetNum)?.mapObject?.name ?? 'Planet')
	);
	let title = $derived(isLabel ? (exists ? 'Map Label' : 'New Map Label') : `Notes: ${targetName}`);

	async function run(action: () => Promise<void>) {
		saving = true;
		error = '';
		try {
			await action();
			onClose();
		} catch (e) {
			error = e instanceof ConnectError ? e.rawMessage : 'Failed to save note';
		} finally {
			saving = false;
		}
	}

	function save() {
		if (saving || $readOnly) {
			return;
		}
		if (!text.trim() && (isLabel || !tags.length)) {
			// clearing a note removes it
			if (exists) {
				remove();
			} else {
				onClose();
			}
			return;
		}
		run(() =>
			savePlayerNote({
				...note,
				text,
				tags: isLabel ? [] : tags,
				// only the text is shown on the map
				showOnMap: isLabel || (isPlanet && showOnMap && !!text.trim())
			})
		);
	}

	function remove() {
		if (saving || $readOnly) {
			return;
		}
		run(() => deletePlayerNote(note));
	}

	function onKeyDown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			e.preventDefault();
			onClose();
		} else if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
			e.preventDefault();
			save();
		}
	}

	onMount(() => {
		// keep the game's hotkeys from firing while this dialog is open
		const originalScope = hotkeys.getScope();
		const scope = 'note';
		hotkeys('Esc', scope, onClose);
		hotkeys.setScope(scope);
		textarea?.focus();

		return () => {
			hotkeys.unbind('Esc', scope, onClose);
			hotkeys.deleteScope(scope);
			hotkeys.setScope(originalScope);
		};
	});
</script>

<div class="flex flex-col gap-2" data-type="note" role="presentation" onkeydown={onKeyDown}>
	<h2 class="text-lg font-semibold text-secondary truncate">{title}</h2>

	{#if error}
		<div class="alert alert-error py-2" data-type="note-error">{error}</div>
	{/if}

	<textarea
		bind:this={textarea}
		bind:value={text}
		maxlength={maxLength}
		readonly={$readOnly}
		rows={isLabel ? 2 : 3}
		aria-label={title}
		placeholder={isLabel ? 'A label only you can see, i.e. "the front"' : 'Notes only you can see'}
		class="textarea textarea-bordered w-full text-base"></textarea>

	{#if !isLabel}
		<div class="text-sm font-semibold">Tags</div>
		<TagPicker bind:tags known={$universe.tags} disabled={$readOnly} />
	{/if}
	{#if isPlanet}
		<label class="label cursor-pointer justify-start gap-2">
			<input
				type="checkbox"
				class="checkbox checkbox-sm"
				bind:checked={showOnMap}
				disabled={$readOnly}
			/>
			<span>Show note on map</span>
		</label>
	{/if}

	<div class="flex flex-row gap-2">
		{#if exists && !$readOnly}
			<button
				type="button"
				class="btn btn-sm btn-outline btn-error"
				disabled={saving}
				onclick={remove}
				data-type="note-delete">Delete</button
			>
		{/if}
		<div class="grow"></div>
		<button type="button" class="btn btn-sm btn-outline" onclick={onClose} data-type="note-cancel"
			>{$readOnly ? 'Close' : 'Cancel'}</button
		>
		{#if !$readOnly}
			<button
				type="button"
				class="btn btn-sm btn-primary"
				disabled={saving}
				onclick={save}
				data-type="note-save">Save</button
			>
		{/if}
	</div>
</div>
