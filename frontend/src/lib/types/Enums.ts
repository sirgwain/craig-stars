import { GameStartMode } from '$lib/types/cs-proto';

import { camelCase, startCase } from 'lodash-es';

/**
 * Converts a protobuf enum value into a PascalCase string without the enum prefix.
 */
export function enumToString<T extends Record<string, string | number>>(
	enumType: T,
	value: T[keyof T]
): string {
	if (value === 0) return ''; // assuming 0 == UNSPECIFIED

	// Get the string key from numeric value
	const key = enumType[value as number] as string | undefined;
	if (!key) return '';

	// Figure out the common prefix from the first non-numeric key
	const sampleKey = Object.keys(enumType).find((k) => isNaN(Number(k)));
	if (!sampleKey) return '';

	const parts = sampleKey.split('_');
	// Remove last part so we get the type prefix only
	const typePrefix = parts.slice(0, -1).join('_') + '_';

	// Remove that prefix from the actual key
	const withoutPrefix = key.startsWith(typePrefix) ? key.substring(typePrefix.length) : key;

	// Convert to PascalCase
	return startCase(camelCase(withoutPrefix));
}

export const GameStartModeFullNames: {
	[key in GameStartMode]: string;
} = {
	[GameStartMode.UNSPECIFIED]: 'Normal',
	[GameStartMode.ACC_BBS]: 'Accelerated BBS Play',
	[GameStartMode.MAX]: 'Max Start'
};


