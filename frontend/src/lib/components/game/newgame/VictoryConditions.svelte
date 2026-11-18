<script lang="ts">
	import type { Snippet } from 'svelte';
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsScore,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionHighestScoreAfterYears,
		VictoryConditionOwnCapitalShips,
		VictoryConditionOwnPlanets,
		VictoryConditionProductionCapacity
	} from '#lib/types/Consts.js';
	import {
		VictoryConditionsSchema,
		type GameSettings,
		type VictoryConditions
	} from '#lib/types/cs-proto.js';
	import { create } from '@bufbuild/protobuf';
	import VictoryConditionCheckbox from './VictoryConditionCheckbox.svelte';
	import VictoryConditionInput from './VictoryConditionInput.svelte';

	type Props = {
		settings: GameSettings;
	};

	let { settings = $bindable() }: Props = $props();

	// Local runes state for reactive bindings
	let victoryConditions: VictoryConditions = $state(
		settings.victoryConditions ?? create(VictoryConditionsSchema)
	);

	$effect(() => {
		settings.victoryConditions = victoryConditions;
	});
</script>

{#snippet criterion(title: string, condition: number, fields: Snippet)}
	<div
		class="grid gap-2 border-b border-base-300 py-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center"
	>
		<label class="flex items-center gap-3"
			><VictoryConditionCheckbox bind:conditions={victoryConditions.conditions} {condition} /><span
				>{title}</span
			></label
		>
		<div class="ml-7 flex flex-wrap items-center gap-2 sm:ml-0">{@render fields()}</div>
	</div>
{/snippet}
<div class="text-sm">
	{#snippet planets()}<VictoryConditionInput
			label="Planets owned"
			bind:value={victoryConditions.ownPlanets}
			min={20}
			max={100}
			unit="%"
		/>of planets{/snippet}
	{@render criterion('Own planets', VictoryConditionOwnPlanets, planets)}
	{#snippet tech()}Level<VictoryConditionInput
			label="Victory tech level"
			bind:value={victoryConditions.attainTechLevel}
			min={8}
			max={26}
		/>in<VictoryConditionInput
			label="Technology fields"
			bind:value={victoryConditions.attainTechLevelNumFields}
			min={2}
			max={6}
		/>fields{/snippet}
	{@render criterion('Attain technology levels', VictoryConditionAttainTechLevels, tech)}
	{#snippet score()}<VictoryConditionInput
			label="Winning score"
			bind:value={victoryConditions.exceedsScore}
			min={1000}
			max={20000}
		/>points{/snippet}
	{@render criterion('Exceed score', VictoryConditionExceedsScore, score)}
	{#snippet second()}<VictoryConditionInput
			label="Lead over second place"
			bind:value={victoryConditions.exceedsSecondPlaceScore}
			min={20}
			max={300}
			unit="%"
		/>ahead{/snippet}
	{@render criterion('Exceed second place score', VictoryConditionExceedsSecondPlaceScore, second)}
	{#snippet production()}<VictoryConditionInput
			label="Production capacity"
			bind:value={victoryConditions.productionCapacity}
			min={10}
			max={500}
			step={10}
		/>thousand resources/yr{/snippet}
	{@render criterion('Production capacity', VictoryConditionProductionCapacity, production)}
	{#snippet ships()}<VictoryConditionInput
			label="Capital ships"
			bind:value={victoryConditions.ownCapitalShips}
			min={10}
			max={300}
			step={10}
		/>ships{/snippet}
	{@render criterion('Own capital ships', VictoryConditionOwnCapitalShips, ships)}
	{#snippet highest()}<VictoryConditionInput
			label="Highest score after years"
			bind:value={victoryConditions.highestScoreAfterYears}
			min={30}
			max={900}
			step={10}
		/>years{/snippet}
	{@render criterion('Highest score after years', VictoryConditionHighestScoreAfterYears, highest)}
	<div class="grid gap-3 pt-3 sm:grid-cols-2">
		<label class="flex flex-col items-start gap-2"
			>Selected criteria required<VictoryConditionInput
				label="Selected criteria required"
				bind:value={victoryConditions.numCriteriaRequired}
				min={1}
				max={7}
			/></label
		>
		<label class="flex flex-col items-start gap-2"
			>Minimum years before victory<VictoryConditionInput
				label="Minimum years before victory"
				bind:value={victoryConditions.yearsPassed}
				min={30}
				max={500}
				step={10}
			/></label
		>
	</div>
</div>
