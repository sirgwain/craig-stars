import { timestampDate, timestampFromDate, type Timestamp } from '@bufbuild/protobuf/wkt';
import { format } from 'date-fns';

// convert an ios time string to a protobuf timestamp
export function isoToTimestamp(iso: string | undefined): Timestamp | undefined {
	if (!iso) {
		return;
	}
	const date = new Date(iso);
	return timestampFromDate(date);
}

export function timestampToDate(timestamp: Timestamp | undefined): Date {
	if (timestamp === undefined) {
		return new Date();
	}
	return timestampDate(timestamp);
}

export function timestampToString(
	timestamp: Timestamp | undefined,
	fmt = 'E, MMM do yyyy hh:mm aaa'
): string {
	if (timestamp === undefined) {
		return '';
	}
	return format(timestampToDate(timestamp), fmt);
}

/**
 * Compare two google.protobuf.Timestamp values.
 *
 * Returns:
 *  - negative if t1 < t2
 *  - zero if equal
 *  - positive if t1 > t2
 *
 * undefined handling:
 *  - configurable via undefinedOrder:
 *      'low'  => undefined < defined (default)
 *      'high' => undefined > defined
 *      'equal'=> undefined == undefined and undefined == any defined (treat all as equal)
 */
export function compare(
	t1: Timestamp | undefined,
	t2: Timestamp | undefined,
	undefinedOrder: 'low' | 'high' | 'equal' = 'low'
): number {
	// Fast-path equal reference
	if (t1 === t2) return 0;

	// Both undefined
	if (t1 === undefined && t2 === undefined) return 0;

	// One undefined
	if (t1 === undefined || t2 === undefined) {
		if (undefinedOrder === 'equal') return 0;
		if (undefinedOrder === 'low') return t1 === undefined ? -1 : 1;
		// 'high'
		return t1 === undefined ? 1 : -1;
	}

	// Both defined: compare seconds first
	const s1 = Number(t1.seconds ?? 0);
	const s2 = Number(t2.seconds ?? 0);
	if (s1 !== s2) return s1 - s2;

	// Seconds equal, compare nanos
	const n1 = t1.nanos ?? 0;
	const n2 = t2.nanos ?? 0;
	return n1 - n2;
}
