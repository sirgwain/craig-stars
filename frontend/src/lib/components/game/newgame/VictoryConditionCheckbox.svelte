<script lang="ts">
	import type { VictoryCondition } from '$lib/types/Game';

	export let conditions: number;
	export let condition: VictoryCondition;
</script>

<input
	class="w-12"
	type="checkbox"
	checked={(conditions & condition) > 0}
	on:change={(e) =>
		e.currentTarget.checked ? (conditions |= condition) : (conditions &= ~condition)}
/>
