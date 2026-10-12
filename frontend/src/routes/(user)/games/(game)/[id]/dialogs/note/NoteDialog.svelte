<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import Note from './Note.svelte';

	// Any part of the game can open this by calling editNote or editLabel on the game context.
	const { editingNote, cancelEditNote } = getGameContext();
</script>

<!-- On a phone the dialog sits at the top and stays short, so all of it is above the keyboard -->
<div
	class="modal modal-top md:modal-middle"
	class:modal-open={!!$editingNote}
	data-type="note-dialog"
>
	<div
		class="modal-box w-full max-w-full max-h-dvh pt-[max(1rem,env(safe-area-inset-top))] md:max-w-lg md:pt-6"
	>
		{#if $editingNote}
			{#key $editingNote}
				<Note note={$editingNote} onClose={cancelEditNote} />
			{/key}
		{/if}
	</div>
</div>
