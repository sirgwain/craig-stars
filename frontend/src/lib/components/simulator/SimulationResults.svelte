<script lang="ts">
	import { outcome, type SimulationRun } from './scenario';
	import SimulationReplay from './SimulationReplay.svelte';
	let { runs }: { runs: SimulationRun[] } = $props();
	let selected = $state(0);
	let showDetails = $state(false);
	let run = $derived(runs[selected] ?? runs[0]);
	const labels = ['A holds the field', 'B holds the field', 'No decisive outcome'];
	let counts = $derived(
		labels.map((label) => runs.filter((r) => outcome(r.response) === label).length)
	);
	function average(playerNum: number, field: 'survivors' | 'retreated' | 'destroyed') {
		return (
			runs.reduce(
				(n, r) =>
					n +
					r.response.stacks
						.filter((s) => s.playerNum === playerNum)
						.reduce((sum, s) => sum + s[field], 0),
				0
			) / runs.length
		).toFixed(1);
	}
	function loss(playerNum: number, field: 'resources' | 'ironium' | 'boranium' | 'germanium') {
		return (
			runs.reduce(
				(n, r) =>
					n +
					r.response.stacks
						.filter((s) => s.playerNum === playerNum)
						.reduce((sum, s) => sum + (s.loss?.[field] ?? 0), 0),
				0
			) / runs.length
		).toFixed(0);
	}
</script>

<div class="border border-secondary bg-base-200 p-3" data-testid="simulation-summary">
	<h2 class="text-xl font-semibold">{runs.length} {runs.length === 1 ? 'battle' : 'battles'}</h2>
	<div class="mt-2 flex h-3 overflow-hidden rounded" aria-hidden="true">
		{#each counts as count, i (i)}<div
				class={i === 0 ? 'bg-info' : i === 1 ? 'bg-error' : 'bg-secondary'}
				style={`width:${(count / runs.length) * 100}%`}
			></div>{/each}
	</div>
	<div class="mt-2 flex flex-wrap justify-between gap-2">
		{#each labels as label, i (label)}<span>{label}: <strong>{counts[i]}</strong></span>{/each}
	</div>
	<button
		class="btn btn-sm btn-outline btn-secondary mt-3 md:hidden"
		aria-expanded={showDetails}
		onclick={() => (showDetails = !showDetails)}>Survivors and losses</button
	>
	<div class="md:block" class:hidden={!showDetails}>
		<div class="mt-3 overflow-x-auto">
			<table class="table table-sm">
				<thead
					><tr
						><th>Side</th><th>Avg. survivors</th><th>Retreated</th><th>Destroyed</th><th
							>Resources lost</th
						><th>Ironium</th><th>Boranium</th><th>Germanium</th></tr
					></thead
				><tbody
					>{#each [1, 2] as n (n)}<tr
							><th>{n === 1 ? 'A' : 'B'}</th><td>{average(n, 'survivors')}</td><td
								>{average(n, 'retreated')}</td
							><td>{average(n, 'destroyed')}</td><td>{loss(n, 'resources')}</td><td
								>{loss(n, 'ironium')}</td
							><td>{loss(n, 'boranium')}</td><td>{loss(n, 'germanium')}</td></tr
						>{/each}</tbody
				>
			</table>
		</div>
		<p class="mt-2 text-secondary">
			Survivors include ships that retreated. Holding the field means only one side has ships
			remaining on the board.
		</p>
	</div>
</div>
<div class="mt-3 flex flex-wrap items-center justify-between gap-2 border border-secondary p-2">
	<label class="flex items-center gap-2"
		>Replay<select aria-label="Battle replay" class="select select-sm" bind:value={selected}
			>{#each runs as r, i (i)}<option value={i}>Battle {i + 1} · {outcome(r.response)}</option
				>{/each}</select
		></label
	><span class="text-secondary">Seed {run.seed.toString()}</span>
</div>
{#key run}<SimulationReplay {run} />{/key}
<details class="mt-3 border border-secondary p-3">
	<summary class="cursor-pointer font-semibold">Stack losses for this battle</summary>
	<div class="overflow-x-auto">
		<table class="table table-sm">
			<thead
				><tr
					><th>Side</th><th>Design / stack</th><th>Survived</th><th>Retreated</th><th>Destroyed</th
					></tr
				></thead
			><tbody
				>{#each run.response.stacks as stack (`${stack.playerNum}-${stack.designNum}`)}<tr
						><td>{stack.playerNum === 1 ? 'A' : 'B'}</td><td
							>{run.response.designs.find(
								(d) => d.playerNum === stack.playerNum && d.num === stack.designNum
							)?.name} #{stack.designNum}</td
						><td>{stack.survivors}</td><td>{stack.retreated}</td><td>{stack.destroyed}</td></tr
					>{/each}</tbody
			>
		</table>
	</div>
</details>
