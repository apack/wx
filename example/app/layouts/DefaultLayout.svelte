<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap" rel="stylesheet">
</svelte:head>

<div class="layout">
  <header>
    <img src="/static/logo.svg" alt="wx" />
  </header>
  <main>
    <slot></slot>
  </main> 
  <footer>
    <a target="_blank" href="https://github.com/apack/wx">
        Github
      </a>
  </footer>
</div>

<style>
  :global(html, body) {
    padding: 0;
    margin: 0;
  }
  :global(html) {
    font-size: 10px;
  }
  :global(body) {
    font-family: 'Roboto', sans-serif;
    font-size: 1.6rem;
  }
  :global(body > #app) {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
  }
  :global(a) {
    color: #C72121;
    text-decoration: none;
  }
  :global(a:hover) {
    text-decoration: underline;
  }
  .layout {
    display: flex;
    flex-direction: column;
    max-width: 60rem;
    padding: 2rem;
  }
  header, footer {
    text-align: center;
    margin: 1.6rem 0;
  }
</style>
