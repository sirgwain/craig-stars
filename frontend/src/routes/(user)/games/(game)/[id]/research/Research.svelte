<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { CommandedPlayer, NextResearchFields, TechFields } from '$lib/types/Player';
	import { Beaker } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	import SpinnerNumberText from '$lib/components/SpinnerNumberText.svelte';
	import Factory from '$lib/components/icons/Factory.svelte';
	import Microscope from '$lib/components/icons/Microscope.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { TechField, TechLevel } from '$lib/types/cs';
	import { startCase } from 'lodash-es';
	import FutureTechs from './FutureTechs.svelte';

	const { player } = getGameContext();

	type Props = {
		onUpdatePlayer?: () => Promise<void>;
	};
	let { onUpdatePlayer }: Props = $props();

	const getLevel = (player: CommandedPlayer, field: TechField | string): number => {
		const f: keyof TechLevel = `${field}`.toLowerCase() as keyof TechLevel;
		return player.techLevels[f] ?? 0;
	};

	let field: keyof TechLevel = $derived(`${$player.researching}`.toLowerCase() as keyof TechLevel);
	let spent = $derived($player.techLevelsSpent[field] ?? 0);

	let leftToSpend = $derived(($player.spec.currentResearchCost ?? 0) - spent);
	let yearsLeft = $derived(
		Math.ceil(leftToSpend / ($player.spec.resourcesPerYearResearchEstimated ?? 0))
	);
</script>

<ItemTitle>Research</ItemTitle>
<div class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full">
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Researching</div>
		<div class="stat-figure"><Icon class="w-8 h-8" src={Beaker} /></div>
		<div class="stat-value">
			{$player.researching}
			{getLevel($player, $player.researching) + 1}
		</div>
		<div class="stat-desc pt-1">
			{spent ?? 0}/{$player.spec.currentResearchCost} resources
			{#if yearsLeft < 100}
				, {yearsLeft.toFixed()}
				{Math.floor(yearsLeft) > 1 ? 'years' : 'year'}
			{/if}
		</div>
	</div>
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Resources Available</div>
		<div class="stat-figure"><Factory class="w-8 h-8 fill-primary" /></div>
		<div class="stat-value">
			{$player.spec.resourcesPerYear ?? 0}
		</div>
	</div>
</div>
<div class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full">
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Spent Last Year</div>
		<div class="stat-figure"><Microscope class="w-8 h-8 fill-primary" /></div>
		<div class="stat-value">
			{$player.researchSpentLastYear ?? 0}
		</div>
	</div>
	<div class="grow stat place-items-center">
		<div class="stat-title">Estimated Spending Next Year</div>
		<div class="stat-figure"><Microscope class="w-8 h-8 fill-warning" /></div>
		<div class="stat-value">
			{$player.spec.resourcesPerYearResearchEstimated ?? 0}
		</div>
	</div>
</div>
<div class="grid grid-cols-1 md:grid-cols-2 justify-center gap-2">
	<div class="w-full">
		<SectionHeader>Tech Levels</SectionHeader>

		<SpinnerNumberText
			class="flex flex-row gap-1 place-content-center text-2xl mb-2"
			bind:value={$player.researchAmount}
			min={0}
			max={100}
			step={1}
			unit="%"
			onChange={onUpdatePlayer}
		>
			{#snippet begin()}
				Research Budget
			{/snippet}
			{#snippet end()}{/snippet}
		</SpinnerNumberText>

		<div class="grid grid-cols-2">
			<div class="text-center">
				Field of Study <div class="divider secondary w-[90%]"></div>
			</div>
			<div class="text-center">
				Current Level <div class="divider secondary w-[90%]"></div>
			</div>
			{#each TechFields as field}
				<div class="form-control">
					<label class="label cursor-pointer">
						<span class="label-text">{startCase(field.toString())}</span>
						<input
							type="radio"
							name="researching"
							value={field}
							class="radio radio-sm checked:bg-primary"
							bind:group={$player.researching}
							onchange={onUpdatePlayer}
						/>
					</label>
				</div>
				<div class="text-center">
					{getLevel($player, field)}
				</div>
			{/each}
		</div>
		<EnumSelect
			name="nextResearchField"
			options={NextResearchFields}
			bind:value={$player.nextResearchField}
			onchange={onUpdatePlayer}
		/>
	</div>

	<div class="w-full">
		<SectionHeader>Expected Research Benefits</SectionHeader>
		<FutureTechs field={$player.researching} />
	</div>
</div>
