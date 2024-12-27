<template>
    <input v-model="companion" placeholder="Find & add companion" required><br>
    <button @click="findCompanion">press</button>
</template>

<script>
import axios from 'axios'
export default {
    name : 'Find',
    data(){
        return {
            companion: '',
        }
    },
    methods: {
        async findCompanion(){
            try {
                const response = await axios.get('api/dialogs', {params: {companion: this.companion}});
                console.log(this.companion + ' is added successfully');
                this.companion = '';
            }catch(error){console.log(error)}
        },
    }

}
</script>