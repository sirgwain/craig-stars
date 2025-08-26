import { RulesSchema, type Rules, type RulesJson } from '$lib/types/cs-proto';
import rulesjson from '$lib/ssr/rules.json';
import { create, fromJson } from '@bufbuild/protobuf';

export const defaultRules: Rules = create(
	RulesSchema,
	fromJson(RulesSchema, rulesjson as unknown as RulesJson)
);
