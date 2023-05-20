import got from 'got';
import { Hono } from 'hono';
import { log } from '../logger';
import { pathBuilder } from '../helpers';
import { cacheHandler } from './database';
import { notFound, urls, config, args } from './objects';

const api = new Hono();

api.notFound((c) => c.json(notFound, 404));

api.post('/auth/login', async (c) => {
	const body = await c.req.json();
	const response: any = await got
		.post(urls.login, config(body))
		.json()
		.catch((err) => log.error(err));

	response.error != null && c.status(401);
	return c.json(response);
});

api.post('/api/user/refresh', async (c) => {
	const body = await c.req.json();
	const token = c.req.header('Authorization');

	const response: any = await got
		.post(urls.refresh, {
			...config(body),
			headers: {
				Authorization: token,
			},
		})
		.json()
		.catch((err) => log.error(err));

	response.error != null && c.status(401);
	return c.json(response);
});

api.get('/client/api/flag/:id', async (c) => {
	const key = c.req.query('key');
	const id = c.req.param('id');
	const response: any = cacheHandler(id, key);

	response != null && response.error != null && c.status(404);
	return response != null ? c.text(response) : c.json(notFound, 404);
});

api.post('/api/user/create', async (c) => {
	const body = await c.req.json();
	const response: any = await got
		.post(urls.create, config(body))
		.json()
		.catch((err) => log.error(err));

	response.created == false && c.status(400);
	return c.json(response);
});

api.post('/api/flag/create', async (c) => {
	const { environment, flag, body } = await c.req.json();
	const builder = pathBuilder(urls.database, `${environment}.flags.${flag}`);

	const response: any = await got
		.put(builder, config(body, true))
		.json()
		.catch((err) => log.error(err));

	response.success == false && c.status(400);
	return c.json(response);
});

export { api };
