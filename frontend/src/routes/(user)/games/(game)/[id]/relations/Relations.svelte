<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { PlayerRelation, type PlayerRelationship } from '$lib/types/cs-proto';
	import { CommandedPlayer } from '$lib/types/Player';

	const { universe } = getGameContext();

	type Props = {
		player: CommandedPlayer;
		onUpdatePlayerRelationships?: (relations: PlayerRelationship[]) => void;
	};
	let { player, onUpdatePlayerRelationships: onUpdatePlayerRelationship }: Props = $props();

	let relations: PlayerRelationship[] = $derived($state.snapshot(player.relations));

	function updateRelationship() {
		onUpdatePlayerRelationship?.(relations);
	}
</script>

<ItemTitle>Relations</ItemTitle>

<div class="flex flex-col justify-between gap-1">
	{#each relations as relation, index (index)}
		{#if player.num != index + 1}
			<SectionHeader>{$universe.getPlayerPluralName(index + 1)}</SectionHeader>
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Friend</span>
					<input
						type="radio"
						name={`player-relation-${index + 1}`}
						class="radio checked:bg-success"
						value={PlayerRelation.FRIEND}
						bind:group={relation.relation}
						onchange={updateRelationship}
					/>
				</label>
			</div>
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Neutral</span>
					<input
						type="radio"
						name={`player-relation-${index + 1}`}
						class="radio checked:bg-info"
						value={PlayerRelation.NEUTRAL}
						bind:group={relation.relation}
						onchange={updateRelationship}
					/>
				</label>
			</div>
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Enemy</span>
					<input
						type="radio"
						name={`player-relation-${index + 1}`}
						class="radio checked:bg-error"
						value={PlayerRelation.ENEMY}
						bind:group={relation.relation}
						onchange={updateRelationship}
					/>
				</label>
			</div>
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Share Map</span>
					<input
						type="checkbox"
						name={`player-relation-${index + 1}-share-map`}
						class="checkbox"
						bind:checked={relation.shareMap}
						onchange={updateRelationship}
					/>
				</label>
			</div>
		{/if}
	{/each}
</div>
