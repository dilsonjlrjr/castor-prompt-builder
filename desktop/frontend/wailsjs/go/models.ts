export namespace main {
	
	export class StepDTO {
	    titulo: string;
	    descricao: string;
	
	    static createFrom(source: any = {}) {
	        return new StepDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.titulo = source["titulo"];
	        this.descricao = source["descricao"];
	    }
	}
	export class GapAnswerDTO {
	    field_id?: string;
	    pergunta: string;
	    resposta: string;
	    role_nome?: string;
	
	    static createFrom(source: any = {}) {
	        return new GapAnswerDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.field_id = source["field_id"];
	        this.pergunta = source["pergunta"];
	        this.resposta = source["resposta"];
	        this.role_nome = source["role_nome"];
	    }
	}
	export class BuildRequestDTO {
	    model_id: string;
	    role_ids: string[];
	    narrativa: string;
	    gap_answers: GapAnswerDTO[];
	    steps: StepDTO[];
	
	    static createFrom(source: any = {}) {
	        return new BuildRequestDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model_id = source["model_id"];
	        this.role_ids = source["role_ids"];
	        this.narrativa = source["narrativa"];
	        this.gap_answers = this.convertValues(source["gap_answers"], GapAnswerDTO);
	        this.steps = this.convertValues(source["steps"], StepDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BuildResultDTO {
	    conteudo: string;
	    caminho: string;
	    erro?: string;
	
	    static createFrom(source: any = {}) {
	        return new BuildResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.conteudo = source["conteudo"];
	        this.caminho = source["caminho"];
	        this.erro = source["erro"];
	    }
	}
	export class CampoDTO {
	    id: string;
	    label: string;
	    tipo: string;
	    obrigatorio: boolean;
	    opcoes?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CampoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.tipo = source["tipo"];
	        this.obrigatorio = source["obrigatorio"];
	        this.opcoes = source["opcoes"];
	    }
	}
	export class FileStatus {
	    arquivo: string;
	    tipo: string;
	    ok: boolean;
	    problema?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.arquivo = source["arquivo"];
	        this.tipo = source["tipo"];
	        this.ok = source["ok"];
	        this.problema = source["problema"];
	    }
	}
	
	export class ModelDTO {
	    id: string;
	    nome: string;
	    descricao: string;
	    campos: CampoDTO[];
	
	    static createFrom(source: any = {}) {
	        return new ModelDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nome = source["nome"];
	        this.descricao = source["descricao"];
	        this.campos = this.convertValues(source["campos"], CampoDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RoleDTO {
	    id: string;
	    nome: string;
	    categoria: string;
	    tom: string;
	    gaps_comuns: string[];
	    habilidades: string[];
	
	    static createFrom(source: any = {}) {
	        return new RoleDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nome = source["nome"];
	        this.categoria = source["categoria"];
	        this.tom = source["tom"];
	        this.gaps_comuns = source["gaps_comuns"];
	        this.habilidades = source["habilidades"];
	    }
	}

}

