import { VectorSchema, type Vector } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

export const emptyVector = () => create(VectorSchema, { x: 0, y: 0 });

export const equal = (v1: Vector, v2: Vector) => v1.x === v2.x && v1.y === v2.y;

// compute the distance between two vectors
export const distance = (v1: Vector, v2: Vector): number =>
	v2 && v1
		? Math.sqrt(
				((v1?.x ?? 0) - (v2?.x ?? 0)) * ((v1?.x ?? 0) - (v2?.x ?? 0)) +
					((v1?.y ?? 0) - (v2?.y ?? 0)) * ((v1?.y ?? 0) - (v2?.y ?? 0))
			)
		: 0;

export const lengthSquared = (v: Vector): number => v.x * v.x + v.y * v.y;

export const normalized = (from: Vector): Vector => {
	const v = create(VectorSchema, { x: from.x, y: from.y });
	const lengthsq = lengthSquared(v);

	if (lengthsq === 0) {
		v.x = 0;
		v.y = 0;
	} else {
		const length = Math.sqrt(lengthsq);
		v.x /= length;
		v.y /= length;
	}
	return v;
};

export const subtract = (from: Vector, to: Vector): Vector => {
	return create(VectorSchema, { x: from.x - to.x, y: from.y - to.y });
};

export const string = (v: Vector | undefined) => `(${v?.x ?? 0}, ${v?.y ?? 0})`;
