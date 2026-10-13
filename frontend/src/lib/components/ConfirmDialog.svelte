<script lang="ts">
	import { ExclamationTriangle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import type { Snippet } from 'svelte';

	type Props = {
		// the dialog closes itself when it's confirmed or canceled
		show: boolean;
		title: string;
		// what is being confirmed, any markup
		children?: Snippet;
		// warn that confirming destroys something, with a warning icon and error colors
		destructive?: boolean;
		confirmText?: string;
		cancelText?: string;
		onConfirm?: () => void;
		onCancel?: () => void;
	};

	let {
		show = $bindable(),
		title,
		children,
		destructive = false,
		confirmText = 'Ok',
		cancelText = 'Cancel',
		onConfirm,
		onCancel
	}: Props = $props();

	const id = $props.id();

	function confirm() {
		show = false;
		onConfirm?.();
	}

	function cancel() {
		show = false;
		onCancel?.();
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (show && e.key === 'Escape') {
			cancel();
		}
	}}
/>

{#if show}
	<div
		class="modal modal-open"
		role="alertdialog"
		aria-modal="true"
		aria-labelledby={`${id}-title`}
		data-type="confirm-dialog"
	>
		<div class="modal-box" class:border-2={destructive} class:border-error={destructive}>
			<h3
				id={`${id}-title`}
				class="flex items-center gap-2 text-lg font-bold"
				class:text-error={destructive}
			>
				{#if destructive}
					<Icon src={ExclamationTriangle} size="32" class="shrink-0" />
				{/if}
				<span class="min-w-0 wrap-break-word">{title}</span>
			</h3>
			<div class="py-4">
				{@render children?.()}
			</div>
			<div class="modal-action">
				<button
					type="button"
					class="btn btn-outline btn-secondary"
					data-type="confirm-dialog-cancel"
					onclick={cancel}>{cancelText}</button
				>
				<button
					type="button"
					class="btn"
					class:btn-error={destructive}
					class:btn-primary={!destructive}
					data-type="confirm-dialog-confirm"
					onclick={confirm}>{confirmText}</button
				>
			</div>
		</div>
	</div>
{/if}
