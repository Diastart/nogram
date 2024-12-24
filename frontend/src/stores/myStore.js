import { defineStore } from 'pinia'

export const useUserInformation = defineStore('userInformation',{
  state: () => ({
    username : '',
  }),
})
