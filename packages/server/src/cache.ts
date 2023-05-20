class Cache extends Map {
	delete = (key) => super.delete(this.preProcess(key));
	get = (key) => super.get(this.preProcess(key));
	has = (key) => super.has(this.preProcess(key));
	set = (key, value) => super.set(this.preProcess(key), value);
	preProcess = (key) => (typeof key === 'object' ? JSON.stringify(key) : key);
}

export { Cache };
