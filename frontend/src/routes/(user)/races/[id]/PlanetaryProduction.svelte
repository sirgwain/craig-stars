<script lang="ts">
	import { PRT, type Race } from '$lib/types/Race';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';

	export let race: Race;

	const updatePopEfficiency = (value: number) => {
		race.popEfficiency = value / 100;
	};
</script>

{#if race.prt === PRT.AR}
	<p>
		<SpinnerNumberText bind:value={race.popEfficiency} step={1} min={7} max={25}>
			<svelte:fragment slot="begin"
				>Annual Resources = Planet Value * sqrt(Population * Energy Tech /</svelte:fragment
			>
			<svelte:fragment slot="end">)</svelte:fragment>
		</SpinnerNumberText>
	</p>
{:else}
	<p>
		<SpinnerNumberText
			step={100}
			value={race.popEfficiency * 100}
			on:change={(e) => updatePopEfficiency(e.detail)}
			min={700}
			max={2500}
		>
			<svelte:fragment slot="begin">One resource is generated each year for every</svelte:fragment>
			<svelte:fragment slot="end">colonists.</svelte:fragment>
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.factoryOutput} step={1} min={5} max={15}>
			<svelte:fragment slot="begin">Every 10 factories produce</svelte:fragment>
			<svelte:fragment slot="end">resources each year.</svelte:fragment>
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.factoryCost} step={1} min={5} max={25}>
			<svelte:fragment slot="begin">Factories require</svelte:fragment>
			<svelte:fragment slot="end">resources to build.</svelte:fragment>
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.numFactories} step={1} min={5} max={25}>
			<svelte:fragment slot="begin">Every 10,000 colonists may operate up to</svelte:fragment>
			<svelte:fragment slot="end">factories.</svelte:fragment>
		</SpinnerNumberText>
	</p>

	<p>
		<input
			class="checkbox checkbox-xs"
			type="checkbox"
			name="factoriesCostLess"
			bind:checked={race.factoriesCostLess}
		/>
		Factories cost 1kT less of Germanium to build
	</p>
	<p>
		<SpinnerNumberText bind:value={race.mineOutput} step={1} min={5} max={25}>
			<svelte:fragment slot="begin">Every 10 mines produce up to</svelte:fragment>
			<svelte:fragment slot="end">kT of each mineral every year.</svelte:fragment>
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.mineCost} step={1} min={2} max={15}>
			<svelte:fragment slot="begin">Mines require</svelte:fragment>
			<svelte:fragment slot="end">resources to build.</svelte:fragment>
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.numMines} step={1} min={5} max={25}>
			<svelte:fragment slot="begin">Every 10,000 colonists may operate up to</svelte:fragment>
			<svelte:fragment slot="end">mines.</svelte:fragment>
		</SpinnerNumberText>
	</p>
{/if}
