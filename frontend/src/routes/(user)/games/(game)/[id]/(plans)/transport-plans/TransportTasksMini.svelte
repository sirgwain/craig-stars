<script lang="ts">
	import type { WaypointTransportTask, WaypointTransportTasks } from '$lib/types/Fleet';
	import TransportActionIcon from './TransportActionIcon.svelte';

	export let transportTasks: WaypointTransportTasks;
</script>

<div class="flex flex-row justify-between mt-1">
	<div class="text-fuel">F</div>
	<TransportActionIcon action={transportTasks.fuel.action} />
	<div class="text-ironium">I</div>
	<TransportActionIcon action={transportTasks.ironium.action} />
	<div class="text-boranium">B</div>
	<TransportActionIcon action={transportTasks.boranium.action} />
	<div class="text-germanium">G</div>
	<TransportActionIcon action={transportTasks.germanium.action} />
	<div class="text-colonists">C</div>
	<TransportActionIcon action={transportTasks.colonists.action} />
</div>
