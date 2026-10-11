<!--
  @component
  The percent of their range that scanners are drawn at.
 -->
<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { clamp } from '#lib/services/Math.js';

	const { settings } = getGameContext();
</script>

<input
	class="input input-sm w-16 px-1"
	type="number"
	title="Scanner Range"
	min={0}
	max={100}
	step={10}
	value={$settings.scannerPercent}
	onchange={(e) => {
		const val = parseInt(e.currentTarget.value);
		if (val) {
			$settings.scannerPercent = clamp(val, 0, 100);
			e.currentTarget.value = `${$settings.scannerPercent}`;
		} else {
			$settings.scannerPercent = 0;
			e.currentTarget.value = '0';
		}
	}}
/>
