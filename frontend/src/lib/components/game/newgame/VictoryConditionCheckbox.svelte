<script lang="ts">
	import type { VictoryCondition } from '$lib/types/cs';

	type Props = {
		conditions: number;
		condition: VictoryCondition;
	};

	let { conditions = $bindable(), condition }: Props = $props();
</script>

<input
	class="w-12"
	type="checkbox"
	checked={(conditions & condition) > 0}
	onchange={(e) =>
		e.currentTarget.checked ? (conditions |= condition) : (conditions &= ~condition)}
/>
