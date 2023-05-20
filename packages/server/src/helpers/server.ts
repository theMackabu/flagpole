import got from 'got';
import { log } from '../logger';
import { cache } from '../api/database';
import { serve } from '@hono/node-server';
import { urls, agent, args } from '../api/objects';

const startServer = async (app) => {
	const data: any = await got
		.get(urls.list, agent)
		.json()
		.catch((err) => log.error(err));
	const items = Object.entries(data.items);

	items.forEach(([key, value]: any) => cache.set(key, value));
	log.info({ items: items.length }, 'fetched inital cache');
	log.info({ port: args.port }, 'server started');

	serve({
		fetch: app.fetch,
		port: args.port,
	});
};

export default startServer;
