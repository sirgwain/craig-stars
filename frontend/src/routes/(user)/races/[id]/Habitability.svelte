<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import { Grav, Rad, Temp } from '$lib/types/Hab';
	import { HabSchema, type Race } from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';
	import Habitation from './Habitation.svelte';

	type Props = {
		race: Race;
	};

	let { race = $bindable() }: Props = $props();

	let habLow = $state(create(HabSchema, race.habLow ?? {}));
	let habHigh = $state(create(HabSchema, race.habHigh ?? {}));

	$effect(() => {
		race.habLow = habLow;
		race.habHigh = habHigh;
	});
</script>

<div class="flex flex-col gap-2">
	<Habitation
		habType={Grav}
		bind:habLow={habLow.grav}
		bind:habHigh={habHigh.grav}
		bind:immune={race.immuneGrav}
	/>
	<Habitation
		habType={Temp}
		bind:habLow={habLow.temp}
		bind:habHigh={habHigh.temp}
		bind:immune={race.immuneTemp}
	/>
	<Habitation
		habType={Rad}
		bind:habLow={habLow.rad}
		bind:habHigh={habHigh.rad}
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
