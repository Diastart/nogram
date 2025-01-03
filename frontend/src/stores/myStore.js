import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserInformation = defineStore('userInformation',{
  state: () => ({
    companions: [],
    messages: [],
    groups: [],
    users: [],
    members: [],
    candidates: [],
    conversationId: 0,
    conversationPhoto: '',
    conversationName: '',
  }),
  actions: {
    async fetchCompanions() {
      try {
        const response = await axios.get('api/companions')
        console.log('Fetching companions...');
        this.companions = response.data
      } catch (error) {console.error('Error fetching companions:', error)}
    },
    async fetchMessages() {
      try{
        const response = await axios.get('api/messages', {params: {conversationId: this.conversationId}});
        console.log('Fetching messages...');
        this.messages = response.data
      } catch (error) {console.error('Error fetching messages:', error)}
    },
    async fetchGroups() {
      try{
        const response = await axios.get('api/groups');
        console.log('Fetching groups...');
        this.groups = response.data;
      }catch(error) {console.error('Error fetching groups:', error)}
    },
    async fetchUsers() {
      try{
        const response = await axios.get('api/users');
        console.log('Fetching users...');
        this.users = response.data;
      }catch(error) {console.error('Error fetching users:', error)}
    },
    async fetchMembers(groupId) {
      try{
        const response = await axios.get('api/members', {params: {groupId: groupId}});
        console.log('Fetching members...');
        this.members = response.data;
        this.candidates = this.companions.filter(c => !this.members.some(m => m.id === c.id))
      }catch(error) {console.error('Error fetching members:', error)}
    },
  }
})
