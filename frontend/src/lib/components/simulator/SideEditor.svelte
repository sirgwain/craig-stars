<script lang="ts">
	import TechAvatar from '#lib/components/tech/TechAvatar.svelte';
	import { techs } from '#lib/services/Stores.js';
	import { getLabelForLRT, getLabelForPRT, lrts } from '#lib/types/Race.js';
	import { enumToString, battleTargetToString } from '#lib/types/Enums.js';
	import {
		BattleTactic,
		BattleTarget,
		BattleSimulationStackSchema,
		Prt,
		type BattleSimulationSide,
		type ShipDesign
	} from '#lib/types/cs-proto.js';
	import { clone } from '@bufbuild/protobuf';
	import { newStack, techFields } from './scenario';
	let {
		side = $bindable(),
		index,
		designs,
		onEdit,
		onNew,
		onImport,
		onUseGame,
		hasGame
	}: {
		side: BattleSimulationSide;
		index: number;
		designs: ShipDesign[];
		onEdit: (stack: number) => void;
		onNew: () => void;
		onImport: () => void;
		onUseGame: () => void;
		hasGame: boolean;
	} = $props();
	const tactics = Object.values(BattleTactic).filter(
		(v): v is BattleTactic => typeof v === 'number' && v > 0
	);
	const targets = Object.values(BattleTarget).filter(
		(v): v is BattleTarget => typeof v === 'number'
	);
	const prts = Object.values(Prt).filter((v): v is Prt => typeof v === 'number' && v > 0);
	let choices = $derived([...side.stacks.flatMap((s) => (s.design ? [s.design] : [])), ...designs]);
	let selectedDesign = $state(0);
	function duplicate(i: number) {
		side.stacks.splice(i + 1, 0, clone(BattleSimulationStackSchema, side.stacks[i]));
	}
</script>

<section
	class="min-w-0 rounded border border-secondary bg-base-200"
	aria-label={`Side ${index === 0 ? 'A' : 'B'}`}
