<script lang="ts">
	import {
		ARM,
		BET,
		CE,
		GR,
		IFE,
		ISB,
		LSP,
		MA,
		NAS,
		NRSE,
		OBRM,
		RS,
		TT,
		UR,
		type Race
	} from '$lib/types/cs';

	type Props = {
		race: Race;
	};

	let { race = $bindable() }: Props = $props();
</script>

<div class="flex flex-row flex-wrap justify-between gap-1">
	<button
		type="button"
		onclick={() => (race.lrts ^= IFE)}
		class:btn-primary={race.lrts & IFE}
		class="btn w-[11rem]">Improved Fuel Efficiency</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= TT)}
		class:btn-primary={race.lrts & TT}
		class="btn w-[11rem]">Total Terraforming</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= ARM)}
		class:btn-primary={race.lrts & ARM}
		class="btn w-[11rem]">Advanced Remote Mining</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= ISB)}
		class:btn-primary={race.lrts & ISB}
		class="btn w-[11rem]">Improved Starbases</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= GR)}
		class:btn-primary={race.lrts & GR}
		class="btn w-[11rem]">Generalized Research</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= UR)}
		class:btn-primary={race.lrts & UR}
		class="btn w-[11rem]">Ultimate Recycling</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= NRSE)}
		class:btn-primary={race.lrts & NRSE}
		class="btn w-[11rem]">No Ramscoop Engines</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= OBRM)}
		class:btn-primary={race.lrts & OBRM}
		class="btn w-[11rem]">Only Basic Remote Mining</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= NAS)}
		class:btn-primary={race.lrts & NAS}
		class="btn w-[11rem]">No Advanced Scanners</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= LSP)}
		class:btn-primary={race.lrts & LSP}
		class="btn w-[11rem]">Low Starting Population</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= BET)}
		class:btn-primary={race.lrts & BET}
		class="btn w-[11rem]">Bleeding Edge Technology</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= RS)}
		class:btn-primary={race.lrts & RS}
		class="btn w-[11rem]">Regenerating Shields</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= MA)}
		class:btn-primary={race.lrts & MA}
		class="btn w-[11rem]">Mineral Alchemy</button
	>

	<button
		type="button"
		onclick={() => (race.lrts ^= CE)}
		class:btn-primary={race.lrts & CE}
		class="btn w-[11rem]">Cheap Engines</button
	>
</div>
