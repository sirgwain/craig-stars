<script lang="ts">
	import type { Race } from '$lib/types/Race';

	export let race: Race;

	const updatePopEfficiency = (value: number) => {
		race.popEfficiency = value / 100;
	};
</script>

<p>tbd</p>
