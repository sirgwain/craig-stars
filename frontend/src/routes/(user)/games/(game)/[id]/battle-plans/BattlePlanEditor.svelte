<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { BattleAttackWho, BattleTactic, BattleTarget } from '$lib/types/Battle';
	import type { BattlePlan } from '$lib/types/Player';
	import { ExclamationTriangle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';

	const dispatch = createEventDispatcher();

	export let plan: BattlePlan;
	export let readonlyName = false;
	export let error: string = '';

	const onSubmit = async () => {
		dispatch('save');
	};
</script>

<form on:submit|preventDefault={onSubmit}>
	{#if error !== ''}
		<div
			class="alert alert-error shadow-lg w-1/2 mx-auto"
			in:fade
			out:fade={{ delay: 5000 }}
			on:introend={(e) => (error = '')}
		>
			<div>
				<Icon src={ExclamationTriangle} size="24" class="hover:stroke-accent" />
				<span>{error}</span>
			</div>
		</div>
	{/if}

	<TextInput name="name" bind:value={plan.name} required disabled={readonlyName} />
	<EnumSelect name="primaryTarget" enumType={BattleTarget} bind:value={plan.primaryTarget} />
	<EnumSelect name="secondaryTarget" enumType={BattleTarget} bind:value={plan.secondaryTarget} />
	<EnumSelect name="tactic" enumType={BattleTactic} bind:value={plan.tactic} />
	<EnumSelect name="attackWho" enumType={BattleAttackWho} bind:value={plan.attackWho} />
</form>
