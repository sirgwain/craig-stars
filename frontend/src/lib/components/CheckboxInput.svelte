<script lang="ts">
	import { startCase } from 'lodash-es';

	type Props = {
		name: string;
		checked: boolean | undefined;
		title?: string | undefined;
		titleClass?: string;
	};

	let {
		name,
		checked = $bindable(),
		titleClass = 'cs-label-text w-32 text-right',
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full grow">
	<div class="cs-form-control">
		<label class="cs-form-label"
			><span class={titleClass}>{title}</span>
			<input class="checkbox ml-2" type="checkbox" {name} bind:checked />
		</label>
	</div>
</div>
