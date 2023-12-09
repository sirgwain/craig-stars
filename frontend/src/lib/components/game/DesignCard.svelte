<script lang="ts">
	import { techs } from '$lib/services/Stores';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import { QuestionMarkCircle, Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import TechAvatar from '../tech/TechAvatar.svelte';
	import Cost from './Cost.svelte';
	import DesignStats from './DesignStats.svelte';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip.svelte';
	import { None } from '$lib/types/MapObject';

	const dispatch = createEventDispatcher();

	export let design: ShipDesign;
	export let href: string;
	export let copyhref: string;

	const deleteDesign = async (design: ShipDesign) => {
		if (design.num != undefined && confirm(`Are you sure you want to delete ${design.name}?`)) {
			dispatch('delete', { design });
		}
	};
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 w-full sm:w-[430px]">
	<figure>
		<div class="border border-secondary bg-black p-1">
			<a class="cs-link" {href}>
				<TechAvatar tech={$techs.getHull(design.hull)} hullSetNumber={design.hullSetNumber} />
			</a>
		</div>
	</figure>
	<div class="card-body">
		<h2 class="card-title">
			<div class="flex flex-row gap-1">
				<div>
					<a class="cs-link" {href}>{design.name}</a>
				</div>
				<div>
					<button
						type="button"
						class="w-full h-full cursor-help"
						on:pointerdown|preventDefault={(e) => onShipDesignTooltip(e, design)}
					>
						<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
					</button>
				</div>
			</div>
		</h2>
		<div class="flex flex-col sm:flex-row justify-between">
			<div class="mr-2">
				<Cost cost={design.spec.cost} />
			</div>
			<DesignStats spec={design.spec} />
		</div>
		<div class="card-actions justify-start">
			<div class="grow">
				{#if design.originalPlayerNum == None}
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
				{/if}
			</div>
			<div class="flex flex-row join">
				{#if !design.cannotDelete}
					<button
						class="btn btn-outline btn-secondary joint-item"
						on:click={(e) => deleteDesign(design)}
					>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				{/if}
				{#if design.originalPlayerNum == None}
					<a class="btn btn-outline btn-secondary joint-item" href={copyhref}>Copy</a>
					<!-- cannotDelete = true is for designs that are reserved for the system -->
					{#if !design.spec?.numInstances && !design.cannotDelete}
						<a class="btn btn-outline btn-secondary joint-item" href={`${href}}/edit`}>Edit</a>
					{/if}
				{/if}
			</div>
		</div>
	</div>
</div>
