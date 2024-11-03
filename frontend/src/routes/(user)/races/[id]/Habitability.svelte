<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import { HabTypes } from '$lib/types/Hab';
	import type { Race } from '$lib/types/Race';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';
	import HabBar from './HabBar.svelte';

	interface Props {
		race: Race;
	}

	let { race = $bindable() }: Props = $props();
</script>

<div class="flex flex-col gap-2">
	<HabBar
		habType={HabTypes.Gravity}
		bind:habLow={race.habLow.grav}
		bind:habHigh={race.habHigh.grav}
		bind:immune={race.immuneGrav}
	/>
	<HabBar
		habType={HabTypes.Temperature}
		bind:habLow={race.habLow.temp}
		bind:habHigh={race.habHigh.temp}
		bind:immune={race.immuneTemp}
	/>
	<HabBar
		habType={HabTypes.Radiation}
		bind:habLow={race.habLow.rad}
		bind:habHigh={race.habHigh.rad}
		bind:immune={race.immuneRad}
	/>
	<SpinnerNumberText min={1} max={20} bind:value={race.growthRate}>
		{#snippet begin()}
				Maximum Colonist Growth Rate Per Year
			{/snippet}
		{#snippet end()}
				%.
			{/snippet}
	</SpinnerNumberText>

	<HabChance {race} />
</div>
