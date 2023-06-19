<script lang="ts">
	import { page } from '$app/stores';
	
	let gameId = parseInt($page.params.id);
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<slot />
</div>
