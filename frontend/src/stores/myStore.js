import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserInformation = defineStore('userInformation',{
  state: () => ({
    username : '',
    companions: [],
    conversationId: 0,
    conversationPhoto: '',
    conversationName: '',
  }),
  actions: {
    async fetchCompanions() {
      try {
        const response = await axios.get('api/companions')
        this.companions = response.data
      } catch (error) {
        console.error('Error fetching companions:', error)
      }
    },
  }
})
