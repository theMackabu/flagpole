import { hello } from '@/index';

test('hello', () => {
	expect(hello('foo')).toEqual('Hello foo');
});
