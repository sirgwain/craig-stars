import { addError, CSError, type ErrorResponse } from './Errors';

export abstract class Service {
	// convert an error response to a CSError and add it to the errors store
	// then throw
	static async throwError(response: Response) {
		let result: ErrorResponse | undefined;
		try {
			result = await response.json();
		} catch {
			// ignore, there might not be a body
		}
		const err = new CSError(result, response.statusText, response.status);
		addError(err);
		console.error(err.status, err);
		throw err;
	}

	static async get<T>(url: string, body?: BodyInit): Promise<T> {
		const response = await fetch(url, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			},
			body: body
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as T;
	}
	static async create<T>(item: T, url: string): Promise<T> {
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(item)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as T;
	}
	static async update<T>(item: T, url: string): Promise<T> {
		const response = await fetch(url, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(item)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as T;
	}

	static async delete(url: string): Promise<void> {
		const response = await fetch(`${url}`, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return Promise.resolve();
	}
}
