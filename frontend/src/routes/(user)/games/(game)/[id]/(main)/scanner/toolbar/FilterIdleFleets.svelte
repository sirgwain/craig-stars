<script lang="ts">
	import IdleFleets from '#lib/components/icons/IdleFleets.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Idle Fleets Filter"
	class:btn-primary={$settings.showIdleFleetsOnly}
	class:btn-ghost={!$settings.showIdleFleetsOnly}
	class="btn btn-ghost fill-current btn-xs size-full p-0"
	onclick={() => ($settings.showIdleFleetsOnly = !$settings.showIdleFleetsOnly)}
	><IdleFleets class="size-6 " /></a
>
