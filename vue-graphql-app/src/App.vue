<script setup>
import UserManagement from './components/UserManagement.vue'
</script>

<template>
  <div id="app">
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0 24px; display: flex; align-items: center;">
        <h1 style="margin: 0; font-size: 20px;">Vue + GraphQL App</h1>
      </a-layout-header>
      <a-layout-content style="padding: 24px; min-height: calc(100vh - 64px); background: #f0f2f5;">
        <UserManagement />
      </a-layout-content>
    </a-layout>
  </div>
</template>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  margin: 0;
  padding: 0;
}
</style>
