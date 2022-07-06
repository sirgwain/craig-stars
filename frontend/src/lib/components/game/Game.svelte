<script lang="ts">
	import CommandPane from '$lib/components/game/command/CommandPane.svelte';
	import Viewport from '$lib/components/game/Scanner.svelte';
</script>

<div class="flex-none w-[31rem]">
	<CommandPane />
</div>
<div class="ml-5 flex-1 h-full border-gray-700 border-2 shadow-sm">
	<Viewport />
</div>
