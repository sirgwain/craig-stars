<script lang="ts">
	import { humanoid } from '$lib/types/Race';
	import { loadWasm, type CS } from '$lib/wasm';

	import { onMount } from 'svelte';

	let cs: CS | undefined = $state();
	onMount(async () => {
		cs = await loadWasm();
	});

	const h = humanoid();
	let points = $state(0);
	$effect(() => {
		if (!cs) {
			// make sure wasm is loaded
			return;
		}

		cs.wasmService.calculateRacePoints({ race: h }).then((resp) => (points = resp.points));
	});
</script>

<h1 class="text-xl">Race {h.name}</h1>
<div class="flex flex-col">
	<div class="flex flex-row gap-1">
		<div>Points</div>
		<div>{points}</div>
	</div>
</div>
