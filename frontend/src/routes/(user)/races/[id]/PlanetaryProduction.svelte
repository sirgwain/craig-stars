<script lang="ts">
	import type { Race } from '$lib/types/Race';
	
	export let race: Race;

	const updatePopEfficiency = (value: number) => {
		race.popEfficiency = value / 100;
	};
</script>

<p>
	One resource is generated each year for every
	<input
		class="input input-bordered w-24 input-sm"
		type="number"
		name="popEfficiency"
		value={race.popEfficiency * 100}
		on:change={(e) => updatePopEfficiency(e.currentTarget.valueAsNumber)}
		step={100}
		min={700}
		max={2500}
	/>
	colonists.
</p>
<p>
	Every 10 factories produce
	<input
		class="input input-bordered w-16 input-sm"
		type="number"
		name="factoryOutput"
		bind:value={race.factoryOutput}
		step={1}
		min={5}
		max={15}
	/>
	resources each year.
</p>
<p>
	Factories require
	<input
		class="input input-bordered w-16 input-sm"
		type="number"
		name="factoryCost"
		bind:value={race.factoryCost}
		step={1}
		min={5}
		max={25}
	/>
	resources to build.
</p>
<p>
	Every 10,000 colonists may operate up to
	<input
		class="input input-bordered w-16 input-sm"
		type="number"
		name="numFactories"
		bind:value={race.numFactories}
		step={1}
		min={5}
		max={25}
	/>
	factories.
</p>

<p>
	<input
		class="input input-bordered input-xs"
		type="checkbox"
		name="factoriesCostLess"
		bind:checked={race.factoriesCostLess}
	/>
	Factories cost 1kT less of Germanium to build
</p>
<p>
	Every 10 mines produce up to
	<input
		class="input input-bordered w-16 input-sm"
		type="number"
		name="mineOutput"
		bind:value={race.mineOutput}
		step={1}
		min={5}
		max={25}
	/>kT of each mineral every year.
</p>
<p>
	Mines require
	<input
		class="input input-bordered w-16 input-sm"
		type="number"
		name="mineCost"
		bind:value={race.mineCost}
		step={1}
		min={2}
		max={15}
	/>
	resources to build.
</p>
<p>
	Every 10,000 colonists may operate up to
	<input
		class="input input-bordered w-16 input-sm"
		type="number"
		name="numMines"
		bind:value={race.numMines}
		step={1}
		min={5}
		max={25}
	/>
	mines.
</p>
