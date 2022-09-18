<script>
  import Button from "./Button.svelte";

  export let count = 0;
</script>

<style>
  .counter {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 20rem;
  }
  .counter > .value {
    font-size: 4rem;
    font-weight: 500;
    height: 4rem;
    margin: .8rem;
  } 
  .counter > .actions {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    margin: .8rem;
  }
</style>

<div class="counter">
  <div class="value">{count}</div>
  <div class="actions">
    <Button on:click={() => count--}>-</Button>
    <Button on:click={() => count++}>+</Button>
  </div>
</div>