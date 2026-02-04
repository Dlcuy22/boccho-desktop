export namespace AnimationEngine {
	
	export class CharacterInfo {
	    name: string;
	    path: string;
	    previewPath: string;
	    frameCount: number;
	
	    static createFrom(source: any = {}) {
	        return new CharacterInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.previewPath = source["previewPath"];
	        this.frameCount = source["frameCount"];
	    }
	}

}

export namespace PackManagement {
	
	export class PackInfo {
	    filePath: string;
	    packName: string;
	    characters: string[];
	    previewImage: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new PackInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.packName = source["packName"];
	        this.characters = source["characters"];
	        this.previewImage = source["previewImage"];
	        this.error = source["error"];
	    }
	}

}

export namespace main {
	
	export class CharacterWindowInfo {
	    id: string;
	    characterName: string;
	    isRunning: boolean;
	    scale: number;
	
	    static createFrom(source: any = {}) {
	        return new CharacterWindowInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.characterName = source["characterName"];
	        this.isRunning = source["isRunning"];
	        this.scale = source["scale"];
	    }
	}

}

