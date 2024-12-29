import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserInformation = defineStore('userInformation',{
  state: () => ({
    username: '',
    companions: [],
    conversationId: 0,
    conversationPhoto: '',
    conversationName: '',
    messages: [],
  }),
  actions: {
    async fetchCompanions() {
      try {
        const response = await axios.get('api/companions')
        this.companions = response.data
      } catch (error) {console.error('Error fetching companions:', error)}
    },
    async fetchMessages() {
      try{
        const response = await axios.get('api/messages', {params: {conversationId: this.conversationId}});
        this.messages = response.data
      } catch (error) {console.error('Error fetching companions:', error)}
    }
  }
})
