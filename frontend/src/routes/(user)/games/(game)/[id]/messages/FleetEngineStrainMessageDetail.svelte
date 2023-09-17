<script lang="ts">
	import type { Message } from '$lib/types/Message';

	export let message: Message;
</script>

{#if message.spec.amount == 1}
	One of your ships in {message.spec.name} was destroyed when its engines over-reacted trying to achieve
	Warp 10.
{:else}
	{message.spec.amount} of your ships in {message.spec.name} were destroyed due to engine strain.
{/if}