>
	<div class="flex items-center justify-between border-b border-secondary p-3">
		<h2 class="text-xl font-semibold text-secondary">
			<span class="mr-2 badge badge-outline">{index === 0 ? 'A' : 'B'}</span>{index === 0
				? 'Your forces'
				: 'Opposing forces'}
		</h2>
		{#if hasGame}<button class="btn btn-sm btn-outline btn-secondary" onclick={onImport}
				>Import fleet</button
			>{/if}
	</div>
	<div class="p-3">
		{#if side.race && side.techLevels}
			<details class="mb-3 border border-base-300 p-2">
				<summary class="cursor-pointer font-semibold"
					>{side.race.name} · Race &amp; technology</summary
				>
				{#if hasGame}<button class="btn btn-sm btn-outline btn-secondary mt-2" onclick={onUseGame}
						>Use current game settings</button
					>{/if}
				<div class="mt-3 grid gap-2 sm:grid-cols-2">
					<label class="flex flex-col gap-1"
						>Race name<input class="input input-sm w-full" bind:value={side.race.name} /></label
					>
					<label class="flex flex-col gap-1"
						>Primary trait<select class="select select-sm w-full" bind:value={side.race.prt}
							>{#each prts as prt (prt)}<option value={prt}>{getLabelForPRT(prt)}</option
								>{/each}</select
						></label
					>
				</div>
				<div class="mt-3 grid grid-cols-3 gap-2 sm:grid-cols-6">
					{#each techFields as field (field)}<label class="flex min-w-0 flex-col gap-1 capitalize"
							>{field.slice(0, 4)}<input
								aria-label={`${index === 0 ? 'A' : 'B'} ${field}`}
								class="input input-sm w-full min-w-0"
								type="number"
								min="0"
								max="26"
								step="1"
								bind:value={side.techLevels[field]}
							/></label
						>{/each}
				</div>
				<details class="mt-3">
					<summary class="cursor-pointer text-secondary">Lesser racial traits</summary>
					<div class="mt-2 grid gap-2 sm:grid-cols-2">
						{#each lrts as lrt (lrt)}<label class="flex items-center gap-2"
								><input
									class="checkbox checkbox-sm checkbox-secondary"
									type="checkbox"
									checked={!!(side.race.lrts & lrt)}
									onchange={() => {
										if (side.race) side.race.lrts ^= lrt;
									}}
								/>{getLabelForLRT(lrt)}</label
							>{/each}
					</div>
				</details>
			</details>
		{/if}
		<div class="mb-2 flex items-center justify-between">
			<h3 class="font-semibold">Ship stacks</h3>
			<span class="text-secondary"
				>{side.stacks.reduce((n, s) => n + (s.quantity || 0), 0)} ships · {side.stacks.length}
				{side.stacks.length === 1 ? 'stack' : 'stacks'}</span
			>
		</div>
		<div class="flex flex-col gap-2">
			{#each side.stacks as stack, i (stack)}
				<div class="border border-base-300 bg-base-100 p-2" data-testid="simulation-stack">
					<div class="flex items-center gap-2">
						<div class="shrink-0 border border-secondary bg-black">
							<TechAvatar
								tech={$techs.getHull(stack.design?.hull ?? '')}
								hullSetNumber={stack.design?.hullSetNumber}
							/>
						</div>
						<div class="min-w-0 grow">
							<div class="font-semibold">{stack.design?.name}</div>
							<div class="text-secondary">{stack.design?.hull}</div>
							<button class="btn btn-xs btn-outline btn-secondary mt-1" onclick={() => onEdit(i)}
								>Edit design</button
							>
						</div>
						<label class="flex w-20 shrink-0 flex-col gap-1"
							>Quantity<input
								aria-label={`Stack ${i + 1} quantity`}
								class="input input-sm w-full"
								type="number"
								min="1"
								max="10000"
								step="1"
								bind:value={stack.quantity}
							/></label
						>
					</div>
					{#if stack.plan}
						<details class="mt-2 border-t border-base-300 pt-2">
							<summary class="cursor-pointer text-secondary"
								>Battle plan · {enumToString(BattleTactic, stack.plan.tactic)}</summary
							>
							<div class="mt-2 grid gap-2 sm:grid-cols-2">
								<label class="flex min-w-0 flex-col gap-1 sm:col-span-2"
									>Tactic<select class="select select-sm w-full" bind:value={stack.plan.tactic}
										>{#each tactics as tactic (tactic)}<option value={tactic}
												>{enumToString(BattleTactic, tactic)}</option
											>{/each}</select
									></label
								>
								<label class="flex min-w-0 flex-col gap-1"
									>Primary target<select
										class="select select-sm w-full"
										bind:value={stack.plan.primaryTarget}
										>{#each targets as target (target)}<option value={target}
												>{battleTargetToString(target)}</option
											>{/each}</select
									></label
								>
								<label class="flex min-w-0 flex-col gap-1"
									>Secondary target<select
										class="select select-sm w-full"
										bind:value={stack.plan.secondaryTarget}
										>{#each targets as target (target)}<option value={target}
												>{battleTargetToString(target)}</option
											>{/each}</select
									></label
								>
							</div>
						</details>
					{/if}
					<details class="mt-2">
						<summary class="cursor-pointer text-secondary"
							>Starting damage · {stack.quantityDamaged} ships at {stack.damage} dp</summary
						>
						<div class="mt-2 grid grid-cols-2 gap-2">
							<label class="flex min-w-0 flex-col gap-1"
								>Damaged ships<input
									class="input input-sm w-full"
									type="number"
									min="0"
									max={stack.quantity}
									step="1"
									bind:value={stack.quantityDamaged}
								/></label
							>
							<label class="flex min-w-0 flex-col gap-1"
								>Damage per ship (dp)<input
									class="input input-sm w-full"
									type="number"
									min="0"
									step="1"
									bind:value={stack.damage}
								/></label
							>
						</div>
					</details>
					<div class="mt-2 flex justify-end gap-2">
						<button
							class="btn btn-xs btn-outline btn-secondary"
							disabled={side.stacks.length >= 64}
							onclick={() => duplicate(i)}>Duplicate stack</button
						><button
							class="btn btn-xs btn-ghost text-error"
							onclick={() => side.stacks.splice(i, 1)}>Remove</button
						>
					</div>
				</div>
			{/each}
		</div>
		<div class="mt-3 flex flex-wrap items-end gap-2">
			<label class="flex min-w-0 grow flex-col gap-1"
				>Design<select
					aria-label="Design to add"
					class="select select-sm w-full"
					bind:value={selectedDesign}
					>{#each choices as design, i (i)}<option value={i}>{design.name} · {design.hull}</option
						>{/each}</select
				></label
			>
			<button
				class="btn btn-sm btn-outline btn-secondary"
				disabled={!choices[selectedDesign] || side.stacks.length >= 64}
				onclick={() => side.stacks.push(newStack(choices[selectedDesign]))}>Add stack</button
			>
			<button
				class="btn btn-sm btn-outline btn-secondary"
				disabled={side.stacks.length >= 64}
				onclick={onNew}>New design</button
			>
		</div>
	</div>
</section>
