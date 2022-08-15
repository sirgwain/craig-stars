export abstract class Service {
	protected async get<T>(url: string, body?: BodyInit): Promise<T> {
		const response = await fetch(url, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			},
			body: body
		});

		if (response.ok) {
			return (await response.json()) as T;
		} else {
			console.error(response);
			throw new Error(`${response}`);
		}
	}
	protected async update<T>(item: T, url: string): Promise<T> {
		const response = await fetch(url, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(item)
		});

		if (response.ok) {
			return (await response.json()) as T;
		} else {
			console.error(response);
		}
		return Promise.resolve(item);
	}


}
