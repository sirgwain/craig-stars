<script lang="ts">
	import Discord from '../icons/Discord.svelte';

	type Props = {
		from?: string;
	};

	let { from = '/' }: Props = $props();
</script>

<a class="btn bg-[#7289da]" id="discord-login" href={`/api/auth/discord/login?from=${from}`}
	>Login&nbsp;
	<Discord class="fill-base-content size-4 " />
</a>
