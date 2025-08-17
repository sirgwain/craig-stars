<script lang="ts">
	import {
		VictoryConditionsSchema,
		type GameSettings,
		type VictoryConditions
	} from '$lib/types/cs-proto';
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsScore,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionHighestScoreAfterYears,
		VictoryConditionOwnCapitalShips,
		VictoryConditionOwnPlanets,
		VictoryConditionProductionCapacity
	} from '$lib/types/Consts';
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

<div>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionOwnPlanets}
		/>
		<span class="label-text mr-1">
			Owns
			<VictoryConditionInput
				bind:value={victoryConditions.ownPlanets}
				min={20}
				max={100}
				unit="%"
			/>
			of all planets.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionAttainTechLevels}
		/>
		<span class="label-text mr-1">
			Attains Tech
			<VictoryConditionInput bind:value={victoryConditions.attainTechLevel} min={8} max={26} />
			in
			<VictoryConditionInput
				bind:value={victoryConditions.attainTechLevelNumFields}
				min={2}
				max={6}
			/>
			fields.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionExceedsScore}
		/>
		<span class="label-text mr-1">
			Exceeds a score of
			<VictoryConditionInput bind:value={victoryConditions.exceedsScore} min={1000} max={20000} />
			.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionExceedsSecondPlaceScore}
		/>
		<span class="label-text mr-1">
			Exceeds second place score by
			<VictoryConditionInput
				bind:value={victoryConditions.exceedsSecondPlaceScore}
				min={20}
				max={300}
				unit="%"
			/>
			.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionProductionCapacity}
		/>
		<span class="label-text mr-1">
			Has a production capacity of
			<VictoryConditionInput
				bind:value={victoryConditions.productionCapacity}
				min={10}
				max={500}
				step={10}
			/>,000 resources/yr.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionOwnCapitalShips}
		/>
		<span class="label-text mr-1">
			Owns
			<VictoryConditionInput
				bind:value={victoryConditions.ownCapitalShips}
				min={10}
				max={300}
				step={10}
			/>
			capital ships.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={victoryConditions.conditions}
			condition={VictoryConditionHighestScoreAfterYears}
		/>
		<span class="label-text mr-1">
			Has the highest score after
			<VictoryConditionInput
				bind:value={victoryConditions.highestScoreAfterYears}
				min={30}
				max={900}
				step={10}
			/>
			years.
		</span>
	</label>
	<label class="label justify-start">
		<span class="label-text mr-1">
			Winner must meet
			<VictoryConditionInput bind:value={victoryConditions.numCriteriaRequired} min={1} max={7} />
			of the above selected criteria.
		</span>
	</label>
	<label class="label justify-start">
		<span class="label-text mr-1">
			At least
			<VictoryConditionInput
				bind:value={victoryConditions.yearsPassed}
				min={30}
				max={500}
				step={10}
			/>
			years must pass before a winner is declared.
		</span>
	</label>
</div>
