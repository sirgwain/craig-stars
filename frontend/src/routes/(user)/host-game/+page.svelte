<script lang="ts">
	import NewGame from '#lib/components/game/newgame/NewGame.svelte';
	import { hostedGameSettings } from '#lib/components/game/newgame/defaults.js';
	import { me } from '#lib/services/Stores.js';
	import { untrack } from 'svelte';

	const settings = untrack(() =>
		hostedGameSettings($me.username ? `${$me.username}'s game` : undefined)
	);
</script>

<NewGame initialSettings={settings} />
