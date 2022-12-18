<script lang="ts">
	import Login from '$lib/components/Login.svelte';
	import Menu from '$lib/components/Menu.svelte';
	import { me } from '$lib/services/Context';
	import '../../css/app.css';
	import '../../css/hulls.css';
	import '../../css/planets.css';
	import '../../css/techs.css';
</script>

<main class="p-3 flex flex-col h-screen">
	<div class="flex-initial">
		<Menu user={$me} />
	</div>
	<div class="flex-1 h-full">
		<slot>This is the main content</slot>
	</div>
</main>
