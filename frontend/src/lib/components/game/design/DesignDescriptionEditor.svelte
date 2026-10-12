<script lang="ts">
	import { addError } from '#lib/services/Errors.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { MaxDescriptionLength } from '#lib/types/Consts.js';
	import type { ShipDesign } from '#lib/types/cs-proto.js';
	import type { ConnectError } from '@connectrpc/connect';
	import { untrack } from 'svelte';

	const { updateDesignDescription, readOnly } = getGameContext();

	type Props = {
		design: ShipDesign;
	};

	let { design }: Props = $props();

	// unlike the rest of a design, the description can be changed while ships of the design exist
	let description = $state(untrack(() => design.description));
	let saving = $state(false);
	let changed = $derived(description.trim() !== design.description);

	async function save() {
		saving = true;
		try {
			const updated = await updateDesignDescription(design.num, description);
			description = updated.description;
		} catch (e) {
			addError(e as ConnectError);
		} finally {
			saving = false;
		}
	}
</script>

{#if $readOnly}
	{#if design.description}
		<p class="whitespace-pre-wrap wrap-break-word mx-2 my-1" data-type="design-description">
			{design.description}
		</p>
	{/if}
{:else}
	<div class="flex flex-col gap-1 mx-2 my-1" data-type="design-description-editor">
		<textarea
			class="textarea w-full text-base"
			name="description"
			aria-label="Design notes"
			rows={2}
			maxlength={MaxDescriptionLength}
			placeholder="What this design is for, only you can see it"
			bind:value={description}></textarea>
		{#if changed}
			<div class="flex flex-row justify-end gap-2">
				<button
					type="button"
					class="btn btn-sm btn-outline"
					disabled={saving}
					onclick={() => (description = design.description)}>Cancel</button
				>
				<button
					type="button"
					class="btn btn-sm btn-primary"
					disabled={saving}
					onclick={save}
					data-type="design-description-save">Save notes</button
				>
			</div>
		{/if}
	</div>
{/if}
