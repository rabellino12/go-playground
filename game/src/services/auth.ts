export function fetchJWT(username: string): Promise<string> {
	return fetch('http://localhost:8080/auth/jwt', {
		body: JSON.stringify({ username }),
		method: 'POST'
	})
		.then(res => {
			const body = res.json();
			return body;
		})
		.then(body => {
			return body.jwt;
		})
		.catch(err => {
			throw err;
		});
}
