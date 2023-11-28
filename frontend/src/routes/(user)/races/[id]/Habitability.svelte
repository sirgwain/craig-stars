<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import { HabTypes } from '$lib/types/Hab';
	import type { Race } from '$lib/types/Race';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';
	import HabBar from './HabBar.svelte';

	export let race: Race;
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
	<SpinnerNumberText min={1} max={100} bind:value={race.growthRate}>
		<svelte:fragment slot="begin">Maxium Colonist Growth Rate Per Year</svelte:fragment>
		<svelte:fragment slot="end">%.</svelte:fragment>
	</SpinnerNumberText>

	<HabChance {race} />
</div>
