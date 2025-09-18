export namespace app {
	
	export class AppInfo {
	    version: string;
	    os: string;
	    arch: string;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	    }
	}
	export class GGUFFile {
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new GGUFFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}

}

export namespace storage {
	
	export class ConnectionMetadata {
	    id: string;
	    name: string;
	    group?: string;
	    provider: string;
	    model: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.group = source["group"];
	        this.provider = source["provider"];
	        this.model = source["model"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

