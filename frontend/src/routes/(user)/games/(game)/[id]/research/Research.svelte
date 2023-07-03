<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { NextResearchField, Player } from '$lib/types/Player';
	import { Beaker } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	import NumberInput from '$lib/components/NumberInput.svelte';
	import Factory from '$lib/components/icons/Factory.svelte';
	import Microscope from '$lib/components/icons/Microscope.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { TechField, type TechLevel } from '$lib/types/TechLevel';
	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';
	import { $enum as eu } from 'ts-enum-util';
	import FutureTechs from './FutureTechs.svelte';

	const dispatch = createEventDispatcher();

	const { game, player, universe } = getGameContext();

	const getLevel = (player: Player, field: TechField | string): number => {
		const f: keyof TechLevel = `${field}`.toLowerCase() as keyof TechLevel;
		return $player.techLevels[f] ?? 0;
	};

	const updatePlayerOrders = async () => {
		dispatch('update-player');
	};

	let spent = 0;
	$: {
		const field: keyof TechLevel = `${$player.researching}`.toLowerCase() as keyof TechLevel;
		spent = $player.techLevelsSpent[field] ?? 0;
	}

	$: leftToSpend = ($player.spec.currentResearchCost ?? 0) - spent;
	$: yearsLeft = (leftToSpend / ($player.spec.resourcesPerYearResearch ?? 0)) + 1;
	$: console.log("leftToSpend", leftToSpend, "spent", spent, "yearsLeft", yearsLeft, "spend per year", $player.spec.resourcesPerYearResearch)
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
			{spent ?? 0}/{$player.spec.currentResearchCost} resources, {yearsLeft.toFixed()}
			{Math.floor(yearsLeft) > 1 ? 'years' : 'year'}
		</div>
	</div>
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Resources Available</div>
		<div class="stat-figure"><Factory class="w-8 h-8" /></div>
		<div class="stat-value">
			{$player.spec.resourcesPerYear}
		</div>
	</div>
</div>
<div class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full">
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Spent Last Year</div>
		<div class="stat-figure"><Microscope class="w-8 h-8" /></div>
		<div class="stat-value">
			{$player.researchSpentLastYear ?? 0}
		</div>
	</div>
	<div class="grow stat place-items-center">
		<div class="stat-title">Spending Next Year</div>
		<div class="stat-figure"><Microscope class="w-8 h-8" /></div>
		<div class="stat-value">
			{$player.spec.resourcesPerYearResearch ?? 0}
		</div>
	</div>
</div>
<div class="grid grid-cols-1 md:grid-cols-2 justify-center gap-2">
	<div class="w-full">
		<SectionHeader>Tech Levels</SectionHeader>

		<NumberInput
			name="researchAmount"
			title="Research Budget"
			bind:value={$player.researchAmount}
			min={0}
			max={100}
			step={1}
			unit="%"
			on:change={updatePlayerOrders}
		/>

		<div class="grid grid-cols-2">
			<div class="text-center">
				Field of Study <div class="divider secondary w-[90%]" />
			</div>
			<div class="text-center">
				Current Level <div class="divider secondary w-[90%]" />
			</div>
			{#each eu(TechField).getKeys() as field}
				<div class="form-control">
					<label class="label cursor-pointer">
						<span class="label-text">{startCase(field.toString())}</span>
						<input
							type="radio"
							name="researching"
							value={field}
							class="radio radio-sm checked:bg-primary"
							bind:group={$player.researching}
							on:change={updatePlayerOrders}
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
			enumType={NextResearchField}
			bind:value={$player.nextResearchField}
			on:change={updatePlayerOrders}
		/>
	</div>

	<div class="w-full">
		<SectionHeader>Expected Research Benefits</SectionHeader>
		<FutureTechs field={$player.researching} />
	</div>
</div>
