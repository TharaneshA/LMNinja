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

export namespace attacks {
	
	export class AttackCategory {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new AttackCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
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
	    gpuLayers?: number;
	
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
	        this.gpuLayers = source["gpuLayers"];
	    }
	}
	export class DashboardStats {
	    modelsScanned: number;
	    totalScans: number;
	    vulnerabilitiesFound: number;
	    overallPassRate: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.modelsScanned = source["modelsScanned"];
	        this.totalScans = source["totalScans"];
	        this.vulnerabilitiesFound = source["vulnerabilitiesFound"];
	        this.overallPassRate = source["overallPassRate"];
	    }
	}
	export class ScanHistoryItem {
	    id: string;
	    targetModelName: string;
	    startTime: string;
	    status: string;
	    totalPrompts: number;
	    vulnerabilitiesFound: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanHistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.targetModelName = source["targetModelName"];
	        this.startTime = source["startTime"];
	        this.status = source["status"];
	        this.totalPrompts = source["totalPrompts"];
	        this.vulnerabilitiesFound = source["vulnerabilitiesFound"];
	    }
	}
	export class ScanResultItem {
	    id: number;
	    prompt: string;
	    response: string;
	    evaluationJson: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResultItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.prompt = source["prompt"];
	        this.response = source["response"];
	        this.evaluationJson = source["evaluationJson"];
	    }
	}
	export class VulnerabilityByModel {
	    modelName: string;
	    vulnCount: number;
	
	    static createFrom(source: any = {}) {
	        return new VulnerabilityByModel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.modelName = source["modelName"];
	        this.vulnCount = source["vulnCount"];
	    }
	}

}

