<script lang="ts">
	import { techs } from '#lib/services/Stores.js';
	import { None } from '#lib/types/Consts.js';
	import { QuestionMarkCircle, Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import ConfirmDialog from '../ConfirmDialog.svelte';
	import TechAvatar from '../tech/TechAvatar.svelte';
	import Cost from './Cost.svelte';
	import DesignStats from './DesignStats.svelte';
	import { onShipDesignTooltip } from './tooltips/ShipDesignTooltip';
	import type { ShipDesign } from '#lib/types/cs-proto.js';

	type Props = {
		design: ShipDesign;
		href: string;
		copyhref: string;
		onDelete?: (design: ShipDesign) => Promise<void>;
		// how many of this design are in production queues, they're removed along with the design
		numQueued?: number;
		// hide the delete/copy/edit actions
		readOnly?: boolean;
	};

	let { design, href, copyhref, onDelete, numQueued = 0, readOnly = false }: Props = $props();

	// deleting a design scraps everything built from it and removes it from every production queue,
	// so that gets a louder confirmation
	let numInstances = $derived(design.spec?.numInstances ?? 0);
	const count = (num: number) =>
		`${num} ${design.spec?.starbase ? 'starbase' : 'ship'}${num === 1 ? '' : 's'}`;
	let instancesText = $derived(count(numInstances));
	let queuedText = $derived(count(numQueued));
	let destructive = $derived(numInstances > 0 || numQueued > 0);
	let confirmingDelete = $state(false);
</script>

<div
	class="card bg-base-200 shadow-sm rounded-xs border-2 border-base-300 pt-2 w-full sm:w-[430px]"
>
	<figure>
		<div class="border border-secondary bg-black p-1">
			<a class="cs-link" {href}>
				<TechAvatar tech={$techs.getHull(design.hull)} hullSetNumber={design.hullSetNumber} />
			</a>
		</div>
	</figure>
	<div class="card-body">
		<h2 class="card-title">
			<div class="flex flex-row gap-1 min-w-0">
				<div class="min-w-0">
					<a class="cs-link" {href}>{design.name}</a>
				</div>
				<div>
					<button
						type="button"
						class="size-full cursor-help"
						onpointerdown={(e) => onShipDesignTooltip(e, design)}
					>
						<Icon src={QuestionMarkCircle} size="16" class="cursor-help inline-block" />
					</button>
				</div>
			</div>
		</h2>
		{#if design.description}
			<p class="text-sm opacity-80 line-clamp-2 wrap-break-word" data-type="design-description">
				{design.description}
			</p>
		{/if}
		<div class="flex flex-col sm:flex-row justify-between">
			<div class="mr-2">
				<Cost cost={design.spec?.cost} />
			</div>
			<DesignStats spec={design.spec} />
		</div>
		<div class="card-actions justify-start">
			<div class="grow">
				<!-- for a design another player gave us, built is how many of them we were given -->
				<div
					class="tooltip"
					data-tip={`${design.spec?.numInstances ?? 0} remaining of ${design.spec?.numBuilt ?? 0} ${
						design.originalPlayerNum == None ? 'built' : 'received'
					}`}
				>
					<span class="btn rounded-lg border border-secondary normal-case">
						{design.spec?.numInstances ?? 0} of {design.spec?.numBuilt ?? 0}
					</span>
				</div>
			</div>
			<div class="flex flex-row join">
				{#if !design.cannotDelete && !readOnly}
					<button
						class="btn btn-outline btn-secondary joint-item"
						onclick={() => (confirmingDelete = true)}
					>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				{/if}
				{#if design.originalPlayerNum == None && !design.mysteryTrader && !readOnly}
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

<ConfirmDialog
	bind:show={confirmingDelete}
	title={`Delete ${design.name}?`}
	{destructive}
	confirmText={numInstances > 0 ? `Delete design and ${instancesText}` : 'Delete design'}
	onConfirm={() => onDelete?.(design)}
>
	{#if destructive}
		<p>Deleting this design will also:</p>
		<ul class="list-disc pl-6 py-2">
			{#if numInstances > 0}
				<li data-type="delete-design-instances">
					destroy the <strong class="text-error">{instancesText}</strong> built from it
				</li>
			{/if}
			{#if numQueued > 0}
				<li data-type="delete-design-queued">
					remove the <strong class="text-error">{queuedText}</strong> queued for production
				</li>
			{/if}
		</ul>
		<p>This cannot be undone.</p>
	{:else}
		<p>Are you sure you want to delete this design? It is currently unused.</p>
	{/if}
</ConfirmDialog>
