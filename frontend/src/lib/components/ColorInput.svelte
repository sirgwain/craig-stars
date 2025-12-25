<script lang="ts">
	import { startCase } from 'lodash-es';

	type Props = {
		name: string;
		value: string | undefined;
		title?: string | undefined;
		titleClass?: string;
		required?: boolean;
	};

	let {
		name,
		value = $bindable(),
		titleClass = 'label-text w-32 text-right',
		required = false,
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full flex-grow">
	<div class="form-control">
		<label class="label"
			><span class={titleClass}>{title}</span>
			<input
				class="input input-bordered ml-2 flex-grow"
				type="color"
				{name}
				{required}
				bind:value
			/>
		</label>
	</div>
</div>
