<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import {
		BattleAttackWho,
		type BattlePlan,
		BattleTactic,
		BattleTarget
	} from '$lib/types/cs-proto';

	type Props = {
		plan: BattlePlan;
	};

	let { plan = $bindable() }: Props = $props();

	// Local runes state for reactive bindings
	let name: string = $state(plan.name ?? '');
	let primaryTarget: BattleTarget = $state(plan.primaryTarget ?? BattleTarget.UNSPECIFIED);
	let secondaryTarget: BattleTarget = $state(plan.secondaryTarget ?? BattleTarget.UNSPECIFIED);
	let tactic: BattleTactic = $state(plan.tactic ?? BattleTactic.UNSPECIFIED);
	let attackWho: BattleAttackWho = $state(plan.attackWho ?? BattleAttackWho.UNSPECIFIED);

	// Sync local state back to plan
	$effect(() => {
		plan.name = name;
		plan.primaryTarget = primaryTarget;
		plan.secondaryTarget = secondaryTarget;
		plan.tactic = tactic;
		plan.attackWho = attackWho;
	});
</script>

<TextInput name="name" bind:value={name} required />
<EnumSelect name="primaryTarget" enumType={BattleTarget} bind:value={primaryTarget} />
<EnumSelect
	name="secondaryTarget"
	enumType={BattleTarget}
	showEmpty={true}
	bind:value={secondaryTarget}
/>
<EnumSelect name="tactic" enumType={BattleTactic} bind:value={tactic} />
<EnumSelect name="attackWho" enumType={BattleAttackWho} bind:value={attackWho} />
