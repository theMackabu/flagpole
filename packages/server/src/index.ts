import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { log } from './logger';
import got from 'got';
import { jwt } from 'hono/jwt';
import { parse } from './cli';
import { logger } from 'hono/logger';

class Cache extends Map {
	constructor(array) {
		super(array);
	}
	delete(key) {
		return super.delete(this.preProcess(key));
	}
	get(key) {
		return super.get(this.preProcess(key));
	}
	has(key) {
		return super.has(this.preProcess(key));
	}
	set(key, value) {
		return super.set(this.preProcess(key), value);
	}

	preProcess(key) {
		if (typeof key === 'object') {
			return JSON.stringify(key);
		} else {
			return key;
		}
	}
}

const app = new Hono();
const cache = new Cache();
const args = parse(process.argv.slice(2));

const notFound = {
	error: {
		code: 404,
		message: 'not found',
	},
};

const throttle = (callback, delay) => {
	let throttleTimeout = null;
	let storedEvent = null;

	const throttledEventHandler = (event) => {
		storedEvent = event;
		const shouldHandleEvent = !throttleTimeout;

		if (shouldHandleEvent) {
			callback(storedEvent);

			storedEvent = null;
			throttleTimeout = setTimeout(() => {
				throttleTimeout = null;
				if (storedEvent) {
					throttledEventHandler(storedEvent);
				}
			}, delay);
		}
	};

	return throttledEventHandler;
};

const fetch = throttle(async (data) => {
	const { address, builder } = data;
	const response = await got
		.get(address, { throwHttpErrors: false })
		.json()
		.catch((err) => log.error(err));

	cache.set(builder, response);
}, 10 * 1000);

const startServer = async () => {
	const data = await got
		.get(urls.list)
		.json()
		.catch((err) => log.error(err));
	const items = Object.keys(data.items);

	items.forEach((key) => cache.set(key, JSON.parse(data.items[key])));
	log.info(`fetched inital cache, ${items.length >= 1 ? items.length + ' item' : items.length + ' items'}`);
	log.info(`started on port ${args.port}`);

	serve({
		fetch: app.fetch,
		port: args.port,
	});
};

const cacheHandler = (id, key) => {
	const address = urls.database + `/${id}.flags.${key}`;
	const builder = id + '.flags.' + key;

	fetch({ address, builder });
	return cache.get(builder);
};

const urls = {
	login: `http://${args.auth}/_/login`,
	create: `http://${args.auth}/_/signup`,
	refresh: `http://${args.auth}/_/refresh`,
	database: `http://${args.db_url}/${args.database}`,
	list: `http://${args.db_url}/list`,
};

const config = (body) => ({ json: body, throwHttpErrors: false });

app.use('*', logger());
app.use('/api/*', jwt({ secret: args.secret }));

app.notFound((c) => c.json(notFound, 404));

app.post('/auth/login', async (c) => {
	const body = await c.req.json();
	const response = await got
		.post(urls.login, config(body))
		.json()
		.catch((err) => log.error(err));

	response.error != null && c.status(401);
	return c.json(response);
});

app.post('/api/user/refresh', async (c) => {
	const body = await c.req.json();
	const token = c.req.header('Authorization');

	const response = await got
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

app.get('/client/api/flag/:id', async (c) => {
	const key = c.req.query('key');
	const id = c.req.param('id');
	const response = cacheHandler(id, key);

	response != null && response.error != null && c.status(404);
	return response != null ? c.json(response) : c.json(notFound, 404);
});

app.post('/api/user/create', async (c) => {
	const body = await c.req.json();
	const response = await got
		.post(urls.create, config(body))
		.json()
		.catch((err) => log.error(err));

	response.created == false && c.status(400);
	return c.json(response);
});

startServer();
