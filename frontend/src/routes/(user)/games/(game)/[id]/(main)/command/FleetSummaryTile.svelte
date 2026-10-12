<script lang="ts">
	import { MaxNameLength } from '#lib/types/Consts.js';
	import type { NextPrevMapObjectProps, RenameFleetProps } from '#lib/services/Events.js';

	import { getGameContext } from '#lib/services/GameContext.js';
	import { getHullIcon } from '#lib/techicon.js';
	import type { ShipDesign } from '#lib/types/cs-proto.js';
	import type { CommandedFleet } from '#lib/types/Fleet.js';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import CommandTile from './CommandTile.svelte';

	const { universe, readOnly, editNote } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		hideTitle?: boolean;
	} & NextPrevMapObjectProps &
		RenameFleetProps;

	let { fleet, hideTitle, onRenameFleet, onNextMapObject, onPreviousMapObject }: Props = $props();

	const design: ShipDesign | undefined = $derived.by(() => {
		if (fleet.tokens.length > 0) {
			const designNum = fleet.tokens[0].designNum;
			return $universe.getDesign(fleet.mapObject.playerNum, designNum);
		}
	});

	let note = $derived($universe.getNote(fleet));

	function rename() {
		let name = prompt('Enter fleet name', fleet.baseName);
		if (!name || name === '') {
			name = $universe.getMyDesign(fleet.tokens[0].designNum)?.name ?? '';
		}
		if (name !== '') {
			onRenameFleet?.({ fleet, name: name.slice(0, MaxNameLength) });
		}
	}
</script>

<CommandTile title={hideTitle ? '' : fleet.mapObject.name}>
	{#snippet actions()}
		{#if !hideTitle && (note || !$readOnly)}
			<button
				type="button"
				class="flex cursor-pointer"
				title={note?.text || (note ? 'edit notes and tags' : 'add notes and tags')}
				aria-label={note ? 'Edit notes and tags' : 'Add notes and tags'}
				data-type="fleet-title-note-button"
				onclick={() => editNote(fleet)}
			>
				<Icon
					src={PencilSquare}
					class={`size-4 hover:stroke-accent ${note ? 'stroke-accent' : ''}`}
				/>
			</button>
		{/if}
	{/snippet}
	<div class="grid grid-cols-2">
		<div class="avatar border border-secondary p-2 bg-black m-auto relative">
			{#if fleet.tokens.reduce((count, t) => count + t.quantity, 0) > 1}
				<div class="absolute -right-2 -top-1 text-xl size-6">+</div>
			{/if}
			<div class="fleet-avatar {getHullIcon(design)} bg-black"></div>
		</div>
		<div class="flex flex-col gap-y-1">
			<button
				onclick={onPreviousMapObject}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Prev</button
			>
			<button
				onclick={onNextMapObject}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Next</button
			>
			<button
				onclick={rename}
				type="button"
				class="btn btn-outline btn-sm normal-case btn-secondary">Rename</button
			>
		</div>
	</div>
</CommandTile>
