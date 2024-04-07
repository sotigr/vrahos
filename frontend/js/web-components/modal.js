class CustomModal extends HTMLElement {
    constructor() {
        super();
        this._modalVisible = false;
        this._modal;
        this.attachShadow({ mode: 'open' });
        this.shadowRoot.innerHTML = `
        <style>
            .__modal {
                display: none;
                position: fixed;
                z-index: 1;
                left: 0; top: 0;
                width: 100%; height: 100%;
                overflow: auto;
                background-color: rgba(0,0,0,0.4); 
                align-items: center;
                justify-content: center;
            }
        </style>
		<div class="__modal"> 
            <slot name="content" /> 
		</div>
		`
    }
 
	connectedCallback() {
		this._modal = this.shadowRoot.querySelector(".__modal");
 	}

    show() {
        this._modalVisible = true;
        this._modal.style.display = 'flex';
    }
    
    hide() {
        this._modalVisible = false;
        this._modal.style.display = 'none'; 
    } 
}
customElements.define('c-modal', CustomModal);