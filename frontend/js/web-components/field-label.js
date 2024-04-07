class FieldLabel extends HTMLElement {
    constructor() {
        super();
        this.classList.add("block", "text-md")
    }
}

customElements.define("c-field-label", FieldLabel);