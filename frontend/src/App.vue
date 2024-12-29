<template>
  <router-view :key="$route.fullPath"/>
</template>

<script>
export default {
  name: 'App',
}
</script>

<style>
html, body {
    margin: 0;
    padding: 0;
    background-color: black;
    color: white;
    min-height: 100vh;
    font-family: 'Impact', fantasy;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
}
</style>
