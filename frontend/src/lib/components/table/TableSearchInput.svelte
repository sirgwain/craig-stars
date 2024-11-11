<script lang="ts">
	type Props = {
		value?: string;
		[key: string]: any;
	};

	let { value = $bindable(''), ...rest }: Props = $props();
</script>

<input
	type="text"
	name="search"
	placeholder="search"
	class="input input-bordered input-sm w-full sm:w-auto"
	autocomplete="off"
	autocorrect="off"
	autocapitalize="off"
	spellcheck="false"
	bind:value
	{...rest}
/>
