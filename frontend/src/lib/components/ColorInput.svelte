<script lang="ts">
	import { startCase } from 'lodash-es';
	export let name: string;
	export let value: string | undefined;

	export let title: string | undefined = undefined;
	export let titleClass = 'label-text w-32 text-right';
	export let required = false;

	$: !title && (title = startCase(name));
    $: console.log(value)
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
