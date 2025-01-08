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
    latestMessages: [],
    latestMessagesOfGroups: [],
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
        console.log('Messages response:', response.data); 
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
    async fetchLatestMessages() {
      try {
          console.log('Fetching latest messages...');
          const response = await axios.post('api/latest/messages', this.companions);
          console.log('Latest messages response:', response.data);
          this.latestMessages = response.data;
      } catch (error) {console.error('Error fetching latest messages:', error);}
    },
    async fetchLatestMessagesOfGroups() {
      try {
          if (!this.groups || this.groups.length === 0) {
              console.log('No groups to fetch messages for');
              return;
          }
          const response = await axios.post('api/latest/messages/groups', this.groups);
          this.latestMessagesOfGroups = response.data;
      } catch (error) {console.error('Error fetching group messages:', error);}
    }
  }
})
