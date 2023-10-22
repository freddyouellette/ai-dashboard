import { createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const messagesSlice = createSlice({
	name: 'messages',
	initialState: {
		messages: [],
	},
	reducers: {
		addMessage: (state, action) => {
			state.messages.push(action.payload)
		},
		setMessages: (state, action) => {
			state.messages = action.payload
		}
	}
});

// thunk
export const sendMessage = (chatId, message) => async dispatch => {
	let newMessageData = {
		chat_id: chatId,
		text: message,
	}
	axios.post('http://localhost:8080/messages', newMessageData)
	.then(response => {
		console.log(response)
		dispatch(messagesSlice.actions.addMessage(response.data))
		axios.get("http://localhost:8080/chats/"+chatId+"/response")
		.then(response => {
			console.log(response)
			dispatch(messagesSlice.actions.addMessage(response.data))
		}, error => console.error(error))
	}, error => console.error(error))
}

// thunk
export const fetchMessages = () => async dispatch => {
	axios.get('http://localhost:8080/messages')
	.then(response => {
		console.log(response)
		dispatch(messagesSlice.actions.setMessages(response.data))
	}, error => console.error(error))
}

export const selectMessages = state => state.messages.messages

export default messagesSlice.reducer;