<script>
  import Counter from "../components/Counter.svelte";
import DefaultLayout from "../layouts/DefaultLayout.svelte";

  export let count = 0;
</script>

<svelte:head>
  <title>Welcome to your wx app!</title>
</svelte:head>

<DefaultLayout>
  <h1>Welcome to your wx app!</h1>
  <h2><b>wx</b> is an open-source framework that allows you to quickly build <b>Go</b> apps that use <b>Svelte</b> for rendering web pages.</h2>
  <div class="counter-wrapper">
    <Counter {count}/>
  </div>
</DefaultLayout>

<style>
  h1 {
    font-size: 3rem;
    font-weight: 500;
    text-align: center;
  }
  h2 {
    font-size: 2rem;
    text-align: center;
    font-weight: 400;
  }
  .counter-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: row;
  }
</style>
