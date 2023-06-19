<script lang="ts">
	import CommandPane from '$lib/components/game/command/CommandPane.svelte';
	import Viewport from '$lib/components/game/Scanner.svelte';
</script>

<div class="flex flex-row h-full w-full">
	<div class="flex-none w-[31rem] h-full">
		<CommandPane />
	</div>
	<div class="grow ml-5 h-full border-gray-700 border-2 shadow-sm">
		<Viewport />
	</div>
</div>
