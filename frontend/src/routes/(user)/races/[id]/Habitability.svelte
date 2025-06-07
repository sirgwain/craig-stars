<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import { Grav, Rad, Temp, type Race } from '$lib/types/cs';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';
	import Habitation from './Habitation.svelte';

	type Props = {
		race: Race;
	};

	let { race = $bindable() }: Props = $props();
</script>

<div class="flex flex-col gap-2">
	<Habitation
		habType={Grav}
		bind:habLow={race.habLow.grav}
		bind:habHigh={race.habHigh.grav}
		bind:immune={race.immuneGrav}
	/>
	<Habitation
		habType={Temp}
		bind:habLow={race.habLow.temp}
		bind:habHigh={race.habHigh.temp}
		bind:immune={race.immuneTemp}
	/>
	<Habitation
		habType={Rad}
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
