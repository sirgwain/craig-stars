<script lang="ts">
    import { run } from 'svelte/legacy';

	import { startCase } from 'lodash-es';

    interface Props {
        name: string;
        checked: boolean | undefined;
        title?: string | undefined;
        titleClass?: string;
    }

    let {
        name,
        checked = $bindable(),
        title = $bindable(undefined),
        titleClass = "label-text w-32 text-right"
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
                class="checkbox ml-2"
                type="checkbox"
                {name}
                bind:checked={checked}
            />
        </label>
    </div>
</div>
