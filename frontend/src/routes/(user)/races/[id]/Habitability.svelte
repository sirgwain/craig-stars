<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import { HabType } from '$lib/types/Hab';
	import type { Race } from '$lib/types/Race';
	import HabBar from './HabBar.svelte';

	export let race: Race;
</script>

<div class="flex flex-col gap-2">
	<HabBar
		habType={HabType.Gravity}
		bind:habLow={race.habLow.grav}
		bind:habHigh={race.habHigh.grav}
		bind:immune={race.immuneGrav}
	/>
	<HabBar
		habType={HabType.Temperature}
		bind:habLow={race.habLow.temp}
		bind:habHigh={race.habHigh.temp}
		bind:immune={race.immuneTemp}
	/>
	<HabBar
		habType={HabType.Radiation}
		bind:habLow={race.habLow.rad}
		bind:habHigh={race.habHigh.rad}
		bind:immune={race.immuneRad}
	/>
	<label>
		Maxium Colonist Growth Rate Per Year
		<input
			class="input input-sm input-bordered"
			type="number"
			name="growthRate"
			min={1}
			max={100}
			bind:value={race.growthRate}
		/></label
	>

	<HabChance {race} />
</div>
