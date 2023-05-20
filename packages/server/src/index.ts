import { Hono } from 'hono';
import { jwt } from 'hono/jwt';
import { logger } from './logger';
import { api } from './api/routes';
import { args } from './api/objects';
import { startServer } from './helpers';

const app = new Hono();

app.use('*', logger());
app.use('/api/*', jwt({ secret: args.secret }));
app.route('/', api);

startServer(app);
