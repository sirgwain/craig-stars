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
		title = startCase(name),
		titleClass = 'label-text w-32 text-right'
	}: Props = $props();
</script>

<div class="w-full flex-grow">
	<div class="form-control">
		<label class="label"
			><span class={titleClass}>{title}</span>
			<input class="checkbox ml-2" type="checkbox" {name} bind:checked />
		</label>
	</div>
</div>
