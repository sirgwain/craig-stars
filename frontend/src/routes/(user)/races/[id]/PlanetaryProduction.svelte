<script lang="ts">
	import { PRT, type Race } from '$lib/types/Race';
	import SpinnerNumberText from '../../../../lib/components/SpinnerNumberText.svelte';

	interface Props {
		race: Race;
	}

	let { race = $bindable() }: Props = $props();

	const updatePopEfficiency = (value: number) => {
		race.popEfficiency = value / 100;
	};
</script>

{#if race.prt === PRT.AR}
	<p>
		<SpinnerNumberText bind:value={race.popEfficiency} step={1} min={7} max={25}>
			{#snippet begin()}
						Annual Resources = Planet Value * sqrt(Population * Energy Tech /
					{/snippet}
			{#snippet end()}
						)
					{/snippet}
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
			{#snippet begin()}
						One resource is generated each year for every
					{/snippet}
			{#snippet end()}
						colonists.
					{/snippet}
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.factoryOutput} step={1} min={5} max={15}>
			{#snippet begin()}
						Every 10 factories produce
					{/snippet}
			{#snippet end()}
						resources each year.
					{/snippet}
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.factoryCost} step={1} min={5} max={25}>
			{#snippet begin()}
						Factories require
					{/snippet}
			{#snippet end()}
						resources to build.
					{/snippet}
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.numFactories} step={1} min={5} max={25}>
			{#snippet begin()}
						Every 10,000 colonists may operate up to
					{/snippet}
			{#snippet end()}
						factories.
					{/snippet}
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
			{#snippet begin()}
						Every 10 mines produce up to
					{/snippet}
			{#snippet end()}
						kT of each mineral every year.
					{/snippet}
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.mineCost} step={1} min={2} max={15}>
			{#snippet begin()}
						Mines require
					{/snippet}
			{#snippet end()}
						resources to build.
					{/snippet}
		</SpinnerNumberText>
	</p>
	<p>
		<SpinnerNumberText bind:value={race.numMines} step={1} min={5} max={25}>
			{#snippet begin()}
						Every 10,000 colonists may operate up to
					{/snippet}
			{#snippet end()}
						mines.
					{/snippet}
		</SpinnerNumberText>
	</p>
{/if}
