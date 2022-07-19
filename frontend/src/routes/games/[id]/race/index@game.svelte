<script lang="ts">
	import { player } from '$lib/services/Context';
</script>

{#if $player}
	<div class="w-full">
		You are playing as the {$player?.race?.pluralName}
	</div>
{/if}
