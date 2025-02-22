<script lang="ts">
	import {
		VictoryConditionAttainTechLevels,
		VictoryConditionExceedsScore,
		VictoryConditionExceedsSecondPlaceScore,
		VictoryConditionHighestScoreAfterYears,
		VictoryConditionOwnCapitalShips,
		VictoryConditionOwnPlanets,
		VictoryConditionProductionCapacity,
		type GameSettings
	} from '$lib/types/cs';
	import VictoryConditionCheckbox from './VictoryConditionCheckbox.svelte';
	import VictoryConditionInput from './VictoryConditionInput.svelte';

	type Props = {
		settings: GameSettings;
	};

	let { settings = $bindable() }: Props = $props();
</script>

<div>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionOwnPlanets}
		/>
		<span class="label-text mr-1">
			Owns
			<VictoryConditionInput
				bind:value={settings.victoryConditions.ownPlanets}
				min={20}
				max={100}
				unit="%"
			/>
			of all planets.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionAttainTechLevels}
		/>
		<span class="label-text mr-1">
			Attains Tech
			<VictoryConditionInput
				bind:value={settings.victoryConditions.attainTechLevel}
				min={8}
				max={26}
			/>
			in
			<VictoryConditionInput
				bind:value={settings.victoryConditions.attainTechLevelNumFields}
				min={2}
				max={6}
			/>
			fields.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionExceedsScore}
		/>
		<span class="label-text mr-1">
			Exceeds a score of
			<VictoryConditionInput
				bind:value={settings.victoryConditions.exceedsScore}
				min={1000}
				max={20000}
			/>
			.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionExceedsSecondPlaceScore}
		/>
		<span class="label-text mr-1">
			Exceeds second place score by
			<VictoryConditionInput
				bind:value={settings.victoryConditions.exceedsSecondPlaceScore}
				min={20}
				max={300}
				unit="%"
			/>
			.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionProductionCapacity}
		/>
		<span class="label-text mr-1">
			Has a production capacity of
			<VictoryConditionInput
				bind:value={settings.victoryConditions.productionCapacity}
				min={10}
				max={500}
				step={10}
			/>,000 resources/yr.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionOwnCapitalShips}
		/>
		<span class="label-text mr-1">
			Owns
			<VictoryConditionInput
				bind:value={settings.victoryConditions.ownCapitalShips}
				min={10}
				max={300}
				step={10}
			/>
			capital ships.
		</span>
	</label>
	<label class="label justify-start">
		<VictoryConditionCheckbox
			bind:conditions={settings.victoryConditions.conditions}
			condition={VictoryConditionHighestScoreAfterYears}
		/>
		<span class="label-text mr-1">
			Has the highest score after
			<VictoryConditionInput
				bind:value={settings.victoryConditions.highestScoreAfterYears}
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
			<VictoryConditionInput
				bind:value={settings.victoryConditions.numCriteriaRequired}
				min={1}
				max={7}
			/>
			of the above selected criteria.
		</span>
	</label>
	<label class="label justify-start">
		<span class="label-text mr-1">
			At least
			<VictoryConditionInput
				bind:value={settings.victoryConditions.yearsPassed}
				min={30}
				max={500}
				step={10}
			/>
			years must pass before a winner is declared.
		</span>
	</label>
</div>
