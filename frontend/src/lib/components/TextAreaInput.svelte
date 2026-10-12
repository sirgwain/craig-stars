<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { HTMLTextareaAttributes } from 'svelte/elements';

	type Props = {
		name: string;
		value: string | undefined;
		title?: string | undefined;
		titleClass?: string;
	} & HTMLTextareaAttributes;

	let {
		name,
		value = $bindable(),
		titleClass = 'cs-label-text w-32 text-right',
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full grow">
	<div class="cs-form-control">
		<label class="cs-form-label items-start!"
			><span class="{titleClass} pt-2 shrink-0">{title}</span>
			<textarea class="textarea ml-2 grow text-base" {name} bind:value {...rest}></textarea>
		</label>
	</div>
</div>
