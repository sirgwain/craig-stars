import { CSError, type ErrorResponse } from './Errors';

export abstract class Service {
	static async get<T>(url: string, body?: BodyInit): Promise<T> {
		const response = await fetch(url, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			},
			body: body
		});

		if (!response.ok) {
			const err = new CSError((await response.json()) as ErrorResponse, response.status);
			console.error(err);
			throw err;
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
			const err = new CSError((await response.json()) as ErrorResponse, response.status);
			console.error(err);
			throw err;
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
			const err = new CSError((await response.json()) as ErrorResponse, response.status);
			console.error(err);
			throw err;
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
			const err = new CSError((await response.json()) as ErrorResponse, response.status);
			console.error(err);
			throw err;
		}
		return Promise.resolve();
	}
}
