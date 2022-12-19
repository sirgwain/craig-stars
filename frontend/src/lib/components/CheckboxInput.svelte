<script lang="ts">
	import { startCase } from 'lodash-es';
	export let name: string;
	export let checked: boolean | undefined;

	export let title: string | undefined = undefined;
	export let titleClass = "label-text w-32 text-right"

	$: !title && (title = startCase(name));
</script>

<div class="w-full flex-grow">
    <div class="form-control">
        <label class="label"
            ><span class={titleClass}>{title}</span>
            <input
                class="checkbox ml-2"
                type="checkbox"
                {name}
                bind:checked={checked}
            />
        </label>
    </div>
</div>
