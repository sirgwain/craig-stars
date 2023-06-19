<script lang="ts">
	import Menu from '$lib/components/Menu.svelte';
	import { me } from '$lib/services/Context';
</script>

<div class="mx-auto md:max-w-2xl">
	<slot>This is the main content</slot>
</div>
