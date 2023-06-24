import { get, writable } from 'svelte/store';

export type ErrorResponse = { status?: string; error: string };

export const errors = writable<CSError[]>([]);

export function addError(err: CSError | undefined) {
	if (err) {
		const errs = get(errors);
		errors.update(() => [...errs, err]);
	}
}

export class CSError implements ErrorResponse {
	statusCode: number | undefined = undefined;
	status = '';
	error = '';

	constructor(data?: ErrorResponse | string, statusCode?: number) {
		if (data as string) {
			this.status = data as string;
		} else {
			Object.assign(this, data);
		}
		this.statusCode = statusCode;
	}

	toString(): string {
		return (this.statusCode ? `${this.statusCode}: ` : '') + this.error;
	}
}
