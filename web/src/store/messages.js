import { createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const messagesSlice = createSlice({
	name: 'messages',
	initialState: {
		messages: [],
		waitingForResponse: false,
	},
	reducers: {
		addMessage: (state, action) => {
			state.messages.push(action.payload)
		},
		setMessages: (state, action) => {
			state.messages = action.payload
		},
		setWaitingForResponse: (state, action) => {
			state.waitingForResponse = action.payload
		},
	}
});

// thunk
export const sendMessage = (chatId, message) => async dispatch => {
	let newMessageData = {
		chat_id: chatId,
		text: message,
	}
	axios.post('http://localhost:8080/api/messages', newMessageData)
	.then(response => {
		console.log(response)
		dispatch(messagesSlice.actions.addMessage(response.data))
		dispatch(messagesSlice.actions.setWaitingForResponse(true))
		axios.get("http://localhost:8080/api/chats/"+chatId+"/response")
		.then(response => {
			console.log(response)
			dispatch(messagesSlice.actions.setWaitingForResponse(false))
			dispatch(messagesSlice.actions.addMessage(response.data))
		}, error => console.error(error))
	}, error => console.error(error))
}

// thunk
export const fetchMessages = () => async dispatch => {
	axios.get('http://localhost:8080/api/messages')
	.then(response => {
		console.log(response)
		dispatch(messagesSlice.actions.setMessages(response.data))
	}, error => console.error(error))
}

export const selectMessages = state => state.messages.messages
export const selectWaitingForResponse = state => state.messages.waitingForResponse

export default messagesSlice.reducer;