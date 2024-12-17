<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { Player, PlayerRelation, type PlayerRelationship } from '$lib/types/Player';

	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { universe } = getGameContext();

	type Props = {
		player: Player;
		onUpdatePlayerRelationships?: (relations: PlayerRelationship[]) => void;
	};
	let { player, onUpdatePlayerRelationships: onUpdatePlayerRelationship }: Props = $props();

	let relations: PlayerRelationship[] = $state([]);

	$effect(() => {
		relations = $state.snapshot(player.relations);
	});

	function updateRelationship() {
		onUpdatePlayerRelationship?.(relations);
	}
</script>

<ItemTitle>Relations</ItemTitle>

<div class="flex flex-col justify-between gap-1">
	{#each relations as relation, index}
		{#if player.num != index + 1}
			<!-- content here -->
			<SectionHeader>{$universe.getPlayerPluralName(index + 1)}</SectionHeader>
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Friend</span>
					<input
						type="radio"
						name={`player-relation-${index + 1}`}
						class="radio checked:bg-success"
						value={PlayerRelation.Friend}
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
						value={PlayerRelation.Neutral}
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
						value={PlayerRelation.Enemy}
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
