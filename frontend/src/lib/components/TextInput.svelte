<script lang="ts">
	import { run } from 'svelte/legacy';

	import { startCase } from 'lodash-es';

	type Props = {
		name: string;
		value: string | undefined;
		title?: string | undefined;
		titleClass?: string;
		required?: boolean;
		disabled?: boolean;
	};

	let {
		name,
		value = $bindable(),
		title = $bindable(undefined),
		titleClass = 'label-text w-32 text-right',
		required = false,
		disabled = false
	}: Props = $props();

	run(() => {
		!title && (title = startCase(name));
	});
</script>

<div class="w-full flex-grow">
	<div class="form-control">
		<label class="label"
			><span class={titleClass}>{title}</span>
			<input
				class="input input-bordered ml-2 flex-grow"
				type="text"
				{name}
				{required}
				{disabled}
				bind:value
			/>
		</label>
	</div>
</div>
