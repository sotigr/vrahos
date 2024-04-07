class Box extends HTMLElement {  
    constructor() { 
        super();
        this.classList.add("bg-white", "block", "rounded-lg", "shadow-xl", "p-5")
    } 
}

customElements.define("c-box", Box);

class BoxInner extends HTMLElement {  
    constructor() { 
        super();
        this.classList.add("bg-slate-300", "block", "rounded-lg", "p-5")
    } 
}

customElements.define("c-box-inner", BoxInner);