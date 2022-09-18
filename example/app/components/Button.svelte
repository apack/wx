<style>
    .button {
        font-size: 2rem;
        font-weight: 500;
        height: 3.2rem;
        padding: 0 1.6rem;
        text-align: center;
        background-color: #C72121;
        border: none;
        color: #FFF;
        border-radius: .3rem;
        cursor: pointer;
    }
    .button:hover {
        background-color: #A61B1B;
    }
    .button:not(:first-child) {
        margin-left: .8rem;
    }
</style>
<button class="button" on:click>
    <slot></slot>
</button>
