<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import Designs from './Designs.svelte';

	const { universe } = getGameContext();
</script>

<Designs designs={$universe.designs} />
