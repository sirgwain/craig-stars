<script lang="ts">
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Idle Fleets Filter"
	class:fill-accent={$settings.showIdleFleetsOnly}
	class:fill-current={!$settings.showIdleFleetsOnly}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={() => ($settings.showIdleFleetsOnly = !$settings.showIdleFleetsOnly)}
	><IdleFleets class="w-6 h-6" /></a
>
