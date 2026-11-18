<script lang="ts">
	import Select from '#lib/components/Select.svelte';
	import type { NewGamePlayer, Race } from '#lib/types/cs-proto.js';

	type Props = { player: NewGamePlayer; races: Race[] };
	let { player = $bindable(), races }: Props = $props();

	$effect(() => {
		if (!player.race && races[0]) player.race = races[0];
	});

	function raceChanged(id: bigint) {
		const race = races.find((r) => r.id === id);
		if (race) player.race = race;
	}
</script>

<Select
	values={races.map((r) => ({ value: r.id, title: r.pluralName }))}
	name="Host"
	value={player.race?.id ?? races[0]?.id}
	onchange={(e) => raceChanged(BigInt(e.currentTarget.value))}
/>
