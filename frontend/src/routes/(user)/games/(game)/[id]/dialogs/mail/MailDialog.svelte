<script lang="ts">
	import { clickOutside } from '#lib/clickOutside.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import Mail from './Mail.svelte';

	// Any part of the game can open this by calling openMail on the game context.
	const { mailTo, closeMail } = getGameContext();
</script>

<!-- On a phone the dialog sits at the top and stays short, so all of it is above the keyboard -->
<div
	class="modal modal-top md:modal-middle"
	class:modal-open={$mailTo !== undefined}
	data-type="mail-dialog"
>
	<div
		class="modal-box w-full max-w-full max-h-dvh pt-[max(1rem,env(safe-area-inset-top))] md:max-w-lg md:pt-6"
		use:clickOutside={closeMail}
	>
		{#if $mailTo !== undefined}
			{#key $mailTo}
				<Mail toPlayerNum={$mailTo} onClose={closeMail} />
			{/key}
		{/if}
	</div>
</div>
