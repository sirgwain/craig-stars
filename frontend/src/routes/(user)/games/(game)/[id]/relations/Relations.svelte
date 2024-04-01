<script lang="ts" context="module">
	export type PlayerUpdateEvent = {
		'update-player': void;
	};
</script>

<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { PlayerRelation } from '$lib/types/Player';

	import { getGameContext } from '$lib/services/GameContext';
	import { createEventDispatcher } from 'svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';

	const dispatch = createEventDispatcher<PlayerUpdateEvent>();

	const { game, player, universe } = getGameContext();

	function updatePlayerOrders() {
		dispatch('update-player');
	}
</script>

<ItemTitle>Relations</ItemTitle>

<div class="flex flex-col justify-between gap-1">
	{#each $player.relations as relation, index}
		{#if $player.num != index + 1}
			<!-- content here -->
			<SectionHeader>{$universe.getPlayerName(index + 1)}</SectionHeader>
			<div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Friend</span>
					<input
						type="radio"
						name={`player-relation-${index + 1}`}
						class="radio checked:bg-success"
						value={PlayerRelation.Friend}
						bind:group={relation.relation}
						on:change={updatePlayerOrders}
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
						on:change={updatePlayerOrders}
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
						on:change={updatePlayerOrders}
					/>
				</label>
			</div>
			<!-- Coming Soon! -->
			<!-- <div class="form-control">
				<label class="label cursor-pointer">
					<span class="label-text">Share Map</span>
					<input
						type="checkbox"
						name={`player-relation-${index + 1}-share-map`}
						class="checkbox"
						bind:checked={relation.shareMap}
						on:change={updatePlayerOrders}
					/>
				</label>
			</div> -->
		{/if}
	{/each}
</div>
