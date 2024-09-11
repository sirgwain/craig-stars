import { get, writable } from 'svelte/store';

export type ErrorResponse = { status?: string; error: string };

export const errors = writable<CSError[]>([]);

export function addError(err: CSError | string | undefined) {
	if (err) {
		const errs = get(errors);
		if (typeof err === 'string') {
			errors.update(() => [...errs, new CSError(undefined, err, 0)]);
		} else {
			errors.update(() => [...errs, err]);
		}
	}
}

export class CSError implements ErrorResponse {
	statusCode: number | undefined = undefined;
	status = '';
	error = '';

	constructor(data: ErrorResponse | undefined, status: string, statusCode: number) {
		this.status = status;
		this.error = status;
		this.statusCode = statusCode;
		if (typeof data === 'string') {
			this.status = this.error = data as string;
		} else {
			Object.assign(this, data);
		}
	}

	toString(): string {
		return (this.statusCode ? `${this.statusCode}: ` : '') + this.error;
	}
}
