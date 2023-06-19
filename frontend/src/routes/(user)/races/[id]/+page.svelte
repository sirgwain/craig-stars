<script lang="ts">
	import { page } from '$app/stores';
	import RaceEditor from '$lib/components/RaceEditor.svelte';

	let id = parseInt($page.params.id);
</script>

<RaceEditor {id} />
