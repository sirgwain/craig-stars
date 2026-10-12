<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import Note from './Note.svelte';

	// Any part of the game can open this by calling editNote or editLabel on the game context.
	const { editingNote, cancelEditNote } = getGameContext();
</script>

<div class="modal" class:modal-open={!!$editingNote} data-type="note-dialog">
	<div class="modal-box max-w-full max-h-max size-full md:max-w-lg md:h-auto">
		{#if $editingNote}
			{#key $editingNote}
				<Note note={$editingNote} onClose={cancelEditNote} />
			{/key}
		{/if}
	</div>
</div>
