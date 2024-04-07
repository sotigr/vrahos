
colors = {
    blue: (b) => b.classList.add("bg-blue-600", "hover:bg-blue-700", "active:bg-blue-800", "color-white"),
    red: (b) => b.classList.add("bg-red-600", "hover:bg-red-700", "active:bg-red-800", "color-white"),
    gray: (b) => b.classList.add("bg-gray-600", "hover:bg-gray-700", "active:bg-gray-800", "color-black"),
}

class Button extends HTMLElement { 
    constructor() {
        super();
        buttonContractor(this)
    }
} 

class ButtonExtended extends HTMLButtonElement {
    constructor() {
        super();
        buttonContractor(this) 
    }
}

function buttonContractor(b) {
    b.classList.add("inline-block", "cursor-pointer", "text-white", "px-5", "py-2", "rounded-3xl", "select-none")

    let color = b.getAttribute("color");
    if (!color) {
        color = "blue"
    }
    if (colors[color]) {
        colors[color](b)
    }
} 
 

customElements.define("c-button", Button); 
customElements.define("ce-button", ButtonExtended, {extends: "button"}); 