<script lang="ts">
	import { page } from '$app/stores';
	import Game from '$lib/components/game/Game.svelte';

	let id = $page.params.id;
</script>

<Game {id} />
