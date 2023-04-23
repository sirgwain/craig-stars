<script lang="ts">
	import { player, techs } from '$lib/services/Context';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { Icon } from '@steeze-ui/svelte-icon';
	import TechAvatar from '../tech/TechAvatar.svelte';
	import Cost from './Cost.svelte';
	import DesignStats from './DesignStats.svelte';
	import { Trash } from '@steeze-ui/heroicons';
	import { DesignService } from '$lib/services/DesignService';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let design: ShipDesign;
	export let gameId: number;

	const deleteDesign = async (design: ShipDesign) => {
		if (design.num != undefined && confirm(`Are you sure you want to delete ${design.name}?`)) {
			const { fleets, starbases } = await DesignService.delete(gameId, design.num);
			if ($player) {
				const p = $player;
				const designs = p.designs.filter((d) => d.num !== design.num);
				player.update(() => ({ ...p, fleets, designs }));
			}
			dispatch('deleted', { design });
		}
	};
</script>

<div class="card bg-base-200 shadow-xl rounded-sm border-2 border-base-300 pt-2">
	<figure>
		<TechAvatar tech={$techs.getHull(design.hull)} hullSetNumber={design.hullSetNumber} />
	</figure>
	<div class="card-body">
		<h2 class="card-title">
			<a class="cs-link" href={`/games/${gameId}/designs/${design.num}`}>{design.name}</a>
		</h2>
		<div class="flex flex-row justify-between">
			<div class="mr-2">
				<Cost cost={design.spec.cost} />
			</div>
			<DesignStats {design} />
		</div>
		<div class="card-actions justify-end">
			<button class="btn" on:click={(e) => deleteDesign(design)}>
				<Icon src={Trash} size="24" class="hover:stroke-accent" />
			</button>
			{#if !design.spec?.numInstances}
				<a class="btn" href={`/games/${gameId}/designs/${design.num}/edit`}>Edit</a>
			{/if}
		</div>
	</div>
</div>
