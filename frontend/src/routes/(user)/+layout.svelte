<script lang="ts">
	import Menu from '$lib/components/Menu.svelte';
	import { me } from '$lib/services/Context';
	import '../../css/app.css';
	import '../../css/hulls.css';
	import '../../css/planets.css';
	import '../../css/mapobjects.css';
	import '../../css/techs.css';
</script>

<main class="p-3 flex flex-col">
	<div class="sticky top-0 z-10">
		<Menu user={$me} />
	</div>
	<div class="w-full mx-auto md:max-w-2xl">
		<slot>This is the main content</slot>
	</div>
</main>
