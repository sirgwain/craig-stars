<script lang="ts">
	import { LRT, type Race } from '$lib/types/Race';

	export let race: Race;
</script>

<div class="flex flex-row flex-wrap justify-between gap-1">
	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.IFE)}
		class:btn-primary={race.lrts & LRT.IFE}
		class="btn w-[11rem]">Improved Fuel Efficiency</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.TT)}
		class:btn-primary={race.lrts & LRT.TT}
		class="btn w-[11rem]">Total Terraforming</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.ARM)}
		class:btn-primary={race.lrts & LRT.ARM}
		class="btn w-[11rem]">Advanced Remote Mining</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.ISB)}
		class:btn-primary={race.lrts & LRT.ISB}
		class="btn w-[11rem]">Improved Starbases</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.GR)}
		class:btn-primary={race.lrts & LRT.GR}
		class="btn w-[11rem]">Generalized Research</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.UR)}
		class:btn-primary={race.lrts & LRT.UR}
		class="btn w-[11rem]">Ultimate Recycling</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.NRSE)}
		class:btn-primary={race.lrts & LRT.NRSE}
		class="btn w-[11rem]">No Ramscoop Engines</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.OBRM)}
		class:btn-primary={race.lrts & LRT.OBRM}
		class="btn w-[11rem]">Only Basic Remote Mining</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.NAS)}
		class:btn-primary={race.lrts & LRT.NAS}
		class="btn w-[11rem]">No Advanced Scanners</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.LSP)}
		class:btn-primary={race.lrts & LRT.LSP}
		class="btn w-[11rem]">Low Starting Population</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.BET)}
		class:btn-primary={race.lrts & LRT.BET}
		class="btn w-[11rem]">Bleeding Edge Technology</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.RS)}
		class:btn-primary={race.lrts & LRT.RS}
		class="btn w-[11rem]">Regenerating Shields</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.MA)}
		class:btn-primary={race.lrts & LRT.MA}
		class="btn w-[11rem]">Mineral Alchemy</button
	>

	<button
		type="button"
		on:click|preventDefault={() => (race.lrts ^= LRT.CE)}
		class:btn-primary={race.lrts & LRT.CE}
		class="btn w-[11rem]">Cheap Engines</button
	>
</div>
