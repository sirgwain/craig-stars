<script lang="ts">
	import { game, techs } from '$lib/services/Stores';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import TechAvatar from '../tech/TechAvatar.svelte';
	import Cost from './Cost.svelte';
	import DesignStats from './DesignStats.svelte';

	const dispatch = createEventDispatcher();

	export let design: ShipDesign;
	export let gameId: number;

	const deleteDesign = async (design: ShipDesign) => {
		if (design.num != undefined && confirm(`Are you sure you want to delete ${design.name}?`)) {
			dispatch('delete', { design });
		}
	};
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 w-full sm:w-[400px]">
	<figure>
		<div class="border border-secondary bg-black p-1">
			<a class="cs-link" href={`/games/${gameId}/designs/${design.num}`}>
				<TechAvatar tech={$techs.getHull(design.hull)} hullSetNumber={design.hullSetNumber} />
			</a>
		</div>
	</figure>
	<div class="card-body">
		<h2 class="card-title">
			<a class="cs-link" href={`/games/${gameId}/designs/${design.num}`}>{design.name}</a>
		</h2>
		<div class="flex flex-row justify-between">
			<div class="mr-2">
				<Cost cost={design.spec.cost} />
			</div>
			<DesignStats spec={design.spec} />
		</div>
		<div class="card-actions justify-start">
			<div class="grow">
				<div
					class="tooltip"
					data-tip={`${design.spec.numInstances ?? 0} remaining of ${
						design.spec.numBuilt ?? 0
					} built`}
				>
					<span class="btn rounded-lg border border-secondary normal-case">
						{design.spec.numInstances ?? 0} of {design.spec.numBuilt ?? 0}
					</span>
				</div>
			</div>
			<div>
				<button class="btn" on:click={(e) => deleteDesign(design)}>
					<Icon src={Trash} size="24" class="hover:stroke-accent" />
				</button>
				{#if !design.spec?.numInstances}
					<a class="btn" href={`/games/${gameId}/designs/${design.num}/edit`}>Edit</a>
				{/if}
			</div>
		</div>
	</div>
</div>
