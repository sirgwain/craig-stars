<script lang="ts">
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Idle Fleets Filter"
	class:btn-primary={$settings.showIdleFleetsOnly}
	class:btn-ghost={!$settings.showIdleFleetsOnly}
	class="btn btn-ghost fill-current btn-xs w-full h-full"
	onclick={() => ($settings.showIdleFleetsOnly = !$settings.showIdleFleetsOnly)}
	><IdleFleets class="w-6 h-6" /></a
>
