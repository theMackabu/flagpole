const log = {
	info: (object: any) => console.log(JSON.stringify({ level: 'info', name: 'server', time: new Date(), message: object })),
	error: (object: any) => console.error(JSON.stringify({ level: 'error', name: 'server', time: new Date(), message: object })),
	server: (method: string, path: string, status?: number, elapsed?: string) => {
		console.log(JSON.stringify({ level: 'info', name: 'server', time: new Date(), method, path: getPathFromURL(path), status, elapsed }));
	},
};

const time = (start: number) => (delta = Date.now() - start);

const getPathFromURL = (url: string, strict: boolean = true): string => {
	const queryIndex = url.indexOf('?', 8);
	const result = url.substring(url.indexOf('/', 8), queryIndex === -1 ? url.length : queryIndex);

	if (strict === false && /.+\/$/.test(result)) {
		return result.slice(0, -1);
	}

	return result;
};

const logger = () => {
	return async (c, next) => {
		const { method } = c.req;
		log.server(method, c.req.url);

		const start = Date.now();
		await next();

		log.server(method, c.req.url, c.res.status, time(start));
	};
};

export { log, logger };
