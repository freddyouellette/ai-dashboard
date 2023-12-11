import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const messagesSlice = createSlice({
	name: 'messages',
	initialState: {
		messages: {},
		messagesLoading: false,
		messagesError: null,
		waitingForResponse: false,
	},
	reducers: {
		addMessage: (state, action) => {
			state.messages.push(action.payload)
		},
		setMessages: (state, action) => {
			state.messages = action.payload
		},
		setMessagesLoading: (state, action) => {
			state.messagesLoading = action.payload
		},
		setMessagesError: (state, action) => {
			state.messagesError = action.payload
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
export const getChatMessages = chat => async dispatch => {
	dispatch(messagesSlice.actions.setMessages({}));
	dispatch(messagesSlice.actions.setMessagesLoading(true));
	dispatch(messagesSlice.actions.setMessagesError(null));
	axios.get(`http://localhost:8080/api/chats/${chat.ID}/messages`)
	.then(response => {
		console.log(response)
		dispatch(messagesSlice.actions.setMessagesLoading(false));
		dispatch(messagesSlice.actions.setMessagesError(null));
		dispatch(messagesSlice.actions.setMessages(response.data))
	}, error => {
		dispatch(messagesSlice.actions.setMessagesLoading(false));
		dispatch(messagesSlice.actions.setMessagesError(error));
		console.error(error)
	});
}

export const selectMessages = createSelector(
	state => state.messages.messages,
	state => state.messages.messagesLoading,
	state => state.messages.messagesError,
	(messages, messagesLoading, messagesError) => ({ messages, messagesLoading, messagesError })
);
export const selectWaitingForResponse = state => state.messages.waitingForResponse

export default messagesSlice.reducer;